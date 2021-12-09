#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function installChaincode() {
    local channelName=$1
    local chaincodeName=$2

    IFS=',' read -r -a orgTypes <<< $3
    local orgNum=$4
    local peerNum=$5
    local chaincodePackagePath="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"

    local maxPeerId=$(($peerNum - 1))
    local maxOrgId=$(($orgNum - 1))

    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrgId); do
            for peerId in $(seq 0 $maxPeerId); do
                local peerName="peer${peerId}.${orgType}${orgId}"
                # infoln "[${orgType}.${orgId}.${peerId}] Install chaincode ${chaincodeName}"
                infoln "Installing chaincode ${chaincodeName} in channel ${channelName} of ${peerName}..."
                selectPeer $orgType $orgId $peerId

                set -x
                # FABRIC_CFG_PATH="${PWD}/{}"
                peer lifecycle chaincode install $chaincodePackagePath >&log.txt
                res=$?
                { set +x; } 2>/dev/null
                cat log.txt
                verifyResult $res "Chaincode installation on ${peerName} has failed"
                successln "Chaincode is installed on ${peerName}"
            done
        done
    done
}

installChaincode $1 $2 $3 $4 $5