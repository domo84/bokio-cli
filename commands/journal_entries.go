package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newJournalEntriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "journal-entries",
		Short: "Manage journal entries",
	}

	cmd.AddCommand(newJournalEntriesListCmd())
	cmd.AddCommand(newJournalEntriesGetCmd())
	cmd.AddCommand(newJournalEntriesCreateCmd())
	cmd.AddCommand(newJournalEntriesReverseCmd())

	return cmd
}

func newJournalEntriesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List journal entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			params := listParamsFromFlags()

			if flagAll {
				entries, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.JournalEntry], error) {
					return state.client.ListJournalEntries(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(entries)
			}

			resp, err := state.client.ListJournalEntries(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newJournalEntriesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a journal entry by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			entry, err := state.client.GetJournalEntry(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(entry)
		},
	}
}

func newJournalEntriesCreateCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a journal entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			if fromFile == "" {
				return fmt.Errorf("--from-file is required (journal entries have complex row structures)")
			}

			var req api.CreateJournalEntryRequest
			if err := loadJSONFile(fromFile, &req); err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			entry, err := state.client.CreateJournalEntry(cmd.Context(), req)
			if err != nil {
				return err
			}
			return state.formatter.Format(entry)
		},
	}

	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file (required)")

	return cmd
}

func newJournalEntriesReverseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reverse <id>",
		Short: "Reverse a journal entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			entry, err := state.client.ReverseJournalEntry(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(entry)
		},
	}
}
