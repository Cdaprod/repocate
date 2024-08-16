#!/usr/bin/env bash

set -euo pipefail

# Source common functions and variables
source "$(dirname "$0")/repocate-common.sh"

# Set BASE_IMAGE to repocate-base
BASE_IMAGE="${BASE_IMAGE:-repocate-base}"

# Build the Docker image if it doesn't exist
if ! docker image inspect "$BASE_IMAGE" > /dev/null 2>&1; then
    echo "Building Docker image $BASE_IMAGE from Dockerfile"
    docker build -t "$BASE_IMAGE" . || error_exit "Failed to build Docker image $BASE_IMAGE"
fi

# Define REPOCATE_WORKSPACE if not set
REPOCATE_WORKSPACE="${REPOCATE_WORKSPACE:-$HOME/Repocate}"
mkdir -p "$REPOCATE_WORKSPACE"

# Function to ensure the user is in the Docker group
ensure_user_in_docker_group() {
    if ! groups $USER | grep -q "\bdocker\b"; then
        echo "User is not in the Docker group. Adding user to Docker group..."
        sudo usermod -aG docker $USER
        echo "User added to Docker group. Please log out and log back in or run 'newgrp docker' to apply the changes."
        exit 1
    fi
}

# Function to create a dynamic branch name based on the current time
create_dynamic_branch() {
    local branch_name="repocate-$(date +%Y%m%d-%H%M%S)"
    git checkout -b "$branch_name"
    log "INFO" "Created and switched to new branch: $branch_name"
}

# Function to dynamically commit and push changes
dynamic_commit_and_push() {
    local commit_message="${1:-"Automated commit by repocate on $(date)"}"
    git add .
    git commit -m "$commit_message"
    git push origin $(git rev-parse --abbrev-ref HEAD)
    log "INFO" "Committed and pushed changes with message: $commit_message"
}

# Function to create a snapshot (Git tag) before significant changes
create_snapshot() {
    local snapshot_name="snapshot-$(date +%Y%m%d-%H%M%S)"
    git tag "$snapshot_name"
    log "INFO" "Created snapshot: $snapshot_name"
}

# Function to find a free port on the host
find_free_port() {
    local port
    while true; do
        port=$(shuf -i 2000-65000 -n 1)  # Generate a random port number between 2000 and 65000
        if ! netstat -tuln | grep -q ":$port "; then  # Check if the port is free
            echo "$port"
            return
        fi
    done
}

# Updated ensure_repo function
ensure_repo() {
    local repo_url=$1
    local repo_name=$(basename "$repo_url" .git)
    local project_dir="$REPOCATE_WORKSPACE/$repo_name/source_code"
    
    mkdir -p "$REPOCATE_WORKSPACE/$repo_name/container_configs"
    mkdir -p "$project_dir"
    
    if [[ ! -d "$project_dir/.git" ]]; then
        log "INFO" "Cloning repository $repo_url"
        git clone --verbose "$repo_url" "$project_dir" || error_exit "Failed to clone repository"
        echo -n "Cloning repository... "
        progress_bar 5 20
    else
        log "INFO" "Updating repository $repo_url"
        (cd "$project_dir" && git pull --verbose) || log "WARN" "Failed to update repository"
        echo -n "Updating repository... "
        progress_bar 2 10
    fi
    
    # Create a dynamic branch after cloning/updating
    (cd "$project_dir" && create_dynamic_branch)
    
    echo "$project_dir"
}

# Function to create and start container with dynamic port and volume management
# init_container() {
#     ensure_user_in_docker_group
#     local repo_url=$1
#     local repo_name
#     local project_dir
#     local container_name
#     local volume_name
#     local port_3000
#     local port_50051
#     local container_id

#     # Extract repository name and sanitize it for use as a container name
#     repo_name=$(basename "$repo_url" .git | tr '[:upper:]' '[:lower:]' | tr -cd '[:alnum:]-')

#     # Set up project directory and container/volume names
#     project_dir=$(ensure_repo "$repo_url")
#     container_name="${repo_name}"
#     volume_name="repocate-${repo_name}-vol"

