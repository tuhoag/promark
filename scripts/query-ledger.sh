#!/bin/bash

. $SCRIPTS_DIR/utils.sh

CC_QUERYBYID_FCN="QueryLedgerById"

function parsePeerConnectionParameters() {
    local orgNum=$1
    local peerNum=$2
    local orgType=$3
   
    PEER_CONN_PARMS=""
    PEERS=""
    local peerNames=""

    infoln "$orgNum ; $peerNum"

    local maxOrgId=$(($orgNum - 1))
    local maxPeerId=$(($peerNum - 1))

    for orgId in $(seq 0 $maxOrgId); do
         infoln $orgId
         for peerId in $(seq 0 $maxPeerId); do
             #for orgType in "adv" "bus"; do
                 local peerName="peer${peerId}.${orgType}${orgId}"
                 infoln "processed $peerName"
                 selectPeer $orgType $orgId $peerId
                 PEERS="$peerNames ${peerName}"
                 PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"
                ## Set path to TLS certificate
                TLSINFO=$(eval echo "--tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE")
                PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
             #done
         done
     done

    infoln "parsePeerConnectionParameters: $PEERS $PEER_CONN_PARMS"
}

function queryById() {
    local chaincodeName=$1
    local channelName=$2
    local orgType=$3
    local orgNum=$4
    local peerNum=$5

    infoln $orgType

    parsePeerConnectionParameters $orgNum $peerNum $orgType

    fcnCall='{"function":"'QueryLedgerById'","Args":["id1"]}'

    infoln "Invoke fcn call:${fcnCall} on peers: $peers"

    set -x
    peer chaincode query --channelID $channelName --name $chaincodeName -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --cafile $ORDERER_CA --tls $peerConnectionParams -c $fcnCall >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution on $peers failed "
    successln "Invoke transaction successful on $peers on channel '$channelName'"
}

queryById $1 $2 $3 $4 $5