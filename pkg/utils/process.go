package utils

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

// GetTimeout returns a timeout duration from a lambda timeout string or a default timeout.
func GetTimeout(lambdaTimeout string, defaultTimeout time.Duration) time.Duration {
	if lambdaTimeout != "" {
		if timeout, err := strconv.Atoi(lambdaTimeout); err == nil {
			defaultTimeout = time.Duration(timeout) * time.Second
		}
	}
	return defaultTimeout
}

// Template takes a template string and data to fill the template
// and returns the generated string or an error.
func Template(body string, data any) (string, error) {
	tmpl, err := template.New("promt").Parse(body)
	if err != nil {
		return "", errors.Wrap(err, "parse teplate failed")
	}
	var content bytes.Buffer
	if err = tmpl.Execute(&content, data); err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}
	return content.String(), nil
}

// CSV takes a string containing a table and returns CSV data as a byte slice or an error.
// The table should be represented as lines separated by newline characters, and each line should
// contain fields separated by the '|' character.
func CSV(body string) ([]byte, error) {
	if strings.TrimSpace(body) == "" {
		return nil, errors.New("empty input")
	}

	lines := strings.Split(strings.TrimSpace(body), "\n")
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
	var columnCount int

	for i, line := range cleanLines {
		fields := strings.Split(line, "|")
		for i, field := range fields {
			fields[i] = strings.TrimSpace(field)
		}

		if i == 0 {
			columnCount = len(fields)
		} else if len(fields) != columnCount {
			return nil, errors.New("inconsistent number of columns")
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
