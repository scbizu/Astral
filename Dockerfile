
FROM  daocloud.io/library/golang:1.8.1

MAINTAINER scnace "scbizu@gmail.com"

ADD . $GOPATH/src/github.com/scbizu/wxgo

RUN go run $GOPATH/src/github.com/scbizu/wxgo/wx.go
