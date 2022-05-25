#!/usr/bin/env python3
from __future__ import print_function
from typing import NewType

import yaml
import sys

# SERVICES = {
#   'web' : 'ubuntu:14.04',
# }

PEER = {
  'peer': 'peer-verifier-base',
}

COUCHDB = {
  'couchdb': 'couchdb:3.1.1',
}

ORDERER = {
  'orderer': 'orderer-base'
}

LOG = {
  'logs': 'logs.dockerfile'
}

EXT = {
  'external' : 'external-service.dockerfile'
}

# define the list of ports
ADV_BASE_PORT=1050
PUB_BASE_PORT=2050
PEER_LOCAL_PORT=7051

ADV_COUCHDB_PORT=5984
PUB_COUCHDB_PORT=6484
COUCHDB_LOCAL_PORT=5984

ORDERER_LOCAL_PORT=7050
ORDERER_LISTEN_PORT=9443

LOG_PORT=5003
EXT_PORT=5000
ADV_WEB_PORT=8500
PUB_WEB_PORT=9000
ADV_GOSSIP_PORT=55000
PUB_GOSSIP_PORT=60000
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

def generateLogService(name, image):
  entry = {'container_name': name}
  entry['networks']=['test']

  entry['build']= {'context':'.',
                    'dockerfile': image,
                  }
  log_port = "{0}:{1}".format(LOG_PORT, LOG_PORT)
  entry['ports']=[log_port]
  entry['volumes']= ['../services/log:/log']
  entry['command']= '/bin/sh run.sh'

  log_env = "LOG_PORT={0}".format(str(LOG_PORT))
  entry['environment']=[log_env]

  return entry

def generateExtService(name, image):
  entry = {'container_name': name}
  entry['networks']=['test']

  entry['build']= {'context':'.',
                  'dockerfile': image,
                }
  ext_port = "{0}:{1}".format(EXT_PORT, EXT_PORT)
  entry['ports']=[ext_port]

  entry['volumes']= ['../services/ext:/code']
  entry['command']= '/bin/sh run.sh'
  api_port = "API_PORT={0}".format(str(EXT_PORT))
  redis_port = "REDIS_PORT={0}".format(str(REDIS_PORT))
  entry['environment']=[api_port,
                        redis_port]

  return entry

def generateOrderer(name, image):
  entry = {'container_name': name}
  entry['networks']=['test']

  # orderer_port=str(ORDERER_LOCAL_PORT)+':'+str(ORDERER_LOCAL_PORT)
  orderer_port="{0}:{1}".format(str(ORDERER_LOCAL_PORT), str(ORDERER_LOCAL_PORT))
  # orderer_listen_map_port='53732'+':'+str(ORDERER_LISTEN_PORT)
  orderer_listen_map_port="53732:{0}".format(str(ORDERER_LISTEN_PORT))
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

  return entry

def generateCouchDB(name, image, orgname, orgid, peerid):
  entry = {'container_name': name,
            # 'network': 'test'
          }
  entry['networks']=['test']

  if orgname == 'adv':
    couchdb_port = ADV_COUCHDB_PORT + (int(orgid)*10)+ int(peerid)
  elif orgname == 'pub':
    couchdb_port = PUB_COUCHDB_PORT + (int(orgid)*10)+ int(peerid)

  map_couchdb_port = str(couchdb_port)+':'+str(COUCHDB_LOCAL_PORT)
  entry['image'] = image
  entry['environment'] = ['COUCHDB_USER=admin',
                          'COUCHDB_PASSWORD=adminpw'
                          ]
  entry['ports']=[map_couchdb_port]

  return entry

