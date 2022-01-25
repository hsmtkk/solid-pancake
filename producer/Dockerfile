FROM golang:1.17 AS builder
COPY . /opt
WORKDIR /opt/producer
RUN go build

FROM gcr.io/distroless/cc-debian11 AS runtime
COPY --from=builder /opt/producer/producer /usr/local/bin/producer
ENTRYPOINT ["/usr/local/bin/producer"]
