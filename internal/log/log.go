package log

import (
    "fmt"
    "log"
    "os"
    "path/filepath" 
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

// Info logs an informational message
func Info(message string) {
    if logLevel == "INFO" || logLevel == "DEBUG" {
        logger.Printf("INFO: %s", message)
    }
}

// Warn logs a warning message
func Warn(message string) {
    if logLevel != "ERROR" {
        logger.Printf("WARN: %s", message)
    }
}

// Error logs an error message and exits the program
func Error(message string) {
    logger.Printf("ERROR: %s", message)
    os.Exit(1)
}