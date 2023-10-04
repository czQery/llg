FROM alpine

LABEL maintainer="czQery"
LABEL authors="czQery"

RUN mkdir /data
WORKDIR /data

COPY /client/dist dist
COPY /server/.config.json config.json
COPY /server/llg llg

CMD ["./llg"]