### builder ###
FROM golang:1.20 as builder

WORKDIR /workspace

# Copy the Go Modules
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
# Build
ARG GOOS=linux
ARG GOARCH=amd64
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -trimpath -tags timetzdata -o dkw cmd/main.go

### runner ###
FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /
EXPOSE 8080

COPY internal internal
COPY --from=builder /workspace/dkw .

ENTRYPOINT ["/dkw"]
