### builder ###
FROM golang:1.20 as builder

WORKDIR /workspace

# Copy the Go Modules
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags tools -a -o initdb tools/initdb/main.go

### runner ###
FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /

COPY --from=builder /workspace/initdb .
COPY internal internal

ENTRYPOINT ["/initdb"]
