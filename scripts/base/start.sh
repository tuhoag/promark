#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function startNetwork() {
    local orgNum=$1
    local peerNum=$2
    local logLevel=$3

    infoln "Starting the network"
    infoln $FABRIC_CFG_PATH


    local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"

    infoln $orgNum
    infoln $peerNum
    infoln $docker_compose_path

    FABRIC_LOG=$logLevel COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} up -d --remove-orphans 2>&1

    docker ps -a
    if [ $? -ne 0 ]; then
        fatalln "Unable to start network"
    fi
}

startNetwork $1 $2 $3