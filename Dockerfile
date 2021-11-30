# sytntax=docker/dockerfile:1

FROM golang:1.16-alpine

ADD . /go/src/smcp
WORKDIR /go/src/smcp
RUN go get smcp
RUN go install
RUN mkdir /go/src/smcp/images
ENTRYPOINT ["/go/bin/smcp"]