#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function runService() {
	local log_level=$1
    local mode=$2

    if test $mode -eq 0; then
        FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} run --service-ports external.promark.com /bin/sh 2>&1
    elif test $mode -eq 1; then
        FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} run --service-ports peer0.bus0.promark.com /bin/sh 2>&1
    elif test $mode -eq 2; then
        FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} run --service-ports verifier2.promark.com /bin/sh 2>&1
    fi
}

runService $1 $2