#     # Check that BASE_IMAGE is set
#     if [[ -z "${BASE_IMAGE:-}" ]]; then
#         error_exit "BASE_IMAGE is not set. Please set it before running the script."
#     fi
#     log "INFO" "Using BASE_IMAGE: $BASE_IMAGE"

#     # Pull the Docker image if it's not available locally
#     if ! docker image inspect "$BASE_IMAGE" > /dev/null 2>&1; then
#         log "INFO" "Pulling Docker image $BASE_IMAGE"
#         docker pull "$BASE_IMAGE" || error_exit "Failed to pull Docker image $BASE_IMAGE"
#     fi

#     # Find free ports to bind to the container
#     port_3000=$(find_free_port)
#     port_50051=$(find_free_port)
#     log "INFO" "Using ports: $port_3000 (for 3000) and $port_50051 (for 50051)"

#     # Check if the Docker volume already exists, and create it if necessary
#     if docker volume ls --format '{{.Name}}' | grep -q "^$volume_name$"; then
#         log "INFO" "Volume $volume_name already exists. Using existing volume."
#     else
#         log "INFO" "Creating Docker volume: $volume_name"
#         docker volume create "$volume_name" || error_exit "Failed to create Docker volume $volume_name"
#     fi

#     # Run the Docker container and capture its ID
#     log "INFO" "Creating new container $container_name"
#     container_id=$(docker run -d \
#         -v "$volume_name:/workspace" \
#         -v "$HOME/.ssh:/root/.ssh:ro" \
#         -v "$HOME/.gitconfig:/root/.gitconfig:ro" \
#         -p "$port_3000:3000" \
#         -p "$port_50051:50051" \
#         --name "$container_name" \
#         "$BASE_IMAGE" \
#         tail -f /dev/null)

#     # Check if the container creation was successful
#     if [ $? -ne 0 ] || [ -z "$container_id" ]; then
#         error_exit "Failed to create container $container_name"
#     fi
#     log "INFO" "Container $container_name created successfully with ID $container_id"

#     # Verify that the container is listed and recognized by Docker
#     if ! docker ps -a --format '{{.Names}}' | grep -q "^$container_name$"; then
#         error_exit "Failed to confirm the creation of container $container_name"
#     fi
#     log "INFO" "Container $container_name created with ports $port_3000:3000 and $port_50051:50051"

#     # Optionally show a progress bar or delay for effect (can be removed if not needed)
#     progress_bar 3 15
# }

