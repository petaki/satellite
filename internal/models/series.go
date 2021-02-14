package models

type SeriesType string

const (
	Day   SeriesType = "day"
	Week  SeriesType = "week"
	Month SeriesType = "month"
)

type Series map[string]float64

type SeriesRepository interface {
	FindCpu(SeriesType) (Series, error)
	FindMemory(SeriesType) (Series, error)
	FindDisk(SeriesType, string) (Series, error)
	FindDiskPaths() ([]string, error)
}