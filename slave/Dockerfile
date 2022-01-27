FROM golang:1.17-alpine as build

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download; \
    go build -ldflags="-w -s" -o slave ./

CMD ["/app/main"]

FROM scratch AS prod

COPY --from=build /app/slave ./

CMD ["/slave"]