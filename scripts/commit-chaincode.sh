#!/bin/bash

. $SCRIPTS_DIR/utils.sh

function commitChaincode() {
    local chaincodeName=$1
    local channelName=$2
    local orgNum=$3
    local peerNum=$4

    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"
    local peer_name="peer${peerId}.${orgType}${orgId}"

    infoln "Commiting chaincode $chaincodeName in channel '$channelName'..."

    parsePeerConnectionParameters $orgNum $peerNum
    # res=$?
    # verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

    # while 'peer chaincode' command can get the orderer endpoint from the
    # peer (if join was successful), let's supply it directly as we know
    # it using the "-o" option
    # set -x
    # peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} $PEER_CONN_PARMS --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
    # res=$?
    # { set +x; } 2>/dev/null
    # cat log.txt
    # verifyResult $res "Chaincode definition commit failed on peer0.org${ORG} on channel '$CHANNEL_NAME' failed"
    # successln "Chaincode definition committed on channel '$CHANNEL_NAME'"
}

function parsePeerConnectionParameters() {
    local orgNum=$1
    local peerNum=$2

    PEER_CONN_PARMS=""
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

                peerNames="$peerNames ${peerName}"
            done

        done

    done

    infoln $peerNames

    # while [ "$#" -gt 0 ]; do
    #     selectPeer $1
    #     PEER="peer0.org$1"
    #     ## Set peer addresses
    #     PEERS="$PEERS $PEER"
    #     PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"
    #     ## Set path to TLS certificate
    #     TLSINFO=$(eval echo "--tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE")
    #     PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"

    #     # infoln "PEERS_CONN_PARMS: ${PEER_CONN_PARMS}"
    #     # infoln "PEERS: ${PEERS}"
    #     # shift by one to get to the next organization
    #     shift
    # done
    # # remove leading space for output
    # PEERS="$(echo -e "$PEERS" | sed -e 's/^[[:space:]]*//')"
}

commitChaincode $1 $2 $3 $4
