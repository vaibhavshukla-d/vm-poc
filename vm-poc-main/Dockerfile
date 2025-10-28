# syntax=docker/dockerfile:1

FROM golang:1.25.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o vm-server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/vm-server .
EXPOSE 8080
ENTRYPOINT ["/app/vm-server"]
