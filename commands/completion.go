package commands

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for bokio.

To load completions:

Bash:
  $ source <(bokio completion bash)
  # To load completions for each session, execute once:
  $ bokio completion bash > /etc/bash_completion.d/bokio

Zsh:
  $ source <(bokio completion zsh)
  # To load completions for each session, execute once:
  $ bokio completion zsh > "${fpath[1]}/_bokio"

Fish:
  $ bokio completion fish | source
  # To load completions for each session, execute once:
  $ bokio completion fish > ~/.config/fish/completions/bokio.fish

PowerShell:
  PS> bokio completion powershell | Out-String | Invoke-Expression
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}
