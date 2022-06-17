#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

infoln "Stopping the network"

function stopNetwork() {
    IFS=',' read -r -a orgTypes <<< $1
    local orgNum=$2
    local peerNum=$3
    local logLevel=$4

    local maxOrgId=$(($orgNum - 1))

    infoln "Starting the network"
    local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"

    local allfiles="-f ${DOCKER_COMPOSE_DIR_PATH}/core.yml"

    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrgId); do
            filepath="${DOCKER_COMPOSE_DIR_PATH}/${orgType}${orgId}-${peerNum}.yml"
            allfiles="${allfiles} -f ${filepath}"
        done
    done

    set -x
    FABRIC_LOG=$logLevel COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION NUM_ORGS=$orgNum NUM_PEERS=$peerNum docker-compose $allfiles down -v 2>&1
    { set +x; } 2>/dev/null

    docker ps -a
    if [ $? -ne 0 ]; then
        fatalln "Unable to start network"
    fi
}

stopNetwork $1 $2 $3
