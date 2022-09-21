package models

// Alarm model.
type Alarm struct {
	CPU    float64 `redis:"cpu"`
	Memory float64 `redis:"memory"`
	Disk   float64 `redis:"disk"`
}

// AlarmRepository type.
type AlarmRepository interface {
	Find() (*Alarm, error)
}
