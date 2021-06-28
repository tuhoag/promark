#!/bin/bash

. $SCRIPTS_DIR/utils.sh


infoln "Cleaning the repository"

$SCRIPTS_DIR/stop.sh

# remove organizations
rm -rf $ORGANIZATION_OUTPUTS

# remove volumes
rm -rf volumes

# remove channels
rm -rf channels

# remove log
rm -rf $PWD/log.txt

# remove the chaincode package file before commit
# rm chaincode/main.tar.gz
