FROM golang:1.23.0-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add gcc sqlite-dev musl-dev

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/payment-system cmd/payment-system/main.go

FROM alpine

RUN apk --no-cache add sqlite

WORKDIR /app

COPY --from=builder /usr/local/src/bin/payment-system .
COPY --from=builder /usr/local/src/config/ ./config/

VOLUME /app/storage

ENV CONFIG_PATH=/app/config/local.yaml

EXPOSE 8080

CMD ["./payment-system"]