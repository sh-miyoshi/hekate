FROM golang:1.15 as builder

WORKDIR /hekate
COPY go.mod .
COPY go.sum .
COPY pkg ./pkg
COPY cmd/hekate ./cmd/hekate
WORKDIR /hekate/cmd/hekate
RUN go build -o hekate-server


FROM ubuntu:18.04

COPY --from=builder /hekate/cmd/hekate/hekate-server /bin/hekate-server
RUN mkdir -p /etc/hekate
COPY build/server/config.yaml /etc/hekate/config/default.yaml
COPY cmd/hekate/_data/testcerts /etc/hekate/secret
COPY cmd/hekate/_data/login /etc/hekate/login

CMD ["hekate-server", "--config", "/etc/hekate/config/default.yaml"]
