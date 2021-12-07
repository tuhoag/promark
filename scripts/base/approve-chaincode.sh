#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function approveForMyOrg() {
    local chaincodeName=$1
    local channelName=$2
    IFS=',' read -r -a orgTypes <<< $3
    local orgNum=$4
    local peerNum=$5
    local sequence=$6
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"

    local maxPeerId=$(($peerNum - 1))
    local maxOrgId=$(($orgNum - 1))

    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrgId); do
            for peerId in $(seq 0 $maxPeerId); do
                local peerName="peer${peerId}.${orgType}${orgId}"

                infoln "Approving chaincode ${chaincodeName} in channel ${channelName} of ${peerName}..."

                selectPeer $orgType $orgId $peerId

                local packageId=$(peer lifecycle chaincode queryinstalled)
                packageId=${packageId%,*}
                packageId=${packageId#*:}
                packageId=${packageId##* }
                infoln "My package id: $packageId"

                set -x
                peer lifecycle chaincode approveformyorg -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --version $sequence --package-id $packageId --sequence $sequence >&log.txt
                res=$?
                { set +x; } 2>/dev/null

                cat log.txt

                verifyResult $res "Chaincode definition approved on ${peerName} on channel '$channelName' failed"
                successln "Chaincode definition approved on ${peerName} on channel '$channelName'"
            done
        done
    done
}

approveForMyOrg $1 $2 $3 $4 $5 $6
