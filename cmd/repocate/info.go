package repocate

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/cdaprod/repocate/internal/container"
)

var ListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all repocate containers.",
    Run: func(cmd *cobra.Command, args []string) {
        containers, err := container.ListContainers()
        if err != nil {
            fmt.Println("Error listing containers:", err)
            return
        }

        for _, c := range containers {
            fmt.Println(c) // Assuming `c` has a `String` method or implement a way to print the container information
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