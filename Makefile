

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o bin/goz .

.PHONY: release
release:
	GOOS=linux GOARCH=arm64 go build -o bin/goz.linux.arm64 .
	GOOS=linux GOARCH=amd64 go build -o bin/goz.linux.amd64 .
	GOOS=windows GOARCH=arm64 go build -o bin/goz.win.arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/goz.win.amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/goz.mac.arm64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/goz.mac.amd64 .

.PHONY: clean
clean:
	rm bin/*
