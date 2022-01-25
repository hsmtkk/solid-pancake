FROM golang:1.17 AS builder
COPY . /opt
WORKDIR /opt/consumer
RUN go build

FROM gcr.io/distroless/cc-debian11 AS runtime
COPY --from=builder /opt/consumer/consumer /usr/local/bin/consumer
ENTRYPOINT ["/usr/local/bin/consumer"]
