FROM golang:1.13.8-alpine3.10

WORKDIR /go/src/app 

RUN apk add make git

RUN go get -u github.com/golang/dep/cmd/dep

ADD . .

RUN make build