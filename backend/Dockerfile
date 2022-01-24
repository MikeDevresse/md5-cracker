FROM golang:1.17-alpine as build

RUN mkdir /app

COPY ./backend /app/

WORKDIR /app

COPY docker-compose.yaml ./
COPY .env ./

RUN go mod download; \
    go build -o main ./

CMD ["/app/main"]

FROM alpine AS prod

RUN apk upgrade --no-cache && \
    apk add docker-cli docker-compose

RUN mkdir /app
WORKDIR /app
COPY --from=build /app/main ./
COPY --from=build /app/docker-compose.yaml ./
COPY --from=build /app/.env ./

CMD ["/app/main"]