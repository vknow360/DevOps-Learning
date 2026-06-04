# Stage 1: build using heavy base image with all required tools
FROM golang:1.25.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o main .

# Stage 2: run using small base image
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/main .

CMD ["./main"]