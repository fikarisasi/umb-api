FROM golang:1.11

USER nobody

RUN mkdir -p /go/src/github.com/openshift/umb_api
WORKDIR /go/src/github.com/openshift/umb_api

COPY . /go/src/github.com/openshift/umb_api
RUN go mod init
RUN go mod tidy
RUN go build

CMD ["./umb_api"]
