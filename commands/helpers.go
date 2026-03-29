package commands

import (
	"encoding/json"
	"os"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	flagPage     int
	flagPageSize int
	flagQuery    string
	flagAll      bool
)

// addListFlags adds pagination and filter flags to a command.
func addListFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&flagPage, "page", 1, "Page number")
	cmd.Flags().IntVar(&flagPageSize, "page-size", 25, "Items per page (max 100)")
	cmd.Flags().StringVarP(&flagQuery, "query", "q", "", "Filter expression")
	cmd.Flags().BoolVar(&flagAll, "all", false, "Fetch all pages")
}

// listParamsFromFlags builds ListParams from the current flag values.
func listParamsFromFlags() api.ListParams {
	return api.ListParams{
		Page:     flagPage,
		PageSize: flagPageSize,
		Query:    flagQuery,
	}
}

// loadJSONFile reads a JSON file and decodes it into the given value.
func loadJSONFile(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
