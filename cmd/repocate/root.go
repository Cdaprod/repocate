package repocate

import (
    "github.com/spf13/cobra"
)

// rootCmd is the root command for the CLI
var rootCmd = &cobra.Command{
    Use:   "repocate",
    Short: "Repocate is a tool for managing development environments using Docker containers.",
    Long:  `Repocate allows you to clone repositories, create isolated development environments, and manage them using Docker containers.`,
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
    displayBanner() // Call the banner display function from banner.go
    return rootCmd.Execute()
}
