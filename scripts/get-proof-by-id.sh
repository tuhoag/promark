#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_GET_PROOF_BY_ID_FCN="GetProofById"

function getData() {
    local channelName=$1
    local chaincodeName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5
    local proofId=$6

    infoln "getData: $1 $2 $3 $4 $5 $6"

    fcnCall='{"function":"'${CC_GET_PROOF_BY_ID_FCN}'","Args":["'${proofId}'"]}'

    # infoln "Invoking Init Chaincode with $@\n"
    parsePeerConnectionParameters $orgTypes $orgNum $peerNum
    res=$?
    verifyResult $res "Invoke transaction failed on channel '$channelName' due to uneven number of peer and org parameters "

    sleep $DELAY
    infoln "Attempting to Query peer0.org${ORG}, Retry after $DELAY seconds."
    set -x
    peer chaincode query --channelID $channelName --name $chaincodeName -c $fcnCall >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    let rc=$res
    cat log.txt
}

getData $1 $2 $3 $4 $5 $6
