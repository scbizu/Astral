
FROM  daocloud.io/library/golang:1.8.1

MAINTAINER scnace "scbizu@gmail.com"

ADD . $GOPATH/src/astral

RUN go run $GOPATH/src/astral/wx.go