init_container() {
    ensure_user_in_docker_group
    local repo_url=$1
    local repo_name
    local project_dir
    local container_name
    local volume_name
    local port_3000
    local port_50051
    local container_id

    # Extract repository name and sanitize it for use as a container name
    repo_name=$(basename "$repo_url" .git | tr '[:upper:]' '[:lower:]' | tr -cd '[:alnum:]-')

    # Set up project directory and container/volume names
    project_dir=$(ensure_repo "$repo_url")
    container_name="${repo_name}"
    volume_name="repocate-${repo_name}-vol"

    # Check that BASE_IMAGE is set
    if [[ -z "${BASE_IMAGE:-}" ]]; then
        error_exit "BASE_IMAGE is not set. Please set it before running the script."
    fi
    log "INFO" "Using BASE_IMAGE: $BASE_IMAGE"

    # Pull the Docker image if it's not available locally
    if ! docker image inspect "$BASE_IMAGE" > /dev/null 2>&1; then
        log "INFO" "Pulling Docker image $BASE_IMAGE"
        docker pull "$BASE_IMAGE" || error_exit "Failed to pull Docker image $BASE_IMAGE"
    fi

    # Find free ports to bind to the container
    port_3000=$(find_free_port)
    port_50051=$(find_free_port)
    log "INFO" "Using ports: $port_3000 (for 3000) and $port_50051 (for 50051)"

    # Check if the Docker volume already exists, and create it if necessary
    if docker volume ls --format '{{.Name}}' | grep -q "^$volume_name$"; then
        log "INFO" "Volume $volume_name already exists. Using existing volume."
    else
        log "INFO" "Creating Docker volume: $volume_name"
        docker volume create "$volume_name" || error_exit "Failed to create Docker volume $volume_name"
    fi

    # Copy project files to the Docker volume
    log "INFO" "Copying project files to Docker volume"
    docker run --rm \
        -v "$volume_name:/workspace" \
        -v "$project_dir:/source" \
        "$BASE_IMAGE" \
        /bin/bash -c "cp -r /source/. /workspace/" || error_exit "Failed to copy files to Docker volume"

    # Run the Docker container and capture its ID
    log "INFO" "Creating new container $container_name"
    container_id=$(docker run -d \
        -v "$volume_name:/workspace" \
        -v "$HOME/.ssh:/root/.ssh:ro" \
        -v "$HOME/.gitconfig:/root/.gitconfig:ro" \
        -p "$port_3000:3000" \
        -p "$port_50051:50051" \
        --name "$container_name" \
        "$BASE_IMAGE" \
        tail -f /dev/null)

    # Check if the container creation was successful
    if [ $? -ne 0 ] || [ -z "$container_id" ]; then
        error_exit "Failed to create container $container_name"
    fi
    log "INFO" "Container $container_name created successfully with ID $container_id"

    # Verify that the container is listed and recognized by Docker
    if ! docker ps -a --format '{{.Names}}' | grep -q "^$container_name$"; then
        error_exit "Failed to confirm the creation of container $container_name"
    fi
    log "INFO" "Container $container_name created with ports $port_3000:3000 and $port_50051:50051"

    # Optionally show a progress bar or delay for effect (can be removed if not needed)
    progress_bar 3 15
}

