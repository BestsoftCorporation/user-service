FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o main .

FROM alpine:latest
RUN apk add --no-cache mongodb-tools

WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080

CMD ["./main"]