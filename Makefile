.PHONY: build
build:
	go build -o build/gb-emu ./cmd/gb-emu/main.go

all:
	go build -o build/gb-emu ./cmd/gb-emu/main.go
	CGO_ENABLED=0 go build -o build/frame-log-compare ./cmd/frame-log-compare/main.go
	CGO_ENABLED=0 go build -o build/lcd-log-compare ./cmd/lcd-log-compare/main.go
	CGO_ENABLED=0 go build -o build/log-compare ./cmd/log-compare/main.go
