# sytntax=docker/dockerfile:1

FROM golang:1.16-alpine

ADD . /go/src/smcp
WORKDIR /go/src/smcp
RUN go get smcp
RUN go install
ENTRYPOINT ["/go/bin/smcp"]