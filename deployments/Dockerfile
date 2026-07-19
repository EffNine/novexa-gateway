# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/gateway ./cmd/gateway

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata sqlite

RUN adduser -D -g '' novexa

WORKDIR /app

COPY --from=builder /app/gateway .

RUN mkdir -p /app/data && chown -R novexa:novexa /app/data

USER novexa

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./gateway"]
