# Start from Ubuntu 22.04 base image
FROM ubuntu:22.04

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Set environment variables
ARG GO111MODULE=on
ARG GOROOT=/usr/local/go
ARG GOPATH=/root/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ARG NVM_DIR=/root/.nvm
ARG NODE_VERSION=16.15.1
ENV LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8

# Install basic tools and dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    curl \
    wget \
    zsh \
    neovim \
    ripgrep \
    fd-find \
    tmux \
    fzf \
    unzip \
    python3 \
    python3-pip \
    protobuf-compiler \
    build-essential \
    software-properties-common \
    locales \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set the locale
RUN locale-gen en_US.UTF-8

# Set up the cache directory
RUN mkdir -p /tmp/.buildx-cache

# Install Go for ARM64 (or the correct architecture)
RUN wget https://golang.org/dl/go1.20.3.linux-arm64.tar.gz \
    && tar -C /usr/local -xzf go1.20.3.linux-arm64.tar.gz \
    && rm go1.20.3.linux-arm64.tar.gz

# Install Go Tools
RUN go install github.com/golang/protobuf/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install golang.org/x/tools/gopls@latest \
    && go install github.com/fatih/gomodifytags@latest \
    && go install github.com/cweill/gotests/gotests@latest \
    && go install github.com/josharian/impl@latest \
    || { echo "Go tools installation failed"; exit 1; }

# Install Node.js and npm using nvm
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash \
    && . $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

# Install global npm packages
RUN . $NVM_DIR/nvm.sh \
    && npm install -g typescript eslint prettier ts-node yarn

# Install Rust (for some Neovim plugins)
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Set up persistent history
RUN mkdir -p /root/.persistent_history \
    && touch /root/.persistent_history/.zsh_history \
    && chmod 644 /root/.persistent_history/.zsh_history \
    && echo 'export HISTFILE=/root/.persistent_history/.zsh_history' >> /root/.zshrc \
    && echo 'export HISTSIZE=10000' >> /root/.zshrc \
    && echo 'export SAVEHIST=10000' >> /root/.zshrc \
    && echo 'setopt SHARE_HISTORY' >> /root/.zshrc \
    && echo 'setopt HIST_IGNORE_ALL_DUPS' >> /root/.zshrc

# Set up Zsh and Oh My Zsh
RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" \
    && git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions \
    && git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

# Verify Zsh installation and set fallback
RUN if command -v zsh > /dev/null; then \
        echo "Zsh installed successfully"; \
    else \
        echo "Zsh not found, falling back to Bash"; \
        SHELL ["/bin/bash", "-c"]; \
        CMD ["/bin/bash"]; \
    fi

# Copy Zsh configuration files
COPY .config/zsh/.zshrc /root/.zshrc
COPY .config/zsh/custom_plugins.zsh /root/.oh-my-zsh/custom/custom_plugins.zsh

# Ensure correct ownership and permissions
RUN chown -R root:root /root/.zshrc /root/.oh-my-zsh/custom/custom_plugins.zsh \
    && chmod 644 /root/.zshrc /root/.oh-my-zsh/custom/custom_plugins.zsh

# Verify Zsh configuration
RUN zsh -c "source /root/.zshrc && echo 'Zsh configuration loaded successfully'"

# Install Neovim plugin manager (vim-plug)
RUN curl -fLo "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim

# Copy Neovim configuration
COPY .config/nvim/init.vim /root/.config/nvim/init.vim

# Install Neovim plugins
RUN nvim --headless +PlugInstall +qall

# Set up working directory
WORKDIR /workspace

# Set Zsh as the default shell
SHELL ["/bin/zsh", "-c"]

# Command to run when starting the container
CMD ["/bin/zsh"]