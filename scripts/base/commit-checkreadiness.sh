
#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function checkCommitReadiness() {
    local channelName=$1
    local chaincodeName=$2
    IFS=',' read -r -a orgTypes <<< $3
    local orgNum=$4
    local peerNum=$5
    local sequence=$6

    local maxPeerId=$(($peerNum - 1))
    local maxOrgId=$(($orgNum - 1))

    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrgId); do
            for peerId in $(seq 0 $maxPeerId); do
                sleep $DELAY
                infoln "Attempting to check the commit readiness of the chaincode definition on peer$peerId.$orgType$orgId, Retry after $DELAY seconds."

                set -x
                peer lifecycle chaincode checkcommitreadiness --channelID $channelName --name $chaincodeName --version $sequence --sequence $sequence --output json >&log.txt

                res=$?
                { set +x; } 2>/dev/null
                let rc=0
                for var in "$4"; do
                    infoln "checkCommitReadiness: var =$var"
                    grep "$var" log.txt &>/dev/null || let rc=1
                done
            done
        done
    done
}

checkCommitReadiness $1 $2 $3 $4 $5 $6
# checkCommitReadiness $1 $2 "adv" "\"adv0MSP\": true" "\"pub0MSP\": false"
# checkCommitReadiness $1 $2 "pub" "\"adv0MSP\": true" "\"pub0MSP\": false"
# checkCommitReadiness $1 $2 "adv" "\"adv0MSP\": true" "\"pub0MSP\": true"
# checkCommitReadiness $1 $2 "pub" "\"adv0MSP\": true" "\"pub0MSP\": true"