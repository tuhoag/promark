#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_QUERY_CUSTOMER_CAMPAIGN_PROOF_FCN="GetProofCustomerCampaign"


function queryProofCustomerCampaign() {
    local chaincodeName=$1
    local channelName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5

    parsePeerConnectionParameters $orgTypes $orgNum $peerNum
    res=$?
    verifyResult $res "Invoke transaction failed on channel '$channelName' due to uneven number of peer and org parameters "

    set -x
    fcn_call0='{"function":"'${CC_QUERY_CUSTOMER_CAMPAIGN_PROOF_FCN}'","Args":["c:001","u:001"]}'
    { set +x; } 2>/dev/null

    set -x
    peer chaincode query --channelID $channelName --name $chaincodeName  -c ${fcn_call0} >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    let rc=$res
    cat log.txt
}

queryProofCustomerCampaign $1 $2 $3 $4 $5