# Start from Ubuntu 22.04 base image
FROM ubuntu:22.04

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Install basic tools and dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends git curl wget zsh && \
    apt-get install -y --no-install-recommends neovim ripgrep fd-find tmux fzf unzip && \
    apt-get install -y --no-install-recommends python3 python3-pip python3-dev protobuf-compiler && \
    apt-get install -y --no-install-recommends build-essential software-properties-common locales ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Set up the cache directory
RUN mkdir -p /tmp/.buildx-cache

# Set environment variables
ARG GO111MODULE=on
ARG GOROOT=/usr/local/go
ARG GOPATH=/root/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ARG NVM_DIR=/root/.nvm
ARG NODE_VERSION=16.15.1
ENV LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8

# Install Go for the correct architecture
ARG GO_VERSION=1.20.3
ARG TARGETARCH
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-${TARGETARCH}.tar.gz \
    && rm go${GO_VERSION}.linux-${TARGETARCH}.tar.gz

# Install Node.js and npm using nvm
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash \
    && . $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

# Install VS Code server and extensions
RUN . $NVM_DIR/nvm.sh \
    && npm install -g @coder/code-server \
    && code-server --install-extension github.copilot \
    && code-server --install-extension ms-python.python \
    && code-server --install-extension VisualStudioExptTeam.vscodeintellicode \
    && code-server --install-extension eamodio.gitlens \
    && code-server --install-extension esbenp.prettier-vscode \
    && code-server --install-extension ms-azuretools.vscode-docker \
    && code-server --install-extension ms-vsliveshare.vsliveshare \
    && code-server --install-extension streetsidesoftware.code-spell-checker \
    && code-server --install-extension coenraads.bracket-pair-colorizer-2

# Mount VS Code settings and extensions
COPY .config/vscode/settings.json /root/.config/vscode/settings.json
COPY .config/vscode/extensions.json /root/.config/vscode/extensions.json

# Final setup for running VS Code server
EXPOSE 3000 8080 9229
CMD ["code-server", "--bind-addr", "0.0.0.0:8080", "."]