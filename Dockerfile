# syntax=docker/dockerfile:1.4
### builder ###
FROM golang:1.20-bullseye as builder

WORKDIR /workspace

# Copy the Go Modules
COPY --link go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .
# Build
ARG GOOS=linux
ARG GOARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -trimpath -tags timetzdata -o dkw cmd/main.go

### runner ###
FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /
EXPOSE 8080

COPY --link internal internal
COPY --link --from=builder /workspace/dkw .

ENTRYPOINT ["/dkw"]
