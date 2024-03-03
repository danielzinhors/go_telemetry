FROM golang:1.21.6 as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server_weather cmd/server/main.go

FROM scratch
COPY --from=builder /app/cmd/server/.env .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server_weather .
CMD ["./server_weather"]
