#!/bin/bash

. $PWD/settings.sh

export CHANNEL_NAME="mychannel"
export LOG_LEVEL=INFO
export FABRIC_LOGGING_SPEC=INFO # log of the host machine
export CAMPAIGN_CHAINCODE_NAME="campaign"
export PROOF_CHAINCODE_NAME="proof"

function initialize() {
    # generate all organizations
    $BASE_SCRIPTS_DIR/gen-orgs.sh

    # generate genesis-block
    $BASE_SCRIPTS_DIR/gen-genesis-block.sh
}

function createChannel() {
    $BASE_SCRIPTS_DIR/gen-channel-tx.sh $CHANNEL_NAME
    $BASE_SCRIPTS_DIR/gen-channel.sh $CHANNEL_NAME "adv" 0
}

function joinChannel() {
     # args: $CHANNEL_NAME <org type> <number of org> <number of peer>
    $BASE_SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "adv,bus" $1 $2
}

function networkUp() {
    $BASE_SCRIPTS_DIR/start.sh $LOG_LEVEL
}

function networkDown() {
    # docker rm -f logspout

    $BASE_SCRIPTS_DIR/stop.sh $LOG_LEVEL
}

function clear() {
    $BASE_SCRIPTS_DIR/clear.sh
}

function monitor() {
    $BASE_SCRIPTS_DIR/monitor.sh
}

function packageChaincode() {
    $BASE_SCRIPTS_DIR/package-chaincode.sh $1 $2
}

function installChaincode() {
    # args: $CHANNEL_NAME $CHAINCODE_NAME  <org name> <org id> <number of peer>
    $BASE_SCRIPTS_DIR/install-chaincode.sh $CHANNEL_NAME $1 "adv,bus" $2 $3
    # $BASE_SCRIPTS_DIR/install-chaincode.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,bus" $1 $2
}

function approveChaincode {
    # args: $CHANNEL_NAME $CHAINCODE_NAME  <org name> <org id> <number of peer> <sequence>
    $BASE_SCRIPTS_DIR/approve-chaincode.sh $CHANNEL_NAME $1 "adv,bus" $2 $3 $4
    # $BASE_SCRIPTS_DIR/approve-chaincode.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,bus" $1 $2 $3
}

function commitChaincode() {
    # args: $CHANNEL_NAME $CHAINCODE_NAME  <number of org> <number of peer>
    $BASE_SCRIPTS_DIR/commit-chaincode.sh $CHANNEL_NAME $1 "adv,bus" $2 $3 $4
    # $BASE_SCRIPTS_DIR/commit-chaincode.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,bus" $1 $2 $3
}

function deployChaincode() {
    # args: $CHAINCODE_NAME $NO_ORG $NO_PEERS $SEQUENCE
    local CHAINCODE_NAME=$1
    local NO_ORG=$2
    local NO_PEERS=$3
    local SEQUENCE=$4

    sleep 1
    packageChaincode $CHAINCODE_NAME $SEQUENCE
    sleep 1
    installChaincode $CHAINCODE_NAME $NO_ORG $NO_PEERS $SEQUENCE
    sleep 1
    approveChaincode $CHAINCODE_NAME $NO_ORG $NO_PEERS $SEQUENCE
    sleep 1
    commitChaincode $CHAINCODE_NAME $NO_ORG $NO_PEERS $SEQUENCE
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

    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache 2>&1

    # FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache verifier1.promark.com 2>&1

    # FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache peer0.bus0.promark.com 2>&1

    # FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache peer0.adv0.promark.com 2>&1
}


function invokeCreateCamp() {
    $SCRIPTS_DIR/create-camp.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,bus" $1 $2
}

function getAllCamp() {
    $SCRIPTS_DIR/get-all-camp.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,bus" $1 $2
}

function deleteCampById() {
    $SCRIPTS_DIR/delete-camp-by-id.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME  "adv,bus" $1 $2
}

function getCampById() {
    $SCRIPTS_DIR/get-campaign-by-id.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,bus" $1 $2
}

function invokeGenerateCustomerProof() {
    $SCRIPTS_DIR/generate-customer-proof.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,bus" $1 $2
}

function invokeAddCustomerProof() {
    $SCRIPTS_DIR/add-proof.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,bus" $1 $2 $3 $4 $5
}

function invokeGetAllCustomerProofs() {
    $SCRIPTS_DIR/get-all-proofs.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,bus" $1 $2
}

function invokeDelProofById() {
    $SCRIPTS_DIR/del-proof-by-id.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,bus" $1 $2 $3
}

function invokeGetCustomerProofById() {
    $SCRIPTS_DIR/get-proof-by-id.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,bus" $1 $2 $3
}

function invokeVerifyProof() {
    $SCRIPTS_DIR/verify-proof.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,bus" $1 $2 $3 $4 $5
}

function getAllCampaignData() {
    $SCRIPTS_DIR/get-all-campaign-data.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" $1 $2

}

function invokeCollectData() {
    $SCRIPTS_DIR/create-data.sh $CHAINCODE_NAME $CHANNEL_NAME "adv,bus" $1 $2
}

# function invokeQueryById() {
#     $SCRIPTS_DIR/query-ledger.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" $1 $2
# }

function runLogService() {
    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} run --service-ports logs.promark.com /bin/sh 2>&1
}

