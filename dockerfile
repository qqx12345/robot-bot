FROM golang:1.24.5-alpine AS builder

WORKDIR /robot-bot

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /robot-bot/myapp .

EXPOSE 8080

CMD ["./myapp"]