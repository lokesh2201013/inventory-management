
FROM golang:1.24.2-alpine AS builder

ENV CGO_ENABLED=0 GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main . && ls -l main


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["./main"]
