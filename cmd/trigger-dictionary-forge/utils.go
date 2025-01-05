package main

import (
	"bytes"
	"encoding/csv"
	"strings"

	"github.com/pkg/errors"
)

// toCSV converts a markdown/text table to CSV format
func toCSV(tableText string) ([]byte, error) {
	lines := strings.Split(strings.TrimSpace(tableText), "\n")
	if len(lines) < 3 {
		return nil, errors.New("invalid table format")
	}

	var cleanLines []string
	for i, line := range lines {
		if i == 1 && strings.Contains(line, "-") {
			continue
		}
		line = strings.Trim(line, "|")
		cleanLines = append(cleanLines, line)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	for _, line := range cleanLines {
		fields := strings.Split(line, "|")
		for i, field := range fields {
			fields[i] = strings.TrimSpace(field)
		}
		if err := writer.Write(fields); err != nil {
			return nil, errors.Wrap(err, "failed to write CSV line")
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, errors.Wrap(err, "failed to flush CSV writer")
	}
	return buf.Bytes(), nil
}
