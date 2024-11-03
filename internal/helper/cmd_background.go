//go:build !windows
// +build !windows

package helper

import (
	"os/exec"
)

func RunCmdBackground(cmd *exec.Cmd) {}
