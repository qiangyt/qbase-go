package comm

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"mvdan.cc/sh/v3/interp"
)

type FnInput func() string

type CommandOutputKind byte

const (
	_ CommandOutputKind = iota
	COMMAND_OUTPUT_KIND_TEXT
	COMMAND_OUTPUT_KIND_VARS
	COMMAND_OUTPUT_KIND_JSON
)

type CommandOutputT struct {
	Kind CommandOutputKind
	Vars map[string]any
	Text string
	Json any
}

type CommandOutput = *CommandOutputT

func ParseCommandOutputP(outputText string) CommandOutput {
	r, err := ParseCommandOutput(outputText)
	if err != nil {
		panic(err)
	}
	return r
}

func ParseCommandOutput(outputText string) (CommandOutput, error) {
	r := &CommandOutputT{Kind: COMMAND_OUTPUT_KIND_TEXT, Text: outputText}

	if strings.HasPrefix(outputText, "$json$\n\n") {
		jsonBody := outputText[len("$json$\n\n"):]

		err := json.Unmarshal([]byte(jsonBody), &r.Json)
		if err != nil {
			return nil, errors.Wrapf(err, "json: %s"+jsonBody)
		}

		r.Kind = COMMAND_OUTPUT_KIND_JSON
		return r, nil
	}

	if strings.HasPrefix(outputText, "$vars$\n\n") {
		varsBody := outputText[len("$vars$\n\n"):]

		r.Vars = Text2Vars(varsBody)
		r.Kind = COMMAND_OUTPUT_KIND_VARS
	}

	return r, nil
}

func Vars2Pair(vars map[string]any) []string {
	if len(vars) == 0 {
		return nil
	}

	r := make([]string, 0, len(vars))
	for k, v := range vars {
		r = append(r, k+"="+cast.ToString(v))
	}
	return r
}

func Text2Vars(text string) map[string]any {
	pairs := Text2Lines(text)
	return Pair2Vars(pairs)
}

func Pair2Vars(pairs []string) map[string]any {
	if len(pairs) == 0 {
		return map[string]any{}
	}

	r := map[string]any{}
	for _, pair := range pairs {
		pair = strings.TrimLeft(pair, "\t \r")
		pos := strings.IndexByte(pair, '=')
		if pos <= 0 {
			continue
		}
		k := pair[:pos]
		if pos == len(pair)-1 {
			r[k] = ""
		} else {
			r[k] = pair[pos+1:]
		}
	}
	return r
}

func openHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return devNull{}, nil
	}
	return interp.DefaultOpenHandler()(ctx, path, flag, perm)
}

func RunShellCommandP(vars map[string]any, dir string, sh string, cmd string, passwordInput FnInput) CommandOutput {
	r, err := RunShellCommand(vars, dir, sh, cmd, passwordInput)
	if err != nil {
		panic(err)
	}
	return r
}

func RunShellCommand(vars map[string]any, dir string, sh string, cmd string, passwordInput FnInput) (CommandOutput, error) {
	if len(sh) == 0 || sh == "gosh" {
		return RunGoshCommand(vars, dir, cmd, passwordInput)
	}

	if IsSudoCommand(cmd) {
		return RunSudoCommand(vars, dir, cmd, passwordInput)
	}
	return RunUserCommand(vars, dir, cmd)
}

func RunUserCommandP(vars map[string]any, dir string, cmd string) CommandOutput {
	r, err := RunUserCommand(vars, dir, cmd)
	if err != nil {
		panic(err)
	}
	return r
}

/*
func RunShellScriptFile(afs afero.Fs, url string, credentials comm.Credentials, timeout time.Duration,
	dir string, sh string) string {

	scriptContent := comm.DownloadText(afs, url, credentials, timeout)
	return RunShellCommand(dir, sh, scriptContent)
}*/

func newExecCommand(vars map[string]any, dir string, cmd string, args ...string) (*exec.Cmd, error) {
	r := exec.Command(cmd, args...)
	env, err := EnvironList(vars)
	if err != nil {
		return nil, err
	}
	r.Env = env
	r.Dir = dir
	return r, nil
}

func RunCommandNoInputP(vars map[string]any, dir string, cmd string, args ...string) CommandOutput {
	r, err := RunCommandNoInput(vars, dir, cmd, args...)
	if err != nil {
		panic(err)
	}
	return r
}

func RunCommandNoInput(vars map[string]any, dir string, cmd string, args ...string) (CommandOutput, error) {
	_cmd, err := newExecCommand(vars, dir, cmd, args...)
	if err != nil {
		return nil, err
	}

	b, err := _cmd.Output()
	if err != nil {
		cli := strings.Join(append([]string{cmd}, args...), " ")
		return nil, errors.Wrapf(err, "get output for command '%s'", cli)
	}

	return ParseCommandOutput(cast.ToString(b))
}

func RunCommandWithInput(vars map[string]any, dir string, cmd string, args ...string) func(...string) (CommandOutput, error) {
	return func(input ...string) (CommandOutput, error) {
		if IsSudoCommand(cmd) && len(input) > 0 {
			cmd = InstrumentSudoCommand(cmd)
		}

		cli := cmd + " " + strings.Join(args, " ")

		_cmd, err := newExecCommand(vars, dir, cmd, args...)
		if err != nil {
			return nil, err
		}

		stdin, err := _cmd.StdinPipe()
		if err != nil {
			return nil, errors.Wrapf(err, "open stdin for command '%s'", cli)
		}
		defer func() {
			if stdin != nil {
				stdin.Close()
				stdin = nil
			}
		}()

		_, err = io.WriteString(stdin, strings.Join(input, " "))
		if err != nil {
			return nil, errors.Wrap(err, "write something to stdin")
		}
		stdin.Close()
		stdin = nil

		b, err := _cmd.Output()
		if err != nil {
			return nil, errors.Wrapf(err, "get output for command '%s'", cli)
		}

		return ParseCommandOutput(cast.ToString(b))
	}
}

func IsSudoCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "sudo ")
}

func InputSudoPassword(passwordInput FnInput) string {
	if passwordInput == nil {
		return ""
	}

	r := passwordInput()
	if len(r) == 0 {
		return ""
	}
	return r
}

func InstrumentSudoCommand(cmd string) string {
	if !IsSudoCommand(cmd) {
		return cmd
	}

	if strings.Contains(cmd, "-S") || strings.Contains(cmd, "--stdin") {
		return cmd
	}

	noSudoCmd := cmd[len("sudo "):]
	return "sudo --stdin " + noSudoCmd
}
