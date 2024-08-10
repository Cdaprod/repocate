#!/usr/bin/env bash

set -euo pipefail

# Source common functions and variables
source "$(dirname "$0")/repocate-common.sh"

# Function to ensure repo is cloned/updated
ensure_repo() {
    local repo_url=$1
    local dir_name=$(basename "$repo_url" .git)
    local full_path="$WORKSPACE_DIR/$dir_name"
    
    mkdir -p "$WORKSPACE_DIR"
    
    if [[ ! -d "$full_path" ]]; then
        log "INFO" "Cloning repository $repo_url"
        git clone "$repo_url" "$full_path" || error_exit "Failed to clone repository"
        echo -n "Cloning repository... "
        progress_bar 5 20
    else
        log "INFO" "Updating repository $repo_url"
        (cd "$full_path" && git pull) || log "WARN" "Failed to update repository"
        echo -n "Updating repository... "
        progress_bar 2 10
    fi
    
    echo "$full_path"
}

# Function to create and start container
create_container() {
    local repo_url=$1
    local dir_path=$(ensure_repo "$repo_url")
    local container_name=$(get_container_name "$repo_url")
    
    if ! docker ps -a --format '{{.Names}}' | grep -q "^$container_name$"; then
        log "INFO" "Creating new container $container_name"
        echo -n "Creating container... "
        docker run -d \
            -v "$dir_path:/workspace" \
            -v "$HOME/.gitconfig:/root/.gitconfig:ro" \
            -v "$HOME/.ssh:/root/.ssh:ro" \
            -p 3000:3000 -p 50051:50051 \
            -e TERM="$TERM" \
            -e GIT_AUTHOR_NAME="$(git config user.name)" \
            -e GIT_AUTHOR_EMAIL="$(git config user.email)" \
            -e GIT_COMMITTER_NAME="$(git config user.name)" \
            -e GIT_COMMITTER_EMAIL="$(git config user.email)" \
            --name "$container_name" \
            "$BASE_IMAGE" \
            tail -f /dev/null > /dev/null 2>&1 || error_exit "Failed to create container"
        progress_bar 3 15
    elif [[ "$(docker ps -q -f name=$container_name)" ]]; then
        log "INFO" "Container $container_name is already running"
    else
        log "INFO" "Starting existing container $container_name"
        echo -n "Starting container... "
        docker start "$container_name" > /dev/null || error_exit "Failed to start container"
        progress_bar 2 10
    fi
    
    echo -n "Entering container... "
    docker exec -it "$container_name" /bin/zsh -c "cd /workspace && /bin/zsh" || error_exit "Failed to exec into container"
}

# Function to enter container
enter_container() {
    local repo_url=$1
    local container_name=$(get_container_name "$repo_url")
    
    if [[ "$(docker ps -q -f name=$container_name)" ]]; then
        log "INFO" "Entering container $container_name"
        docker exec -it "$container_name" /bin/zsh -c "cd /workspace && /bin/zsh" || error_exit "Failed to exec into container"
    else
        log "WARN" "Container $container_name is not running. Use 'repocate create $repo_url' first."
    fi
}

# Function to stop container
stop_container() {
    local repo_url=$1
    local container_name=$(get_container_name "$repo_url")
    
    if [[ "$(docker ps -q -f name=$container_name)" ]]; then
        log "INFO" "Stopping container $container_name"
        echo -n "Stopping container... "
        docker stop "$container_name" > /dev/null || error_exit "Failed to stop container"
        progress_bar 2 10
        echo "${GREEN}Container stopped successfully${RESET}"
    else
        log "WARN" "Container $container_name is not running"
    fi
}

# Function to rebuild container
rebuild_container() {
    local repo_url=$1
    local container_name=$(get_container_name "$repo_url")
    
    if [[ "$(docker ps -a -q -f name=$container_name)" ]]; then
        log "INFO" "Removing existing container $container_name"
        echo -n "Removing existing container... "
        docker rm -f "$container_name" > /dev/null || error_exit "Failed to remove container"
        progress_bar 2 10
    fi
    
    create_container "$repo_url"
}

# Function to list containers
list_containers() {
    log "INFO" "Listing all repocate containers"
    docker ps -a --filter "name=repocate-" --format "table ${BLUE}{{.Names}}${RESET}\t${GREEN}{{.Status}}${RESET}\t${YELLOW}{{.Ports}}${RESET}"
}

# Function to show version
show_version() {
    echo "${GREEN}Repocate version $VERSION${RESET}"
}

# Function to show usage
usage() {
    cat << EOF
${BLUE}Usage:${RESET} repocate <command> [<repo-url>]

${YELLOW}Commands:${RESET}
  ${GREEN}create <repo-url>${RESET}  Clone repo and create/start dev container
  ${GREEN}enter <repo-url>${RESET}   Enter dev container for repo
  ${GREEN}stop <repo-url>${RESET}    Stop dev container for repo
  ${GREEN}rebuild <repo-url>${RESET} Rebuild dev container for repo
  ${GREEN}list${RESET}               List all dev containers
  ${GREEN}version${RESET}            Show version information
  ${GREEN}help${RESET}               Show this help message

For more information, see ${BLUE}https://github.com/yourusername/repocate${RESET}
EOF
}

# Check prerequisites
check_and_install_prerequisites

# Main script logic
case ${1:-} in
    create)
        [[ $# -eq 2 ]] || error_exit "The 'create' command requires a repository URL"
        create_container "$2"
        ;;
    enter)
        [[ $# -eq 2 ]] || error_exit "The 'enter' command requires a repository URL"
        enter_container "$2"
        ;;
    stop)
        [[ $# -eq 2 ]] || error_exit "The 'stop' command requires a repository URL"
        stop_container "$2"
        ;;
    rebuild)
        [[ $# -eq 2 ]] || error_exit "The 'rebuild' command requires a repository URL"
        rebuild_container "$2"
        ;;
    list)
        list_containers
        ;;
    version)
        show_version
        ;;
    help)
        usage
        ;;
    *)
        usage
        exit 1
        ;;
esac