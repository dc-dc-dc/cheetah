package csv

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CsvReader struct {
	src    io.Reader
	reader *bufio.Reader
	line   int
}

func NewCsvReader(src io.Reader) *CsvReader {
	return &CsvReader{
		src:    src,
		reader: bufio.NewReader(src),
	}
}

func (c *CsvReader) Header() (map[string]int, error) {
	if c.line != 0 {
		return nil, fmt.Errorf("header already read")
	}
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

func (c *CsvReader) LineNumber() int {
	return c.line
}

func (c *CsvReader) NextLine() ([]string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	c.line += 1
	return strings.Split(strings.Trim(line, "\n"), ","), nil
}
