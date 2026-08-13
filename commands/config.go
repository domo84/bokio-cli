package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/domo84/bokio-cli/internal/config"
	"github.com/domo84/bokio-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configKey describes a supported setting. This table is the single source of truth
// for validation, shell completion and help text.
type configKey struct {
	Name        string
	Env         string
	Default     string
	Description string
	// Secret marks values that should be masked when listed.
	Secret bool
}

var configKeys = []configKey{
	{
		Name: "company_id", Env: "BOKIO_COMPANY_ID",
		Description: "Default company ID, so -c can be omitted",
	},
	{
		Name: "output_format", Env: "BOKIO_OUTPUT", Default: "table",
		Description: "Output format: table or json",
	},
	{
		Name: "client_id", Env: "BOKIO_CLIENT_ID",
		Description: "OAuth client ID",
	},
	{
		Name: "client_secret", Env: "BOKIO_CLIENT_SECRET", Secret: true,
		Description: "OAuth client secret",
	},
	{
		Name: "redirect_port", Env: "BOKIO_REDIRECT_PORT", Default: "8585",
		Description: "Local port for the OAuth callback",
	},
}

func configKeyNames() []string {
	names := make([]string, 0, len(configKeys))
	for _, key := range configKeys {
		names = append(names, key.Name)
	}
	return names
}

func isValidKey(key string) bool {
	for _, k := range configKeys {
		if k.Name == key {
			return true
		}
	}
	return false
}

// configFilePath is the default location of the config file. --config-dir moves the
// directory, so this is what help text should advertise.
func configFilePath() string {
	return filepath.Join(config.DefaultConfigDir(), "config.yaml")
}

// configKeysHelp renders the supported keys as an aligned block for help text.
func configKeysHelp() string {
	var b strings.Builder
	b.WriteString("Supported keys:\n")

	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	for _, key := range configKeys {
		description := key.Description
		if key.Default != "" {
			description += fmt.Sprintf(" (default %s)", key.Default)
		}
		fmt.Fprintf(w, "  %s\t%s\t%s\n", key.Name, description, key.Env)
	}
	w.Flush()

	return b.String()
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long: fmt.Sprintf(`Manage configuration.

Settings are stored in a YAML file and can be overridden per invocation by BOKIO_*
environment variables and by command-line flags.

  File: %s
        (use --config-dir to keep it somewhere else)

Precedence, highest first: flags, environment variables, config file, defaults.

%s
Examples:
  bokio config set company_id 3fa85f64-5717-4562-b3fc-2c963f66afa6
  bokio config get company_id
  bokio config list`, configFilePath(), configKeysHelp()),
	}

	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigListCmd())

	return cmd
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <property>",
		Short: "Get a config value",
		Long: "Print the current value of a config key, including any override from a\n" +
			"BOKIO_* environment variable.\n\n" + configKeysHelp(),
		Args:      cobra.ExactArgs(1),
		ValidArgs: configKeyNames(),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := initState(cmd)
			if err != nil {
				return err
			}

			key := args[0]
			if !isValidKey(key) {
				return unknownKeyError(key)
			}

			val := viper.Get(key)
			if val == nil || val == "" {
				fmt.Println("(not set)")
			} else {
				fmt.Println(val)
			}
			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <property> <value>",
		Short: "Set a config value",
		Long: "Write a config key to the config file at\n  " + configFilePath() +
			"\n\n" + configKeysHelp(),
		Args:      cobra.ExactArgs(2),
		ValidArgs: configKeyNames(),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initState(cmd)
			if err != nil {
				return err
			}

			key := args[0]
			value := args[1]

			if !isValidKey(key) {
				return unknownKeyError(key)
			}

			viper.Set(key, value)

			if err := viper.Unmarshal(state.cfg); err != nil {
				return err
			}

			if err := config.Save(state.configDir, state.cfg); err != nil {
				return err
			}

			fmt.Printf("%s = %s\n", key, value)
			return nil
		},
	}
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all config keys and their current values",
		Long: `List every supported config key with its current value and where that value
came from.

Sources are "env" (a BOKIO_* variable), "file" (the config file) or "default".
Note that "bokio config set" writes a full snapshot of every key, so once anything
has been saved the remaining keys report "file" even though they still hold their
default values.

The --company-id and --output flags override these at run time and are not
reflected here. Secrets are masked; use "bokio config get" to read one.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initState(cmd)
			if err != nil {
				return err
			}

			settings := make([]output.ConfigSetting, 0, len(configKeys))
			for _, key := range configKeys {
				settings = append(settings, output.ConfigSetting{
					Key:    key.Name,
					Value:  configDisplayValue(key),
					Source: configValueSource(key),
				})
			}
			return state.formatter.Format(settings)
		},
	}
}

// configDisplayValue renders a key's current value, masking secrets so that a listing
// can be shared without leaking one.
func configDisplayValue(key configKey) string {
	value := viper.GetString(key.Name)
	if value == "" {
		return "(not set)"
	}
	if key.Secret {
		return "(set)"
	}
	return value
}

// configValueSource reports which layer supplied the current value.
func configValueSource(key configKey) string {
	if _, ok := os.LookupEnv(key.Env); ok {
		return "env"
	}
	if viper.InConfig(key.Name) {
		return "file"
	}
	return "default"
}

func unknownKeyError(key string) error {
	return fmt.Errorf("unknown config key %q. Valid keys: %s",
		key, strings.Join(configKeyNames(), ", "))
}
