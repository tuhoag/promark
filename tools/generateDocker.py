#!/usr/bin/env python3
from __future__ import print_function
from typing import NewType

import yaml
import sys

# SERVICES = {
#   'web' : 'ubuntu:14.04',
# }

SERVICES = {
  'peer': 'peer-verifier-base',
  'couchdb': 'couchdb:3.1.1',
}

# NW = {'test':''}

ORDERER = {
  'orderer': 'orderer-base'
}

LOG = {
  'log': 'logs.dockerfile'
}

EXT = {
  'external' : 'external-service.dockerfile'
}

# define the list of ports
ADV_BASE_PORT=1050
BUS_BASE_PORT=2050
PEER_LOCAL_PORT=7051

ADV_COUCHDB_PORT=5984
BUS_COUCHDB_PORT=6484
COUCHDB_LOCAL_PORT=5984

ORDERER_LOCAL_PORT=7050
ORDERER_LISTEN_PORT=9443

LOG_PORT=5003
EXT_PORT=5000
ADV_WEB_PORT=8500
BUS_WEB_PORT=9000
ADV_GOSSIP_PORT=55000
BUS_GOSSIP_PORT=60000
LOCAL_GOSSIP_PORT=9443
REDIS_PORT=6379

#define the specific name
org_suffix='promark.com'
var_suffix='${PROJECT_NAME}.com'
compose_suffix='${COMPOSE_PROJECT_NAME}.com'

# COMPOSITION = {'version': '2', 'network': 'test', 'services': {}}
COMPOSITION = {'services': {}}

HEADER = {'version': '2', 'networks': 'test'}

# NETWORK ={'network':{}}

