//go:build linux
// +build linux

package comm

func RunSudoCommand(vars map[string]string, dir string, cmd string, passwordInput FnInput) (CommandOutput, error) {
	password := passwordInput()
	return RunCommandWithInput(vars, dir, "sh", cmd)(password)
}

func RunUserCommand(vars map[string]string, dir string, cmd string) (CommandOutput, error) {
	return RunCommandNoInput(vars, dir, "sh", cmd)
}
