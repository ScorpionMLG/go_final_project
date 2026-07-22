FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/main /app/main

COPY --from=builder /app/web /app/web

EXPOSE 7540

ENV TODO_PORT=7540

CMD ["/app/main"]