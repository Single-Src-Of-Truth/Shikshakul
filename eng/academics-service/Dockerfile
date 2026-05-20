FROM golang:1.26-alpine AS builder

WORKDIR /workspace

COPY lib/core-go ./lib/core-go
COPY eng/academics-service ./eng/academics-service

RUN go work init && \
    go work use ./lib/core-go && \
    go work use ./eng/academics-service

WORKDIR /workspace/eng/academics-service
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /workspace/eng/academics-service/main .

EXPOSE 9001

CMD ["./main"]
