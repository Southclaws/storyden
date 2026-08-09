//go:build windows

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"unsafe"

	"golang.org/x/sys/windows"
)

type managedProcess struct {
	cmd *exec.Cmd
	job windows.Handle
}

func newRunContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt)
}

func startManagedProcess(cmd *exec.Cmd) (*managedProcess, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("create process job: %w", err)
	}

	var limits windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err := windows.SetInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&limits)),
		uint32(unsafe.Sizeof(limits)),
	); err != nil {
		windows.CloseHandle(job)
		return nil, fmt.Errorf("configure process job: %w", err)
	}

	if err := cmd.Start(); err != nil {
		windows.CloseHandle(job)
		return nil, err
	}

	process, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(cmd.Process.Pid))
	if err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		windows.CloseHandle(job)
		return nil, fmt.Errorf("open managed process: %w", err)
	}
	defer windows.CloseHandle(process)

	if err := windows.AssignProcessToJobObject(job, process); err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		windows.CloseHandle(job)
		return nil, fmt.Errorf("assign process to job: %w", err)
	}

	return &managedProcess{cmd: cmd, job: job}, nil
}

func stopManagedProcess(name string, process *managedProcess) {
	if process == nil || process.cmd == nil || process.cmd.Process == nil {
		return
	}
	defer windows.CloseHandle(process.job)

	log.Println("Stopping", name, "...")
	if err := windows.TerminateJobObject(process.job, 1); err != nil {
		log.Println("failed to stop", name, "process job:", err)
	}
	if err := process.cmd.Wait(); err != nil {
		log.Println(name, "exited with:", err)
	}
}
