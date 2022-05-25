#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_COLLECT_PROOF_FCN="AddCustomerProofCampaign"


function addCustomerProof() {
    local channelName=$1
    local chaincodeName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5
    local proofId=$6
    local comm=$7
    local rsStr=$8

    infoln "createCamp: $1 $2 $3 $4 $5 $6 $7 $8"

    #TODO: need to use the list of orgType
    parsePeerConnectionParameters $orgTypes $orgNum $peerNum
    res=$?
    verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

    set -x
    fcn_call0='{"function":"'${CC_COLLECT_PROOF_FCN}'","Args":["'${proofId}'","'${comm}'","'${rsStr}'"]}'
    # fcn_call0='{"function":"'${CC_CREATE_FCN}'","Args":["c:001","campaign100","Adv0","Pub0","http://peer0.adv0.promark.com:8500;http://peer0.pub0.promark.com:9000"]}'
    # fcn_call0='{"function":"CollectCustomerProofCampaign","Args":["p001","p001","p001"]}'

    infoln "invoke fcn call:${fcn_call0}"

    peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $peerConnectionParams -c ${fcn_call0} >&log.txt


    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution failed"
    successln "Invoke transaction successful on channel '$channelName'"

}

addCustomerProof $1 $2 $3 $4 $5 $6 $7 $8