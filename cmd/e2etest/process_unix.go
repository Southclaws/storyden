//go:build darwin || linux

package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const processStopTimeout = 10 * time.Second

type managedProcess struct {
	cmd *exec.Cmd
}

func newRunContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

func startManagedProcess(cmd *exec.Cmd) (*managedProcess, error) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &managedProcess{cmd: cmd}, nil
}

func stopManagedProcess(name string, process *managedProcess) {
	stopManagedProcessWithin(name, process, processStopTimeout)
}

func stopManagedProcessWithin(name string, process *managedProcess, timeout time.Duration) {
	if process == nil || process.cmd == nil || process.cmd.Process == nil {
		return
	}

	log.Println("Stopping", name, "...")

	processGroupID := process.cmd.Process.Pid
	if err := syscall.Kill(-processGroupID, syscall.SIGTERM); err != nil && err != syscall.ESRCH {
		log.Println("failed to stop", name, "process group:", err)
	}

	done := make(chan error, 1)
	go func() { done <- process.cmd.Wait() }()

	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	var waitErr error
	for processGroupRunning(processGroupID) {
		select {
		case err := <-done:
			waitErr = err
			done = nil
		case <-deadline.C:
			log.Println(name, "did not stop in time; killing")
			if err := syscall.Kill(-processGroupID, syscall.SIGKILL); err != nil && err != syscall.ESRCH {
				log.Println("failed to kill", name, "process group:", err)
			}
			if done != nil {
				waitErr = <-done
			}
			if waitErr != nil {
				log.Println(name, "exited with:", waitErr)
			}
			return
		case <-ticker.C:
		}
	}

	if done != nil {
		waitErr = <-done
	}
	if waitErr != nil {
		log.Println(name, "exited with:", waitErr)
	}
}

func processGroupRunning(processGroupID int) bool {
	err := syscall.Kill(-processGroupID, 0)
	return err == nil || err == syscall.EPERM
}
