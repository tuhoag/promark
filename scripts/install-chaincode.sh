#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function installChaincode() {
    local chaincode_name=$1
    local channel_name=$2
    local org_type=$3
    local org_id=$4
    local peer_id=$5
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincode_name}.tar.gz"
    local peer_name="peer${peer_id}.${org_type}${org_id}"

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
}

installChaincode $1 $2 $3 $4 $5
