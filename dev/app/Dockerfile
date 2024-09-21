FROM golang:1.22-alpine3.20

RUN apk add -U sudo
RUN go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY dev/app/entrypoint.sh ./dev/app/
RUN chmod +x ./dev/app/entrypoint.sh

ENTRYPOINT ["./dev/app/entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]
