#!/bin/bash

. $PWD/settings.sh
. $BASE_SCRIPTS_DIR/utils.sh

export CHANNEL_NAME="mychannel"
export LOG_LEVEL=INFO
export FABRIC_LOGGING_SPEC=INFO # log of the host machine
export CAMPAIGN_CHAINCODE_NAME="campaign"
export PROOF_CHAINCODE_NAME="proof"
export POC_CHAINCODE_NAME="poc"

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
    $BASE_SCRIPTS_DIR/join-channel.sh $CHANNEL_NAME "adv,pub" $1 $2
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

function clearDocker() {
    docker rm -f $(docker ps -aq) || true
    docker rmi -f $(docker images -a -q) || true
    docker volume rm $(docker volume ls)  || true
}

function monitor() {
    $BASE_SCRIPTS_DIR/monitor.sh
}

function packageChaincode() {
    $BASE_SCRIPTS_DIR/package-chaincode.sh $1 $2
}

function installChaincode() {
    # args: $CHANNEL_NAME $CHAINCODE_NAME  <org name> <org id> <number of peer>
    $BASE_SCRIPTS_DIR/install-chaincode.sh $CHANNEL_NAME $1 "adv,pub" $2 $3
    # $BASE_SCRIPTS_DIR/install-chaincode.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,pub" $1 $2
}

function approveChaincode {
    # approveChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE

    if [ $1 = "poc" ]; then
        # POLICY=--signature-policy="OR('adv0MSP.member','pub0MSP.member')"
        getSingleRandomOrgPolicy "adv,pub" $2 $3
        # echo "returned policy: $policy"
        policy=--signature-policy=${policy}
        # POLICY="--signature-policy=${getSingleRandomOrgPolicy "adv,pub" $2 $3}"
        # local policy=--signature-policy="OR('adv0MSP.member','pub0MSP.member')"
    # else
    #     POLICY=--signature-policy "OR(adv0MSP.member, pub0MSP.member)"
        # $BASE_SCRIPTS_DIR/approve-chaincode.sh $CHANNEL_NAME $1 "adv,pub" $2 $3 $4
    fi

    echo $policy
    $BASE_SCRIPTS_DIR/approve-chaincode.sh $CHANNEL_NAME $1 "adv,pub" $2 $3 $4 $policy
    # $BASE_SCRIPTS_DIR/approve-chaincode.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,pub" $1 $2 $3
}

function commitChaincode() {
    # args: $CHANNEL_NAME $CHAINCODE_NAME  <number of org> <number of peer>
    if [ $1 = "poc" ]; then
        # POLICY=--signature-policy="OR('adv0MSP.member','pub0MSP.member')"
        getSingleRandomOrgPolicy "adv,pub" $2 $3
        policy=--signature-policy=${policy}
        # local policy=--signature-policy="OR('adv0MSP.member','pub0MSP.member')"
    # else
    #     POLICY=--signature-policy "OR(adv0MSP.member, pub0MSP.member)"
        # $BASE_SCRIPTS_DIR/approve-chaincode.sh $CHANNEL_NAME $1 "adv,pub" $2 $3 $4
    fi

    # echo $POLICY

    $BASE_SCRIPTS_DIR/commit-chaincode.sh $CHANNEL_NAME $1 "adv,pub" $2 $3 $4 $policy
    # $BASE_SCRIPTS_DIR/commit-chaincode.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,pub" $1 $2 $3
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
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 1 $1 $2
}

function runVerifier2Service() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 2
}

function runClientService() {
    $SCRIPTS_DIR/external-service.sh $LOG_LEVEL 3 $1 $2
}

function buildExternalService() {
    local orgNum=$1
    local peerNum=$2

    local docker_compose_path="${DOCKER_COMPOSE_DIR_PATH}/docker-compose-${orgNum}-${peerNum}.yml"

    FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${docker_compose_path} build --no-cache 2>&1

    # FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache verifier1.promark.com 2>&1

    # FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache peer0.pub0.promark.com 2>&1

    # FABRIC_LOG=$LOG_LEVEL COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} build --no-cache peer0.adv0.promark.com 2>&1
}


function invokeCreateCamp() {
    local orgNum=$1
    local peerNum=$2
    local numVerifiers=1
    local deviceIdsStr=$3

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum campaign add $numVerifiers $deviceIdsStr
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/create-camp.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,pub" $1 $2
}

function getAllCamp() {
    local orgNum=$1
    local peerNum=$2

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum campaign all
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/get-all-camp.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,pub" $1 $2
}

