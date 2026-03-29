package main

import (
	"os"

	"github.com/domo84/bokio-cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
