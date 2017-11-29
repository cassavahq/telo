// +build windows,amd64

package api

import (
	"os"
	"syscall"
)

var signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
