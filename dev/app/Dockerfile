FROM golang:1.20-alpine3.17

RUN go install github.com/cosmtrek/air@v1.43.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
