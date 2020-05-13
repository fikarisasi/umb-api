FROM golang:latest

USER nobody

WORKDIR /go/src/umb_api
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# COPY . /go/src/github.com/openshift/umb_api
# RUN go env
# RUN go env -w GO111MODULE=on
RUN go build

CMD ["./umb_api"]
