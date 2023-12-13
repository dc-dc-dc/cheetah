package csv

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CsvHeader = map[string]int

type CsvReader struct {
	src    io.Reader
	reader *bufio.Reader
	line   int
	header CsvHeader
}

func NewCsvReader(src io.Reader) *CsvReader {
	return &CsvReader{
		src:    src,
		reader: bufio.NewReader(src),
	}
}

func (c *CsvReader) Header() (CsvHeader, error) {
	if c.header != nil {
		return c.header, nil
	}
	if c.line != 0 {
		return nil, fmt.Errorf("header already read")
	}
	line, err := c.NextLine()
	if err != nil {
		return nil, err
	}
	c.header = make(CsvHeader)
	for i, col := range line {
		c.header[strings.ToLower(col)] = i
	}
	return c.header, nil
}

func (c *CsvReader) LineNumber() int {
	return c.line
}

func (c *CsvReader) NextLine() ([]string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimRight(line, "\n")
	if len(line) == 0 || line == "" {
		return nil, io.EOF
	}
	c.line += 1
	return strings.Split(line, ","), nil
}
