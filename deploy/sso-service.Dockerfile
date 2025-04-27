FROM golang:1.24.2 AS builder
WORKDIR /app

COPY . .

WORKDIR /app/sso-service/cmd

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main
RUN chmod +x main

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sso-service/cmd/main .
EXPOSE 8080
CMD ["./main"]