# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /server .
EXPOSE 8080
ENTRYPOINT ["./server"]
