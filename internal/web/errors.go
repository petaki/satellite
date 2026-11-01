package web

import "errors"

var (
	// ErrBadStatusCode error.
	ErrBadStatusCode = errors.New("web: bad status code")

	// ErrWebhook error.
	ErrWebhook = errors.New("web: webhook failed")
)
