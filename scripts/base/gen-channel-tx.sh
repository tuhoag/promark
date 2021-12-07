#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function createChannelTx() {
#   if [ ! -d "channels" ]; then
#     mkdir channels
#   fi

    channelName=$1
	getChannelTxPath $channelName
	# channel_tx_path=$?

	set -x
	configtxgen -profile TwoOrgsChannel -outputCreateChannelTx $channelTxPath -channelID $channelName -configPath $CONFIG_PATH
	res=$?
	{ set +x; } 2>/dev/null

    verifyResult $res "Failed to generate channel configuration transaction..."
}



# channelName=$1

createChannelTx $1