#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function approveForMyOrg() {
    local chaincodeName=$1
    local channelName=$2
    local orgType=$3
    local orgId=$4
    local peerId=$5
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"
    local peer_name="peer${peerId}.${orgType}${orgId}"

    infoln "Approving chaincode ${chaincodeName} in channel ${channelName} of ${peer_name}..."

    selectPeer $orgType $orgId $peerId

    local packageId=$(peer lifecycle chaincode queryinstalled)
    packageId=${packageId%,*}
    packageId=${packageId#*:}
    packageId=${packageId##* }
    infoln "My package id: $packageId"

    set -x
    peer lifecycle chaincode approveformyorg -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --version 1.0 --package-id $packageId --sequence 1 >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    cat log.txt

    verifyResult $res "Chaincode definition approved on ${peer_name} on channel '$channelName' failed"
    successln "Chaincode definition approved on ${peer_name} on channel '$channelName'"
}

approveForMyOrg $1 $2 $3 $4 $5
