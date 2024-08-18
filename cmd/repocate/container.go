package container

import (
    "fmt"
    "github.com/spf13/cobra"
    "repocate/internal/container"
    "repocate/internal/config"
    "repocate/internal/log"
)

var CreateCmd = &cobra.Command{
    Use:   "create [repository URL or name]",
    Short: "Clone a repo and create/start a development container.",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        repoInput := args[0]
        config.LoadConfig()
        log.SetupLogger()

        repoName, err := container.ResolveRepoName(repoInput)
        if err != nil {
            log.Error(err)
            return
        }

        if !container.IsRepoCloned(config.WorkspaceDir, repoName) {
            err = container.CloneRepository(config.WorkspaceDir, repoInput)
            if err != nil {
                log.Error(err)
                return
            }
        }

        err = container.InitContainer(config.WorkspaceDir, repoName)
        if err != nil {
            log.Error(err)
        }
    },
}

// Similarly define EnterCmd, StopCmd, RebuildCmd in this file.

var EnterCmd = &cobra.Command{
    Use:   "enter [repository URL or name]",
    Short: "Enter the development container for a specific repo.",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        repoInput := args[0]
        config.LoadConfig()
        log.SetupLogger()

        repoName, err := container.ResolveRepoName(repoInput)
        if err != nil {
            log.Error(err)
            return
        }

        err = container.EnterContainer(config.WorkspaceDir, repoName)
        if err != nil {
            log.Error(err)
        }
    },
}

// StopCmd and RebuildCmd follow the same pattern.