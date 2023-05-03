FROM golang:1.20-alpine3.17

RUN apk add -U sudo
RUN go install github.com/cosmtrek/air@v1.43.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY dev/app/entrypoint.sh ./dev/app/
RUN chmod +x ./dev/app/entrypoint.sh

ENTRYPOINT ["./dev/app/entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]
