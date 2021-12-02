#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function createChannelTx() {
#   if [ ! -d "channels" ]; then
#     mkdir channels
#   fi

    channel_name=$1
	getChannelTxPath $channel_name
	# channel_tx_path=$?

	set -x
	configtxgen -profile TwoOrgsChannel -outputCreateChannelTx $channel_tx_path -channelID $channel_name -configPath $CONFIG_PATH
	res=$?
	{ set +x; } 2>/dev/null

    verifyResult $res "Failed to generate channel configuration transaction..."
}



channel_name=$1

createChannelTx $channel_name