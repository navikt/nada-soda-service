FROM golang:1.21 AS builder

WORKDIR /src
COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download

COPY main.go main.go
COPY pkg pkg

RUN go build -o nada-soda-service .

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /src/nada-soda-service /app/nada-soda-service

CMD ["/app/nada-soda-service"]
