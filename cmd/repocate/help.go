package help

import (
    "github.com/spf13/cobra"
)

var HelpCmd = &cobra.Command{
    Use:   "help",
    Short: "Show help information.",
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
    },
}