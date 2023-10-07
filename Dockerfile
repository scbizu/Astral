
FROM golang:1.21 AS BUILDER

WORKDIR /project/Astral

ADD . /project/Astral

RUN go build -o astral .

FROM golang:1.21

WORKDIR /Astral

COPY --from=BUILDER /project/Astral/astral /Astral/astral

EXPOSE 8443

ENTRYPOINT [ "/Astral/astral" ]
