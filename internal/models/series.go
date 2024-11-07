package models

// SeriesType type.
type SeriesType string

const (
	// Last5Minutes series.
	Last5Minutes SeriesType = "last_5_minutes"

	// Last15Minutes series.
	Last15Minutes SeriesType = "last_15_minutes"

	// Last30Minutes series.
	Last30Minutes SeriesType = "last_30_minutes"

	// Last1Hour series.
	Last1Hour SeriesType = "last_1_hour"

	// Last3Hours series.
	Last3Hours SeriesType = "last_3_hours"

	// Last6Hours series.
	Last6Hours SeriesType = "last_6_hours"

	// Last12Hours series.
	Last12Hours SeriesType = "last_12_hours"

	// Last24Hours series.
	Last24Hours SeriesType = "last_24_hours"

	// Last2Days series.
	Last2Days SeriesType = "last_2_days"

	// Last7Days series.
	Last7Days SeriesType = "last_7_days"

	// Last30Days series.
	Last30Days SeriesType = "last_30_days"
)

var SeriesTypes = []map[string]interface{}{
	{
		"name":  "Last 5 minutes",
		"value": Last5Minutes,
	},
	{
		"name":  "Last 15 minutes",
		"value": Last15Minutes,
	},
	{
		"name":  "Last 30 minutes",
		"value": Last30Minutes,
	},
	{
		"name":  "Last 1 hour",
		"value": Last1Hour,
	},
	{
		"name":  "Last 3 hours",
		"value": Last3Hours,
	},
	{
		"name":  "Last 6 hours",
		"value": Last6Hours,
	},
	{
		"name":  "Last 12 hours",
		"value": Last12Hours,
	},
	{
		"name":  "Last 24 hours",
		"value": Last24Hours,
	},
	{
		"name":  "Last 2 days",
		"value": Last2Days,
	},
	{
		"name":  "Last 7 days",
		"value": Last7Days,
	},
	{
		"name":  "Last 30 days",
		"value": Last30Days,
	},
}

// Value type.
type Value struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

// Series type.
type Series []Value

// ProcessValue type.
type ProcessValue struct {
	Name string  `json:"name"`
	X    int64   `json:"x"`
	Y    float64 `json:"y"`
}

// ProcessSeries type.
type ProcessSeries []ProcessValue

// SeriesRepository type.
type SeriesRepository interface {
	FindCPU(Probe, SeriesType) (Series, Series, Series, ProcessSeries, ProcessSeries, ProcessSeries, error)
	FindMemory(Probe, SeriesType) (Series, Series, Series, ProcessSeries, ProcessSeries, ProcessSeries, error)
	FindLoad(Probe, SeriesType) (Series, Series, Series, error)
	FindDisk(Probe, SeriesType, string) (Series, Series, Series, error)
	FindDiskPaths(Probe) ([]string, error)
	ChunkSize(seriesType SeriesType) int
}
