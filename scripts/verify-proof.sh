#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_VERIFY_CAMPAIGN_PROOF_FCN="VerifyCampaignProof"


function generateProofCustomerCampaign() {
    local channelName=$1
    local chaincodeName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5

    parsePeerConnectionParameters $orgTypes $orgNum $peerNum
    res=$?
    verifyResult $res "Invoke transaction failed on channel '$channelName' due to uneven number of peer and org parameters "

    set -x
    fcn_call0='{"function":"'${CC_VERIFY_CAMPAIGN_PROOF_FCN}'","Args":["c:001","p:001"]}'
    { set +x; } 2>/dev/null

    set -x
    peer chaincode query --channelID $channelName --name $chaincodeName  -c ${fcn_call0} >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    let rc=$res
    cat log.txt
}

generateProofCustomerCampaign $1 $2 $3 $4 $5