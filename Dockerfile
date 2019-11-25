FROM golang:1.12.9 as builder

WORKDIR /myapp
COPY . .
WORKDIR /myapp/cmd/server
RUN go build -o jwt-server


FROM ubuntu:18.04

WORKDIR /myapp
COPY --from=builder /myapp/cmd/server/jwt-server /myapp
COPY --from=builder /myapp/cmd/server/config.yaml /myapp
RUN mkdir -p /myapp/db

CMD ["./jwt-server"]