def servicize(name, image, org, n):
  orgid = org[3:]
  orgname = org[:3]
  print(orgid, orgname)

  entry = {'container_name': name,
            # 'network': 'test'
          }
  entry['networks']=['test']

  if name.startswith('couchdb'):
    if orgname == 'adv':
      couchdb_port = ADV_COUCHDB_PORT + (int(orgid)*100)+ int(n)
    elif orgname == 'bus':
      couchdb_port = BUS_COUCHDB_PORT + (int(orgid)*100)+ int(n)

    map_couchdb_port = str(couchdb_port)+':'+str(COUCHDB_LOCAL_PORT)
    entry['image'] = image
    entry['environment'] = ['COUCHDB_USER=admin',
                            'COUCHDB_PASSWORD=adminpw'
                           ]
    entry['ports']=[map_couchdb_port]

  elif name.startswith('orderer'):
    orderer_port=str(ORDERER_LOCAL_PORT)+':'+str(ORDERER_LOCAL_PORT)
    orderer_listen_map_port='53732'+':'+str(ORDERER_LISTEN_PORT)
    orderer_listen_port=str(ORDERER_LISTEN_PORT)
    
    entry['extends']= {'file':'docker-compose-base.yml',
                       'service': image, 
                      }

    entry['ports']=[orderer_port,
                    orderer_listen_map_port,
                  ]

    # envStr= 'CORE_OPERATIONS_LISTENADDRESS=orderer.${PROJECT_NAME}.com'+':'+str(orderer_listen_port)
    entry['environment']=['CORE_OPERATIONS_LISTENADDRESS=orderer.${PROJECT_NAME}.com'+':'+str(orderer_listen_port)]

    entry['volumes']=['../channels/genesis.block:/var/hyperledger/orderer/orderer.genesis.block',
                      '../organizations/ordererOrganizations/${PROJECT_NAME}.com/orderers/orderer.${PROJECT_NAME}.com/msp:/var/hyperledger/orderer/msp',
                      '../organizations/ordererOrganizations/${PROJECT_NAME}.com/orderers/orderer.${PROJECT_NAME}.com/tls/:/var/hyperledger/orderer/tls',
    ]
  elif name.startswith('peer'):
    if orgname == 'adv':
      port = ADV_BASE_PORT + (int(orgid)*100)+ int(n)
      webPort = ADV_WEB_PORT + (int(orgid)*100)+ int(n)
      gossipPort= ADV_GOSSIP_PORT + (int(orgid)*100)+ int(n)
    elif orgname == 'bus':
      port = BUS_BASE_PORT + (int(orgid)*100)+ int(n)
      webPort = BUS_WEB_PORT + (int(orgid)*100)+ int(n)
      gossipPort= BUS_GOSSIP_PORT + (int(orgid)*100)+ int(n)
    
    mapPort ="{0}:{1}".format(port, PEER_LOCAL_PORT)
    mapWebPort="{0}:{1}".format(webPort, webPort)
    mapGossipPort="{0}:{1}".format(gossipPort, LOCAL_GOSSIP_PORT)

    entry['extends']= {'file':'docker-compose-base.yml', 
                       'service': image, 
                      }
    entry['ports']=[mapPort,
                    mapWebPort,
                    mapGossipPort,
                    ]

    # for environment part of each peer
    # - CORE_PEER_ID=peer0.bus0.${PROJECT_NAME}.com
    # - CORE_PEER_ADDRESS=peer0.bus0.${PROJECT_NAME}.com:7051
    # - CORE_PEER_LISTENADDRESS=0.0.0.0:7051
    # - CORE_PEER_CHAINCODEADDRESS=peer0.bus0.${PROJECT_NAME}.com:7052
    # - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
    # - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.bus0.${PROJECT_NAME}.com:7051
    # - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.bus0.${PROJECT_NAME}.com:7051
    # - CORE_PEER_LOCALMSPID=bus0MSP
    # - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
    # - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb1.bus0.promark.com:5984
    # - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin
    # - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw
  
    peerName="peer{0}.{1}.{2}".format(n, org, var_suffix)
    couchdbName="couchdb{0}.{1}.{2}".format(n, org, org_suffix)
    corePeerId="CORE_PEER_ID={0}".format(peerName)
    corePeerAdd="CORE_PEER_ADDRESS={0}:{1}".format(peerName, PEER_LOCAL_PORT)
    corePeerListenAdd="CORE_PEER_LISTENADDRESS=0.0.0.0:{0}".format(PEER_LOCAL_PORT)
    corePeerChaincodeAdd="CORE_PEER_CHAINCODEADDRESS={}:7052".format(peerName)
    corePeerChaincodeListenAdd="CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052"
    corePeerGossipExtEndpoint="CORE_PEER_GOSSIP_EXTERNALENDPOINT={0}:{1}".format(peerName, PEER_LOCAL_PORT)
    corePeerGossipBootstrap="CORE_PEER_GOSSIP_BOOTSTRAP={0}:{1}".format(peerName, PEER_LOCAL_PORT)
    corePeerLocalMspId="CORE_PEER_LOCALMSPID={0}MSP".format(org)
    coreLedgerStateDB="CORE_LEDGER_STATE_STATEDATABASE=CouchDB"
    coreLedgerStateCouchDBConfigAdd="CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS={0}:{1}".format(couchdbName, COUCHDB_LOCAL_PORT)
    coreLedgerStateCouchDBUsername="CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin"
    coreLedgerStateCoudchDBPassword="CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw"
    
    verPort="VER_PORT={0}".format(webPort) 
    verName="VER_NAME={0}.log".format(name)

    entry['environment']=[corePeerId,
                          corePeerAdd,
                          corePeerListenAdd,
                          corePeerChaincodeAdd,
                          corePeerChaincodeListenAdd,
                          corePeerGossipExtEndpoint,
                          corePeerGossipBootstrap,
                          corePeerLocalMspId,
                          coreLedgerStateDB,
                          coreLedgerStateCouchDBConfigAdd,
                          coreLedgerStateCouchDBUsername,
                          coreLedgerStateCoudchDBPassword,
                          verPort,
                          verName,
    ]
    # - /var/run/docker.sock:/host/var/run/docker.sock
    # - ../organizations/peerOrganizations/bus0.${PROJECT_NAME}.com/peers/peer0.bus0.${PROJECT_NAME}.com/msp:/etc/hyperledger/fabric/msp
    # - ../organizations/peerOrganizations/bus0.${PROJECT_NAME}.com/peers/peer0.bus0.${PROJECT_NAME}.com/tls:/etc/hyperledger/fabric/tls

    volumnSock="/var/run/docker.sock:/host/var/run/docker.sock"
    volumnMsp="../organizations/peerOrganizations/{0}.{1}/peers/{2}/msp:/etc/hyperledger/fabric/msp".format(org, var_suffix, peerName)
    volumnTls="../organizations/peerOrganizations/{0}.{1}/peers/{2}/tls:/etc/hyperledger/fabric/tls".format(org, var_suffix, peerName)
    
    entry['volumes']=[volumnSock,
                      volumnMsp,
                      volumnTls,
    ]

    # - orderer.${COMPOSE_PROJECT_NAME}.com
    # - couchdb1.bus0.promark.com
    ordererName="orderer.{0}".format(compose_suffix)

    entry['depends_on']=[ordererName,
                        couchdbName]

  return entry

