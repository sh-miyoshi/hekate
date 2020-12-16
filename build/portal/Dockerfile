FROM node:14 as builder

WORKDIR /hekate
COPY ./cmd/portal ./cmd/portal
WORKDIR /hekate/cmd/portal
RUN npm install --production


FROM node:14-slim

EXPOSE 3000
WORKDIR /hekate

COPY cmd/hekate/_data/testcerts /hekate/secret

COPY --from=builder /hekate/cmd/portal /hekate/portal
COPY build/portal/nuxt.config.js /hekate/portal/nuxt.config.js

COPY build/portal/run.sh /hekate/run.sh


CMD ["./run.sh"]
