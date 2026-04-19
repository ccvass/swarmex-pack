FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /swarmex-pack ./cmd

FROM alpine:3.21
RUN apk add --no-cache docker-cli
COPY --from=build /swarmex-pack /usr/local/bin/swarmex-pack
ENTRYPOINT ["swarmex-pack"]
