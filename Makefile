.DEFAULT_GOAL := build
SHELL := /bin/bash

fmt:
	go fmt ./...
.PHONY:fmt

vet: fmt
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	go vet ./...
	shadow ./...
.PHONY:vet

build: vet
	go build -o clnr -ldflags "-w -s" main.go
.PHONY:build

install:
	mkdir /usr/local/clnr
	cp -r locales /usr/local/clnr
	cp clnr /usr/local/clnr
	ln -sf /usr/local/clnr/clnr /usr/bin/clnr
.PHONY:install

uninstall:
	rm -r /usr/local/clnr
	rm /usr/bin/clnr
.PHONY:uninstall