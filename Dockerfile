FROM golang:1.12.9 as builder

WORKDIR /myapp
COPY . .
RUN go build -o jwt-server


FROM ubuntu:18.04

WORKDIR /myapp
COPY --from=builder /myapp/jwt-server /myapp

CMD ["./jwt-server"]
