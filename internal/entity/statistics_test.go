package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountByMonthToResponse(t *testing.T) {
	cs := CountByMonths{
		CountByMonth{Month: "m", Count: 10},
		CountByMonth{Month: "x", Count: 5},
	}
	ms := []CountByCatAndMonth{
		{CategoryId: 1, CountByMonth: CountByMonth{Month: "m", Count: 4}},
		{CategoryId: 2, CountByMonth: CountByMonth{Month: "m", Count: 6}},
		{CategoryId: 1, CountByMonth: CountByMonth{Month: "x", Count: 5}},
	}
	stats := cs.ToResponse(ms)

	assert.Len(t, stats.Statistics, 2)
	assert.Equal(t, stats.Statistics[0].Month, "m")
	assert.EqualValues(t, stats.Statistics[0].DreamCount, 10)
	assert.Len(t, stats.Statistics[0].Categories, 2)
	assert.EqualValues(t, stats.Statistics[0].Categories[0].CategoryId, 1)
	assert.EqualValues(t, stats.Statistics[0].Categories[0].Count, 4)
	assert.EqualValues(t, stats.Statistics[0].Categories[1].CategoryId, 2)
	assert.EqualValues(t, stats.Statistics[0].Categories[1].Count, 6)

	assert.Equal(t, stats.Statistics[1].Month, "x")
	assert.EqualValues(t, stats.Statistics[1].DreamCount, 5)
	assert.Len(t, stats.Statistics[1].Categories, 1)

}
