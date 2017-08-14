
FROM  daocloud.io/library/golang:1.8.1

MAINTAINER scnace "scbizu@gmail.com"

ADD . $GOPATH/src/Astral

RUN go install Astral

ENTRYPOINT $GOPATH/bin/Astral
