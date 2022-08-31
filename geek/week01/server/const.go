package server

import (
	"os"
	"runtime"
	"syscall"
	"time"
)

const (
	DefaultShutdownTimeout = time.Second * 30
	DefaultWaitTimeout     = time.Second * 10
	DefaultCbTimeout       = time.Second * 3
)

// Signals 信号数组map
var Signals = map[string][]os.Signal{
	"darwin": {os.Interrupt, os.Kill,
		syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGILL, syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM},
	"linux": {os.Interrupt, os.Kill,
		syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL,
		syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM},
	"windows": {os.Interrupt, os.Kill,
		syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGILL, syscall.SIGABRT, syscall.SIGTERM},
}[runtime.GOOS]
