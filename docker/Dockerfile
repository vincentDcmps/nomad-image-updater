FROM golang:alpine AS builder
COPY . /app 
WORKDIR /app
RUN apk add --no-cache make
RUN make build

FROM alpine

COPY --from=builder /app/dist/nomad-image-updater /app/nomad-image-updater
VOLUME [ "/src" ]
WORKDIR "/src"
CMD ["/app/nomad-image-updater","update", "."]
