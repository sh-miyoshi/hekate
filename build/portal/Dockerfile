FROM node as builder

WORKDIR /myapp
COPY . .
WORKDIR /myapp/cmd/portal
RUN npm install
RUN npm run build

FROM node:13-alpine

WORKDIR /myapp
ENV HOST=0.0.0.0

ADD cmd/portal/package.json ./
ADD cmd/portal/nuxt.config.js ./

COPY --from=builder /myapp/cmd/portal/node_modules ./node_modules/
COPY --from=builder /myapp/cmd/portal/.nuxt ./.nuxt/
COPY --from=builder /myapp/cmd/portal/static ./static/

EXPOSE 3000
CMD ["npm", "run", "start"]
