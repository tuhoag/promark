FROM hyperledger/fabric-peer:2.2

RUN apk update
RUN apk upgrade
RUN apk add --no-cache git
RUN apk --update add redis
RUN apk add go
# RUN go get github.com/bwesterb/go-ristretto
# RUN go get gopkg.in/redis.v4

ENV REDIS_URL 127.0.0.1:6379