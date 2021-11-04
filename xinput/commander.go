package xinput

import "os/exec"

type Commander interface {
	Command(name string, arg ...string)
	Output() ([]byte, error)
}

type ExecCommander struct {
	cmd *exec.Cmd
}

func NewExecCommander() Commander {
	return &ExecCommander{}
}

func (e *ExecCommander) Command(name string, arg ...string) {
	e.cmd = exec.Command(name, arg...)
}

func (e *ExecCommander) Output() ([]byte, error) {
	return e.cmd.Output()
}
