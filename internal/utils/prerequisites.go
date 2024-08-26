package prerequisites

import (
    "fmt"
    "os/exec"
    "github.com/cdaprod/repocate/internal/log"
)

// CheckAndInstall checks for prerequisites like Docker and Git
func CheckAndInstall() {
    checkCommand("docker", "--version")
    checkCommand("git", "--version")
}

// checkCommand checks if a command is available
func checkCommand(command string, args ...string) {
    cmd := exec.Command(command, args...)
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("%s is not installed. Please install it and try again.", command))
    }
}