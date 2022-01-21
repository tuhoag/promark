#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function generateOrgs() {
    local orgNum=$1
    local peerNum=$2

    local configPath="${ORG_CONFIG_PATH}/crypto-config-${orgNum}-${peerNum}.yaml"
    local outputPath=$CREDENTIALS_OUTPUTS

    infoln "Config path: $configPath"
    infoln "Output path: $outputPath"

    set -x
    cryptogen generate --config=$configPath --output=$outputPath
    res=$?
    { set +x; } 2>/dev/null

    if [ $res -ne 0 ]; then
        fatalln "Failed to generate certificates..."
    fi
}


# check if cryptogen is accessible
set -x
which cryptogen
{ set +x; } 2>/dev/null
if [ "$?" -ne 0 ]; then
    fatalln "cryptogen tool not found."
fi

generateOrgs $1 $2

