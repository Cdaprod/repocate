#!/usr/bin/env bash

VERSION="1.2.0"
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/repocate"
CONFIG_FILE="$CONFIG_DIR/config"
LOG_FILE="$CONFIG_DIR/repocate.log"

# Color definitions
if command -v tput > /dev/null && tput setaf 1 > /dev/null 2>&1; then
    RED=$(tput setaf 1)
    GREEN=$(tput setaf 2)
    YELLOW=$(tput setaf 3)
    BLUE=$(tput setaf 4)
    RESET=$(tput sgr0)
else
    RED=""
    GREEN=""
    YELLOW=""
    BLUE=""
    RESET=""
fi

# Ensure config directory exists
mkdir -p "$CONFIG_DIR"

# Load configuration
if [[ -f "$CONFIG_FILE" ]]; then
    source "$CONFIG_FILE"
else
    BASE_IMAGE=${BASE_IMAGE:-"repocate-env"}
    WORKSPACE_DIR=${WORKSPACE_DIR:-"$HOME/repocate-workspaces"}
    echo "BASE_IMAGE=$BASE_IMAGE" > "$CONFIG_FILE"
    echo "WORKSPACE_DIR=$WORKSPACE_DIR" >> "$CONFIG_FILE"
fi

# Logging function
log() {
    local level=$1
    shift
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [$level] $*" >> "$LOG_FILE"
    case "$level" in
        ERROR)   echo "${RED}ERROR: $*${RESET}" >&2 ;;
        WARN)    echo "${YELLOW}WARN: $*${RESET}" ;;
        INFO)    echo "${GREEN}INFO: $*${RESET}" ;;
        DEBUG)   [[ "${DEBUG:-}" == "true" ]] && echo "${BLUE}DEBUG: $*${RESET}" ;;
    esac
}

# Error handling function
error_exit() {
    log "ERROR" "$1"
    exit 1
}

# Progress bar function
progress_bar() {
    local duration=$1
    local steps=$2
    local sleep_time=$(bc <<< "scale=3; $duration / $steps")
    local progress=0
    echo -n "["
    for ((i=0; i<steps; i++)); do
        echo -n " "
    done
    echo -n "]"
    echo -ne "\r["
    for ((i=0; i<=steps; i++)); do
        progress=$((100 * i / steps))
        echo -ne "${GREEN}#${RESET}"
        sleep "$sleep_time"
    done
    echo -ne "\r[${GREEN}"
    for ((i=0; i<steps; i++)); do
        echo -n "#"
    done
    echo -ne "${RESET}] ${progress}%\n"
}

# Check and install prerequisites
check_and_install_prerequisites() {
    local prerequisites=("docker" "git" "wget" "bc")
    for prereq in "${prerequisites[@]}"; do
        if ! command -v "$prereq" &> /dev/null; then
            log "WARN" "$prereq is not installed. Attempting to install..."
            if command -v apt-get &> /dev/null; then
                sudo apt-get update && sudo apt-get install -y "$prereq" || error_exit "Failed to install $prereq"
            elif command -v yum &> /dev/null; then
                sudo yum install -y "$prereq" || error_exit "Failed to install $prereq"
            elif command -v brew &> /dev/null; then
                brew install "$prereq" || error_exit "Failed to install $prereq"
            else
                error_exit "Unsupported package manager. Please install $prereq manually."
            fi
            log "INFO" "$prereq installed successfully"
        fi
    done
}

# Function to get container name from repo URL
get_container_name() {
    basename=$(basename "$1" .git)
    echo "repocate-${basename}"
}