# Function to enter container
enter_container() {
    ensure_user_in_docker_group
    local repo_url=$1
    local repo_name=$(basename "$repo_url" .git | tr '[:upper:]' '[:lower:]' | tr -cd '[:alnum:]-')
    local container_name="${repo_name}"

    # Check if the container is running, and if not, attempt to start it
    if [[ "$(docker ps -q -f name=$container_name)" ]]; then
        log "INFO" "Entering container $container_name"
        docker exec -it "$container_name" /bin/bash -c "cd /workspace && exec /bin/zsh || exec /bin/bash" || error_exit "Failed to exec into container"
    else
        log "WARN" "Container $container_name is not running. Starting container..."
        if docker start "$container_name" > /dev/null; then
            log "INFO" "Successfully started container $container_name"
            docker exec -it "$container_name" /bin/bash -c "cd /workspace && exec /bin/zsh || exec /bin/bash" || error_exit "Failed to exec into container"
        else
            error_exit "Failed to start container $container_name"
        fi
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

# Function to rebuild container with snapshot and dynamic branching
rebuild_container() {
    local repo_url=$1
    local container_name=$(get_container_name "$repo_url")
    
    # Check if the provided argument is a valid repository URL or a container name
    if [[ "$repo_url" == *"/"* ]]; then
        # Assuming it's a URL
        create_snapshot
        stop_container "$repo_url"
        echo "INFO: Removing existing container $container_name"
        docker rm -f "$container_name" > /dev/null || error_exit "Failed to remove container"
        init_container "$repo_url"
    else
        # Treat as a container name, skip the repository cloning step
        echo "INFO: Rebuilding container $container_name without cloning a repository"
        stop_container "$repo_url"
        docker rm -f "$container_name" > /dev/null || error_exit "Failed to remove container"
        init_container "$repo_url" # You may need to modify init_container function to handle this case
    fi
}

# Function to list containers
list_containers() {
    log "INFO" "Listing all repocate containers"
    docker ps -a --filter "name=repocate-" --format "table ${BLUE}{{.Names}}${RESET}\t${GREEN}{{.Status}}${RESET}\t${YELLOW}{{.Ports}}${RESET}"
}

# Function to stop all containers
stop_all_containers() {
    log "INFO" "Stopping all repocate containers..."
    docker stop $(docker ps -a -q --filter "name=repocate-") > /dev/null || log "WARN" "No containers found to stop."
    echo "${GREEN}All containers stopped successfully${RESET}"
}

# Function to clean up stopped containers
cleanup_containers() {
    log "INFO" "Cleaning up all stopped repocate containers..."
    docker rm $(docker ps -a -q --filter "status=exited" --filter "name=repocate-") > /dev/null || log "WARN" "No stopped containers found to clean up."
    echo "${GREEN}All stopped containers cleaned up successfully${RESET}"
}

# Function to show version
show_version() {
    echo "${GREEN}Repocate version $VERSION${RESET}"
}

# Function to better log error_exit
error_exit() {
    log "ERROR" "$1"
    echo "Last 10 Docker logs:"  # Example for Docker-specific error logging
    docker logs --tail 10 "$container_name" || echo "No logs available or container not found"
    exit 1
}

# Function to show usage
usage() {
    cat << "EOF"

______                           _       
| ___ \                         | |      
| |_/ /___ _ __   ___   ___ __ _| |_ ___ 
|    // _ \ '_ \ / _ \ / __/ _` | __/ _ \
| |\ \  __/ |_) | (_) | (_| (_| | ||  __/
\_| \_\___| .__/ \___/ \___\__,_|\__\___|
          | |                            
          |_|                            

By: David Cannan (Cdaprod)

EOF

    cat << EOF
${BLUE}Usage:${RESET} repocate <command> [<repo-url> | <repo_name>]

${YELLOW}Commands:${RESET}
  ${GREEN}create <repo-url>${RESET}  Clone repo and create/start dev container
  ${GREEN}enter <repo-url>${RESET}   Enter dev container for repo
  ${GREEN}stop <repo-url>${RESET}    Stop dev container for repo
  ${GREEN}stop-all${RESET}           Stop all repocate containers
  ${GREEN}cleanup${RESET}            Clean up all stopped repocate containers
  ${GREEN}rebuild <repo-url>${RESET} Rebuild dev container for repo
  ${GREEN}list${RESET}               List all dev containers
  ${GREEN}version${RESET}            Show version information
  ${GREEN}help${RESET}               Show this help message

${YELLOW}Advanced Usage:${RESET}
  ${GREEN}snapshot${RESET}           Create a Git snapshot before major changes
  ${GREEN}branch${RESET}             Create and switch to a dynamic branch
  ${GREEN}commit${RESET}             Automatically commit and push changes with a dynamic message
  ${GREEN}rollback${RESET}           Revert to the last known good commit
  ${GREEN}volume${RESET}             Dynamically create and manage Docker volumes

${CYAN}Tip:${RESET} After the first time you use \`repocate <repo-url>\`, you can simply use \`repocate <repo_name>\` for quicker access!

For more information, visit ${BLUE}https://github.com/Cdaprod/repocate${RESET}
EOF
}

# Check prerequisites
check_and_install_prerequisites
# Other functions (enter_container, stop_container, etc.) remain unchanged

# Main script logic
case ${1:-} in
    create)
        [[ $# -eq 2 ]] || error_exit "The 'create' command requires a repository URL"
        init_container "$2"
        ;;
    enter)
        [[ $# -eq 2 ]] || error_exit "The 'enter' command requires a repository URL"
        enter_container "$2"
        ;;
    stop)
        [[ $# -eq 2 ]] || error_exit "The 'stop' command requires a repository URL"
        stop_container "$2"
        ;;
    stop-all)
        stop_all_containers
        ;;
    cleanup)
        cleanup_containers
        ;;
    rebuild)
        [[ $# -eq 2 ]] || error_exit "The 'rebuild' command requires a repository URL"
        rebuild_container "$2"
        ;;
    list)
        ensure_user_in_docker_group
        list_containers
        ;;
    version)
        show_version
        ;;
    help)
        usage
        ;;
    *)
        # Default behavior: enter the container for the given repo_url
        if [[ $# -eq 1 ]]; then
            enter_container "$1"
        else
            usage
            exit 1
        fi
        ;;
esac