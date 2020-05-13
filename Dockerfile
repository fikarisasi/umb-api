FROM golang:1.11

USER nobody

RUN mkdir -p /go/src/github.com/fikarisasi/umb-api
WORKDIR /go/src/github.com/fikarisasi/umb-api

COPY . /go/src/github.com/fikarisasi/umb-api
RUN go build

CMD ["./umb-api"]
