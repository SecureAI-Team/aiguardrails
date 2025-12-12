FROM golang:1.22-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOTOOLCHAIN=go1.24.2
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o aiguardrails ./cmd/api

FROM alpine:3.19
WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /app/aiguardrails /app/aiguardrails
COPY migrations /app/migrations
ENV PORT=8080
EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/aiguardrails"]

