FROM golang:1.6.1-alpine
ADD . /code
WORKDIR /code
#ENV APP_USER app
#ENV APP_HOME /code
RUN apk update
RUN apk upgrade
RUN apk add --no-cache git
RUN apk --update add redis 
RUN go get github.com/bwesterb/go-ristretto
RUN go get gopkg.in/redis.v4
ENV REDIS_URL redis:6379
EXPOSE $API_PORT

