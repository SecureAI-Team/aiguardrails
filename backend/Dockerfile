FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o aiguardrails ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /app/aiguardrails /app/aiguardrails
COPY migrations /app/migrations
ENV PORT=8080 TZ=Asia/Shanghai
EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/aiguardrails"]


