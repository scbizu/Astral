
FROM  daocloud.io/library/golang:1.8.1

MAINTAINER scnace "scbizu@gmail.com"

ADD . $GOPATH/src/github.com/scbizu/astral

RUN go run $GOPATH/src/github.com/scbizu/astral/wx.go
