package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func processTemplate(ctx context.Context, body string, data any) (string, error) {
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

func processCSV(body string) ([]byte, error) {
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
