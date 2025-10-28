//go:build !windows
// +build !windows

package proc

import (
	"syscall"
)

func ProcAttributes() *syscall.SysProcAttr {
	return nil
}
