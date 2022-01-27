FROM golang:1.17-alpine as build

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN apk add binutils; \
    go mod download; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o slave ./; \
    strip slave

CMD ["/app/main"]

FROM scratch AS prod

COPY --from=build /app/slave ./

CMD ["/slave"]