#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function joinChannel() {
    local channel_name=$1
    local org_type=$2
    local org_id=$3
    local peerNum=$4

    local org_name="${org_type}${org_id}"
    
    local maxPeerId=$(($peerNum - 1))

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
}

joinChannel $1 $2 $3 $4