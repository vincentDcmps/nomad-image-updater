FROM golang AS builder
COPY . /app 
WORKDIR /app
RUN make build

FROM alpine

RUN apk add --no-cache libc6-compat
COPY --from=builder /app/dist/nomad-image-updater /app/nomad-image-updater
VOLUME [ "/src" ]
CMD ["/app/nomad-image-updater", "/src"]
