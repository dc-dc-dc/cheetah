package csv

import (
	"bufio"
	"io"
	"strings"
)

type CsvReader struct {
	src    io.Reader
	reader *bufio.Reader
}

func NewCsvReader(src io.Reader) *CsvReader {
	return &CsvReader{
		src:    src,
		reader: bufio.NewReader(src),
	}
}

func (c *CsvReader) Header() (map[string]int, error) {
	line, err := c.NextLine()
	if err != nil {
		return nil, err
	}
	header := make(map[string]int)
	for i, col := range line {
		header[strings.ToLower(col)] = i
	}
	return header, nil
}

func (c *CsvReader) NextLine() ([]string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.Trim(line, "\n"), ","), nil
}
