#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function joinChannel() {
    local channelName=$1
    IFS=',' read -r -a orgTypes <<< $2
    # local org_type=$2
    local orgNum=$3
    local peerNum=$4
    local orgName=""

    local maxPeerId=$(($peerNum - 1))
    local maxOrgId=$(($orgNum - 1))

    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrgId); do
            for peerId in $(seq 0 $maxPeerId); do
    # for org_id in $(seq 0 $maxOrgId); do
                orgName="peer${orgId}.${orgType}${orgId}"
    #     for peer_id in $(seq 0 $maxPeerId); do
                selectPeer $orgType $orgId $peerId

                infoln "Joining Channel $channelName from Org $orgName's peer$peerId"

                getBlockPath $channelName

                set -x
                peer channel join -b $blockPath
                res=$?
                { set +x; } 2>/dev/null
                #cat log.txt
            done
        done
    done
}

joinChannel $1 $2 $3 $4