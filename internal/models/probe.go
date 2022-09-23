package models

// Probe type.
type Probe string

// ProbeRepository type.
type ProbeRepository interface {
	FindAll() ([]Probe, error)
}
