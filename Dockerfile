FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server_app ./server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server_app ./server
COPY static ./static

EXPOSE 8080

CMD ["./server"]