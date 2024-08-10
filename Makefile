PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man/man1

.PHONY: all install uninstall

all: repocate repocate-common.sh repocate.1

install: repocate repocate-common.sh repocate.1
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 repocate $(DESTDIR)$(BINDIR)/repocate
	install -m 644 repocate-common.sh $(DESTDIR)$(BINDIR)/repocate-common.sh
	install -d $(DESTDIR)$(MANDIR)
	install -m 644 repocate.1 $(DESTDIR)$(MANDIR)/repocate.1

uninstall:
	rm -