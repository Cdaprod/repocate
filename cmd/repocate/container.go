package repocate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	cont "github.com/cdaprod/repocate/internal/container" // Alias to avoid conflict
	"github.com/cdaprod/repocate/internal/utils"
	"github.com/cdaprod/repocate/internal/config"
	"github.com/cdaprod/repocate/internal/log"
	"github.com/cdaprod/repocate/internal/git"
)

// Helper function to load config and set up logging
func initializeEnvironment() {
	config.LoadConfig()
	log.SetupLogger()
}

type CommandFactory struct{}

func NewCommandFactory() *CommandFactory {
	return &CommandFactory{}
}

func (f *CommandFactory) AddContainerCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		CreateCmd,
		CloneCmd,
		EnterCmd,
		StopCmd,
		RebuildCmd,
	)
}

var CreateCmd = &cobra.Command{
	Use:   "create [repository URL or name]",
	Short: "Clone a repo and create/start a development container.",
	Args:  cobra.ExactArgs(1),
	Run:   runCreateCommand,
}

func runCreateCommand(cmd *cobra.Command, args []string) {
	repoInput := args[0]
	initializeEnvironment()

	repoName, err := utils.ExtractRepoName(repoInput)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to extract repo name: %s", err))
		return
	}

	repoPath := filepath.Join(config.WorkspaceDir, repoName)

	if !utils.IsRepoCloned(config.WorkspaceDir, repoName) {
		err = git.CloneRepository(config.WorkspaceDir, repoInput)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to clone repository: %s", err))
			return
		}
	}

	dockerfilePath := filepath.Join(repoPath, "Dockerfile")

	if repoName == "repocate-default" {
		log.Info("Using Dockerfile.multiarch for repocate-default.")
		// Add logic here to ensure the Dockerfile.multiarch is used
	} else if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		log.Error("Dockerfile not found in the cloned repository.")
		return
	} else {
		log.Info(fmt.Sprintf("Using Dockerfile for repository %s.", repoName))
	}

	err = cont.InitContainer(config.WorkspaceDir, repoName)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to initialize container: %s", err))
		return
	}

	log.Info("Project environment created successfully.")
}

var CloneCmd = &cobra.Command{
	Use:   "clone [repository URL]",
	Short: "Clone a repository.",
	Args:  cobra.ExactArgs(1),
	Run:   runCloneCommand,
}

func runCloneCommand(cmd *cobra.Command, args []string) {
	repoURL := args[0]
	initializeEnvironment()

	err := git.CloneRepository(config.WorkspaceDir, repoURL)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to clone repository: %s", err))
		return
	}

	log.Info("Repository cloned successfully.")
}

var EnterCmd = &cobra.Command{
	Use:   "enter [repository URL or name]",
	Short: "Enter the development container for a specific repo.",
	Args:  cobra.ExactArgs(1),
	Run:   runEnterCommand,
}

func runEnterCommand(cmd *cobra.Command, args []string) {
	repoInput := args[0]
	initializeEnvironment()

	repoName, err := utils.ExtractRepoName(repoInput)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = cont.EnterContainer(config.WorkspaceDir, repoName)
	if err != nil {
		log.Error(err.Error())
	}
}

var StopCmd = &cobra.Command{
	Use:   "stop [repository URL or name]",
	Short: "Stop the development container for a specific repo.",
	Args:  cobra.ExactArgs(1),
	Run:   runStopCommand,
}

func runStopCommand(cmd *cobra.Command, args []string) {
	repoInput := args[0]
	initializeEnvironment()

	repoName, err := utils.ExtractRepoName(repoInput)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = cont.StopContainer(config.WorkspaceDir, repoName)
	if err != nil {
		log.Error(err.Error())
	}
}

var RebuildCmd = &cobra.Command{
	Use:   "rebuild [repository URL or name]",
	Short: "Rebuild the development container for a specific repo.",
	Args:  cobra.ExactArgs(1),
	Run:   runRebuildCommand,
}

func runRebuildCommand(cmd *cobra.Command, args []string) {
	repoInput := args[0]
	initializeEnvironment()

	repoName, err := utils.ExtractRepoName(repoInput)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = cont.RebuildContainer(config.WorkspaceDir, repoName)
	if err != nil {
		log.Error(err.Error())
	}
}