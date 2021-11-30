# sytntax=docker/dockerfile:1

FROM golang:1.16-alpine

ADD . /go/src/smcp
WORKDIR /go/src/smcp
RUN go get smcp
RUN go install
RUN mkdir /home/gokalp
RUN mkdir /home/gokalp/Pictures
RUN mkdir /home/gokalp/Pictures/detected
ENTRYPOINT ["/go/bin/smcp"]