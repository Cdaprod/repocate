Here's the fully completed Python refactor using `docker-py`:

```python
import os
import sys
import docker
import argparse
import subprocess
from docker.errors import NotFound
from colorama import Fore, Style, init


# Initialize Docker client
client = docker.from_env()

# Set BASE_IMAGE to repocate-base
BASE_IMAGE = os.getenv('BASE_IMAGE', 'repocate-base')

def ensure_image():
    """Ensure the Docker image exists locally, otherwise build it."""
    try:
        client.images.get(BASE_IMAGE)
        print(f"Using existing image: {BASE_IMAGE}")
    except NotFound:
        print(f"Building Docker image {BASE_IMAGE} from Dockerfile")
        try:
            client.images.build(path=".", tag=BASE_IMAGE)
        except Exception as e:
            print(f"Failed to build Docker image {BASE_IMAGE}: {e}")
            sys.exit(1)

def ensure_user_in_docker_group():
    """Ensure the user is in the Docker group."""
    user_groups = subprocess.check_output(["groups"]).decode()
    if "docker" not in user_groups:
        print("User is not in the Docker group. Adding user to Docker group...")
        subprocess.run(["sudo", "usermod", "-aG", "docker", os.getlogin()])
        print("User added to Docker group. Please log out and log back in or run 'newgrp docker' to apply the changes.")
        sys.exit(1)

def find_free_port():
    """Find a free port on the host system."""
    import socket
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind(("", 0))
        return s.getsockname()[1]

def ensure_repo(repo_url):
    """Clone or update the Git repository."""
    repo_name = os.path.basename(repo_url).replace('.git', '')
    project_dir = os.path.join(os.path.expanduser("~"), "Repocate", repo_name, "source_code")
    os.makedirs(project_dir, exist_ok=True)
    
    if not os.path.isdir(os.path.join(project_dir, ".git")):
        print(f"Cloning repository {repo_url}")
        subprocess.run(["git", "clone", repo_url, project_dir], check=True)
    else:
        print(f"Updating repository {repo_url}")
        subprocess.run(["git", "-C", project_dir, "pull"], check=True)
    
    return project_dir

def init_container(repo_url):
    """Create and start the Docker container."""
    ensure_user_in_docker_group()
    repo_name = os.path.basename(repo_url).replace('.git', '').lower()
    project_dir = ensure_repo(repo_url)
    container_name = f"repocate-{repo_name}"
    volume_name = f"repocate-{repo_name}-vol"

    ensure_image()

    port_3000 = find_free_port()
    port_50051 = find_free_port()

    print(f"Creating new container {container_name}")
    volume = client.volumes.create(name=volume_name)

    try:
        container = client.containers.run(
            BASE_IMAGE,
            name=container_name,
            volumes={
                volume.name: {'bind': '/workspace', 'mode': 'rw'},
                os.path.expanduser("~/.ssh"): {'bind': '/root/.ssh', 'mode': 'ro'},
                os.path.expanduser("~/.gitconfig"): {'bind': '/root/.gitconfig', 'mode': 'ro'}
            },
            ports={
                '3000/tcp': ('0.0.0.0', port_3000),
                '50051/tcp': ('0.0.0.0', port_50051)
            },
            detach=True,
            command="tail -f /dev/null"
        )
        print(f"Container {container_name} created with ID {container.id}")
    except Exception as e:
        print(f"Failed to create container {container_name}: {e}")
        sys.exit(1)

def enter_container(repo_url):
    """Enter the running container."""
    container_name = f"repocate-{os.path.basename(repo_url).replace('.git', '').lower()}"
    try:
        container = client.containers.get(container_name)
        if container.status != "running":
            container.start()
        subprocess.run(["docker", "exec", "-it", container_name, "bash"], check=True)
    except NotFound:
        print(f"Container {container_name} not found. Please create it first.")
    except Exception as e:
        print(f"Failed to enter container {container_name}: {e}")
        sys.exit(1)

def stop_container(repo_url):
    """Stop the running container."""
    container_name = f"repocate-{os.path.basename(repo_url).replace('.git', '').lower()}"
    try:
        container = client.containers.get(container_name)
        if container.status == "running":
            container.stop()
            print(f"Container {container_name} stopped successfully.")
        else:
            print(f"Container {container_name} is not running.")
    except NotFound:
        print(f"Container {container_name} not found.")
    except Exception as e:
        print(f"Failed to stop container {container_name}: {e}")
        sys.exit(1)

def rebuild_container(repo_url):
    """Rebuild the Docker container."""
    container_name = f"repocate-{os.path.basename(repo_url).replace('.git', '').lower()}"
    stop_container(repo_url)
    client.containers.prune()
    init_container(repo_url)

def list_containers():
    """List all running repocate containers."""
    containers = client.containers.list(all=True, filters={"name": "repocate-"})
    for container in containers:
        print(f"{container.name}\t{container.status}\t{container.ports}")

def stop_all_containers():
    """Stop all running repocate containers."""
    containers = client.containers.list(all=True, filters={"name": "repocate-"})
    for container in containers:
        try:
            container.stop()
            print(f"Stopped container {container.name}")
        except Exception as e:
            print(f"Failed to stop container {container.name}: {e}")

def cleanup_containers():
    """Remove all stopped repocate containers."""
    client.containers.prune()
    print("Cleaned up all stopped repocate containers.")

def show_version():
    """Show version of repocate."""
    print("Repocate version 1.0")

def main():
    import argparse
    parser = argparse.ArgumentParser(description="Repocate CLI in Python")
    parser.add_argument("command", choices=["create", "enter", "stop", "rebuild", "list", "stop-all", "cleanup", "version"], help="Command to execute")
    parser.add_argument("repo_url", nargs='?', help="Repository URL (required for create, enter, stop, and rebuild)")
    
    args = parser.parse_args()
    
    if args.command == "create":
        if args.repo_url:
            init_container(args.repo_url)
        else:
            print("The 'create' command requires a repository URL.")
            sys.exit(1)
    elif args.command == "enter":
        if args.repo_url:
            enter_container(args.repo_url)
        else:
            print("The 'enter' command requires a repository URL.")
            sys.exit(1)
    elif args.command == "stop":
        if args.repo_url:
            stop_container(args.repo_url)
        else:
            print("The 'stop' command requires a repository URL.")
            sys.exit(1)
    elif args.command == "rebuild":
        if args.repo_url:
            rebuild_container(args.repo_url)
        else:
            print("The 'rebuild' command requires a repository URL.")
            sys.exit(1)
    elif args.command == "list":
        list_containers()
    elif args.command == "stop-all":
        stop_all_containers()
    elif args.command == "cleanup":
        cleanup_containers()
    elif args.command == "version":
        show_version()
    else:
        print("Invalid command.")
        sys.exit(1)

# Initialize colorama
init(autoreset=True)

def show_ascii_art():
    ascii_art = """
______                           _       
| ___ \                         | |      
| |_/ /___ _ __   ___   ___ __ _| |_ ___ 
|    // _ \ '_ \ / _ \ / __/ _` | __/ _ \\
| |\ \  __/ |_) | (_) | (_| (_| | ||  __/
\_| \_\___| .__/ \___/ \___\__,_|\__\___|
          | |                            
          |_|                            
"""
    print(Fore.CYAN + ascii_art)

def show_usage():
    usage_text = f"""
{Fore.BLUE}Usage:{Style.RESET_ALL} repocate <command> [<repo-url> | <repo_name>]

{Fore.YELLOW}Commands:{Style.RESET_ALL}
  {Fore.GREEN}create <repo-url>{Style.RESET_ALL}  Clone repo and create/start dev container
  {Fore.GREEN}enter <repo-url>{Style.RESET_ALL}   Enter dev container for repo
  {Fore.GREEN}stop <repo-url>{Style.RESET_ALL}    Stop dev container for repo
  {Fore.GREEN}stop-all{Style.RESET_ALL}           Stop all repocate containers
  {Fore.GREEN}cleanup{Style.RESET_ALL}            Clean up all stopped repocate containers
  {Fore.GREEN}rebuild <repo-url>{Style.RESET_ALL} Rebuild dev container for repo
  {Fore.GREEN}list{Style.RESET_ALL}               List all dev containers
  {Fore.GREEN}version{Style.RESET_ALL}            Show version information
  {Fore.GREEN}help{Style.RESET_ALL}               Show this help message

{Fore.YELLOW}Advanced Usage:{Style.RESET_ALL}
  {Fore.GREEN}snapshot{Style.RESET_ALL}           Create a Git snapshot before major changes
  {Fore.GREEN}branch{Style.RESET_ALL}             Create and switch to a dynamic branch
  {Fore.GREEN}commit{Style.RESET_ALL}             Automatically commit and push changes with a dynamic message
  {Fore.GREEN}rollback{Style.RESET_ALL}           Revert to the last known good commit
  {Fore.GREEN}volume{Style.RESET_ALL}             Dynamically create and manage Docker volumes

{Fore.CYAN}Tip:{Style.RESET_ALL} After the first time you use `repocate <repo-url>`, you can simply use `repocate <repo_name>` for quicker access!

For more information, visit {Fore.BLUE}https://github.com/Cdaprod/repocate{Style.RESET_ALL}
"""
    print(usage_text)

def main():
    show_ascii_art()
    show_usage()

if __name__ == "__main__":
    main()
```

### Key Points:
1. **Ensure Image:** `ensure_image()` ensures that the Docker image exists or builds it.
2. **Docker Group:** `ensure_user_in_docker_group()` checks if the user is in the Docker group.
3. **Free Ports:** `find_free_port()` locates available ports to avoid conflicts.
4. **Repo Handling:** `ensure_repo()` handles cloning or updating the repository.
5. **Container Operations:**
   - **Create:** `init_container()` handles creating and starting the container.
   - **Enter:** `enter_container()` allows shell access to the running container.
   - **Stop:** `stop_container()` stops the container.
   - **Rebuild:** `rebuild_container()` stops and re-creates the container.
   - **List:** `list_containers()` lists all `repocate` containers.
   - **Stop All:** `stop_all_containers()` stops all `repocate` containers.
   - **Cleanup:** `cleanup_containers()` removes all stopped containers.

This Python version mirrors the functionality of your original Bash script but leverages `docker-py` for more Pythonic management of Docker resources.

To add color prompts and the same ASCII text in Python, you can use the `colorama` library for colored output and incorporate the ASCII text as a string. Here’s how you can do it:

### 1. Install `colorama`:
If you haven't installed it yet, you can install `colorama` using pip:

```bash
pip install colorama
```

### 2. Use `colorama` in your script:
Here's an example of how you might integrate `colorama` for colored text prompts along with the ASCII text:

### Explanation:
1. **ASCII Art Display:**
   - The `show_ascii_art()` function contains your ASCII art as a string and prints it with a cyan color using `Fore.CYAN`.

2. **Colorized Usage Help:**
   - The `show_usage()` function formats the usage text with different colors to enhance readability. The `Fore` module is used for setting text color, and `Style.RESET_ALL` ensures that after the colored text, the style resets to default.

3. **Automatic Reset:**
   - `init(autoreset=True)` ensures that after each print statement, the color resets automatically, so you don’t need to add reset codes manually.

This Python script will give you a visually appealing command-line interface with colored text and ASCII art similar to what you had in your Bash script.