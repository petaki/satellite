package models

import "errors"

var (
	// ErrNoRecord error.
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidLimit error.
	ErrInvalidLimit = errors.New("models: invalid limit")
)
