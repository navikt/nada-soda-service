FROM golang:1.21

WORKDIR /src
COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download
COPY . .
RUN go build -o nada-soda-service .

CMD ["/src/nada-soda-service"]
