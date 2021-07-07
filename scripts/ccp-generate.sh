#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

# function json_ccp {
#     local PP=$(one_line_pem $4)
#     local CP=$(one_line_pem $5)
#     sed -e "s/\${ORG}/$1/" \
#         -e "s/\${P0PORT}/$2/" \
#         -e "s/\${CAPORT}/$3/" \
#         -e "s#\${PEERPEM}#$PP#" \
#         -e "s#\${CAPEM}#$CP#" \
#         organizations/ccp-template.json
# }

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        $CONFIG_PATH/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

org=adv
ORG=0
P0PORT=7051
CAPORT=7054
PEERPEM=organizations/peerOrganizations/adv0.promark.com/tlsca/tlsca.adv0.promark.com-cert.pem
CAPEM=organizations/peerOrganizations/adv0.promark.com/ca/ca.adv0.promark.com-cert.pem

# echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/org1.example.com/connection-org1.json
echo "$(yaml_ccp $org$ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/adv0.promark.com/connection-adv0.yaml

org=bus
ORG=0
P0PORT=9051
CAPORT=8054
PEERPEM=organizations/peerOrganizations/bus0.promark.com/tlsca/tlsca.bus0.promark.com-cert.pem
CAPEM=organizations/peerOrganizations/bus0.promark.com/ca/ca.bus0.promark.com-cert.pem

# echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/org2.example.com/connection-org2.json
echo "$(yaml_ccp $org$ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/bus0.promark.com/connection-bus0.yaml
