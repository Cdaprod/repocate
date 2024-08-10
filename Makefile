# Define default installation paths
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man/man1
REPOCATE_DIR ?= ~/Repocate

# Define source and target files
SOURCES = repocate repocate-common.sh repocate.1
TARGETS = $(addprefix $(DESTDIR)$(BINDIR)/, $(basename $(SOURCES)))

.PHONY: all install uninstall clean setup_env

# Default target that compiles everything
all: setup_env $(SOURCES)

# Install binaries and man pages
install: $(SOURCES)
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 repocate $(DESTDIR)$(BINDIR)/repocate
	install -m 644 repocate-common.sh $(DESTDIR)$(BINDIR)/repocate-common.sh
	install -d $(DESTDIR)$(MANDIR)
	install -m 644 repocate.1 $(DESTDIR)$(MANDIR)/repocate.1

# Setup environment in ~/Repocate directory
setup_env:
	@echo "Setting up the Repocate environment..."
	# Create Repocate directory structure
	mkdir -p $(REPOCATE_DIR)/.config/zsh
	mkdir -p $(REPOCATE_DIR)/.config/nvim
	mkdir -p $(REPOCATE_DIR)/workspace
	# Move configuration files into Repocate directory
	mv .zshrc $(REPOCATE_DIR)/.config/zsh/.zshrc || true
	mv init.vim $(REPOCATE_DIR)/.config/nvim/init.vim || true
	@echo "Environment setup complete in $(REPOCATE_DIR)"

# Uninstall binaries and man pages
uninstall:
	rm -f $(DESTDIR)$(BINDIR)/repocate
	rm -f $(DESTDIR)$(BINDIR)/repocate-common.sh
	rm -f $(DESTDIR)$(MANDIR)/repocate.1

# Clean up any temporary files and the Repocate directory
clean:
	rm -f $(TARGETS)
	rm -rf $(REPOCATE_DIR)