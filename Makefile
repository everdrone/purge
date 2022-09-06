all: build

build:
	tinygo build

build2:
	go build

mkdest:
	sudo mkdir -p /usr/local/bin

install: build mkdest
	sudo cp ./purge /usr/local/bin/purge

install2: build2 mkdest
	sudo cp ./purge /usr/local/bin/purge

clean:
	rm -f ./purge

uninstall:
	sudo rm -f /usr/local/bin/purge

.PHONY: all build install clean uninstall