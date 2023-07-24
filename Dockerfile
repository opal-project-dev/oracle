FROM golang:1.17-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/opal-project-dev/oracle
COPY . .

RUN go mod tidy
RUN go mod vendor

RUN GOOS=linux go build -o /usr/local/bin/chainlink-price-feed github.com/opal-project-dev/oracle


FROM alpine:3.9

RUN apk add --no-cache ca-certificates
COPY --from=buildbase /usr/local/bin/chainlink-price-feed /usr/local/bin/chainlink-price-feed


ENTRYPOINT ["chainlink-price-feed"]
