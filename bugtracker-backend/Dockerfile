# ------------ Build Stage ------------
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/bugtracker/main.go


# ------------ Run Stage ------------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

ENV PORT=10000

EXPOSE 10000

CMD ["./server"]
