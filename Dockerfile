
FROM  daocloud.io/library/golang:1.8.1

MAINTAINER scnace "scbizu@gmail.com"

ADD . $GOPATH/src/github.com/scbizu/Astral

RUN go install Astral

ENTRYPOINT $GOPATH/bin/Astral
