FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY auth ./auth
COPY mshop-transactions-service ./mshop-transactions-service

WORKDIR /app/mshop-transactions-service

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o transactions-service cmd/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/mshop-transactions-service/transactions-service .

EXPOSE 8080

CMD ["./transactions-service"]
