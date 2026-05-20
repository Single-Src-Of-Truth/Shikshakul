FROM golang:1.26-alpine AS builder

WORKDIR /workspace

COPY lib/core-go ./lib/core-go
COPY eng/utility-service ./eng/utility-service

RUN go work init && \
    go work use ./lib/core-go && \
    go work use ./eng/utility-service

WORKDIR /workspace/eng/utility-service
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /workspace/eng/utility-service/internal/templates ./internal/templates
COPY --from=builder /workspace/eng/utility-service/main .

EXPOSE 9004

CMD ["./main"]
