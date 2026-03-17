package service

import (
	"net/http"
	"time"
)

// HTTPClient function.
func HTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
