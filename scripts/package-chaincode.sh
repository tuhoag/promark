#!/bin/bash

. $SCRIPTS_DIR/utils.sh


infoln "Deploying CC"


function packageChaincode() {
    local chaincode_name=$1
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincode_name}.tar.gz"
    local chaincode_label="${chaincode_name}_1.0"

    infoln "Packaging chaincode $chaincode_name"

    infoln "Vendoring Go dependencies at $CHAINCODE_SRC_PATH"
    pushd $CHAINCODE_SRC_PATH
    GO111MODULE=on go mod vendor
    popd
    successln "Finished vendoring Go dependencies"

    # if [ ! -d $CC_PACKAGE_FOLDER_OUTPUT ]; then
    #     mkdir $CC_PACKAGE_FOLDER_OUTPUT
    # fi

    set -x
    peer lifecycle chaincode package $chaincode_package_path --path $CHAINCODE_PACKAGE_DIR --lang $CHAINCODE_LANGUAGE --label $chaincode_label >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    cat log.txt
    verifyResult $res "Chaincode packaging has failed"
    successln "Chaincode is packaged"
}

packageChaincode $1
