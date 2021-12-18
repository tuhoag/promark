#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_CREATE_FCN="CreateCampaign"

function createCamp() {
    local channelName=$1
    local chaincodeName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5

    infoln
    infoln "createCamp: channelName: ${channelName} - chaincodeName: ${chaincodeName}"

    #TODO: need to use the list of orgType
    parsePeerConnectionParameters $orgTypes $orgNum $peerNum
    res=$?
    verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

    set -x
    fcn_call0='{"function":"'${CC_CREATE_FCN}'","Args":["c002","campaign001","Adv0","Bus0","peer0.adv0.promark.com:5000;peer0.bus0.promark.com:5000"]}'

    # fcn_call1='{"function":"'${CC_CREATE_FCN}'","Args":["id11","campaign11","Adv1","Bus1","http://peer0.adv1.promark.com:8510","http://peer0.bus1.promark.com:9010"]}'

    # fcn_call2='{"function":"'${CC_CREATE_FCN}'","Args":["id12","campaign12","Adv2","Bus2","http://peer0.adv2.promark.com:8520","http://peer0.bus2.promark.com:9020"]}'

    # fcn_call3='{"function":"'${CC_CREATE_FCN}'","Args":["id13","campaign13","Adv3","Bus3","http://peer0.adv3.promark.com:8530","http://peer0.bus3.promark.com:9030"]}'

    infoln "invoke fcn call:${fcn_call0}"
    # peer lifecycle chaincode commit -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME  --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --tls $peerConnectionParams --version $sequence --sequence $sequence

    peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $peerConnectionParams -c ${fcn_call0} >&log.txt

    # peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call1} >&log.txt

    # peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call2} >&log.txt

    # peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call3} >&log.txt

    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution failed"
    successln "Invoke transaction successful on channel '$channelName'"

}

createCamp $1 $2 $3 $4 $5
