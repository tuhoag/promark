# ProMark: Privacy and Transparency-Aware Proximity Advertising Platform

## Installation
- Docker@3: [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)
- Docker-Compose@1.28.4: [https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)
- python@3.8: [https://www.python.org/downloads/](https://www.python.org/downloads/)
- node.js: [https://nodejs.org/en/download/](https://nodejs.org/en/download/)
- Go@1.18: [https://go.dev/doc/install](https://go.dev/doc/install)
- Hyperledger Fabric@2.2: This project includes the binaries of Hyperledger Fabric@2.2 in the `bin` folder. If you don't want to use them, feel free to download them from [the official Hyperledger Fabric Samples](https://github.com/hyperledger/fabric-samples).

## Structure
This project includes the following folders & files:
- `bin`: Hyperledger Fabric binaries.
- `config`: configurations of Hyperleder Fabric network (`core.yaml`, `configtx.yaml`) and digital certificates (`crypto-config.yaml`). `network` folder contains the settings of various network architecture.
- `docker`: docker images for ProMark peers, services (e.g., crypto, log, prometheus), and templates used to generate docker-compose files for various network architecture.
- `exp_data`: experimental results and figures visualizing the results.
- `scripts`: `base` contains scripts to interact with Hyperledger Fabric binaries. The scripts can be used for various Hyperledger Fabric. The remaining scripts in `scripts` use those in `base` to deploy ProMark.
- `services`: this folder is deprecated and will be removed in the next version.
- `src`: contains ProMark primary source code
    - `application` (Node.js): contains source codes of the client application that interact with Hyperledger Fabric chaincode (smart contract),
    - `chaincodes` (Go): contains ProMark chaincodes,
    - `caliper` (Node.js): contains experimental simulation in Caliper,
    - `ext` (Go): contains the crypto service,
    - `verifier` (Go): contains source code of peer verifiers,
    - `log` (Go): contains log services,
    - `visualization` (Python): visualizes experimental results.
- `docker.py`: source code to generate docker-compose files.
- `main.sh`: the main script to interact with ProMark.
- `settings.sh`: all ProMark environment variables.

## Run the network
All the scripts required to interact with ProMark is in `main.sh`. However, the simplest way to run ProMark is to restart the whole network and deploy all the chaincodes.
- Restart the network: `$./main.sh restart <N_ORGS> <N_PEERS> all`, where `N_ORGS` and `N_PEERS` are the number of orgnizations per organization type and the number of peers per organization.
- Test all chaincodes: `$./main.sh test <N_ORGS> <N_PEERS>` to execute the test case that creates campaigns, generates tokens and their one-time versions, and add the token transactions to the ledger.

## Run the simulation
The subcommand `evaluate` in `main.sh` contains all commands used to evaluate most important chaincodes.
- Evaluate the `CreateCampaign` chaincode: `$./main.sh evaluate <N_ORGS> <N_PEERS> campaign create`
- Evaluate the `GeneratePoC` chaincode: `$./main.sh evaluate <N_ORGS> <N_PEERS> proof gen`
- Evaluate the `GenerateTPoC` chaincode: `$./main.sh evaluate <N_ORGS> <N_PEERS> proof verifytpoc`
- Evaluate the `AddCampaignTokenTransaction` chaincode: `$./main.sh evaluate <N_ORGS> <N_PEERS> proof add`


## Monitoring
To make it easier to debug ProMark, I created a script to monitor logs of all peers in ProMark network: `$./main.sh monitor`.

## Ports
We export ProMark peers' ports to make it easier to debug ProMark. There are two types of organizations: advertisers (adv) and publisher (pub) with some additional services such as CA, couchdb, verifier,... Their ports can be calcualted as follows:
- Peer `<peerId>.<orgType><orgId>.promark.com`: orgBasePort + orgId * 100 + peerId * 10, where orgBasePort of advertisers is 5000 and that of publishers is 6000.
- CouchDB of peer `<peerId>.<orgType><orgId>.promark.com`: Peer port + 1.
- API of peer `<peerId>.<orgType><orgId>.promark.com`: Peer Port + 2.

`docker.py` contains source codes to calculate the ports for the docker compose file.

## Citation
```
@ARTICLE{promark,
  author={Hoang, Anh-Tu and Carminati, Barbara and Ferrari, Elena},
  journal={IEEE Transactions on Dependable and Secure Computing}, 
  title={ProMark: Ensuring Transparency and Privacy-Awareness in Proximity Marketing Advertising Campaigns}, 
  year={2024},
  pages={1-12},
  doi={10.1109/TDSC.2024.3478049}}
```
