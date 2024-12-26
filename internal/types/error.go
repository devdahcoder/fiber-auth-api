package types

import "fmt"

var (
	ErrInvalidInput  = fmt.Errorf("invalid input parameters")
	ErrDuplicateUser = fmt.Errorf("user already exists")
)