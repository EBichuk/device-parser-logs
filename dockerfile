FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app cmd/main.go

RUN adduser -D -u 1000 appuser

RUN mkdir -p /app/data/device-reports /app/device-storage

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./app"]