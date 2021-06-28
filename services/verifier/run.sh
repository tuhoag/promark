#!/bin/sh

redis-server --daemonize yes
go run /src/ver.go &
peer node start
