package utils

import "fmt"

var (
	ErrRegisterExistedUsername = fmt.Errorf("ErrRegisterExistedUsername")
	ErrRegisterExistedEmail    = fmt.Errorf("ErrRegisterExistedEmail")
	ErrDAOQueryFailed          = fmt.Errorf("ErrDAOQueryFailed")
)
