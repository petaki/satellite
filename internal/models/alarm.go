package models

// Alarm type.
type Alarm struct {
	CPU    float64 `redis:"cpu"`
	Memory float64 `redis:"memory"`
	Disk   float64 `redis:"disk"`
	Load   float64 `redis:"load"`
}

// AlarmRepository type.
type AlarmRepository interface {
	Find(Probe) (*Alarm, error)
}
