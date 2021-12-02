#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function installChaincode() {
    local chaincode_name=$1
    local channel_name=$2
    local org_type=$3
    local orgNum=$4
    # local peer_id=$5
    local peerNum=$5
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincode_name}.tar.gz"

    local maxPeerId=$(($peerNum - 1))
    local maxOrgId=$(($orgNum - 1))

    for org_id in $(seq 0 $maxOrgId); do
        for peer_id in $(seq 0 $maxPeerId); do
            local peer_name="peer${peer_id}.${org_type}${org_id}"
            infoln "[${org_type}.${org_id}.${peer_id}] Install chaincode ${chaincode_name}"
            infoln "Installing chaincode ${chaincode_name} in channel ${channel_name} of ${peer_name}..."
            selectPeer $org_type $org_id $peer_id

            set -x
            # FABRIC_CFG_PATH="${PWD}/{}"
            peer lifecycle chaincode install $chaincode_package_path >&log.txt
            res=$?
            { set +x; } 2>/dev/null
            cat log.txt
            verifyResult $res "Chaincode installation on ${peer_name} has failed"
            successln "Chaincode is installed on ${peer_name}"
        done
    done
}

installChaincode $1 $2 $3 $4 $5
