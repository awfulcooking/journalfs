.PHONY: default install

now := $(shell date +'%Y-%m-%d%S%S%z')
version := $(shell git describe --always --dirty)

PREFIX ?= /usr/local

default: journalfs

journalfs:
	        go build -ldflags '-X main.buildTime=${now} -X main.version=${version}' ./cmd/journalfs

install: journalfs
	        install -m 755 journalfs $(PREFIX)/bin
