FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY cmd/setup.env ./setup.env

RUN CGO_ENABLED=0 GOOS=linux go build -o account-service cmd/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/account-service .
COPY --from=builder /app/setup.env /root/setup.env 

EXPOSE 8080

CMD ["./account-service"]
