DOCKER_ORGANIZATION := ducampsv
DOCKER_IMAGE:= nomad-image-updater


build:
	go build -o dist/nomad-image-updater main.go

format:
	go fmt ./...

dockerbuild:
	docker build  -t $(DOCKER_ORGANIZATION)/$(DOCKER_IMAGE) . -f docker/Dockerfile
