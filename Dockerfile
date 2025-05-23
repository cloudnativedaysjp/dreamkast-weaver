# syntax=docker/dockerfile:1.4
### builder ###
FROM golang:1.24-bullseye as builder

WORKDIR /workspace

# Copy the Go Modules
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

# Build
ARG GOOS=linux
ARG GOARCH=amd64
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -trimpath -tags timetzdata -o dkw cmd/dkw/main.go

### runner ###
FROM gcr.io/distroless/static-debian11:nonroot

LABEL org.opencontainers.image.authors="Yoshitaka Fujii, Hiroki Okui, Kohei Ota"
LABEL org.opencontainers.image.url="https://github.com/cloudnativedaysjp/dreamkast-weaver"
LABEL org.opencontainers.image.source="https://github.com/cloudnativedaysjp/dreamkast-weaver/blob/main/Dockerfile"
EXPOSE 8080

WORKDIR /

COPY --link internal internal
COPY --link --from=builder /workspace/dkw .

ENTRYPOINT ["/dkw","serve"]