def generatePeer (name, image, orgname, orgid, peerid):
  entry = {'container_name': name,
            # 'network': 'test'
          }
  entry['networks']=['test']

  if orgname == 'adv':
    port = ADV_BASE_PORT + (int(orgid)*10)+ int(peerid)
    webPort = ADV_WEB_PORT + (int(orgid)*10)+ int(peerid)
    gossipPort= ADV_GOSSIP_PORT + (int(orgid)*10)+ int(peerid)
  elif orgname == 'pub':
    port = PUB_BASE_PORT + (int(orgid)*10)+ int(peerid)
    webPort = PUB_WEB_PORT + (int(orgid)*10)+ int(peerid)
    gossipPort= PUB_GOSSIP_PORT + (int(orgid)*10)+ int(peerid)

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
  # - CORE_PEER_ID=peer0.pub0.${PROJECT_NAME}.com
  # - CORE_PEER_ADDRESS=peer0.pub0.${PROJECT_NAME}.com:7051
  # - CORE_PEER_LISTENADDRESS=0.0.0.0:7051
  # - CORE_PEER_CHAINCODEADDRESS=peer0.pub0.${PROJECT_NAME}.com:7052
  # - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
  # - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.pub0.${PROJECT_NAME}.com:7051
  # - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.pub0.${PROJECT_NAME}.com:7051
  # - CORE_PEER_LOCALMSPID=pub0MSP
  # - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
  # - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb1.pub0.promark.com:5984
  # - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin
  # - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw

  peerName="peer{0}.{1}{2}.{3}".format(peerid, orgname, orgid, var_suffix)
  couchdbName="couchdb{0}.{1}{2}.{3}".format(peerid, orgname, orgid, org_suffix)
  corePeerId="CORE_PEER_ID={0}".format(peerName)
  corePeerAdd="CORE_PEER_ADDRESS={0}:{1}".format(peerName, PEER_LOCAL_PORT)
  corePeerListenAdd="CORE_PEER_LISTENADDRESS=0.0.0.0:{0}".format(PEER_LOCAL_PORT)
  corePeerChaincodeAdd="CORE_PEER_CHAINCODEADDRESS={}:7052".format(peerName)
  corePeerChaincodeListenAdd="CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052"
  corePeerGossipExtEndpoint="CORE_PEER_GOSSIP_EXTERNALENDPOINT={0}:{1}".format(peerName, PEER_LOCAL_PORT)
  corePeerGossipBootstrap="CORE_PEER_GOSSIP_BOOTSTRAP={0}:{1}".format(peerName, PEER_LOCAL_PORT)
  corePeerLocalMspId="CORE_PEER_LOCALMSPID={0}{1}MSP".format(orgname, orgid)
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
  # - ../organizations/peerOrganizations/pub0.${PROJECT_NAME}.com/peers/peer0.pub0.${PROJECT_NAME}.com/msp:/etc/hyperledger/fabric/msp
  # - ../organizations/peerOrganizations/pub0.${PROJECT_NAME}.com/peers/peer0.pub0.${PROJECT_NAME}.com/tls:/etc/hyperledger/fabric/tls

  volumnSock="/var/run/docker.sock:/host/var/run/docker.sock"
  volumnMsp="../organizations/peerOrganizations/{0}{1}.{2}/peers/{3}/msp:/etc/hyperledger/fabric/msp".format(orgname, orgid, var_suffix, peerName)
  volumnTls="../organizations/peerOrganizations/{0}{1}.{2}/peers/{3}/tls:/etc/hyperledger/fabric/tls".format(orgname, orgid, var_suffix, peerName)

  entry['volumes']=[volumnSock,
                    volumnMsp,
                    volumnTls,
  ]

  # - orderer.${COMPOSE_PROJECT_NAME}.com
  # - couchdb1.pub0.promark.com
  ordererName="orderer.{0}".format(compose_suffix)

  entry['depends_on']=[ordererName,
                      couchdbName]

  return entry

# The arguments to run this file is:
# <number of peer> <org_name1> <num of org> <org_name2> <num of org>
def main():
  # Dictionary Methods
  orgs = {}.fromkeys(['adv', 'pub'], 0)
  peerNumber = sys.argv[1]

  if sys.argv[2] == 'adv':
    orgs['adv'] = sys.argv[3]
  elif sys.argv[2] == 'pub':
    orgs['pub'] = sys.argv[3]

  if sys.argv[4] == 'adv':
    orgs['adv'] = sys.argv[5]
  elif sys.argv[4] == 'pub':
    orgs['pub'] = sys.argv[5]
  print(orgs)

  with open('docker-compose.yml', 'w') as f:
      f.write(yaml.dump(HEADER, default_flow_style=False, indent=4))

  for name, image in EXT.items():
    name="{0}.{1}".format(name, org_suffix)
    COMPOSITION['services'][name] = generateExtService(name, image)

  for name, image in LOG.items():
    name="{0}.{1}".format(name, org_suffix)
    COMPOSITION['services'][name] = generateLogService(name, image)

  for name, image in ORDERER.items():
    name="{0}.{1}".format(name, org_suffix)
    COMPOSITION['services'][name] = generateOrderer(name, image)

  for org in orgs:
    print(orgs.get(org), org)

    for n in range(0, int(orgs.get(org))):
      print(n)
      for i in range (0, int(peerNumber)):
        for couch, image in COUCHDB.items():
          couch="{0}{1}.{2}{3}.{4}".format(couch, i, org, n, org_suffix)
          print(couch+'\n')
          COMPOSITION['services'][couch] = generateCouchDB(couch, image, org, n, i)

    for n in range(0, int(orgs.get(org))):
      print(n, org)
      for i in range (0, int(peerNumber)):
        for peer, image in PEER.items():
          peer="{0}{1}.{2}{3}.{4}".format(peer, i, org, n, org_suffix)
          print(peer+'\n')
          COMPOSITION['services'][peer] = generatePeer(peer, image, org, n, i)

  # print(yaml.dump(COMPOSITION, default_flow_style=False, indent=4), end='')
  with open('docker-compose.yml', 'a') as f:
        f.write(yaml.dump(COMPOSITION, default_flow_style=False, indent=4))


if __name__ == '__main__':
  print('Generate docker file')
  main()


