FROM golang:1.11

USER nobody

RUN mkdir -p /go/src/github.com/openshift/umb-api
WORKDIR /go/src/github.com/openshift/umb-api

COPY . /go/src/github.com/openshift/umb-api
RUN go build

CMD ["./umb-api"]
