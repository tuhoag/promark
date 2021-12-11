#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_VERIFY_CAMPAIGN_PROOF_FCN="VerifyCampaignProof"


function verifyProof() {
    local channelName=$1
    local chaincodeName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5
    local camId=$6
    local proofId=$7

    parsePeerConnectionParameters $orgTypes $orgNum $peerNum
    res=$?
    verifyResult $res "Invoke transaction failed on channel '$channelName' due to uneven number of peer and org parameters "

    set -x
    fcn_call0='{"function":"'${CC_VERIFY_CAMPAIGN_PROOF_FCN}'","Args":["'${camId}'","'${proofId}'"]}'
    { set +x; } 2>/dev/null

    set -x
    peer chaincode query --channelID $channelName --name $chaincodeName  -c ${fcn_call0} >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    let rc=$res
    cat log.txt
}

verifyProof $1 $2 $3 $4 $5 $6 $7