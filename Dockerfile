# Debe coincidir con la version de github.com/a-h/templ en go.mod.
ARG TEMPL_VERSION=v0.3.1100

# 1) Generacion: compila los .templ a codigo Go (views/*_templ.go).
#    La imagen de templ no trae shell, por eso el RUN va en forma "exec".
FROM ghcr.io/a-h/templ:${TEMPL_VERSION} AS generate
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

# 2) Compilacion
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY --from=generate /app .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server_app ./server/main.go

# 3) Imagen final
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server_app ./server
COPY static ./static

EXPOSE 8080

CMD ["./server"]
