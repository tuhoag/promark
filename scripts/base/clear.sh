#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh


infoln "Cleaning the repository"


function clear() {
    local orgNum=$1
    local peerNum=$2
    local logLevel=$3
    local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"

    # infoln $docker_compose_path

    $BASE_SCRIPTS_DIR/stop.sh $orgNum $peerNum $logLevel

    # remove organizations
    rm -rf $CREDENTIALS_OUTPUTS

    # remove volumes
    rm -rf volumes

    # remove channels
    rm -rf ${CHANNEL_PATH}

    # remove log
    rm -rf $PWD/log.txt

    # remove the chaincode package file before commit
    # rm chaincode/main.tar.gz

    # docker rm -f $(docker ps -aq) || true
    # docker rmi -f $(docker images -a -q) || true
    # docker volume rm $(docker volume ls)  || true

    docker ps -a
    if [ $? -ne 0 ]; then
        fatalln "Unable to start network"
    fi
}

clear $1 $2 $3

