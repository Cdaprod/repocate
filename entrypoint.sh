#!/bin/sh

# Base directory for configurations
CONFIG_DIR="/root/.config"

# Symlink function to handle different configurations
symlink_configs() {
  local source_dir=$1
  local target_dir=$2

  if [ ! -d "$target_dir" ]; then
    mkdir -p "$target_dir"
  fi

  for file in $(ls "$source_dir"); do
    ln -sf "$source_dir/$file" "$target_dir/$file"
  done
}

# Setup Zsh configuration
if [ -d "$CONFIG_DIR/zsh" ]; then
  echo "Setting up Zsh configuration..."
  ln -sf "$CONFIG_DIR/zsh/custom_plugins.zsh" "/custom_plugins.zsh"
fi

# Ensure .zshrc exists
if [ ! -f "/root/.zshrc" ]; then
  touch /root/.zshrc
fi

# Setup Neovim configuration
if [ -d "$CONFIG_DIR/nvim" ]; then
  echo "Setting up Neovim configuration..."
  symlink_configs "$CONFIG_DIR/nvim" "/root/.config/nvim"
fi

# Setup VSCode configuration
if [ -d "$CONFIG_DIR/vscode" ]; then
  echo "Setting up VSCode configuration..."
  symlink_configs "$CONFIG_DIR/vscode" "/root/.config/Code/User"
fi

# Setup MetaGPT configuration
if [ -d "$CONFIG_DIR/metagpt" ]; then
  echo "Setting up MetaGPT configuration..."
  mkdir -p "/root/.metagpt"
  symlink_configs "$CONFIG_DIR/metagpt/config" "/root/.metagpt/config"
  symlink_configs "$CONFIG_DIR/metagpt/prompts" "/root/.metagpt/prompts"
  symlink_configs "$CONFIG_DIR/metagpt/roles" "/root/.metagpt/roles"
  symlink_configs "$CONFIG_DIR/metagpt/tools" "/root/.metagpt/tools"
  symlink_configs "$CONFIG_DIR/metagpt/workspace" "/root/.metagpt/workspace"
fi

# Setup Repocate configuration
if [ -d "$CONFIG_DIR/repocate" ]; then
  echo "Setting up Repocate configuration..."
  mkdir -p "/root/.repocate"
  ln -sf "$CONFIG_DIR/repocate/repocate.json" "/root/.repocate/repocate.json"
fi

echo "Configuration setup complete."

# Execute the command passed to the entrypoint
exec "$@"