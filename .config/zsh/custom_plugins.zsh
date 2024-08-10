# Load plugins
plugins=(
  git
  docker
  zsh-syntax-highlighting
  zsh-autosuggestions
  npm
  nvm
  zoxide
  history-substring-search
)

# Load the plugins using Oh My Zsh
for plugin ($plugins); do
  source $ZSH/plugins/$plugin/$plugin.zsh
done

# Other custom plugin configurations or additions can go here