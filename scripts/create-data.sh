#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_CREATE_FCN="AddCollectedData"
CC_READ_ALL_FCN="GetAllCollectedData"

function parsePeerConnectionParameters() {
    local orgTypes=$1
    local orgNum=$2
    local peerNum=$3


    PEER_CONN_PARMS=""
    PEERS=""
    local peerNames=""

    infoln "$orgNum ; $peerNum"

    local maxOrgId=$(($orgNum - 1))
    local maxPeerId=$(($peerNum - 1))

    for orgId in $(seq 0 $maxOrgId); do
        infoln $orgId
        for peerId in $(seq 0 $maxPeerId); do
            for orgType in "adv" "pub"; do
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


function addData() {
    local chaincodeName=$1
    local channelName=$2
    local orgNum=$5
    local peerNum=$6

    infoln "addData: $1 $2 $3 $4 $5 $6"

    #TODO: need to use the list of orgType
    parsePeerConnectionParameters $orgNum $peerNum
    res=$?
    verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

    set -x
    fcn_call='{"function":"'${CC_CREATE_FCN}'","Args":["id6","user6","2","kPEjPY0+xY9am+XnKy5Lj9QJVWZ8GTynTX4LrT+eliU=","RL0Peky2vLZ1khRy9zQo4Bw26tzd28cGruHBacC23Qs=;g6Xt3smh3u3BrdAYYkPrfz8ddpIkX2GMgChSEfHyEQw=","http://peer0.pub0.promark.com:9000;http://peer0.adv0.promark.com:8500"]}'

    infoln "invoke fcn call:${fcn_call}"
    peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call} >&log.txt

    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution failed"
    successln "Invoke transaction successful on channel '$channelName'"

}

addData $1 $2 $3 $4 $5 $6
getData $1 $2

