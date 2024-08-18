package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "repocate/cmd/repocate"
)

func main() {
    if err := repocate.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}