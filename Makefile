# Makefile
# Test goxterm
# By J. Stuart McMurray
# Created 20241128
# Last Modified 20241128

all: test

test:
	go test ./...
	go vet ./...
	staticcheck ./...
