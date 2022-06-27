#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function approveForMyOrg() {
    local channelName=$1
    local chaincodeName=$2
    IFS=',' read -r -a orgTypes <<< $3
    local orgNum=$4
    local peerNum=$5
    local sequence=$6
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"
    # local policy=--signature-policy="OR('adv0MSP.member','pub0MSP.member')"
    local policy=$7


    echo "policy: $policy"
    local maxPeerId=$(($peerNum - 1))
    local maxOrgId=$(($orgNum - 1))

    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrgId); do
            # for peerId in $(seq 0 $maxPeerId); do
                local peerId=0
                local peerName="peer${peerId}.${orgType}${orgId}"

                infoln "Approving chaincode ${chaincodeName} in channel ${channelName} of ${peerName}..."

                selectPeer $orgType $orgId 0

                packageName="${chaincodeName}_${sequence}"
                getPackageId $packageName

                # packageId=$(getPackageId $packageName)
                infoln "My package id: $packageId"

                set -x
                peer lifecycle chaincode approveformyorg -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --version $sequence --package-id $packageId --sequence $sequence $policy >&log.txt
                res=$?
                { set +x; } 2>/dev/null

                cat log.txt

                verifyResult $res "Chaincode definition approved on ${peerName} on channel '$channelName' failed"
                successln "Chaincode definition approved on ${peerName} on channel '$channelName'"
            # done
        done
    done
}

approveForMyOrg $1 $2 $3 $4 $5 $6 $7
