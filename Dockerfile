FROM golang:1.19 AS builder

WORKDIR /app/
COPY . .

ARG VERSION

RUN go mod download
RUN go build -o bins/terradrift-server -v -ldflags="-X main.version=$VERSION" ./terradrift-server
RUN go build -o bins/terradrift-cli -v -ldflags="-X main.version=$VERSION" ./terradrift-cli

FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/bins /bin/

ENTRYPOINT ["/bin/terradrift-server"]
