package config

import (
	"strings"

	"github.com/petaki/satellite/internal/models"
)

func parseSeriesButtons(value string) []models.SeriesType {
	var sb []models.SeriesType
	segments := strings.SplitSeq(value, ",")

	for segment := range segments {
		st := models.SeriesType(strings.TrimSpace(segment))

		for _, current := range models.SeriesTypes {
			if st == current["value"].(models.SeriesType) {
				sb = append(sb, st)

				break
			}
		}
	}

	sb = sb[:min(4, len(sb))]

	return sb
}
