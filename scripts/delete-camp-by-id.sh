#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_DELETE_BY_ID_FCN="DeleteCampaignByID"

function deleteById() {
    local channelName=$1
    local chaincodeName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5

    infoln $orgType

    parsePeerConnectionParameters $orgTypes $orgNum $peerNum

    fcnCall='{"function":"'${CC_DELETE_BY_ID_FCN}'","Args":["c:001"]}'

    infoln "Invoke fcn call:${fcnCall} on peers: $peers"

    set -x
    peer chaincode invoke --channelID $channelName --name $chaincodeName -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --cafile $ORDERER_CA --tls $peerConnectionParams -c $fcnCall >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution on $peers failed "
    successln "Invoke transaction successful on $peers on channel '$channelName'"
}

deleteById $1 $2 $3 $4 $5