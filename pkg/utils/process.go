package utils

import (
	"io"
	"strconv"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

// GetTimeout returns a timeout duration based on the provided lambda timeout string.
// If lambdaTimeout is non-empty and can be converted to an integer, the value is interpreted
// as a number of seconds and converted to a time.Duration. If the conversion fails or the string
// is empty, the defaultTimeout is returned.
//
// Parameters:
//   - lambdaTimeout: A string representing the timeout in seconds.
//   - defaultTimeout: The default time.Duration to return if lambdaTimeout is empty or invalid.
//
// Returns:
//   - time.Duration: The resulting timeout duration.
func GetTimeout(lambdaTimeout string, defaultTimeout time.Duration) time.Duration {
	if lambdaTimeout != "" {
		if timeout, err := strconv.Atoi(lambdaTimeout); err == nil {
			defaultTimeout = time.Duration(timeout) * time.Second
		}
	}
	return defaultTimeout
}

// TemplateFromReaderToWriter reads a template from an io.Reader, parses it, and executes the template
// with the provided data, writing the output to the specified io.Writer.
//
// Parameters:
//   - w: The io.Writer where the executed template output will be written.
//   - r: The io.Reader from which the template is read.
//   - data: The data to be applied to the template during execution.
//
// Returns:
//   - error: An error if reading, parsing, or executing the template fails; otherwise, nil.
func TemplateFromReaderToWriter(w io.Writer, r io.Reader, data any) error {
	templateBytes, err := io.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "read template failed")
	}

	tmpl, err := template.New("prompt").Parse(string(templateBytes))
	if err != nil {
		return errors.Wrap(err, "parse template failed")
	}

	return tmpl.Execute(w, data)
}
