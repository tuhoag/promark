#!/bin/bash

C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

# println echos string
function println() {
    echo -e "$1"
}

# errorln echos i red color
function errorln() {
    println "${C_RED}${1}${C_RESET}"
}

# successln echos in green color
function successln() {
    println "${C_GREEN}${1}${C_RESET}"
}

# infoln echos in blue color
function infoln() {
    println "${C_BLUE}${1}${C_RESET}"
}

# warnln echos in yellow color
function warnln() {
    println "${C_YELLOW}${1}${C_RESET}"
}

# fatalln echos in red color and exits with fail status
function fatalln() {
    errorln "$1"
    exit 1
}

function verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}

function selectPeer() {
    local org_type=$1
    local org_id=$2
    local peer_id=$3

    # calculate port
    if [ $org_type = "adv" ]; then
        local base_port=$ADV_BASE_PORT
    elif [ $org_type = "bus" ]; then
        local base_port=$BUS_BASE_PORT
    else
        errorln "Org type $org_type is unsupported."
    fi

    local port=$(($base_port + $org_id * 10 + $peer_id))
    local org_name="$org_type$org_id"
    local org_domain=$org_name.$PROJECT_NAME.com
    local peer_domain=peer$peer_id.$org_domain

    # infoln "Selecting organization $org_name's peer$peer_id with port $port"
    infoln "selectPeer peer${peer_id}.${org_type}${org_id}:${port}"

    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="${org_name}MSP"
    export CORE_PEER_ADDRESS=0.0.0.0:${port}
    export PEER_ORG_CA=${ORGANIZATION_OUTPUTS}/peerOrganizations/$org_domain/peers/$peer_domain/tls/ca.crt
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER_ORG_CA
    export CORE_PEER_MSPCONFIGPATH=${ORGANIZATION_OUTPUTS}/peerOrganizations/$org_domain/users/Admin@$org_domain/msp
}

function getChannelTxPath() {
    local orgNum=$1
    local peerNum=$2
    local channelName=$3
    channelTxPath=$CHANNEL_PATH/${channelName}-${orgNum}-${peerNum}.tx
    # return $channel_tx_path
}

function getBlockPath() {
    local orgNum=$1
    local peerNum=$2
    local channelName=$3
    blockPath="${CHANNEL_PATH}/${channelName}-genesis-${orgNum}-${peerNum}.block"
    # return $blockPath
}

# function parsePeerConnectionParameters() {
#     local orgNum=$1
#     local peerNum=$2

#     PEER_CONN_PARMS=""
#     PEERS=""
#     local peerNames=""

#     infoln "$orgNum ; $peerNum"

#     local maxOrgId=$(($orgNum - 1))
#     local maxPeerId=$(($peerNum - 1))

#     for orgId in $(seq 0 $maxOrgId); do
#          infoln $orgId
#          for peerId in $(seq 0 $maxPeerId); do
#              for orgType in "adv" "bus"; do
#                  local peerName="peer${peerId}.${orgType}${orgId}"
#                  infoln "processed $peerName"
#                  selectPeer $orgType $orgId $peerId
#                  PEERS="$peerNames ${peerName}"
#                  PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"
#                 ## Set path to TLS certificate
#                 TLSINFO=$(eval echo "--tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE")
#                 PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
#              done
#          done
#      done

#     infoln "parsePeerConnectionParameters1: $PEERS $PEER_CONN_PARMS"
# }

function parsePeerConnectionParameters() {
    IFS=',' read -r -a orgTypes <<< $1
    local maxOrdId=$(($2 - 1))
    local maxPeerId=$(($3 - 1))

    # echo $orgTypes
    peerConnectionParams=""
    peers=""
    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrdId); do
            for peerId in $(seq 0 $maxOrdId); do
                selectPeer $orgType $orgId $peerId

                peers="$peers $CORE_PEER_ADDRESS"
                peerConnectionParams="$peerConnectionParams --peerAddresses $CORE_PEER_ADDRESS"
                peerConnectionParams="$peerConnectionParams --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE"
            done
        done
    done

    infoln "parsePeerConnectionParameters: $PEERS $PEER_CONN_PARMS"
}

export -f errorln
export -f successln
export -f infoln
export -f warnln
