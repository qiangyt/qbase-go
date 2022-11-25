//go:build windows
// +build windows

package comm

import "github.com/pkg/errors"

func RunSudoCommand(vars map[string]string, dir string, cmd string, passwordInput FnInput) (CommandOutput, error) {
	return nil, errors.New("todo")
}

func RunUserCommand(vars map[string]string, dir string, cmd string) (CommandOutput, error) {
	return RunCommandNoInput(vars, dir, "sh", cmd)
}