function deleteCampById() {
    local orgNum=$1
    local peerNum=$2
    local camId=$3

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum campaign del $camId
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/delete-camp-by-id.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME  "adv,pub" $1 $2
}

function getCampById() {
    local orgNum=$1
    local peerNum=$2
    local camId=$3

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum campaign get $camId
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/get-campaign-by-id.sh $CHANNEL_NAME $CAMPAIGN_CHAINCODE_NAME "adv,pub" $1 $2
}

function invokeGenerateCustomerProof() {
    local orgNum=$1
    local peerNum=$2
    local camId=$3
    local userId=$4

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum proof gen $camId $userId
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/generate-customer-proof.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,pub" $1 $2
}

function invokeAddCustomerProof() {
    local orgNum=$1
    local peerNum=$2
    local comm=$3
    local rsStr=$4

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum proof add $comm $rsStr
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/add-proof.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,pub" $1 $2 $3 $4 $5
}

function invokeGetAllCustomerProofs() {
    local orgNum=$1
    local peerNum=$2

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum proof all
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/get-all-proofs.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,pub" $1 $2
}

function invokeDelProofById() {
    local orgNum=$1
    local peerNum=$2
    local proofId=$3

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum proof del $proofId
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/del-proof-by-id.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,pub" $1 $2 $3
}

function invokeGetCustomerProofById() {
    local orgNum=$1
    local peerNum=$2
    local proofId=$3

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum proof get $proofId
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/get-proof-by-id.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,pub" $1 $2 $3
}

function invokeVerifyProof() {
    local orgNum=$1
    local peerNum=$2
    local camId=$3
    local proofId=$4

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum proof verify $camId $proofId
    { set +x; } 2>/dev/null
    popd

    # $SCRIPTS_DIR/verify-proof.sh $CHANNEL_NAME $PROOF_CHAINCODE_NAME "adv,pub" $1 $2 $3 $4 $5
}

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

# function clearNetwork {
#     docker rm -f $(docker ps -aq) || true
#     docker rmi -f $(docker images -a -q) || true
#     docker volume rm $(docker volume ls)  || true
# }

function addGoPath {
    export PATH=$PATH:./go/bin
}

function evaluate {
    local orgNum=$1
    local peerNum=$2
    local benchmarkName=$3

    local networkConfigPath="${CALIPER_DIR_PATH}/config/networkConfig-${orgNum}-${peerNum}.yaml"
    local benchmarksPath="${CALIPER_DIR_PATH}/benchmarks/${benchmarkName}.yaml"


    pushd $CALIPER_DIR_PATH
    set -x
    npx caliper launch manager --caliper-workspace $CALIPER_DIR_PATH --caliper-networkconfig $networkConfigPath --caliper-benchconfig $benchmarksPath  --caliper-fabric-gateway-enabled --caliper-flow-only-test
    { set +x; } 2>/dev/null
    popd
}

function testPromark {
    local orgNum=$1
    local peerNum=$2
    # local numVerifiers=1

    pushd $CLIENT_DIR_PATH
    set -x
    node main.js $orgNum $peerNum "test"
    { set +x; } 2>/dev/null
    popd
}

MODE=$1
NO_ORGS=$2
NO_PEERS=$3
# addGoPath

if [ $MODE = "restart" ]; then
    SUB_MODE=$4

    if [ $SUB_MODE = "network" ]; then
        networkDown $NO_ORGS $NO_PEERS
        clear $NO_ORGS $NO_PEERS
        initialize $NO_ORGS $NO_PEERS
        networkUp $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "channel" ]; then
        networkDown $NO_ORGS $NO_PEERS
        clear $NO_ORGS $NO_PEERS
        initialize $NO_ORGS $NO_PEERS
        networkUp $NO_ORGS $NO_PEERS

        sleep 1
        createChannel $NO_ORGS $NO_PEERS
        sleep 2
        joinChannel $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "all" ]; then
        networkDown $NO_ORGS $NO_PEERS
        clear $NO_ORGS $NO_PEERS
        initialize $NO_ORGS $NO_PEERS
        networkUp $NO_ORGS $NO_PEERS

        sleep 1
        createChannel $NO_ORGS $NO_PEERS
        sleep 2
        joinChannel $NO_ORGS $NO_PEERS

        sleep 2
        deployChaincode $CAMPAIGN_CHAINCODE_NAME $NO_ORGS $NO_PEERS 1
        sleep 2
        deployChaincode $PROOF_CHAINCODE_NAME $NO_ORGS $NO_PEERS 1
        sleep 2
        deployChaincode $POC_CHAINCODE_NAME $NO_ORGS $NO_PEERS 1
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "build" ]; then
    SUB_MODE=$4

    if [ $SUB_MODE = "all" ]; then
        buildAllImages $NO_ORGS $NO_PEERS
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi

