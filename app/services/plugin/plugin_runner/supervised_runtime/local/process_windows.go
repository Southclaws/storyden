//go:build windows

package local

import (
	"os"
	"os/exec"
)

func configureProcessCommand(cmd *exec.Cmd) {}

func interruptProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return os.ErrProcessDone
	}
	return cmd.Process.Signal(os.Interrupt)
}

func killProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return os.ErrProcessDone
	}
	return cmd.Process.Kill()
}
