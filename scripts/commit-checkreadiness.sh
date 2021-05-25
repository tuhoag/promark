
#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function checkCommitReadiness() {
    local chaincodeName=$1
    local channelName=$2
    local orgType=$3

    local maxOrgId=0
    local maxPeerId=0

    #shift 1

    for orgId in $(seq 0 $maxOrgId); do
        for peerId in $(seq 0 $maxPeerId); do
            #for orgType in "adv" "bus"; do
            selectPeer $orgType $orgId $peerId
            echo $orgType $orgId $peerId
            #done
        done
    done

    local rc=1
    local COUNTER=1
    infoln "Checking the commit readiness of the chaincode definition on peer$peerId.$orgType$orgId on channel '$channelName'..."

    while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
        sleep $DELAY
        infoln "Attempting to check the commit readiness of the chaincode definition on peer$peerId.$orgType$orgId, Retry after $DELAY seconds."

        set -x
        peer lifecycle chaincode checkcommitreadiness --channelID $channelName --name $chaincodeName --version 1.0 --sequence 1 --output json >&log.txt

        res=$?
        { set +x; } 2>/dev/null
        let rc=0
        for var in "$4"; do
            echo "checkCommitReadiness: var =$var"
            grep "$var" log.txt &>/dev/null || let rc=1
        done
        COUNTER=$(expr $COUNTER + 1)
        echo "checkCommitReadiness: rc =$rc"
    done

    cat log.txt
    if test $rc -eq 0; then
        echo "checkCommitReadiness: res =$res"
        infoln "Checking the commit readiness of the chaincode definition successful on peer$peerId.$orgType$orgId} on channel '$channelName'"
    else
        fatalln "After $MAX_RETRY attempts, Check commit readiness result on peer$peerId.$orgType$orgId is INVALID!"
    fi
}

checkCommitReadiness $1 $2 "adv" "\"adv0MSP\": true" "\"bus0MSP\": false"
checkCommitReadiness $1 $2 "bus" "\"adv0MSP\": true" "\"bus0MSP\": false"
checkCommitReadiness $1 $2 "adv" "\"adv0MSP\": true" "\"bus0MSP\": true"
checkCommitReadiness $1 $2 "bus" "\"adv0MSP\": true" "\"bus0MSP\": true"