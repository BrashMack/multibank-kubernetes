FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags=netgo \
    -o /app/bin/backend \
    ./cmd/vtbhackback/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
COPY --from=builder /app/bin/backend /app/backend

ENV TZ=Europe/Moscow

WORKDIR /app
USER 65534:65534
EXPOSE 8080

ENTRYPOINT ["/app/backend"]
