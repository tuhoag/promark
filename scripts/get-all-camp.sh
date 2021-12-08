#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

CC_READ_ALL_FCN="GetAllAssets"

function getData () {
  local chaincodeName=$1
  local channelName=$2
  local orgTypes=$3
  local orgNum=$4
  local peerNum=$5

  infoln "getData: $1 $2 $3 $4 $5"

  # infoln "Invoking Init Chaincode with $@\n"
  parsePeerConnectionParameters $orgTypes $orgNum $peerNum
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

getData $1 $2 $3 $4 $5