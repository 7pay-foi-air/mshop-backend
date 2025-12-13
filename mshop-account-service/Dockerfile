FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY mshop-account-service ./mshop-account-service

WORKDIR /app/mshop-account-service

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o account-service cmd/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/mshop-account-service/account-service .

EXPOSE 8080

CMD ["./account-service"]