def generatePeer(name, image, org, n):
  orgid = org[3:]
  orgname = org[:3]
  print(orgid, orgname)

  entry = {'container_name': name,
          }

  if name.startswith('peer'):
    if orgname == 'adv':
      port = ADV_BASE_PORT + (int(orgid)*100)+ int(n)
    elif orgname == 'bus':
      port = BUS_BASE_PORT + (int(orgid)*100)+ int(n)
    
    mapPort = str(port) + ':' + str(PEER_LOCAL_PORT)

    entry['extends']= [{'file':'docker-compose-base.yml'}, 
                        {'service': image}
                      ]
    entry['ports']=[mapPort]

    # for environment part of each peer
    # - CORE_PEER_ID=peer0.bus0.${PROJECT_NAME}.com
    # - CORE_PEER_ADDRESS=peer0.bus0.${PROJECT_NAME}.com:7051
    # - CORE_PEER_LISTENADDRESS=0.0.0.0:7051
    # - CORE_PEER_CHAINCODEADDRESS=peer0.bus0.${PROJECT_NAME}.com:7052
    # - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
    # - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.bus0.${PROJECT_NAME}.com:7051
    # - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.bus0.${PROJECT_NAME}.com:7051
    # - CORE_PEER_LOCALMSPID=bus0MSP
    peerName="peer{0}.{1}.{2}".format(n, org, var_suffix)
    corePeerId="CORE_PEER_ID={0}".format(peerName)
    corePeerAdd="CORE_PEER_ADDRESS={0}:{1}".format(peerName, PEER_LOCAL_PORT)
    corePeerListenAdd="CORE_PEER_LISTENADDRESS=0.0.0.0:{0}".format(PEER_LOCAL_PORT)
    corePeerChaincodeAdd="CORE_PEER_CHAINCODEADDRESS={}:7052".format(peerName)
    corePeerChaincodeListenAdd="CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052"
    corePeerGossipExtEndpoint="CORE_PEER_GOSSIP_EXTERNALENDPOINT={0}:{1}".format(peerName, PEER_LOCAL_PORT)
    corePeerGossipBootstrap="CORE_PEER_GOSSIP_BOOTSTRAP={0}:{1}".format(peerName, PEER_LOCAL_PORT)
    corePeerLocalMspId="CORE_PEER_LOCALMSPID={0}MSP".format(org)

    entry['environment']=[corePeerId,
                          corePeerAdd,
                          corePeerListenAdd,
                          corePeerChaincodeAdd,
                          corePeerChaincodeListenAdd,
                          corePeerGossipExtEndpoint,
                          corePeerGossipBootstrap,
                          corePeerLocalMspId,
                        ]
  return entry

# The arguments to run this file is: number of peers/org org1 org2 ...
def main():
  n = len(sys.argv)
  print("\nNumber of peers:", n)

  orgs = []
  for i in range(n):
    if i >=2:
      orgs.append(sys.argv[i])
  print(orgs)

  # for name, image in SERVICES.items():
  #   COMPOSITION1['services'][name] = servicize(name, image)
  with open('docker-compose.yml', 'w') as f:
      f.write(yaml.dump(HEADER, default_flow_style=False, indent=4))

  for name, image in ORDERER.items():
    # name = name +str(org_suffix)
    name="{0}.{1}".format(name, org_suffix)
    COMPOSITION['services'][name] = servicize(name, image, '', '')

  for n in range(0, int(sys.argv[1])):
    for org in orgs:
      for name, image in SERVICES.items():
        # name = name + str(n)+'.'+org+str(org_suffix)
        name="{0}{1}.{2}.{3}".format(name, n, org, org_suffix)
        print(name)
        COMPOSITION['services'][name] = servicize(name, image, org, str(n))

  # print(yaml.dump(COMPOSITION, default_flow_style=False, indent=4), end='')
  with open('docker-compose.yml', 'a') as f:
      f.write(yaml.dump(COMPOSITION, default_flow_style=False, indent=4))


if __name__ == '__main__':
  print('Generate docker file')
  main()


