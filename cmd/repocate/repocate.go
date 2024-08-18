package repocate

import (
    "github.com/spf13/cobra"
    "repocate/cmd/repocate/container"
    "repocate/cmd/repocate/info"
    "repocate/cmd/repocate/help"
)

// rootCmd is the root command for the CLI
var rootCmd = &cobra.Command{
    Use:   "repocate",
    Short: "Repocate is a tool for managing development environments using Docker containers.",
    Long:  `Repocate allows you to clone repositories, create isolated development environments, and manage them using Docker containers.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(container.CreateCmd)
    rootCmd.AddCommand(container.EnterCmd)
    rootCmd.AddCommand(container.StopCmd)
    rootCmd.AddCommand(container.RebuildCmd)
    rootCmd.AddCommand(info.ListCmd)
    rootCmd.AddCommand(info.VersionCmd)
    rootCmd.AddCommand(help.HelpCmd)
}