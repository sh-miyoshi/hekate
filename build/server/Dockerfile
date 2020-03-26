FROM golang:1.13 as builder

WORKDIR /myapp
COPY go.mod .
COPY go.sum .
COPY pkg ./pkg
COPY cmd/hekate ./cmd/hekate
WORKDIR /myapp/cmd/hekate
RUN go build -o hekate-server


FROM ubuntu:18.04

WORKDIR /myapp
COPY --from=builder /myapp/cmd/hekate/hekate-server /myapp/server/hekate-server
COPY build/allinone/config.yaml /myapp/server/config.yaml
COPY cmd/hekate/_data/testcerts /myapp/server/secret
COPY cmd/hekate/_data/login /myapp/server/login

CMD ["./hekate-server", "--help"]
