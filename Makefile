SERVICE := "zenith"
GOOS = darwin
GOARCH = amd64

server:
	cd internal/server && go run .

client:
	cd internal/client && go run .

.PHONY: server client