#!/bin/sh

redis-server --daemonize yes
cd /src/
go mod tidy
go run verifierService.go &
cd /opt/gopath/src/github.com/hyperledger/fabric/peer/
peer node start
