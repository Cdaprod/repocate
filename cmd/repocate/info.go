package info

import (
    "fmt"
    "github.com/spf13/cobra"
    "repocate/internal/container"
)

var ListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all repocate containers.",
    Run: func(cmd *cobra.Command, args []string) {
        err := container.ListContainers()
        if err != nil {
            fmt.Println("Error listing containers:", err)
        }
    },
}

var VersionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show version information.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Repocate version 1.0.0")
    },
}