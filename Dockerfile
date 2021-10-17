# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.16-alpine

RUN set -ex && \
  apk add --no-cache \
      gcc \
      musl-dev && \
  rm -rf /var/cache/apk/*

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the go files
COPY *.go ./
ADD cmd ./cmd

# Build
WORKDIR /app/cmd/bibfuse
RUN go build -o /docker-bibfuse

# Copy the config file
WORKDIR /app
COPY bibfuse.toml ./

ENTRYPOINT ["/docker-bibfuse"]
CMD ["-config", "/app/bibfuse.toml"]
