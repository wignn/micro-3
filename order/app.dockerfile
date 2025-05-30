FROM golang:1.24-alpine3.21 AS build
RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /go/src/github.com/wignn/grpc-graphql

COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order

RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order

FROM alpine:3.21

WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]