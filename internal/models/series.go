package models

// WatcherType type.
type WatcherType string

const (
	// CPU watcher.
	CPU WatcherType = "cpu"

	// Memory watcher.
	Memory WatcherType = "memory"

	// Disk watcher.
	Disk WatcherType = "disk"
)

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
	X int64          `json:"x"`
	Y float64        `json:"y"`
	P []ProcessValue `json:"p"`
}

// ProcessValue type.
type ProcessValue struct {
	Name string  `json:"name"`
	Y    float64 `json:"y"`
}

// Series type.
type Series []Value

// SeriesRepository type.
type SeriesRepository interface {
	FindCPU(Probe, SeriesType) (Series, Series, Series, error)
	FindMemory(Probe, SeriesType) (Series, Series, Series, error)
	FindDisk(Probe, SeriesType, string) (Series, Series, Series, error)
	FindDiskPaths(Probe) ([]string, error)
	ChunkSize(seriesType SeriesType) int
}
