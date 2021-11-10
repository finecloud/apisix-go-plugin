package common

import "errors"

var (
	ErrConnClosed = errors.New("The connection is closed")
)
