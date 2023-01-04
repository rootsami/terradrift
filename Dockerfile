FROM golang:1.19 AS builder

WORKDIR /app/
COPY . .

RUN make deps
RUN make build


FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/bins /bin/

ENTRYPOINT ["/bin/terradrift-server"]
