package repocate

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

// VersionCmd shows the version of the tool
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Repocate version 1.0.0")
	},
}

// Display the ASCII art banner and general information
func displayBanner() {
	fmt.Println(color.Cyan(`
  ______                           _       
 | ___ \                         | |      
 | |_/ /___ _ __   ___   ___ __ _| |_ ___ 
 |    // _ \ '_ \ / _ \ / __/ _` + "`" + ` | __/ _ \
 | |\ \  __/ |_) | (_) | (_| (_| | ||  __/
 \_| \_\___| .__/ \___/ \___\__,_|\__\___|
           | |                            
           |_|`))

	fmt.Println(color.HiMagentaString("By: David Cannan aka Cdaprod"))
	color.Green("\nRepocate allows you to clone repositories, create isolated development environments, and manage them using Docker containers.")
	color.Green("It now supports dynamic plugin loading and a flexible registry system.\n")

	color.Blue("\nUsage:")
	fmt.Println("  repocate [command]")

	color.Blue("\nAvailable Commands:")
	displayCommands()

	color.Blue("\nFlags:")
	fmt.Println("  -h, --help   help for repocate")

	fmt.Println(color.GreenString("\nUse \"repocate [command] --help\" for more information about a command."))
}

// displayCommands lists the available commands and their descriptions
func displayCommands() {
	commands := []struct {
		name        string
		description string
	}{
		{"clone", "Clone a repository"},
		{"create", "Clone a repo and create/start a development container"},
		{"enter", "Enter the development container for a specific repo"},
		{"help", "Show help information"},
		{"list", "List all repocate containers"},
		{"rebuild", "Rebuild the development container for a specific repo"},
		{"stop", "Stop the development container for a specific repo"},
		{"version", "Show version information"},
		{"register", "Register a new plugin"},
		{"unregister", "Unregister a plugin"},
	}

	for _, cmd := range commands {
		fmt.Printf("  %-12s %s\n", cmd.name, cmd.description)
	}
}