#!/bin/bash

. $SCRIPTS_DIR/utils.sh

CC_INIT_FCN="InitLedger"

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

    infoln "parsePeerConnectionParameters1: $PEERS $PEER_CONN_PARMS"
}

function invokeInitCC() {
  echo "invokeInitCC: $1 $2 $3 $4"

  local chaincodeName=$1
  local channelName=$2
  local orgNum=$3
  local peerNum=$4

  # infoln "Invoking Init Chaincode with $@\n"
  parsePeerConnectionParameters $3 $4
  res=$?
  verifyResult $res "Invoke transaction failed on channel '$channelName' due to uneven number of peer and org parameters "

  # while 'peer chaincode' command can get the orderer endpoint from the
  # peer (if join was successful), let's supply it directly as we know
  # it using the "-o" option
  set -x
  fcn_call='{"function":"'${CC_INIT_FCN}'","Args":[]}'
  infoln "invoke fcn call:${fcn_call}"
  peer chaincode invoke -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS -c ${fcn_call} >&log.txt
  
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Invoke execution on $PEERS failed "
  successln "Invoke transaction successful on $PEERS on channel '$channelName'"
}

invokeInitCC $1 $2 $3 $4