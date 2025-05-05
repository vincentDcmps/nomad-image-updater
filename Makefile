build:
	go build -o dist/nomad-image-updater cmd/main.go

format:
	go fmt ./...
