FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN go build -o belsmapi ./cmd/be-lotsanmateo-api/

FROM debian:12-slim

WORKDIR /app

COPY --from=builder /app/belsmapi .

ENV PORT=":8082"
ARG PORT=$PORT

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082
EXPOSE 8083

CMD ["./belsmapi"]