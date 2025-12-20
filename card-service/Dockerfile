# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app ./cmd/card-service

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /build/app .

COPY config.yaml .

EXPOSE 50052

CMD ["./app"]

