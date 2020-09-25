FROM golang:1.15 as server-builder

WORKDIR /hekate
COPY go.mod .
COPY go.sum .
COPY pkg ./pkg
COPY cmd/hekate ./cmd/hekate
WORKDIR /hekate/cmd/hekate
RUN go build -o hekate-server


FROM node:14 as front-builder

WORKDIR /hekate
COPY ./cmd/portal ./cmd/portal
WORKDIR /hekate/cmd/portal
RUN npm install
RUN npm run build


FROM ubuntu:18.04

RUN apt update -y && apt install -y nodejs npm
WORKDIR /hekate
RUN mkdir -p server && mkdir -p portal
COPY cmd/hekate/_data/testcerts /hekate/secret

COPY --from=server-builder /hekate/cmd/hekate/hekate-server /hekate/server/hekate-server
COPY build/allinone/config.yaml /hekate/server/config.yaml
COPY cmd/hekate/_data/login /hekate/server/login

COPY ./cmd/portal /hekate/portal/
COPY --from=front-builder /hekate/cmd/portal/node_modules /hekate/portal/node_modules
COPY --from=front-builder /hekate/cmd/portal/.nuxt /hekate/portal/.nuxt
COPY --from=front-builder /hekate/cmd/portal/static /hekate/portal/static
COPY build/allinone/nuxt.config.js /hekate/portal/nuxt.config.js

COPY build/allinone/run.sh .

CMD ["./run.sh"]
