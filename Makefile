
VERSION = $(shell git describe --tags --always)
FLAGS = -ldflags "\
  -X main.VERSION=$(VERSION) \
"

GOBUILD = $(GOCMD) build $(FLAGS)


.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build $(FLAGS) -o bin/goz .

.PHONY: release
release: clean
	GOOS=linux   GOARCH=arm64 go build $(FLAGS) -o bin/goz.linux.arm64 .
	GOOS=linux   GOARCH=amd64 go build $(FLAGS) -o bin/goz.linux.amd64 .
	GOOS=windows GOARCH=arm64 go build $(FLAGS) -o bin/goz.win.arm64 .
	GOOS=windows GOARCH=amd64 go build $(FLAGS) -o bin/goz.win.amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(FLAGS) -o bin/goz.mac.arm64 .
	GOOS=darwin  GOARCH=amd64 go build $(FLAGS) -o bin/goz.mac.amd64 .
	md5sum bin/* > bin/checksum

.PHONY: clean
clean:
	rm bin/*
