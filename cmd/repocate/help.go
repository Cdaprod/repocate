// cmd/repocate/help
package repocate

import (
    "github.com/spf13/cobra"
    "github.com/faith/color"
    "fmt"
)

var HelpCmd = &cobra.Command{
    Use:   color.BlueString("help"),
    Short: color.MagentaString("Show help information."),
    Run: func(cmd *cobra.Command, args []string) {
        // Custom help message with colors
        fmt.Println(color.CyanString("Repocate CLI Tool"))
        fmt.Println(color.GreenString("Available Commands:"))
        
        for _, c := range cmd.Commands() {
            fmt.Printf("%s\t%s\n", color.BlueString(c.Use), c.Short)
        }

        fmt.Println(color.GreenString("Flags:"))
        fmt.Println("  -h, --help   ", color.YellowString("help for repocate"))
        
        fmt.Println(color.GreenString("Use"), color.BlueString("repocate [command] --help"), color.GreenString("for more information about a command."))
    },
}