
FROM  golang:1.11

ADD . project/Astral

RUN export GO11MODULE="on" && cd project/Astral && go mod download && go build .


