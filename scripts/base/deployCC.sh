#!/bin/bash

. $BASE_SCRIPTS_DIR/utils.sh


infoln "Deploying CC"


# function deployChaincode() {
#     # generate packages
#     packageChaincode

#     # # install chaincode
#     # installChaincode 1
#     # installChaincode 2

#     # queryInstalled 1
#     # queryInstalled 2

#     # approveForMyOrg 1

#     # ## check whether the chaincode definition is ready to be committed
#     # ## expect org1 to have approved and org2 not to
#     # checkCommitReadiness 1 "\"Org1MSP\": true" "\"Org2MSP\": false"
#     # checkCommitReadiness 2 "\"Org1MSP\": true" "\"Org2MSP\": false"

#     # approveForMyOrg 2

#     # ## check whether the chaincode definition is ready to be committed
#     # ## expect them both to have approved
#     # checkCommitReadiness 1 "\"Org1MSP\": true" "\"Org2MSP\": true"
#     # checkCommitReadiness 2 "\"Org1MSP\": true" "\"Org2MSP\": true"

#     # # now that we know for sure both orgs have approved, commit the definition
#     # commitChaincodeDefinition 1 2

#     # queryCommitted 1
#     # queryCommitted 2
# }