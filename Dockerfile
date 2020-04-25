FROM golang:1.14.1-alpine3.11

ADD . /home/goreadme
WORKDIR /home/goreadme
RUN go install ./cmd/goreadme

FROM alpine:3.11
RUN apk add git curl jq
COPY --from=0 /go/bin/goreadme /bin/goreadme

ADD .github/entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]