#!/bin/bash

. $PWD/settings.sh
. $BASE_SCRIPTS_DIR/utils.sh

export CHANNEL_NAME="mychannel"
export LOG_LEVEL=INFO
export FABRIC_LOGGING_SPEC=INFO # log of the host machine
export CAMPAIGN_CHAINCODE_NAME="campaign"
export PROOF_CHAINCODE_NAME="proof"

function initialize() {
    # generate all organizations
    $BASE_SCRIPTS_DIR/gen-orgs.sh $1 $2

    # generate genesis-block
    $BASE_SCRIPTS_DIR/gen-genesis-block.sh $1 $2 $CHANNEL_NAME
}

function createChannel() {
    $BASE_SCRIPTS_DIR/gen-channel-tx.sh $1 $2 $CHANNEL_NAME
    $BASE_SCRIPTS_DIR/gen-channel.sh $1 $2 $CHANNEL_NAME "adv" 0
}

function joinChannel() {
     # args: $CHANNEL_NAME <org type> <number of org> <number of peer>
    $BASE_SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "adv,bus" $1 $2
}

function networkUp() {
    $BASE_SCRIPTS_DIR/start.sh $1 $2 $LOG_LEVEL
}

function networkDown() {
    # docker rm -f logspout

    $BASE_SCRIPTS_DIR/stop.sh $1 $2 $LOG_LEVEL
}

function clear() {
    $BASE_SCRIPTS_DIR/clear.sh $1 $2 $LOG_LEVEL
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
    # args: $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    local CHAINCODE_NAME=$1
    local NO_ORGS=$2
    local NO_PEERS=$3
    local SEQUENCE=$4

    sleep 1
    packageChaincode $CHAINCODE_NAME $SEQUENCE
    sleep 1
    installChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    sleep 1
    approveChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    sleep 1
    commitChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
}

function runExternalService() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 0 $1 $2
}

function runVerifier1Service() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 1
}

function runVerifier2Service() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 2
}

function buildExternalService() {
    local orgNum=$1
    local peerNum=$2

    local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"

    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} build --no-cache 2>&1

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
    local orgNum=$1
    local peerNum=$2

    local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"

    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} run --service-ports logs.promark.com /bin/sh 2>&1
}

function buildAllImages() {
    local orgNum=$1
    local peerNum=$2

    local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"

    infoln $docker_compose_path

    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} build 2>&1
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
    NO_ORGS=$2
    NO_PEERS=$3
    # SEQUENCE=$4

    networkDown $NO_ORGS $NO_PEERS
    clear $NO_ORGS $NO_PEERS
    initialize $NO_ORGS $NO_PEERS
    networkUp $NO_ORGS $NO_PEERS

    sleep 1
    createChannel $NO_ORGS $NO_PEERS
    sleep 2
    joinChannel $NO_ORGS $NO_PEERS

    sleep 1
    deployChaincode $CAMPAIGN_CHAINCODE_NAME $NO_ORGS $NO_PEERS 1
    sleep 15
    deployChaincode $PROOF_CHAINCODE_NAME $NO_ORGS $NO_PEERS 1

elif [ $MODE = "build" ]; then
    SUB_MODE=$2
    NO_ORGS=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "all" ]; then
        buildAllImages $NO_ORGS $NO_PEERS
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "clean" ]; then
    clearNetwork
elif [ $MODE = "path" ]; then
    addGoPath
elif [ $MODE = "init" ]; then
    NO_ORGS=$2
    NO_PEERS=$3
    clear $NO_ORGS $NO_PEERS
    sleep 10
    # rm chaincode/main.tar.gz
    initialize $NO_ORGS $NO_PEERS
    sleep 10
    networkUp $NO_ORGS $NO_PEERS
elif [ $MODE = "clear" ]; then
    NO_ORGS=$2
    NO_PEERS=$3
    clear $NO_ORGS $NO_PEERS
elif [ $MODE = "up" ]; then
    NO_ORGS=$2
    NO_PEERS=$3
    networkUp $NO_ORGS $NO_PEERS
elif [ $MODE = "down" ]; then
    NO_ORGS=$2
    NO_PEERS=$3
    networkDown $NO_ORGS $NO_PEERS
elif [ $MODE = "monitor" ]; then
    monitor
elif [ $MODE = "channel" ]; then

    SUB_MODE=$2
    NO_ORGS=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "create" ]; then
        createChannel $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "join" ]; then
        joinChannel $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "all" ]; then
        createChannel $NO_ORGS $NO_PEERS
        sleep 10
        joinChannel $NO_ORGS $NO_PEERS
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
        NO_ORGS=$4
        NO_PEERS=$5
        SEQUENCE=$6

        installChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "approve" ]; then
        NO_ORGS=$4
        NO_PEERS=$5
        SEQUENCE=$6

        approveChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    elif [ $SUB_MODE = "commit" ]; then
        NO_ORGS=$4
        NO_PEERS=$5
        SEQUENCE=$6

        commitChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    elif [ $SUB_MODE = "deploy" ]; then
        NO_ORGS=$4
        NO_PEERS=$5
        SEQUENCE=$6

        deployChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "camp" ]; then
    SUB_MODE=$2
    NO_ORGS=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "add" ]; then
        invokeCreateCamp $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "all" ]; then
        getAllCamp $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "del" ]; then
        deleteCampById $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "get" ]; then
        getCampById $NO_ORGS $NO_PEERS
    # elif [ $SUB_MODE = "query" ]; then
    #     invokeQueryById $NO_ORGS $NO_PEERS
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "proof" ]; then
    SUB_MODE=$2
    NO_ORGS=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "gen" ]; then
        invokeGenerateCustomerProof $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "add" ]; then
        proofId=$5
        comm=$6
        rsStr=$7
        invokeAddCustomerProof $NO_ORGS $NO_PEERS $proofId $comm $rsStr
    elif [ $SUB_MODE = "all" ]; then
        invokeGetAllCustomerProofs $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "del" ]; then
        proofId=$5
        invokeDelProofById $NO_ORGS $NO_PEERS $proofId
    elif [ $SUB_MODE = "get" ]; then
        proofId=$5
        invokeGetCustomerProofById $NO_ORGS $NO_PEERS $proofId
    elif [ $SUB_MODE = "verify" ]; then
        camId=$5
        proofId=$6
        invokeVerifyProof $NO_ORGS $NO_PEERS $camId $proofId
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "data" ]; then
    SUB_MODE=$2
    NO_ORGS=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "add" ]; then
        invokeCollectData $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "get" ]; then
        getAllCampaignData $NO_ORGS $NO_PEERS
    fi
elif [ $MODE = "service" ]; then
    SUB_MODE=$2
    NO_ORGS=$3
    NO_PEERS=$4

    if [ $SUB_MODE = "ext" ]; then
        runExternalService $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "ver1" ]; then
        runVerifier1Service
    elif [ $SUB_MODE = "ver2" ]; then
        runVerifier2Service
    elif [ $SUB_MODE = "log" ]; then
        runLogService
    elif [ $SUB_MODE = "build" ]; then
        buildExternalService $NO_ORGS $NO_PEERS
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
else
    errorln "Unsupported $MODE command."
fi