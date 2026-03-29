package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newConnectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connections",
		Short: "Manage API connections",
	}

	cmd.AddCommand(newConnectionsListCmd())
	cmd.AddCommand(newConnectionsGetCmd())
	cmd.AddCommand(newConnectionsDeleteCmd())

	return cmd
}

func newConnectionsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List connections",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}

			params := listParamsFromFlags()

			if flagAll {
				connections, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.Connection], error) {
					return state.client.ListConnections(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(connections)
			}

			resp, err := state.client.ListConnections(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newConnectionsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a connection by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			conn, err := state.client.GetConnection(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(conn)
		},
	}
}

func newConnectionsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a connection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := state.client.DeleteConnection(cmd.Context(), args[0]); err != nil {
				return err
			}
			fmt.Println("Connection deleted.")
			return nil
		},
	}
}
