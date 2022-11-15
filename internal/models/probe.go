package models

import "time"

// Probe type.
type Probe string

// ProbeRepository type.
type ProbeRepository interface {
	FindAll() ([]Probe, error)
	FindLatestValues(Probe, int) ([]interface{}, *time.Time, error)
	HasHeartbeat(Probe) (bool, error)
	SetHeartbeat(Probe, int) error
	Delete(Probe) error
}
