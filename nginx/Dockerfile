ARG BACKEND_URL
ARG BACKEND_PATH

FROM node:17.2 as build

ENV REACT_APP_BACKEND_URL = $BACKEND_URL
ENV REACT_APP_BACKEND_PATH = $BACKEND_PATH

WORKDIR /app

COPY ./frontend/package.json .
COPY ./frontend/yarn.lock .

RUN yarn install

COPY ./frontend .

RUN yarn build

FROM nginx:alpine

COPY --from=build /app/build /app
COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf