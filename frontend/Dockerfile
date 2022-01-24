FROM node:17.2-alpine as dev

WORKDIR /app

COPY ./package.json .
COPY ./yarn.lock .

CMD yarn install && yarn dev-server

FROM dev as prod

COPY . .

CMD yarn build
