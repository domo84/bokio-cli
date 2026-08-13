package output

import (
	"io"
	"os"
)

// Formatter defines how to output data to the user.
type Formatter interface {
	Format(data any) error
}

// ConfigSetting is one configuration key rendered for display. Source is where the
// current value came from: "env", "file" or "default".
type ConfigSetting struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

// NewFormatter creates the appropriate formatter based on format string.
func NewFormatter(format string) Formatter {
	return NewFormatterWithWriter(format, os.Stdout)
}

// NewFormatterWithWriter creates a formatter that writes to the given writer.
func NewFormatterWithWriter(format string, w io.Writer) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{Writer: w}
	default:
		return &TableFormatter{Writer: w}
	}
}
