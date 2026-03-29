package commands

import (
	"strconv"

	"github.com/spf13/cobra"
)

func newAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "View chart of accounts",
	}

	cmd.AddCommand(newAccountsListCmd())
	cmd.AddCommand(newAccountsGetCmd())

	return cmd
}

func newAccountsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			accounts, err := state.client.ListAccounts(cmd.Context())
			if err != nil {
				return err
			}
			return state.formatter.Format(accounts)
		},
	}
}

func newAccountsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <accountNumber>",
		Short: "Get an account by number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			num, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			account, err := state.client.GetAccount(cmd.Context(), num)
			if err != nil {
				return err
			}
			return state.formatter.Format(account)
		},
	}
}
