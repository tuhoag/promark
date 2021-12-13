#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh


function createChannel() {
    local orgNum=$1
    local peerNum=$2
    local channelName=$3
    local orgType=$4
    local orgId=$5

    selectPeer $orgType $orgId 0

    println "Generating channel tx..."
    getChannelTxPath $orgNum $peerNum $channelName
    getBlockPath $orgNum $peerNum $channelName

    println "Creating channel..."
    set -x
    peer channel create -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME -c $channelName -f $channelTxPath --outputBlock $blockPath --tls --cafile $ORDERER_CA
    res=$?
    { set +x; } 2>/dev/null

	# cat log.txt
	# verifyResult $res "Channel creation failed"
}

createChannel $1 $2 $3 $4 $5