#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function startNetwork() {
    local log_level=$1

    infoln "Starting the network"
    infoln $FABRIC_CFG_PATH

    FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} up -d 2>&1

    docker ps -a
    if [ $? -ne 0 ]; then
        fatalln "Unable to start network"
    fi
}

startNetwork $1