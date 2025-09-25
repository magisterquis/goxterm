# Makefile
# Test goxterm
# By J. Stuart McMurray
# Created 20241128
# Last Modified 20250926

all: test

test:
	go test -timeout 3s ./...
	go vet ./...
	staticcheck ./...
