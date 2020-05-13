FROM golang:latest

USER nobody

WORKDIR /go/src/umb_api

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# COPY . /go/src/github.com/openshift/umb_api
# RUN go env
# RUN go env -w GO111MODULE=on
RUN go build

CMD ["./umb_api"]
