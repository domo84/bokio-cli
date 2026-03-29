package output

import (
	"encoding/json"
	"io"
)

// JSONFormatter outputs data as formatted JSON.
type JSONFormatter struct {
	Writer io.Writer
}

func (f *JSONFormatter) Format(data any) error {
	enc := json.NewEncoder(f.Writer)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
