FROM golang:1.15 as go-builder

WORKDIR /hekate
COPY go.mod .
COPY go.sum .
COPY pkg ./pkg
COPY cmd/portal-building ./cmd/portal-building
COPY cmd/hekate ./cmd/hekate
WORKDIR /hekate/cmd/hekate
RUN go build -o hekate-server
WORKDIR /hekate/cmd/portal-building
RUN go build -o portal-building

FROM node:14 as front-builder

WORKDIR /hekate
COPY ./cmd/portal ./cmd/portal
WORKDIR /hekate/cmd/portal
RUN npm install --production


FROM mongo as mongo-installer


FROM ubuntu:18.04

RUN apt update -y && apt install -y curl
RUN curl -fsSL https://deb.nodesource.com/setup_14.x | bash -
RUN apt install -y nodejs
WORKDIR /hekate
RUN mkdir -p server && mkdir -p portal && mkdir -p mongo
COPY cmd/hekate/_data/testcerts /hekate/secret

COPY --from=mongo-installer /usr/bin/mongod mongo/mongod

COPY --from=go-builder /hekate/cmd/hekate/hekate-server /hekate/server/hekate-server
COPY build/allinone/config.yaml /hekate/server/config.yaml
COPY cmd/hekate/_data/login /hekate/server/login

COPY --from=front-builder /hekate/cmd/portal /hekate/portal
COPY build/allinone/nuxt.config.js /hekate/portal/nuxt.config.js
COPY --from=go-builder /hekate/cmd/portal-building/portal-building /hekate/portal/portal-building

COPY build/allinone/run.sh .

CMD ["./run.sh"]
