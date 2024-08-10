# Path to your Oh My Zsh installation.
export ZSH="$HOME/.oh-my-zsh"

# Set name of the theme to load
ZSH_THEME="robbyrussell"

# Uncomment the following line to enable command auto-correction.
ENABLE_CORRECTION="true"

# Load custom plugins and configurations
source "$ZSH_CUSTOM/custom_plugins.zsh"

# Enable color support for 'ls' and add handy aliases
if [ -x /usr/bin/dircolors ]; then
    test -r ~/.dircolors && eval "$(dircolors -b ~/.dircolors)" || eval "$(dircolors -b)"
    alias ls='ls --color=auto'
    alias ll='ls -alF -h'
    alias la='ls -A'
    alias l='ls -CF'
    alias ttt='tree -a -L 3'
    alias tt='tree -a -L 2'
    alias t='tree -a -L 1'

    alias grep='grep --color=auto'
    alias fgrep='fgrep --color=auto'
    alias egrep='egrep --color=auto'
fi

# Alias for 'cfg' with custom Git directory and work-tree
alias cfg='/usr/bin/git --git-dir=$HOME/.cfg/ --work-tree=$HOME'

# Alert alias for notifying after a long-running command
alias alert='notify-send --urgency=low -i "$([ $? = 0 ] && echo terminal || echo error)" "$(history | tail -n1 | sed -e '\''s/^\s*[0-9]\+\s*//;s/[;&|]\s*alert$//'\'')"'

# Load Oh My Zsh
source $ZSH/oh-my-zsh.sh