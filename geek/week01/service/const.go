package service

import (
	"time"
)

const (
	DefaultShutdownTimeout = time.Second * 30
	DefaultWaitTimeout     = time.Second * 10
	DefaultCbTimeout       = time.Second * 3
)
