#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function commitChaincode() {
    # local chaincodeName=$1
    # local channelName=$2
    # local orgType=$3
    # local orgNum=$4
    # local peerNum=$5
    # local sequence=$6


    local chaincodeName=$1
    local channelName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5
    local sequence=$6

    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"

    infoln "Commiting chaincode $chaincodeName in channel '$channelName'..."

    parsePeerConnectionParameters $orgTypes $orgNum $peerNum

    res=$?
    verifyResult $res "Invoke transaction failed on channel '$channelName' due to uneven number of peer and org parameters "

    # while 'peer chaincode' command can get the orderer endpoint from the
    # peer (if join was successful), let's supply it directly as we know
    # it using the "-o" option
    set -x
    peer lifecycle chaincode commit -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME  --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --tls $peerConnectionParams --version $sequence --sequence $sequence

    # peer lifecycle chaincode commit -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS --version $sequence --sequence $sequence >&log.txt

    #peer lifecycle chaincode commit -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --version 1.0 --package-id $packageId --sequence 1 >&log.txt

    res=$?

    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Chaincode definition commit failed on peer.org${ORG} on channel '$channelName' failed"
    successln "Chaincode definition committed on channel '$channelName'"

    peer lifecycle chaincode querycommitted --channelID $channelName --name $chaincodeName

}


commitChaincode $1 $2 $3 $4 $5 $6