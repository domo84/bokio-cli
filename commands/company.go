package commands

import (
	"github.com/spf13/cobra"
)

func newCompanyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "company",
		Short: "Get company information",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			info, err := state.client.GetCompanyInfo(cmd.Context())
			if err != nil {
				return err
			}
			return state.formatter.Format(info)
		},
	}
}
