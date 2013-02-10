PREFIX=/home/brian/go-tracker
OUTPUT=go-tracker

build:
	go build

install: build
	mkdir -p $(PREFIX) 
	mkdir -p $(PREFIX)/config
	cp -fR config/ $(PREFIX)/config/
	mkdir -p $(PREFIX)/bin
	cp -f $(OUTPUT) $(PREFIX)/bin

clean:
	go clean


