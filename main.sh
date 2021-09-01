#!/bin/bash

. $PWD/settings.sh

export CHANNEL_NAME="mychannel"
export LOG_LEVEL=INFO
export FABRIC_LOGGING_SPEC=DEBUG
export CHAINCODE_NAME="campaign"

function initialize() {
    # generate all organizations
    $SCRIPTS_DIR/gen-orgs.sh

    # generate genesis-block
    $SCRIPTS_DIR/gen-genesis-block.sh
}

function createChannel() {
    $SCRIPTS_DIR/gen-channel-tx.sh $CHANNEL_NAME
    $SCRIPTS_DIR/gen-channel.sh $CHANNEL_NAME "adv" 0
}

function joinChannel() {
     # args: $CHANNEL_NAME <org type> <number of org> <number of peer>
    $SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "adv" 1 1
    # $SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "adv" 1 2
    $SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "bus" 1 1
    # $SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "bus" 0 2
}

function networkUp() {
    $SCRIPTS_DIR/start.sh $LOG_LEVEL
}

function networkDown() {
    # docker rm -f logspout

    $SCRIPTS_DIR/stop.sh $LOG_LEVEL
}

function clear() {
    $SCRIPTS_DIR/clear.sh
}

function monitor() {
    $SCRIPTS_DIR/monitor.sh
}

function packageChaincode() {
    $SCRIPTS_DIR/package-chaincode.sh $CHAINCODE_NAME
}

function installChaincode() {
    
    # args: $CHAINCODE_NAME $CHANNEL_NAME <org name> <org id> <number of peer>
    $SCRIPTS_DIR/install-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1
    # $SCRIPTS_DIR/install-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 2

    $SCRIPTS_DIR/install-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "bus" 1 1
    # $SCRIPTS_DIR/install-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "bus" 1 2

}

function approveChaincode {
    # args: $CHAINCODE_NAME $CHANNEL_NAME <org name> <org id> <number of peer>
    $SCRIPTS_DIR/approve-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1
    # $SCRIPTS_DIR/approve-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 2

    $SCRIPTS_DIR/approve-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "bus" 1 1
    # $SCRIPTS_DIR/approve-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "bus" 1 2

    $SCRIPTS_DIR/commit-checkreadiness.sh $CHAINCODE_NAME $CHANNEL_NAME
}

function commitChaincode() {
    # args: $CHAINCODE_NAME $CHANNEL_NAME <number of org> <number of peer>
    $SCRIPTS_DIR/commit-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME 1 1
}

function invokeInitLedger() {
    # args: $CHAINCODE_NAME $CHANNEL_NAME <number of org> <number of peer>
    $SCRIPTS_DIR/init-ledger.sh $CHAINCODE_NAME $CHANNEL_NAME 1 1
}   

function invokeCreateCamp() {
    $SCRIPTS_DIR/create-camp.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" "bus" 1 1
} 

function invokeCollectData() {
    $SCRIPTS_DIR/create-data.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" "bus" 1 1
}

function runExternalService() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 0
}

function runVerifier1Service() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 1
}

function runVerifier2Service() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 2
}

function buildExternalService() {

    # FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache external.promark.com 2>&1

    # FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache verifier1.promark.com 2>&1
    
    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache peer0.bus0.promark.com 2>&1

    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache peer0.adv0.promark.com 2>&1
}

function invokeQueryById() {
    $SCRIPTS_DIR/query-ledger.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1
}

function deleteCampById() {
    $SCRIPTS_DIR/delete-camp-byId.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1 
}

function getAllCamp() {
    $SCRIPTS_DIR/getAllCamp.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1
}

function getAllCampaignData() {
    $SCRIPTS_DIR/getAllCampaignData.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1

}

function runLogService() {
    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} run --service-ports logs.promark.com /bin/sh 2>&1
}

function clearNetwork {
  docker rm -f $(docker ps -aq) || true
  docker rmi -f $(docker images -a -q) || true
  docker volume rm $(docker volume ls)  || true
}

MODE=$1

if [ $MODE = "restart" ]; then
    networkDown
    clear
    initialize
    networkUp
    createChannel
    joinChannel
elif [ $MODE = "clean" ]; then
    clearNetwork
elif [ $MODE = "init" ]; then
    clear 
    sleep 10
    # rm chaincode/main.tar.gz
    initialize
    sleep 10
    networkUp
elif [ $MODE = "clear" ]; then
    clear
elif [ $MODE = "up" ]; then
    networkUp
elif [ $MODE = "monitor" ]; then
    monitor
elif [ $MODE = "channel" ]; then

    SUB_MODE=$2

    if [ $SUB_MODE = "create" ]; then
        createChannel
    elif [ $SUB_MODE = "join" ]; then
        joinChannel
    elif [ $SUB_MODE = "all" ]; then
        createChannel
        sleep 10
        joinChannel
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "chaincode" ]; then
    # deployCC "campaign"

    SUB_MODE=$2

    if [ $SUB_MODE = "package" ]; then
        packageChaincode
    elif [ $SUB_MODE = "install" ]; then
        installChaincode
    elif [ $SUB_MODE = "approve" ]; then
        approveChaincode
    elif [ $SUB_MODE = "commit" ]; then
        commitChaincode
    elif [ $SUB_MODE = "all" ]; then
        packageChaincode
        sleep 10
        installChaincode
        sleep 15
        approveChaincode
        sleep 15
        commitChaincode
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "trans" ]; then
    SUB_MODE=$2
    
    if [ $SUB_MODE = "init" ]; then
        invokeInitLedger
    elif [ $SUB_MODE = "add" ]; then
        invokeCreateCamp
    elif [ $SUB_MODE = "query" ]; then
        invokeQueryById
    elif [ $SUB_MODE = "get" ]; then
        getAllCamp
    elif [ $SUB_MODE = "delete" ]; then
        deleteCampById
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "camp" ]; then
    SUB_MODE=$2
    
    if [ $SUB_MODE = "add" ]; then
        invokeCollectData
    elif [ $SUB_MODE = "get" ]; then
        getAllCampaignData
    fi
elif [ $MODE = "service" ]; then
    SUB_MODE=$2

    if [ $SUB_MODE = "run" ]; then
        runExternalService
    elif [ $SUB_MODE = "run1" ]; then
        runVerifier1Service
    elif [ $SUB_MODE = "run2" ]; then
        runVerifier2Service
    elif [ $SUB_MODE = "log" ]; then
        runLogService
    elif [ $SUB_MODE = "build" ]; then
        buildExternalService
    else 
        echo "Unsupported $MODE $SUB_MODE command."
    fi
else
    echo "Unsupported $MODE command."
fi