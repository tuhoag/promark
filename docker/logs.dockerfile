FROM golang:1.6.1-alpine
ADD . /log
WORKDIR /log

RUN apk update
RUN apk upgrade
RUN apk add --no-cache git

EXPOSE $LOG_PORT

