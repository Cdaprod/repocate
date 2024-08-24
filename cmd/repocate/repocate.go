package repocate

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/fatih/color"
    "github.com/schollz/progressbar/v3"
    "github.com/cdaprod/repocate/internal/container"
)

// rootCmd is the root command for the CLI
var rootCmd = &cobra.Command{
    Use:   "repocate",
    Short: "Repocate is a tool for managing development environments using Docker containers.",
    Long:  `Repocate allows you to clone repositories, create isolated development environments, and manage them using Docker containers.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Check and run default container logic
        handleDefaultContainer()
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
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

    // Coloring usage and available commands
    color.Green("\nRepocate allows you to clone repositories, create isolated development environments, and manage them using Docker containers.")
    color.Blue("\n\nUsage:")
    fmt.Println("  repocate [command]")

    color.Blue("\n\nAvailable Commands:")
    fmt.Println(fmt.Sprintf("  %-12s %s", "clone", "Clone a repository."))
    fmt.Println(fmt.Sprintf("  %-12s %s", "completion", "generate the autocompletion script for the specified shell"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "create", "Clone a repo and create/start a development container."))
    fmt.Println(fmt.Sprintf("  %-12s %s", "enter", "Enter the development container for a specific repo."))
    fmt.Println(fmt.Sprintf("  %-12s %s", "help", "Show help information."))
    fmt.Println(fmt.Sprintf("  %-12s %s", "list", "List all repocate containers."))
    fmt.Println(fmt.Sprintf("  %-12s %s", "rebuild", "Rebuild the development container for a specific repo."))
    fmt.Println(fmt.Sprintf("  %-12s %s", "stop", "Stop the development container for a specific repo."))
    fmt.Println(fmt.Sprintf("  %-12s %s", "version", "Show version information."))

    color.Blue("\n\nFlags:")
    fmt.Println("  -h, --help   help for repocate")

    fmt.Println(color.GreenString("\nUse \"repocate [command] --help\" for more information about a command."))

    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(CreateCmd)
    rootCmd.AddCommand(EnterCmd)
    rootCmd.AddCommand(StopCmd)
    rootCmd.AddCommand(RebuildCmd)
    rootCmd.AddCommand(CloneCmd)
    rootCmd.AddCommand(ListCmd)
    rootCmd.AddCommand(VersionCmd)
    rootCmd.AddCommand(HelpCmd)
}

// Function to handle the default container logic
func handleDefaultContainer() {
    color.Cyan("Checking for 'repocate-default' container...")

    bar := progressbar.NewOptions(100,
        progressbar.OptionSetDescription("[green]Checking container status...[reset]"),
        progressbar.OptionSetWriter(os.Stderr),
        progressbar.OptionShowBytes(false),
        progressbar.OptionShowCount(),
        progressbar.OptionOnCompletion(func() {
            fmt.Fprint(os.Stderr, "\n")
        }),
    )

    // Simulate progress for checking container status
    for i := 0; i < 100; i++ {
        bar.Add(1)
        time.Sleep(10 * time.Millisecond)
    }

    // Check if the default container exists
    containerExists, err := container.CheckContainerExists("repocate-default")
    if err != nil {
        fmt.Println(color.RedString("Error checking container: %s", err))
        os.Exit(1)
    }

    if !containerExists {
        color.Yellow("Default container 'repocate-default' not found. Creating it now...")

        // Show progress for container creation
        bar = progressbar.NewOptions(100,
            progressbar.OptionSetDescription("[green]Creating container...[reset]"),
            progressbar.OptionSetWriter(os.Stderr),
            progressbar.OptionShowBytes(false),
            progressbar.OptionShowCount(),
            progressbar.OptionOnCompletion(func() {
                fmt.Fprint(os.Stderr, "\n")
            }),
        )

        for i := 0; i < 100; i++ {
            bar.Add(1)
            time.Sleep(10 * time.Millisecond)
        }

        err := container.CreateAndStartContainer("repocate-default")
        if err != nil {
            fmt.Println(color.RedString("Error creating default container: %s", err))
            os.Exit(1)
        }
        fmt.Println(color.GreenString("Default container 'repocate-default' created and started."))
    } else {
        color.Green("Default container 'repocate-default' exists. Executing into it now...")

        bar = progressbar.NewOptions(100,
            progressbar.OptionSetDescription("[green]Executing into container...[reset]"),
            progressbar.OptionSetWriter(os.Stderr),
            progressbar.OptionShowBytes(false),
            progressbar.OptionShowCount(),
            progressbar.OptionOnCompletion(func() {
                fmt.Fprint(os.Stderr, "\n")
            }),
        )

        for i := 0; i < 100; i++ {
            bar.Add(1)
            time.Sleep(10 * time.Millisecond)
        }

        err := container.ExecIntoContainer("repocate-default")
        if err != nil {
            fmt.Println(color.RedString("Error executing into default container: %s", err))
            os.Exit(1)
        }
    }
}