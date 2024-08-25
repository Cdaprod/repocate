package repocate

import (
    "fmt"
    "github.com/fatih/color"
)

// PrintBanner displays the ASCII art banner and additional information
func PrintBanner() {
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
    color.Blue("\n\nUsage:")
    fmt.Println("  repocate [command]")

    color.Blue("\n\nAvailable Commands:")
    fmt.Println(fmt.Sprintf("  %-12s %s", "start", "Initialize and start the default Repocate container"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "clone", "Clone a repository"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "create", "Clone a repo and create/start a development container"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "enter", "Enter the development container for a specific repo"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "help", "Show help information"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "list", "List all repocate containers"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "rebuild", "Rebuild the development container for a specific repo"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "stop", "Stop the development container for a specific repo"))
    fmt.Println(fmt.Sprintf("  %-12s %s", "version", "Show version information"))

    color.Blue("\n\nFlags:")
    fmt.Println("  -h, --help   help for repocate")

    fmt.Println(color.GreenString("\nUse \"repocate [command] --help\" for more information about a command."))
}
