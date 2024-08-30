package utils

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/fatih/color"
)

var (
    logFile  *os.File
    logger   *log.Logger
    logLevel = "INFO"
)

// SetupLogger initializes the logging system
func SetupLogger() {
    var err error
    logFile, err = os.OpenFile(getLogFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Could not open log file: %s\n", err)
        os.Exit(1)
    }

    logger = log.New(logFile, "", log.LstdFlags)
}

// getLogFilePath returns the path to the log file
func getLogFilePath() string {
    logDir := filepath.Join(os.Getenv("HOME"), ".config", "repocate", "logs")
    if _, err := os.Stat(logDir); os.IsNotExist(err) {
        os.MkdirAll(logDir, 0755)
    }
    return filepath.Join(logDir, "repocate.log")
}

// EnsureLoggerInitialized ensures that the logger is initialized
func ensureLoggerInitialized() {
    if logger == nil {
        SetupLogger()
    }
}

// Info logs an informational message with a green color
func Info(message string) {
    ensureLoggerInitialized()
    if logLevel == "INFO" || logLevel == "DEBUG" {
        coloredMessage := color.New(color.FgGreen).SprintFunc()(fmt.Sprintf("INFO: %s", message))
        logger.Println(coloredMessage)
    }
}

// Warn logs a warning message with a yellow color
func Warn(message string) {
    ensureLoggerInitialized()
    if logLevel != "ERROR" {
        coloredMessage := color.New(color.FgYellow).SprintFunc()(fmt.Sprintf("WARN: %s", message))
        logger.Println(coloredMessage)
    }
}

// Error logs an error message with a red color and exits the program
func Error(message string) {
    ensureLoggerInitialized()
    coloredMessage := color.New(color.FgRed).SprintFunc()(fmt.Sprintf("ERROR: %s", message))
    logger.Println(coloredMessage)
    os.Exit(1)
}