package models

type SeriesType string

const (
	Day   SeriesType = "day"
	Week  SeriesType = "week"
	Month SeriesType = "month"
)

type Value struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

type Series []Value

type SeriesRepository interface {
	FindCpu(SeriesType) (Series, error)
	FindMemory(SeriesType) (Series, error)
	FindDisk(SeriesType, string) (Series, error)
	FindDiskPaths() ([]string, error)
}
