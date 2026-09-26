package commands

import (
	"fmt"
	"os"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/domo84/bokio-cli/internal/auth"
	"github.com/domo84/bokio-cli/internal/config"
	"github.com/domo84/bokio-cli/internal/output"
	"github.com/spf13/cobra"
)

// rootState holds shared state across all commands.
type rootState struct {
	cfg       *config.Config
	client    *api.Client
	formatter output.Formatter
	store     *auth.TokenStore
	configDir string
}

var (
	flagCompanyID    string
	flagOutputFormat string
	flagConfigDir    string
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bokio",
		Short: "CLI for the Bokio accounting API",
		Long: fmt.Sprintf(`A command-line interface for interacting with the Bokio API (https://api.bokio.se/v1/).

Configuration:
  Settings resolve in this order, highest first:
    1. command-line flags
    2. BOKIO_* environment variables
    3. the config file
    4. built-in defaults

  Config file: %s
  Run "bokio config --help" for the list of settings, or "bokio config list" to see
  their current values.`, configFilePath()),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVarP(&flagCompanyID, "company-id", "c", "", "Company ID (env BOKIO_COMPANY_ID, config company_id)")
	cmd.PersistentFlags().StringVarP(&flagOutputFormat, "output", "o", "", "Output format: table, json (env BOKIO_OUTPUT, config output_format) (default table)")
	cmd.PersistentFlags().StringVar(&flagConfigDir, "config-dir", "", fmt.Sprintf("Config directory (default %s)", config.DefaultConfigDir()))

	return cmd
}

func initState(cmd *cobra.Command) (*rootState, error) {
	cfg, err := config.Load(flagConfigDir)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// Flag overrides
	if flagCompanyID != "" {
		cfg.CompanyID = flagCompanyID
	}
	if flagOutputFormat != "" {
		cfg.OutputFormat = flagOutputFormat
	}

	configDir := flagConfigDir
	if configDir == "" {
		configDir = config.DefaultConfigDir()
	}

	store := auth.NewTokenStore(configDir)

	state := &rootState{
		cfg:       cfg,
		store:     store,
		configDir: configDir,
		formatter: output.NewFormatter(cfg.OutputFormat),
	}

	return state, nil
}

func initStateWithClient(cmd *cobra.Command) (*rootState, error) {
	state, err := initState(cmd)
	if err != nil {
		return nil, err
	}

	token, err := auth.GetToken(state.store)
	if err != nil {
		return nil, err
	}

	state.client = api.NewClient(token, state.cfg.CompanyID)
	return state, nil
}

func requireCompanyID(state *rootState) error {
	if state.cfg.CompanyID == "" {
		return fmt.Errorf("company ID is required. Pass --company-id, set BOKIO_COMPANY_ID, " +
			"or persist it with: bokio config set company_id <id>")
	}
	return nil
}

// Execute runs the root command.
func Execute() error {
	root := newRootCmd()

	// Register all commands
	root.AddCommand(newAuthCmd())
	root.AddCommand(newCompanyCmd())
	root.AddCommand(newCustomersCmd())
	root.AddCommand(newSuppliersCmd())
	root.AddCommand(newItemsCmd())
	root.AddCommand(newInvoicesCmd())
	root.AddCommand(newJournalEntriesCmd())
	root.AddCommand(newCreditNotesCmd())
	root.AddCommand(newUploadsCmd())
	root.AddCommand(newBankPaymentsCmd())
	root.AddCommand(newAccountsCmd())
	root.AddCommand(newFiscalYearsCmd())
	root.AddCommand(newSIECmd())
	root.AddCommand(newConnectionsCmd())
	root.AddCommand(newConfigCmd())
	root.AddCommand(newCompletionCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}
	return nil
}
