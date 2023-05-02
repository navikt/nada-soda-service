FROM golang:1.20 as builder

WORKDIR /src
COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download
COPY . .
RUN go build -o nada-soda-service .

FROM alpine:3
WORKDIR /app
COPY --from=builder /src/nada-soda-service /app/nada-soda-service

CMD ["/app/nada-soda-service"]
