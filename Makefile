all: build

build:
	tinygo build

install: build
	sudo cp ./purge /usr/local/bin/purge

clean:
	rm -f ./purge

uninstall:
	sudo rm -f /usr/local/bin/purge

.PHONY: all build install clean uninstall