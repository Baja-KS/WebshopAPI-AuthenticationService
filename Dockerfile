FROM golang:1.16 AS builder

WORKDIR /app

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

FROM alpine:latest as production

WORKDIR /app

COPY --from=builder /app .

CMD ["./app"]
