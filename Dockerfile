
FROM golang:1.13 AS BUILDER

WORKDIR /project/Astral

ADD . /project/Astral

RUN export GO11MODULE="on" && go build -o astral .

FROM golang:1.13

WORKDIR /Astral

COPY --from=BUILDER /project/Astral/astral /Astral/astral


