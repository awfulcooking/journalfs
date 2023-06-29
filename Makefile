.PHONY: default install

now := $(shell date +'%Y-%m-%d%S%S%z')
version := $(shell git describe --always --dirty --tags)

PREFIX ?= /usr/local

default: journalfs

journalfs:
	        go build -ldflags '-X main.buildTime=${now} -X main.version=${version}' ./cmd/journalfs

install: journalfs
	        install -Dm 755 journalfs $(PREFIX)/bin/journalfs
	        install -Dm 755 contrib/journalfs.service $(PREFIX)/usr/lib/systemd/system/journalfs.service
