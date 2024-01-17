FROM alpine

LABEL maintainer="czQery"
LABEL authors="czQery"

RUN mkdir /data
WORKDIR /data

COPY /frontend/dist dist
COPY /backend/.config.json config.json
COPY /backend/llg llg

CMD ["./llg"]