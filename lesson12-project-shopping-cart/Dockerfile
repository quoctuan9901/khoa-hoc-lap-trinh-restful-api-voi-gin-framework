# ======= STAGE 1: BUILD GO BINARY =========
FROM golang:1.24.5-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

# ======= STAGE 2: RUNTIME CONTAINER =========
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .

COPY .env .

EXPOSE 8080

CMD ["./main"]