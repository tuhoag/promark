#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function runService() {
	local log_level=$1

    #FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} run --service-ports external-db.promark.com & 2>&1

    FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} run --service-ports external.promark.com /bin/sh 2>&1
}

runService $1