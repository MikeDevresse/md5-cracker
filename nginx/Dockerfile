FROM node:17.2 as build

WORKDIR /app

COPY ./frontend/package.json .
COPY ./frontend/yarn.lock .

RUN yarn install

COPY ./frontend .

RUN yarn build

FROM nginx:alpine

COPY --from=build /app/build /app
COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf