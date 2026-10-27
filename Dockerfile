# syntax=docker/dockerfile:1

FROM golang:1.26.10-alpine3.24 AS build

# Fail rather than silently selecting a different Go toolchain.
ENV GOTOOLCHAIN=local CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the go files
COPY *.go ./
COPY cmd ./cmd
COPY bibfuse.toml ./

# Build
RUN go test ./... && go vet ./... && \
    go build -o /out/bibfuse ./cmd/bibfuse && \
    go version -m /out/bibfuse

FROM alpine:3.24
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/bibfuse /usr/local/bin/bibfuse
COPY bibfuse.toml ./

ENTRYPOINT ["bibfuse"]
CMD ["-config", "/app/bibfuse.toml"]
