
FROM  golang:1.10

MAINTAINER scnace "scbizu@gmail.com"

ADD . $GOPATH/src/github.com/scbizu/Astral

RUN cd $GOPATH/src/github.com/scbizu/Astral && go install

ENTRYPOINT $GOPATH/bin/Astral -p

EXPOSE 8080
