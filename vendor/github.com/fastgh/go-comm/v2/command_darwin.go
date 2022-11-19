//go:build darwin
// +build darwin

package comm

import (
	"fmt"
	"strings"
)

func RunSudoCommand(vars map[string]any, dir string, cmd string, passwordInput FnInput) (CommandOutput, error) {
	return RunAppleScript(vars, passwordInput(), dir, cmd)
}

func RunAppleScriptP(vars map[string]any, adminPassword string, dir string, script string) CommandOutput {
	r, err := RunAppleScript(vars, adminPassword, dir, script)
	if err != nil {
		panic(err)
	}
	return r
}

func RunAppleScript(vars map[string]any, adminPassword string, dir string, script string) (CommandOutput, error) {
	subArgs := []string{fmt.Sprintf(`do shell script "%s"`, script)}

	if len(adminPassword) > 0 {
		subArgs = append(subArgs, fmt.Sprintf(`password "%s"`, adminPassword))
	}
	subArgs = append(subArgs, "with administrator privileges")

	return RunCommandNoInput(vars, dir, "osascript", "-e", strings.Join(subArgs, " "))
}

func RunUserCommand(vars map[string]any, dir string, cmd string) (CommandOutput, error) {
	return RunCommandNoInput(vars, dir, "open", cmd)
}
