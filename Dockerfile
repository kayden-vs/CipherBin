FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /main ./cmd/web

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser -D -u 1000 appuser
WORKDIR /home/appuser

COPY --from=builder /main . 
COPY --from=builder /app/ui ./ui

# TLS removed - Render handles this
USER appuser

EXPOSE 8080
CMD ["./main"]