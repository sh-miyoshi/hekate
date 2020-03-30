FROM golang:1.13 as server-builder

WORKDIR /myapp
COPY go.mod .
COPY go.sum .
COPY pkg ./pkg
COPY cmd/hekate ./cmd/hekate
WORKDIR /myapp/cmd/hekate
RUN go build -o hekate-server


FROM node:13.10 as front-builder

WORKDIR /myapp
COPY ./cmd/portal ./cmd/portal
WORKDIR /myapp/cmd/portal
RUN npm install


FROM ubuntu:18.04

RUN apt update -y && apt install -y nodejs npm
WORKDIR /myapp
RUN mkdir -p server && mkdir -p portal

COPY --from=server-builder /myapp/cmd/hekate/hekate-server /myapp/server/hekate-server
COPY build/allinone/config.yaml /myapp/server/config.yaml
COPY cmd/hekate/_data/testcerts /myapp/server/secret
COPY cmd/hekate/_data/login /myapp/server/login

COPY ./cmd/portal /myapp/portal/
COPY --from=front-builder /myapp/cmd/portal/node_modules /myapp/portal/node_modules

COPY build/allinone/run.sh .

CMD ["./run.sh"]