FROM golang:1.13.1-alpine
RUN mkdir /src
RUN mkdir /src/ext
WORKDIR /src/ext
#ENV APP_USER app
#ENV APP_HOME /code
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates
RUN apk add --no-cache git
RUN git config --global http.sslverify false
RUN apk --update add redis
RUN go get github.com/bwesterb/go-ristretto
RUN go get gopkg.in/redis.v4
RUN go get github.com/gorilla/mux
ENV REDIS_URL redis:6379
EXPOSE $API_PORT

