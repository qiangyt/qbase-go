package comm

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func RunGoshCommandP(vars map[string]any, dir string, cmd string, passwordInput FnInput) CommandOutput {
	r, err := RunGoshCommand(vars, dir, cmd, passwordInput)
	if err != nil {
		panic(err)
	}
	return r
}

func RunGoshCommand(vars map[string]any, dir string, cmd string, passwordInput FnInput) (CommandOutput, error) {
	var stdin io.Reader
	if IsSudoCommand(cmd) {
		password := InputSudoPassword(passwordInput)
		if len(password) > 0 {
			stdin = strings.NewReader(password + "\n")
			cmd = InstrumentSudoCommand(cmd)
		}
	}

	sf, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return nil, errors.Wrapf(err, "parse command: \n%s", cmd)
	}

	out := strings.Builder{}

	envList, err := EnvironList(vars)
	if err != nil {
		return nil, err
	}
	environ := append(os.Environ(), envList...)

	opts := []interp.RunnerOption{
		interp.Params("-e"),
		interp.Env(expand.ListEnviron(environ...)),
		interp.ExecHandler(GoshExecHandler(6 * time.Second)),
		interp.OpenHandler(openHandler),
		interp.StdIO(stdin, &out, &out),
	}
	if len(dir) > 0 {
		opts = append(opts, interp.Dir(dir))
	}

	var runner *interp.Runner
	if runner, err = interp.New(opts...); err != nil {
		return nil, errors.Wrapf(err, "create runner for command: \n%s", cmd)
	}

	if err = runner.Run(context.TODO(), sf); err != nil {
		return nil, errors.Wrapf(err, "run command: \n%s", cmd)
	}

	return ParseCommandOutput(out.String())
}
