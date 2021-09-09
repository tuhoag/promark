#!/bin/bash

. $SCRIPTS_DIR/utils.sh

CC_CREATE_FCN="CreateCampaign"
CC_READ_ALL_FCN="GetAllAssets"

function parsePeerConnectionParameters() {
    local orgNum=$1
    local peerNum=$2
   
    PEER_CONN_PARMS=""
    PEERS=""
    local peerNames=""

    infoln "$orgNum ; $peerNum"

    local maxOrgId=$(($orgNum - 1))
    local maxPeerId=$(($peerNum - 1))

    for orgId in $(seq 0 $maxOrgId); do
         infoln $orgId
         for peerId in $(seq 0 $maxPeerId); do
             for orgType in "adv" "bus"; do
                 local peerName="peer${peerId}.${orgType}${orgId}"
                 infoln "processed $peerName"
                 selectPeer $orgType $orgId $peerId
                 PEERS="$peerNames ${peerName}"
                 PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"
                ## Set path to TLS certificate
                TLSINFO=$(eval echo "--tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE")
                PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
             done
         done
     done

    infoln "parsePeerConnectionParameters: $PEERS $PEER_CONN_PARMS"
}

function getData () {
  local chaincodeName=$1
  local channelName=$2

  # infoln "Invoking Init Chaincode with $@\n"
  parsePeerConnectionParameters 1 1
  res=$?
  verifyResult $res "Invoke transaction failed on channel '$channelName' due to uneven number of peer and org parameters "

  sleep $DELAY
  infoln "Attempting to Query peer0.org${ORG}, Retry after $DELAY seconds."
  set -x
  peer chaincode query --channelID $channelName --name $chaincodeName  -c '{"Args":["'${CC_READ_ALL_FCN}'"]}' >&log.txt
  res=$?
  { set +x; } 2>/dev/null

  let rc=$res
  cat log.txt
}


function createCamp() {
    local chaincodeName=$1
    local channelName=$2
    local orgNum=$5
    local peerNum=$6

    echo "createCamp: $1 $2 $3 $4 $5 $6"

    #TODO: need to use the list of orgType 
    parsePeerConnectionParameters $orgNum $peerNum
    res=$?
    verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

    set -x
    fcn_call0='{"function":"'${CC_CREATE_FCN}'","Args":["id21","campaign21","Adv0","Bus0","http://peer0.adv0.promark.com:8500","http://peer0.bus0.promark.com:9000"]}'

    # fcn_call1='{"function":"'${CC_CREATE_FCN}'","Args":["id11","campaign11","Adv1","Bus1","http://peer0.adv1.promark.com:8510","http://peer0.bus1.promark.com:9010"]}'

    # fcn_call2='{"function":"'${CC_CREATE_FCN}'","Args":["id12","campaign12","Adv2","Bus2","http://peer0.adv2.promark.com:8520","http://peer0.bus2.promark.com:9020"]}'

    # fcn_call3='{"function":"'${CC_CREATE_FCN}'","Args":["id13","campaign13","Adv3","Bus3","http://peer0.adv3.promark.com:8530","http://peer0.bus3.promark.com:9030"]}'

    infoln "invoke fcn call:${fcn_call}"
    peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call0} >&log.txt

    # peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call1} >&log.txt
    
    # peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call2} >&log.txt

    # peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call3} >&log.txt

    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution failed"
    successln "Invoke transaction successful on channel '$channelName'"

}

createCamp $1 $2 $3 $4 $5 $6
getData $1 $2

