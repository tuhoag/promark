#!/usr/bin/env python3
from __future__ import print_function

import yaml
import sys

# SERVICES = {
#   'web' : 'ubuntu:14.04',
# }

SERVICES = {
  'peer': 'peer-verifier-base',
  'couchdb': 'couchdb:3.1.1',
}

ORDERER = {
  'orderer': 'orderer-base'
}


ADV_BASE_PORT=1050
BUS_BASE_PORT=2050
PEER_LOCAL_PORT=7051
ORDERER_PORT=7050
#port=$(($base_port + $org_id * 100 + $peer_id))


# COMPOSITION = {'version': '2', 'network': 'test', 'services': {}}
COMPOSITION = {'services': {}}

HEADER = {'version': '2', 'network': 'test'}

def servicize(name, image, org, n):
  entry = {'container-name': name,
             'network': 'test',
          }
  if name.startswith('couchdb'):
    entry['image'] = image
    entry['environment'] = ['COUCHDB_USER=admin',
                            'COUCHDB_PASSWORD=adminpw'
                           ]
  elif name.startswith('peer'):
    orgid = org[3:]
    orgname = org[:3]
    print(orgid, orgname)

    if orgname == 'adv':
      port = ADV_BASE_PORT + (int(orgid)*100)+ int(n)
    elif orgname == 'bus':
      port = BUS_BASE_PORT + (int(orgid)*100)+ int(n)
    
    mapPort = str(port) + ':' + str(PEER_LOCAL_PORT)

    entry['extends']= [{'file':'docker-compose-base.yml'}, 
                       {'service': image}
                      ]
    entry['ports']=[mapPort,
                   ]
  elif name.startswith('orderer'):
    entry['extends']= [{'file':'docker-compose-base.yml'}, 
                       {'service': image}
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
  with open('docker-compose.yaml', 'w') as f:
      f.write(yaml.dump(HEADER, default_flow_style=False, indent=4))

  for name, image in ORDERER.items():
    COMPOSITION['services'][name] = servicize(name, image, '', '')

  for n in range(0, int(sys.argv[1])):
    for org in orgs:
      for name, image in SERVICES.items():
        name = name + str(n)+'.'+org+'.'+'promark.com'
        print(name)
        COMPOSITION['services'][name] = servicize(name, image, org, str(n))

  # print(yaml.dump(COMPOSITION, default_flow_style=False, indent=4), end='')
  with open('docker-compose.yaml', 'a') as f:
      f.write(yaml.dump(COMPOSITION, default_flow_style=False, indent=4))


if __name__ == '__main__':
  print('Generate docker file')
  main()


