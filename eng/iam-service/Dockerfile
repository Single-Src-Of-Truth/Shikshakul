FROM golang:1.26-alpine AS builder

WORKDIR /workspace

COPY lib/core-go ./lib/core-go
COPY eng/iam-service ./eng/iam-service

RUN go work init && \
    go work use ./lib/core-go && \
    go work use ./eng/iam-service

WORKDIR /workspace/eng/iam-service
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /workspace/eng/iam-service/main .

EXPOSE 9000

CMD ["./main"]
