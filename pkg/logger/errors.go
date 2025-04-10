package logger

import "fmt"

var (
	ErrSetLogLevelFailed = fmt.Errorf("ErrSetLogLevelFailed")
	ErrInitLoggerFailed  = fmt.Errorf("ErrInitLoggerFailed")
)
