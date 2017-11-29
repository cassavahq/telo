// +build linux,amd64

package api

import (
	"os"
	"syscall"
)

var signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGKILL,
	syscall.SIGTERM,
	syscall.SIGSTOP,
	syscall.SIGHUP,
}
