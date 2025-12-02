FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY auth ./auth
COPY mshop-articles-service ./mshop-articles-service

WORKDIR /app/mshop-articles-service

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o articles-service cmd/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/mshop-articles-service/articles-service .

EXPOSE 8080

CMD ["./articles-service"]
