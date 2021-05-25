#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function joinChannel() {
    local channel_name=$1
    local org_type=$2
    local org_id=$3
    local peer_id=$4

    local org_name="${org_type}${org_id}"
    selectPeer $org_type $org_id $peer_id

    infoln "Joining Channel $channel_name from Org $org_name's peer$peer_id"


    getBlockPath $channel_name

    set -x
    peer channel join -b $block_path
    res=$?
    { set +x; } 2>/dev/null
    #cat log.txt

}

joinChannel $1 $2 $3 $4