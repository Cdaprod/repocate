package utils

import (
    "fmt"
    "time"
)

// ProgressBar displays a simple progress bar
func ProgressBar(duration time.Duration, steps int) {
    sleepTime := duration / time.Duration(steps)
    fmt.Print("[")
    for i := 0; i < steps; i++ {
        time.Sleep(sleepTime)
        fmt.Print("=")
    }
    fmt.Println("]")
}