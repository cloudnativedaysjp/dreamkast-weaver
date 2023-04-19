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

COPY --from=builder /workspace/serve .
USER 65532:65532

ENTRYPOINT ["/serve"]