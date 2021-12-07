#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh


function createChannel() {
    local channelName=$1
    local orgType=$2
    local orgId=$3

    selectPeer $orgType $orgId 0

    println "Generating channel tx..."
    getChannelTxPath $channelName
    getBlockPath $channelName

    println "Creating channel..."
    set -x
    peer channel create -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME -c $channelName -f $channelTxPath --outputBlock $blockPath --tls --cafile $ORDERER_CA
    res=$?
    { set +x; } 2>/dev/null

	# cat log.txt
	# verifyResult $res "Channel creation failed"
}

createChannel $1 $2 $3