
FROM  daocloud.io/library/golang:1.8.1

MAINTAINER scnace "scbizu@gmail.com"

ADD . $GOPATH/src/github.com/scbizu/Astral

RUN cd $GOPATH/src/github.com/scbizu/Astral && go install

ENTRYPOINT $GOPATH/bin/Astral
