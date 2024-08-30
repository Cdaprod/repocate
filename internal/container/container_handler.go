package container

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// HandleDefaultContainer handles the default container initialization and startup
func HandleDefaultContainer() {
	color.Cyan("Initializing and starting the 'repocate-default' container...")

	showProgress("Checking container status...", 100)

	// Initialize the default container if not exists or ensure it is running
	err := InitRepocateDefaultContainer()
	if err != nil {
		fmt.Println(color.RedString("Error initializing 'repocate-default' container: %s", err))
		os.Exit(1)
	}

	color.Green("Checking status of the 'repocate-default' container...")

	// Check if the container is running
	isRunning, err := IsContainerRunning("repocate-default")
	if err != nil {
		fmt.Println(color.RedString("Error checking container status: %s", err))
		os.Exit(1)
	}

	if !isRunning {
		color.Yellow("Container 'repocate-default' is not running. Starting it now...")

		err := StartContainer("repocate-default")
		if err != nil {
			fmt.Println(color.RedString("Error starting container: %s", err))
			os.Exit(1)
		}
	}

	color.Green("'repocate-default' container is ready.")
}

// showProgress is a placeholder function. Implement it based on your needs.
func showProgress(message string, milliseconds int) {
	// Implement progress display logic here
}