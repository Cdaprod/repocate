package repocate

import (
	"fmt"
	"os"
	"time"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"github.com/cheggaaa/pb/v3"
	"github.com/cdaprod/repocate/internal/container"
	"github.com/cdaprod/repocate/internal/config"
	"github.com/cdaprod/repocate/pkg/plugin"
	"github.com/cdaprod/repocate/pkg/tailscale"
	"github.com/cdaprod/repocate/pkg/vault"
)

var reg = plugin.NewRegistry()

func init() {
	cobra.OnInitialize(initConfig)

	// Register commands
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(CreateCmd)
	rootCmd.AddCommand(EnterCmd)
	rootCmd.AddCommand(StopCmd)
	rootCmd.AddCommand(RebuildCmd)
	rootCmd.AddCommand(CloneCmd)
	rootCmd.AddCommand(ListCmd)
	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(HelpCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(unregisterCmd)

	// Register plugins dynamically without direct import
	reg.RegisterPlugin(&tailscale.TailscalePlugin{})
	reg.RegisterPlugin(&vault.VaultPlugin{})

	// Inject dependencies into plugins
	reg.InjectDependencies("tailscale", config.GetString("TailscaleAuthKey"), config.GetString("TailscaleServiceName"))
	reg.InjectDependencies("vault", config.GetString("VaultAddress"))
}

func initConfig() {
	if err := config.LoadConfig(""); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	displayBanner()
	return rootCmd.Execute()
}

func displayBanner() {
	// Display ASCII art banner with color
	color.Cyan(`  ______                           _       
 | ___ \                         | |      
 | |_/ /___ _ __   ___   ___ __ _| |_ ___ 
 |    // _ \ '_ \ / _ \ / __/ _` + "`" + ` | __/ _ \
 | |\ \  __/ |_) | (_) | (_| (_| | ||  __/
 \_| \_\___| .__/ \___/ \___\__,_|\__\___|
           | |                            
           |_|`)

	fmt.Println(color.HiMagentaString("By: David Cannan aka Cdaprod"))

	color.Green("\nRepocate allows you to clone repositories, create isolated development environments, and manage them using Docker containers.")
	color.Green("It now supports dynamic plugin loading and a flexible registry system.")

	color.Blue("\n\nUsage:")
	fmt.Println("  repocate [command]")

	color.Blue("\n\nAvailable Commands:")
	displayCommands()

	color.Blue("\n\nFlags:")
	fmt.Println("  -h, --help   help for repocate")

	fmt.Println(color.GreenString("\nUse \"repocate [command] --help\" for more information about a command."))
}

func displayCommands() {
	commands := []struct {
		name        string
		description string
	}{
		{"clone", "Clone a repository"},
		{"create", "Clone a repo and create/start a development container"},
		{"enter", "Enter the development container for a specific repo"},
		{"help", "Show help information"},
		{"list", "List all repocate containers and plugins"},
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

// rootCmd is the root command for the CLI
var rootCmd = &cobra.Command{
	Use:   "repocate",
	Short: "Repocate is a tool for managing development environments using Docker containers.",
	Long:  `Repocate allows you to clone repositories, create isolated development environments, and manage them using Docker containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Display help information when no subcommand is provided
		cmd.Help()
	},
}

// startCmd is the command to initialize and start the default container
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Initialize and start the default Repocate container",
	Run: func(cmd *cobra.Command, args []string) {
		container.HandleDefaultContainer()
	},
}

// Add other command definitions here (CreateCmd, EnterCmd, StopCmd, etc.)

// showProgress displays a simple progress bar
func showProgress(message string, milliseconds int) {
	fmt.Print(message)
	bar := pb.New(100)
	bar.SetMaxWidth(80)
	bar.Start()
	for i := 0; i < 100; i++ {
		bar.Increment()
		time.Sleep(time.Duration(milliseconds/100) * time.Millisecond)
	}
	bar.Finish()
	fmt.Println()
}

// Add registerCmd, listCmd, and unregisterCmd definitions here

func init() {
	// Add any additional initialization if needed
}