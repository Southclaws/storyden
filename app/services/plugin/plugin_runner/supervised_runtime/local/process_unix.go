//go:build !windows

package local

import (
	"os"
	"os/exec"
	"syscall"
)

func configureProcessCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func interruptProcess(cmd *exec.Cmd) error {
	return signalProcessGroup(cmd, os.Interrupt)
}

func killProcess(cmd *exec.Cmd) error {
	return signalProcessGroup(cmd, os.Kill)
}

func signalProcessGroup(cmd *exec.Cmd, signal os.Signal) error {
	if cmd == nil || cmd.Process == nil {
		return os.ErrProcessDone
	}

	pid := cmd.Process.Pid
	if pid <= 0 {
		return cmd.Process.Signal(signal)
	}

	if sig, ok := signal.(syscall.Signal); ok {
		if err := syscall.Kill(-pid, sig); err == nil {
			return nil
		}
	}

	return cmd.Process.Signal(signal)
}
