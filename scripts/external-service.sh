#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function runService() {
	local log_level=$1
    local mode=$2
    local orgNum=$3
    local peerNum=$4
    local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"


    if test $mode -eq 0; then
        FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} run --service-ports external.promark.com /bin/sh 2>&1
    elif test $mode -eq 1; then
        FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} run --service-ports peer0.bus0.promark.com /bin/sh 2>&1
    elif test $mode -eq 2; then
        FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} run --service-ports peer0.adv0.promark.com /bin/sh 2>&1
    elif test $mode -eq 3; then
        FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} run --service-ports client.promark.com /bin/sh 2>&1
    fi
}

runService $1 $2 $3 $4