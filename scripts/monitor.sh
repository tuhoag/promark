#!/bin/bash

# This script uses the logspout and http stream tools to let you watch the docker containers
# in action.
#
# More information at https://github.com/gliderlabs/logspout/tree/master/httpstream

. $SCRIPTS_DIR/utils.sh

echo Starting monitoring on all containers on the network $NETWORK_NAME

echo $NETWORK_NAME
echo $LOGSPOUT_PORT

docker kill logspout 2> /dev/null 1>&2 || true
docker rm logspout 2> /dev/null 1>&2 || true

docker run -d --name="logspout" \
	--volume=/var/run/docker.sock:/var/run/docker.sock \
	--publish=127.0.0.1:${LOGSPOUT_PORT}:80 \
	--network  ${NETWORK_NAME} \
	gliderlabs/logspout
sleep 3
curl http://127.0.0.1:${LOGSPOUT_PORT}/logs