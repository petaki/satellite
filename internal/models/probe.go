package models

import "time"

// Probe type.
type Probe string

// ProbeSummary type.
type ProbeSummary struct {
	Name      string  `json:"name"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Load1     float64 `json:"load1"`
	Load5     float64 `json:"load5"`
	Load15    float64 `json:"load15"`
	HasBeat   bool    `json:"hasBeat"`
	CPUAlarm  float64 `json:"cpuAlarm"`
	MemAlarm  float64 `json:"memAlarm"`
	LoadAlarm float64 `json:"loadAlarm"`
}

// ProbeRepository type.
type ProbeRepository interface {
	FindAll() ([]Probe, error)
	FindLatestValues(Probe, int) ([]any, *time.Time, error)
	HasHeartbeat(Probe) (bool, error)
	SetHeartbeat(Probe, int) error
	Delete(Probe) error
}
