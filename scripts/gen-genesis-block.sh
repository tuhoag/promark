#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function generateGenesisBlock() {
    which configtxgen
    if [ "$?" -ne 0 ]; then
        fatalln "configtxgen tool not found."
    fi

    if [ ! -d $CHANNEL_PATH ]; then
        mkdir $CHANNEL_PATH
    fi

    infoln "Generating Orderer Genesis block"

    #   cp ./config/configtx.yaml $OUTPUTS/configtx.yaml

    set -x
    configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock $CHANNEL_PATH/genesis.block -configPath $CONFIG_PATH
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
        fatalln "Failed to generate orderer genesis block..."
    fi
}

generateGenesisBlock
