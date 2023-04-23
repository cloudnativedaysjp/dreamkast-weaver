### builder ###
FROM golang:1.20 as builder

WORKDIR /workspace

# Copy the Go Modules
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags timetzdata -a -o serve cmd/serve/main.go

### runner ###
FROM alpine:3.17.3
WORKDIR /
EXPOSE 8080

# NOTE: Requires home directory at runtime
#   see: https://github.com/ServiceWeaver/weaver/blob/v0.4.0/internal/files/file.go#L111
ARG UID=65532
RUN adduser -D -u ${UID} dk-weaver
COPY --from=builder /workspace/serve .
USER ${UID}

ENTRYPOINT ["/serve"]
