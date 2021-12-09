#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh

function packageChaincode() {
    local chaincode_name=$1
    local sequence=$2
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincode_name}.tar.gz"
    local chaincode_label="${chaincode_name}_${sequence}"
    local chaincode_package_src_path="${CHAINCODE_SRC_PATH}/${chaincode_name}"
    infoln "Packaging chaincode $chaincode_name"

    # infoln "Vendoring Go dependencies at $chaincode_package_src_path"
    # pushd $chaincode_package_src_path
    # GO111MODULE=on go mod vendor
    # popd
    # successln "Finished vendoring Go dependencies"

    # if [ ! -d $CC_PACKAGE_FOLDER_OUTPUT ]; then
    #     mkdir $CC_PACKAGE_FOLDER_OUTPUT
    # fi

    set -x
    peer lifecycle chaincode package $chaincode_package_path --path $chaincode_package_src_path --lang $CHAINCODE_LANGUAGE --label $chaincode_label >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    cat log.txt
    verifyResult $res "Chaincode packaging has failed"
    successln "Chaincode is packaged"
}

packageChaincode $1 $2