# elif [ $MODE = "clean" ]; then
#     clearNetwork
# elif [ $MODE = "path" ]; then
#     addGoPath
elif [ $MODE = "init" ]; then
    clear $NO_ORGS $NO_PEERS
    sleep 10
    # rm chaincode/main.tar.gz
    initialize $NO_ORGS $NO_PEERS
    sleep 10
    networkUp $NO_ORGS $NO_PEERS
elif [ $MODE = "clear" ]; then
    clearDocker
    # sleep 2
    clear $NO_ORGS $NO_PEERS
elif [ $MODE = "up" ]; then
    networkUp $NO_ORGS $NO_PEERS
elif [ $MODE = "down" ]; then
    networkDown $NO_ORGS $NO_PEERS
elif [ $MODE = "monitor" ]; then
    monitor
elif [ $MODE = "channel" ]; then
    SUB_MODE=$4

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
    SUB_MODE=$4
    CHAINCODE_NAME=$5

    if [ $SUB_MODE = "package" ]; then
        SEQUENCE=$6

        packageChaincode $CHAINCODE_NAME $SEQUENCE
    elif [ $SUB_MODE = "install" ]; then
        SEQUENCE=$6

        installChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "approve" ]; then
        SEQUENCE=$6

        approveChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    elif [ $SUB_MODE = "commit" ]; then
        SEQUENCE=$6

        commitChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    elif [ $SUB_MODE = "deploy" ]; then
        SEQUENCE=$6

        deployChaincode $CHAINCODE_NAME $NO_ORGS $NO_PEERS $SEQUENCE
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "campaign" ]; then
    SUB_MODE=$4

    if [ $SUB_MODE = "add" ]; then
        invokeCreateCamp $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "all" ]; then
        getAllCamp $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "del" ]; then
        camId=$5
        deleteCampById $NO_ORGS $NO_PEERS $camId
    elif [ $SUB_MODE = "get" ]; then
        camId=$5
        getCampById $NO_ORGS $NO_PEERS $camId
    # elif [ $SUB_MODE = "query" ]; then
    #     invokeQueryById $NO_ORGS $NO_PEERS
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "proof" ]; then
    SUB_MODE=$4

    if [ $SUB_MODE = "gen" ]; then
        camId=$5
        userId=$6
        invokeGenerateCustomerProof $NO_ORGS $NO_PEERS $camId $userId
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

elif [ $MODE = "service" ]; then
    SUB_MODE=$4

    if [ $SUB_MODE = "ext" ]; then
        runExternalService $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "ver1" ]; then
        runVerifier1Service $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "ver2" ]; then
        runVerifier2Service $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "log" ]; then
        runLogService $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "client" ]; then
        runClientService $NO_ORGS $NO_PEERS
    elif [ $SUB_MODE = "build" ]; then
        buildExternalService $NO_ORGS $NO_PEERS
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "eval" ]; then
    SUB_MODE=$4
    SUB_SUB_MODE=$5

    # evaluate $NO_ORGS $NO_PEERS "CreateCampaign"

    if [ $SUB_MODE = "campaign" ]; then
        if [ $SUB_SUB_MODE = "create" ]; then
            evaluate $NO_ORGS $NO_PEERS "CreateCampaign"
        else
            errorln "Unsupported $MODE $SUB_MODE $SUB_SUB_MODE command."
        fi
    elif [ $SUB_MODE = "proof" ]; then
        if [ $SUB_SUB_MODE = "gen" ]; then
            evaluate $NO_ORGS $NO_PEERS "GeneratePoC"
        elif [ $SUB_SUB_MODE = "add" ]; then
            evaluate $NO_ORGS $NO_PEERS "AddCampaignTokenTransaction"
        elif [ $SUB_SUB_MODE = "verifypoc" ]; then
            evaluate $NO_ORGS $NO_PEERS "VerifyPoC"
        elif [ $SUB_SUB_MODE = "verifytpoc" ]; then
            evaluate $NO_ORGS $NO_PEERS "VerifyTPoC"
        else
            errorln "Unsupported $MODE $SUB_MODE $SUB_SUB_MODE command."
        fi
    else
        errorln "Unsupported $MODE $SUB_MODE command."
    fi
elif [ $MODE = "test" ]; then
    testPromark $NO_ORGS $NO_PEERS

else
    errorln "Unsupported $MODE command."
fi