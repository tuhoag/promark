#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function generateGenesisBlock() {
    local orgNum=$1
    local peerNum=$2
    local channelName=$3

    which configtxgen
    if [ "$?" -ne 0 ]; then
        fatalln "configtxgen tool not found."
    fi

    if [ ! -d $CHANNEL_PATH ]; then
        infoln "creating folder ${CHANNEL_PATH}"
        mkdir $CHANNEL_PATH
    fi

    infoln "Generating Orderer Genesis block"

    #   cp ./config/configtx.yaml $OUTPUTS/configtx.yaml
    getBlockPath $orgNum $peerNum $channelName

    infoln "Blockpath: ${blockPath}"
    set -x
    configtxgen -profile "${orgNum}${peerNum}OrdererGenesis" -channelID system-channel -outputBlock $blockPath -configPath $CONFIG_PATH
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
        fatalln "Failed to generate orderer genesis block..."
    fi
}

generateGenesisBlock $1 $2 $3
