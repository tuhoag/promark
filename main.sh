#!/bin/bash

. $PWD/settings.sh

export CHANNEL_NAME="mychannel"
export LOG_LEVEL=INFO
export FABRIC_LOGGING_SPEC=DEBUG
export CHAINCODE_NAME="main"

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
    $SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "adv" 0 0
    $SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "bus" 0 0
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
    $SCRIPTS_DIR/install-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 0 0
    $SCRIPTS_DIR/install-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "bus" 0 0
}

function approveChaincode {

    $SCRIPTS_DIR/approve-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 0 0
    $SCRIPTS_DIR/approve-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME "bus" 0 0

    $SCRIPTS_DIR/commit-checkreadiness.sh $CHAINCODE_NAME $CHANNEL_NAME
}

function commitChaincode() {
    $SCRIPTS_DIR/commit-chaincode.sh $CHAINCODE_NAME $CHANNEL_NAME 1 1
}

function invokeInitLedger() {
    $SCRIPTS_DIR/init-ledger.sh $CHAINCODE_NAME $CHANNEL_NAME 1 1
}   

function invokeCreateCamp() {
    $SCRIPTS_DIR/create-camp.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" "bus" 1 1
} 

function createExternalService() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL
}

function buildExternalService() {

    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache external.promark.com 2>&1
}

MODE=$1

if [ $MODE = "restart" ]; then
    networkDown
    clear
    initialize
    networkUp
    createChannel
    joinChannel

elif [ $MODE = "init" ]; then
    initialize
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
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "trans" ]; then
    SUB_MODE=$2

    if [ $SUB_MODE = "init" ]; then
        invokeInitLedger
    elif [ $SUB_MODE = "add" ]; then
        invokeCreateCamp
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "service" ]; then
    SUB_MODE=$2

    if [ $SUB_MODE = "run" ]; then
        createExternalService
    elif [ $SUB_MODE = "build" ]; then
        buildExternalService
    else 
        echo "Unsupported $MODE $SUB_MODE command."
    fi
else
    echo "Unsupported $MODE command."
fi