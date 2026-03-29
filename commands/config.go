package commands

import (
	"fmt"
	"strings"

	"github.com/domo84/bokio-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validConfigKeys = []string{
	"auth_mode",
	"company_id",
	"output_format",
	"client_id",
	"client_secret",
	"redirect_port",
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigSetCmd())

	return cmd
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "get <property>",
		Short:     "Get a config value",
		Args:      cobra.ExactArgs(1),
		ValidArgs: validConfigKeys,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := initState(cmd)
			if err != nil {
				return err
			}

			key := args[0]
			if !isValidKey(key) {
				return fmt.Errorf("unknown config key %q. Valid keys: %s", key, strings.Join(validConfigKeys, ", "))
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
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initState(cmd)
			if err != nil {
				return err
			}

			key := args[0]
			value := args[1]

			if !isValidKey(key) {
				return fmt.Errorf("unknown config key %q. Valid keys: %s", key, strings.Join(validConfigKeys, ", "))
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

func isValidKey(key string) bool {
	for _, k := range validConfigKeys {
		if k == key {
			return true
		}
	}
	return false
}

