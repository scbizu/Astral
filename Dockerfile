
FROM golang:1.12 AS BUILDER

ADD . /project/Astral

RUN export GO11MODULE="on" && cd /project/Astral && go mod download && go build .

FROM alpine:lastest

WORKDIR /Astral

COPY --from=BUILDER /project/Astral/Astral .


