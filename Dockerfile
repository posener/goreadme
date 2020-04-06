FROM golang:1.14.1-alpine3.11

ADD . /home/goreadme
WORKDIR /home/goreadme

RUN go install ./cmd/goreadme

FROM alpine:3.11
COPY --from=0 /go/bin/goreadme /bin/goreadme

RUN echo -e "#! /bin/sh\ngoreadme \$@ > README.md" > /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENV GOREADME_DEBUG 1
ENTRYPOINT [ "/entrypoint.sh" ]