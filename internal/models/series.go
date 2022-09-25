package models

// SeriesType type.
type SeriesType string

const (
	// Day series.
	Day SeriesType = "day"

	// Week series.
	Week SeriesType = "week"

	// Month series.
	Month SeriesType = "month"
)

// Value type.
type Value struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

// Series type.
type Series []Value

// SeriesRepository type.
type SeriesRepository interface {
	FindCPU(Probe, SeriesType) (Series, Series, Series, error)
	FindMemory(Probe, SeriesType) (Series, Series, Series, error)
	FindDisk(Probe, SeriesType, string) (Series, Series, Series, error)
	FindDiskPaths(Probe) ([]string, error)
}
