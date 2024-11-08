#!/bin/bash

# environment variables
export PROJECT_DIR=$PWD
export SCRIPTS_DIR=$PROJECT_DIR/scripts
export SCENARIOS_DIR=$PROJECT_DIR/scenarios
export BASE_SCRIPTS_DIR=$PROJECT_DIR/scripts/base
export FABRIC_BIN_PATH=$PROJECT_DIR/bin
export CONFIG_PATH=$PROJECT_DIR/config
export CREDENTIALS_OUTPUTS=$PROJECT_DIR/credentials
export ORG_CONFIG_PATH=$CONFIG_PATH
export CHANNEL_PATH=$PROJECT_DIR/channels
export DOCKER_COMPOSE_DIR_PATH=$PROJECT_DIR/docker
export FABRIC_CFG_PATH=$CONFIG_PATH
export CHAINCODE_SRC_PATH=$PROJECT_DIR/src/chaincodes
export CALIPER_DIR_PATH=$PROJECT_DIR/src/caliper
export CLIENT_DIR_PATH=$PROJECT_DIR/src/application

export FABRIC_VERSION=2.2
export PROJECT_NAME=promark

export NETWORK_NAME="${PROJECT_NAME}net"
export LOGSPOUT_PORT=3004
export ADV_BASE_PORT=5000
export PUB_BASE_PORT=6000

export ORDERER_ADDRESS=0.0.0.0:7050
export ORDERER_HOSTNAME=orderer.$PROJECT_NAME.com
export ORDERER_CA=$CREDENTIALS_OUTPUTS/ordererOrganizations/$PROJECT_NAME.com/orderers/$ORDERER_HOSTNAME/msp/tlscacerts/tlsca.$PROJECT_NAME.com-cert.pem

export PATH=$FABRIC_BIN_PATH:$PATH

export CHAINCODE_LANGUAGE=golang
export CHAINCODE_PACKAGE_DIR=$CHAINCODE_SRC_PATH/packages

# Adding for chaincodeCheckReadiness function
export MAX_RETRY="2"
export DELAY="10"

#For external service
export EXTERNAL_SERVCE_PORT=3002
