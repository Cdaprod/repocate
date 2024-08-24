#!/bin/bash

# Used in Dockerfile.minimal

# Set the base directory for configuration plugins
CONFIG_BASE="/root/.config"

# Ensure the base configuration directory exists
mkdir -p $CONFIG_BASE

# Function to handle plugin directory binding
bind_plugin_directory() {
    local plugin_dir=$1
    local target_dir=$CONFIG_BASE/$(basename $plugin_dir)

    # Check if the plugin directory exists in the base directory
    if [ ! -d "$target_dir" ]; then
        echo "Creating plugin directory: $target_dir"
        mkdir -p $target_dir
    fi

    # Sync plugin directory with target
    echo "Syncing plugin directory: $plugin_dir to $target_dir"
    cp -r $plugin_dir/* $target_dir/
}

# Iterate over all directories in .config and bind them
for plugin in /workspace/.config/*; do
    if [ -d "$plugin" ]; then
        bind_plugin_directory $plugin
    fi
done

# Setup and bind each plugin config
for dir in metagpt nvim vscode zsh; do
    if [ -d "/workspace/.config/$dir" ]; then
        bind_plugin_directory "/workspace/.config/$dir"
    fi
done

# Start the shell
exec "$@"