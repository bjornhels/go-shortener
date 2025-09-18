FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/index.html .

RUN chown -R appuser:appgroup /root

USER appuser

EXPOSE 80

ENV BASE_URL=http://localhost
ENV PORT=80

CMD ["./main"]
