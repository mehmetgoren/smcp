FROM golang:1.18.4

RUN apt-get install -y tzdata

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -buildvcs=false -o main .

CMD ["/app/main"]
