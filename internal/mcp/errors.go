package mcp

import "errors"

var (
	// ErrRequired error.
	ErrRequired = errors.New("mcp: missing required parameter")

	// ErrNotFound error.
	ErrNotFound = errors.New("mcp: not found")

	// ErrInvalid error.
	ErrInvalid = errors.New("mcp: invalid value")

	// ErrFind error.
	ErrFind = errors.New("mcp: failed to find")

	// ErrDelete error.
	ErrDelete = errors.New("mcp: failed to delete")

	// ErrMarshal error.
	ErrMarshal = errors.New("mcp: failed to marshal")
)
