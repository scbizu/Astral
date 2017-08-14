
FROM  daocloud.io/library/golang:1.8.1

MAINTAINER scnace "scbizu@gmail.com"

ADD . $GOPATH/src/astral

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep init

RUN go run $GOPATH/src/astral/wx.go
