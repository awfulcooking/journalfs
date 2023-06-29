.PHONY: default install

now := $(shell date +'%Y-%m-%d%S%S%z')
version := $(shell git describe --always --dirty)

PREFIX ?= /usr/local

default: journalfs

journalfs:
	        go build -o journalfs cmd/journalfs

install: journalfs
	        install -m 755 journalfs $(PREFIX)/bin
