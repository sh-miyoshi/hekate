FROM node:14 as builder

WORKDIR /hekate
COPY ./cmd/portal ./cmd/portal
WORKDIR /hekate/cmd/portal
RUN npm install
RUN npm run build

FROM node:14-slim

EXPOSE 3000
WORKDIR /hekate

COPY ./cmd/portal /hekate/
COPY --from=builder /hekate/cmd/portal/node_modules /hekate/node_modules
COPY --from=builder /hekate/cmd/portal/.nuxt /hekate/.nuxt
COPY --from=builder /hekate/cmd/portal/static /hekate/static

CMD ["npm", "run", "start"]
