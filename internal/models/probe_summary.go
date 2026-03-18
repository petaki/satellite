package models

// ProbeSummary type.
type ProbeSummary struct {
	Name      string  `json:"name"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Load1     float64 `json:"load1"`
	Load5     float64 `json:"load5"`
	Load15    float64 `json:"load15"`
	IsActive  bool    `json:"isActive"`
	CPUAlarm  float64 `json:"cpuAlarm"`
	MemAlarm  float64 `json:"memAlarm"`
	LoadAlarm float64 `json:"loadAlarm"`
}
