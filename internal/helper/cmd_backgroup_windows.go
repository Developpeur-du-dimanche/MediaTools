//go:build windows
// +build windows

package helper

import (
	"os/exec"
	"syscall"
)

func RunCmdBackground(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}