function buildAllImages() {
    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build 2>&1
}

function clearNetwork {
    docker rm -f $(docker ps -aq) || true
    docker rmi -f $(docker images -a -q) || true
    docker volume rm $(docker volume ls)  || true
}

function addGoPath {
    export PATH=$PATH:./go/bin
}

MODE=$1
# addGoPath

if [ $MODE = "restart" ]; then
    NO_ORG=$2
    NO_PEERS=$3
    # SEQUENCE=$4

    networkDown
    clear
    initialize
    networkUp

    sleep 1
    createChannel
    sleep 1
    joinChannel $NO_ORG $NO_PEERS

    sleep 1
    deployChaincode $CAMPAIGN_CHAINCODE_NAME $NO_ORG $NO_PEERS 1
    sleep 10
    deployChaincode $PROOF_CHAINCODE_NAME $NO_ORG $NO_PEERS 1

elif [ $MODE = "build" ]; then
    SUB_MODE=$2

    if [ $SUB_MODE = "all" ]; then
        buildAllImages
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "clean" ]; then
    clearNetwork
elif [ $MODE = "path" ]; then
    addGoPath
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
elif [ $MODE = "down" ]; then
    networkDown
elif [ $MODE = "monitor" ]; then
    monitor
elif [ $MODE = "channel" ]; then

    SUB_MODE=$2
    NO_ORG=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "create" ]; then
        createChannel
    elif [ $SUB_MODE = "join" ]; then
        joinChannel $NO_ORG $NO_PEERS
    elif [ $SUB_MODE = "all" ]; then
        createChannel
        sleep 10
        joinChannel $NO_ORG $NO_PEERS
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "chaincode" ]; then
    SUB_MODE=$2
    CHAINCODE_NAME=$3

    if [ $SUB_MODE = "package" ]; then
        SEQUENCE=$4

        packageChaincode $CHAINCODE_NAME $SEQUENCE
    elif [ $SUB_MODE = "install" ]; then
        NO_ORG=$4
        NO_PEERS=$5
        SEQUENCE=$6

        installChaincode $CHAINCODE_NAME $NO_ORG $NO_PEERS
    elif [ $SUB_MODE = "approve" ]; then
        NO_ORG=$4
        NO_PEERS=$5
        SEQUENCE=$6

        approveChaincode $CHAINCODE_NAME $NO_ORG $NO_PEERS $SEQUENCE
    elif [ $SUB_MODE = "commit" ]; then
        NO_ORG=$4
        NO_PEERS=$5
        SEQUENCE=$6

        commitChaincode $CHAINCODE_NAME $NO_ORG $NO_PEERS $SEQUENCE
    elif [ $SUB_MODE = "deploy" ]; then
        NO_ORG=$4
        NO_PEERS=$5
        SEQUENCE=$6

        deployChaincode $CHAINCODE_NAME $NO_ORG $NO_PEERS $SEQUENCE
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "camp" ]; then
    SUB_MODE=$2
    NO_ORG=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "add" ]; then
        invokeCreateCamp $NO_ORG $NO_PEERS
    elif [ $SUB_MODE = "all" ]; then
        getAllCamp $NO_ORG $NO_PEERS
    elif [ $SUB_MODE = "del" ]; then
        deleteCampById $NO_ORG $NO_PEERS
    elif [ $SUB_MODE = "get" ]; then
        getCampById $NO_ORG $NO_PEERS
    # elif [ $SUB_MODE = "query" ]; then
    #     invokeQueryById $NO_ORG $NO_PEERS
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "proof" ]; then
    SUB_MODE=$2
    NO_ORG=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "gen" ]; then
        invokeGenerateCustomerProof $NO_ORG $NO_PEERS
    elif [ $SUB_MODE = "add" ]; then
        proofId=$5
        comm=$6
        rsStr=$7
        invokeAddCustomerProof $NO_ORG $NO_PEERS $proofId $comm $rsStr
    elif [ $SUB_MODE = "all" ]; then
        invokeGetAllCustomerProofs $NO_ORG $NO_PEERS
    elif [ $SUB_MODE = "del" ]; then
        proofId=$5
        invokeDelProofById $NO_ORG $NO_PEERS $proofId
    elif [ $SUB_MODE = "get" ]; then
        proofId=$5
        invokeGetCustomerProofById $NO_ORG $NO_PEERS $proofId
    elif [ $SUB_MODE = "verify" ]; then
        camId=$5
        proofId=$6
        invokeVerifyProof $NO_ORG $NO_PEERS $camId $proofId
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "data" ]; then
    SUB_MODE=$2
    NO_ORG=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "add" ]; then
        invokeCollectData $NO_ORG $NO_PEERS
    elif [ $SUB_MODE = "get" ]; then
        getAllCampaignData $NO_ORG $NO_PEERS
    fi
elif [ $MODE = "service" ]; then
    SUB_MODE=$2

    if [ $SUB_MODE = "ext" ]; then
        runExternalService
    elif [ $SUB_MODE = "ver1" ]; then
        runVerifier1Service
    elif [ $SUB_MODE = "ver2" ]; then
        runVerifier2Service
    elif [ $SUB_MODE = "log" ]; then
        runLogService
    elif [ $SUB_MODE = "build" ]; then
        buildExternalService
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
else
    errorln "Unsupported $MODE command."
fi