#!/usr/bin/env python3
from __future__ import print_function
from typing import NewType

import yaml
import sys

# SERVICES = {
#   'web' : 'ubuntu:14.04',
# }

CHANNELPEER = {
  'peers': {}
}

# COUCHDB = {
#   'couchdb': 'couchdb:3.1.1',
# }

ORDERER_ITEMS = {
   'orderer.promark.com': {}
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
channel_name = 'mychannel'

HEADER = {'name': 'Caliper test', 'version': '1.0'}

CALIPER = {'caliper': {}}

CHANNELS = {'channels': {}}

ORDERERS = {'orderers': {}}

CLIENTS = {'clients': {}}

ORGANIZATIONS = {'organizations': {}}

PEERS = {'peers': {}}

def generateCaliper():
  entry = {'blockchain': 'fabric'}

  return entry

def generateChannels():
  settingVar = "true"
  entry = {}
  entry['created'] = settingVar
  entry['orderers'] = ['orderer.promark.com']

  entry['contracts'] = {'id':'campaign', 'version': "1.0", 'language':'golang'}

  return entry

def generateOrderer(name, image):
  url = "grpcs://0.0.0.0:7050"

  entry = {}
  entry['url']= url
  entry['tlsCACerts']={'path':'./../organizations/ordererOrganizations/promark.com/tlsca/tlsca.promark.com-cert.pem'}
  entry['grpcOptions'] = {'ssl-target-name-override':name}
  return entry

def generateChannelPeers():
  entry = {}
  return entry

def generateChannelPeer():
  settingVar = "true"
  entry = {}
  entry= {'endorsingPeer': settingVar,
          'chaincodeQuery': "true",
          'ledgerQuery': 'true',
          'eventSource': 'true',}
  return entry

# def generateClientsPeer():
#   entry = {'client': {}}
#   return entry

def generateClients(org, n):
  orgName = "{0}{1}".format(org, n)
  path = '../organizations/peerOrganizations/{}.promark.com/users/Admin@{}.promark.com/msp/admincerts'.format(orgName, orgName)
  cryptoPath= "../organizations/peerOrganizations/{0}.promark.com/users/Admin@{1}.promark.com/msp".format(orgName, orgName)
  privateKeyPath= "../organizations/peerOrganizations/{0}.promark.com/users/Admin@{1}.promark.com/msp/keystore/priv_sk".format(orgName, orgName)
  clientSignedCertPath = "../organizations/peerOrganizations/{0}.promark.com/users/Admin@{1}.promark.com/msp/signcerts/Admin@{2}.promark.com-cert.pem".format(orgName, orgName, orgName)

  entry = {'organization':orgName}
  entry['credentialStore'] ={'path': path,
                              'cryptoStore': {'path': cryptoPath},
                            }
  entry['clientPrivateKey']={'path':privateKeyPath,                             
                          }
  
  entry['clientSignedCert']= {'path':clientSignedCertPath}
  return entry

def generateOrg(orgName, peerNumber):
  mspID = "{0}MSP".format(orgName)
  peerArray = []
  adminPrivateKeyPath = "../organizations/peerOrganizations/{0}.promark.com/users/Admin@{1}.promark.com/msp/keystore/priv_sk".format(orgName, orgName)
  signedCertPath = "../organizations/peerOrganizations/{0}.promark.com/users/Admin@{1}.promark.com/msp/signcerts/Admin@{2}.promark.com-cert.pem".format(orgName, orgName, orgName)

  for i in range (0, int(peerNumber)):
    peerName="peer{0}.{1}.{2}".format(i, orgName, org_suffix)
    peerArray.append(peerName)

  entry = {'mspid':mspID}
  entry['peers']= peerArray
  entry['adminPrivateKey'] = {'path':adminPrivateKeyPath}
  entry['signedCert'] = {'path':signedCertPath}

  return entry

def generatePeer(peerName, peerID, org, orgID):
  orgName = "{0}{1}".format(org, orgID)
  if org == 'adv':
    port = ADV_BASE_PORT + (int(orgID)*100)+ int(peerID)
  elif org == 'bus':
    port = BUS_BASE_PORT + (int(orgID)*100)+ int(peerID)
  url = "grpcs://0.0.0.0:{}".format(port)
  tlsCACertsPath = "./../organizations/peerOrganizations/{0}.promark.com/tlsca/tlsca.{1}.promark.com-cert.pem".format(orgName, orgName)

  entry = {}

  entry['url'] = url
  entry['tlsCACerts'] = {'path':tlsCACertsPath}
  entry['grpcOptions'] = {'ssl-target-name-override': peerName}
  return entry

# The arguments to run this file is:
# <number of peer> <org_name1> <num of org> <org_name2> <num of org> 
def main():
  # Dictionary Methods
  orgs = {}.fromkeys(['adv', 'bus'], 0)
  peerNumber = sys.argv[1]

  if sys.argv[2] == 'adv':
    orgs['adv'] = sys.argv[3]
  elif sys.argv[2] == 'bus':
    orgs['bus'] = sys.argv[3]
  
  if sys.argv[4] == 'adv':
    orgs['adv'] = sys.argv[5]
  elif sys.argv[4] == 'bus':
    orgs['bus'] = sys.argv[5]
  print(orgs)

  # Header
  with open('networkConfig.yaml', 'w') as f:
      f.write(yaml.dump(HEADER, default_flow_style=False, indent=4))

  # caliper
  CALIPER['caliper']  = generateCaliper()

  # channels
  CHANNELS['channels'][channel_name] = generateChannels()
  CHANNELS['channels'][channel_name]['peers'] = generateChannelPeers()
  for org in orgs:
    print(orgs.get(org), org)
    for n in range(0, int(orgs.get(org))):
      print(n)
      for peer in CHANNELPEER.items():
        for i in range (0, int(peerNumber)):
          peerName="peer{0}.{1}{2}.{3}".format(i, org, n, org_suffix)
          CHANNELS['channels'][channel_name]['peers'][peerName]= generateChannelPeer()

  # orderers
  for name, image in ORDERER_ITEMS.items():
    ORDERERS['orderers'][name] = generateOrderer(name, image)

  # clients
  for org in orgs:
    print(orgs.get(org), org)
    for n in range(0, int(orgs.get(org))):
      print(n)
      # for peer in CHANNELPEER.items():
      for i in range (0, int(peerNumber)):
        peerName="peer{0}.{1}{2}.{3}".format(i, org, n, org_suffix)
        CLIENTS['clients'][peerName]= {}
        CLIENTS['clients'][peerName]['client']= generateClients(org, n)

  # orgs
  for org in orgs:
    print(orgs.get(org), org)
    for n in range(0, int(orgs.get(org))):
      orgName ="{0}{1}".format(org, n)
      ORGANIZATIONS['organizations'][orgName] =generateOrg(orgName, peerNumber)

  # peers
  for org in orgs:
    print(orgs.get(org), org)
    for n in range(0, int(orgs.get(org))):
      print(n)
      for i in range (0, int(peerNumber)):
        peerName="peer{0}.{1}{2}.{3}".format(i, org, n, org_suffix)
        PEERS['peers'][peerName]=generatePeer(peerName, i, org, n)

  # print(yaml.dump(COMPOSITION, default_flow_style=False, indent=4), end='')
  with open('networkConfig.yaml', 'a') as f:
      f.write(yaml.dump(CALIPER, default_flow_style=False, indent=4))
  with open('networkConfig.yaml', 'a') as f:
      f.write(yaml.dump(CHANNELS, default_flow_style=False, indent=4))
  with open('networkConfig.yaml', 'a') as f:
      f.write(yaml.dump(ORDERERS, default_flow_style=False, indent=4))
  with open('networkConfig.yaml', 'a') as f:
      f.write(yaml.dump(CLIENTS, default_flow_style=False, indent=4))
  with open('networkConfig.yaml', 'a') as f:
      f.write(yaml.dump(ORGANIZATIONS, default_flow_style=False, indent=4))
  with open('networkConfig.yaml', 'a') as f:
      f.write(yaml.dump(PEERS, default_flow_style=False, indent=4))
  
if __name__ == '__main__':
  print('Generate Caliper networkConfig file')
  main()


