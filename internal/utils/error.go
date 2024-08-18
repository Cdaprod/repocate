package utils

import "fmt"

// CheckError handles errors and panics if any error is encountered
func CheckError(err error) {
    if err != nil {
        panic(fmt.Sprintf("An error occurred: %s", err))
    }
}