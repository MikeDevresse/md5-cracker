FROM golang:1.17-alpine as build

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download; \
    CGO_ENABLED=0 go build -o main ./

CMD ["/app/main"]

FROM busybox AS prod

RUN mkdir /app
WORKDIR /app
COPY --from=build /app/main ./

CMD ["/app/main"]