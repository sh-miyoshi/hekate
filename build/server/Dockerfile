FROM golang:1.13 as builder

WORKDIR /myapp
COPY . .
WORKDIR /myapp/cmd/hekate
RUN go build -o hekate


FROM ubuntu:18.04

WORKDIR /myapp
COPY --from=builder /myapp/cmd/hekate/hekate /myapp
COPY --from=builder /myapp/cmd/hekate/config.yaml /myapp

CMD ["./hekate"]
