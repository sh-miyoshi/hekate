FROM golang:1.13 as builder

WORKDIR /myapp
COPY . .
WORKDIR /myapp/cmd/jwt-server
RUN go build -o jwt-server


FROM ubuntu:18.04

WORKDIR /myapp
COPY --from=builder /myapp/cmd/jwt-server/jwt-server /myapp
COPY --from=builder /myapp/cmd/jwt-server/config.yaml /myapp

CMD ["./jwt-server"]
