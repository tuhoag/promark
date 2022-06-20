#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

infoln "Stopping the network"

function stopNetwork() {
    IFS=',' read -r -a orgTypes <<< $1
    local orgNum=$2
    local peerNum=$3
    local logLevel=$4

    local maxOrgId=$(($orgNum - 1))
    local maxPeerId=$(($peerNum - 1))

    infoln "Stopping the network"
    # local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"
    set -x
    local coreServicesPath="-f ${DOCKER_COMPOSE_DIR_PATH}/core.yml"
    FABRIC_LOG=$logLevel COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION NUM_ORGS=$orgNum NUM_PEERS=$peerNum docker-compose $coreServicesPath down -v 2>&1
    { set +x; } 2>/dev/null

    # local allfiles="-f ${DOCKER_COMPOSE_DIR_PATH}/core.yml"

    # for orgType in ${orgTypes[@]}; do
    #     for orgId in $(seq 0 $maxOrgId); do
    #         filepath="${DOCKER_COMPOSE_DIR_PATH}/${orgType}${orgId}-${peerNum}.yml"
    #         allfiles="${allfiles} -f ${filepath}"
    #     done
    # done

    baseAdvPort=5000
    basePubPort=6000
    peerPortStep=10
    orgPortStep=100
    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrgId); do
            for peerId in $(seq 0 $maxPeerId); do
                filepath="-f ${DOCKER_COMPOSE_DIR_PATH}/peer.yml"
                orgName="${orgType}${orgId}"
                if [ $orgType = "adv" ]; then
                    peerPort=$((baseAdvPort+peerPortStep*peerId))
                else
                    peerPort=$((basePubPort+peerPortStep*peerId))
                fi
                dbPort=$((peerPort+2))
                apiPort=$((peerPort+1))
                peerName="peer${peerId}.${orgName}"
                set -x
                FABRIC_LOG=$logLevel COMPOSE_PROJECT_NAME=$peerName PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION PEER_ID=$peerId ORG_NAME=$orgName PEER_PORT=$peerPort DB_PORT=$dbPort API_PORT=$apiPort docker-compose  $filepath -p $peerName down -v 2>&1
                { set +x; } 2>/dev/null
            done
        done
    done

    docker ps
}

stopNetwork $1 $2 $3
