package config

import "errors"

var (
	// ErrInvalid error.
	ErrInvalid = errors.New("config: invalid value")

	// ErrInvalidHeartbeatWait error.
	ErrInvalidHeartbeatWait = errors.New("config: invalid heartbeat wait")

	// ErrInvalidHeartbeatSleep error.
	ErrInvalidHeartbeatSleep = errors.New("config: invalid heartbeat sleep")

	// ErrInvalidHeartbeatWebhookHeader error.
	ErrInvalidHeartbeatWebhookHeader = errors.New("config: invalid heartbeat webhook header")

	// ErrMissingHeartbeatWebhookURL error.
	ErrMissingHeartbeatWebhookURL = errors.New("config: missing heartbeat webhook url")
)
