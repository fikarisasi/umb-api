FROM golang:1.14.2

ENV GO111MODULE=on
WORKDIR /go/src/umb_api
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build 

CMD ["./umb_api"]
