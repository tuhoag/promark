#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function joinChannel() {
    local channel_name=$1
    local org_type=$2
    local orgNum=$3
    local peerNum=$4
    local org_name=""


    local maxPeerId=$(($peerNum - 1))
    local maxOrgId=$(($orgNum - 1))

    for org_id in $(seq 0 $maxOrgId); do
        org_name="${org_type}${org_id}"
        for peer_id in $(seq 0 $maxPeerId); do
            selectPeer $org_type $org_id $peer_id

            infoln "Joining Channel $channel_name from Org $org_name's peer$peer_id"

            getBlockPath $channel_name

            set -x
            peer channel join -b $block_path
            res=$?
            { set +x; } 2>/dev/null
            #cat log.txt
        done
    done
}

joinChannel $1 $2 $3 $4