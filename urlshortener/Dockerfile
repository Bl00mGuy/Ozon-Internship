FROM golang:1.23-alpine AS builder

# Установка Buf (если требуется генерация в контейнере)
RUN apk add --no-cache curl
RUN curl -sSL https://github.com/bufbuild/buf/releases/download/v1.0.0-rc10/buf-Linux-x86_64 -o /usr/local/bin/buf && \
    chmod +x /usr/local/bin/buf

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o urlshortener ./cmd/app

FROM scratch
COPY --from=builder /app/urlshortener /urlshortener
EXPOSE 8080 9090
ENTRYPOINT ["/urlshortener"]
