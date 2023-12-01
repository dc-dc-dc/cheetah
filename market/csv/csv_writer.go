package csv

import (
	"bufio"
	"io"
)

type CsvWriter struct {
	dst    *io.Writer
	writer bufio.Writer
}

func NewCsvWriter(dst *io.Writer) *CsvWriter {
	return &CsvWriter{
		dst:    dst,
		writer: *bufio.NewWriter(*dst),
	}
}

func (w *CsvWriter) Write(elements []string) error {
	for i, element := range elements {
		if i > 0 {
			w.writer.WriteString(",")
		}
		w.writer.WriteString(element)
	}
	w.writer.WriteString("\n")
	return w.writer.Flush()
}
