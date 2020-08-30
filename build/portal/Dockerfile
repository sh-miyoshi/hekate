FROM node:13.10 as builder

WORKDIR /hekate
COPY . .
WORKDIR /hekate/cmd/portal
RUN npm install

FROM node:13-alpine

EXPOSE 3000
WORKDIR /hekate

COPY ./cmd/portal /hekate/
COPY --from=builder /hekate/cmd/portal/node_modules /hekate/node_modules

COPY build/portal/run.sh /hekate/run.sh

CMD ["./run.sh"]
