package fake

import "errors"

var (
	// ErrUnknownCommand error.
	ErrUnknownCommand = errors.New("fake: unknown command")
)
