#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_INIT_FCN="InitLedger"
CC_READ_ALL_FCN="GetAllAssets"

function invokeInitCC() {
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

function getData () {
  local chaincodeName=$1
  local channelName=$2
  local orgNum=$3
  local peerNum=$4

  # infoln "Invoking Init Chaincode with $@\n"
  parsePeerConnectionParameters $3 $4
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

invokeInitCC $1 $2 $3 $4
getData $1 $2 $3 $4