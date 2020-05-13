FROM golang:1.11.1-alpine3.8 as build-env

USER nobody

RUN mkdir /umb_api
WORKDIR /umb_api
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# COPY . /go/src/github.com/openshift/umb_api
# RUN go env
# RUN go env -w GO111MODULE=on
RUN go build

CMD ["./umb_api"]
