PREFIX=/home/brian/go-tracker
OUTPUT=go-tracker

build:
	go build

install:
	mkdir -p $(PREFIX) 
	mkdir -p $(PREFIX)/config
	cp -Rn config/ $(PREFIX)/config/
	mkdir -p $(PREFIX)/bin
	cp -f $(OUTPUT) $(PREFIX)/bin

clean:
	go clean


