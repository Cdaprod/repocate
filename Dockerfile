# Start from Ubuntu 22.04 base image
FROM ubuntu:22.04

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Set environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/root/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION=16.15.1
ENV LANG=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8

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

# Install Go
RUN wget https://golang.org/dl/go1.17.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz \
    && rm go1.17.linux-amd64.tar.gz

# Install Go tools
RUN go install github.com/golang/protobuf/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install golang.org/x/tools/gopls@latest \
    && go install github.com/fatih/gomodifytags@latest \
    && go install github.com/cweill/gotests/gotests@latest \
    && go install github.com/josharian/impl@latest

# Install Node.js and npm using nvm
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash \
    && . $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

# Install global npm packages
RUN . $NVM_DIR/nvm.sh \
    && npm install -g typescript \
    && npm install -g eslint \
    && npm install -g prettier \
    && npm install -g ts-node \
    && npm install -g yarn

# Install Rust (for some Neovim plugins)
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Set up Zsh and Oh My Zsh
RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" \
    && git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions \
    && git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

# Copy Zsh configuration files
COPY .config/zsh/.zshrc /root/.zshrc
COPY .config/zsh/custom_plugins.zsh /root/.oh-my-zsh/custom/custom_plugins.zsh

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