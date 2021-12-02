#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function commitChaincode() {
    local chaincodeName=$1
    local channelName=$2
    local orgNum=$3
    local peerNum=$4

    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"

    infoln "Commiting chaincode $chaincodeName in channel '$channelName'..."

    parsePeerConnectionParameters $orgNum $peerNum

    res=$?
    verifyResult $res "Invoke transaction failed on channel '$channelName' due to uneven number of peer and org parameters "

    # while 'peer chaincode' command can get the orderer endpoint from the
    # peer (if join was successful), let's supply it directly as we know
    # it using the "-o" option
     set -x
    peer lifecycle chaincode commit -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName $PEER_CONN_PARMS --version 1.0 --sequence 1 >&log.txt

    #peer lifecycle chaincode commit -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --version 1.0 --package-id $packageId --sequence 1 >&log.txt

     res=$?

     { set +x; } 2>/dev/null
     cat log.txt
     verifyResult $res "Chaincode definition commit failed on peer0.org${ORG} on channel '$channelName' failed"
     successln "Chaincode definition committed on channel '$channelName'"

    peer lifecycle chaincode querycommitted --channelID $channelName --name $chaincodeName

}

function parsePeerConnectionParameters() {
    local orgNum=$1
    local peerNum=$2

    PEER_CONN_PARMS=""
    PEERS=""
    local peerNames=""

    infoln "$orgNum ; $peerNum"

    local maxOrgId=$(($orgNum - 1))
    local maxPeerId=$(($peerNum - 1))

    for orgId in $(seq 0 $maxOrgId); do
         infoln $orgId
         for peerId in $(seq 0 $maxPeerId); do
             for orgType in "adv" "bus"; do
                 local peerName="peer${peerId}.${orgType}${orgId}"
                 infoln "processed $peerName"
                 selectPeer $orgType $orgId $peerId
                 PEERS="$peerNames ${peerName}"
                 PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"
                ## Set path to TLS certificate
                TLSINFO=$(eval echo "--tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE")
                PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
             done
         done
     done

    infoln "parsePeerConnectionParameters1: $PEERS $PEER_CONN_PARMS"
}

commitChaincode $1 $2 $3 $4