package csv_test

import (
	"strings"
	"testing"

	"github.com/dc-dc-dc/cheetah/market/csv"
	"github.com/stretchr/testify/assert"
)

func TestCsv(t *testing.T) {
	reader := csv.NewCsvReader(strings.NewReader(""))
	_, err := reader.Header()
	assert.Error(t, err)

	reader = csv.NewCsvReader(strings.NewReader("Date,Open,High,Low\n0,0,0,0\n\n"))
	header, err := reader.Header()
	assert.NoError(t, err)
	assert.Equal(t, csv.CsvHeader{"date": 0, "open": 1, "high": 2, "low": 3}, header)
	assert.Equal(t, 1, reader.LineNumber())
	_, err = reader.Header()
	assert.Error(t, err)

	line, err := reader.NextLine()
	assert.NoError(t, err)
	assert.Equal(t, []string{"0", "0", "0", "0"}, line)
	assert.Equal(t, 2, reader.LineNumber())
	_, err = reader.NextLine()
	assert.Error(t, err)
	_, err = reader.NextLine()
	assert.Error(t, err)
}
