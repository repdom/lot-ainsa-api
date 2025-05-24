FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o belsmapi ./cmd/be-lotsanmateo-api/

FROM debian:12-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/belsmapi .

ENV PORT=8080
ARG PORT=$PORT

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082
EXPOSE 8083

CMD ["./belsmapi"]