#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function createChannelTx() {
#   if [ ! -d "channels" ]; then
#     mkdir channels
#   fi
	local orgNum=$1
    local peerNum=$2
    local channelName=$3
	getChannelTxPath $orgNum $peerNum $channelName
	# channel_tx_path=$?
	local newConfigPath="${CONFIG_PATH}/configtx.yaml"
	infoln $newConfigPath
	FABRIC_CFG_PATH=$newConfigPath
	set -x
	configtxgen -profile "${orgNum}Channel" -outputCreateChannelTx $channelTxPath -channelID $channelName -configPath $CONFIG_PATH
	res=$?
	{ set +x; } 2>/dev/null

    verifyResult $res "Failed to generate channel configuration transaction..."
}



# channelName=$1

createChannelTx $1 $2 $3