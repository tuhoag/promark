#!/bin/sh

redis-server --daemonize yes
cd /src/
go mod tidy
go run ver.go &
cd /opt/gopath/src/github.com/hyperledger/fabric/peer/
peer node start
