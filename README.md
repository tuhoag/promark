# Advertising Blockchain

## Run the network
- Initialize credentials by using cryptogen: `$ ./main.sh init`
- Create channel and let peer.adv0 and peer0.bus0 join the created channel: `$ ./main.sh channel all 1 1`
- Package chaincode and deploy it in peer0.adv0 and peer0.bus0: `$ ./main.sh chaincode all 1 1`

`chmod -R a+rwx promark/`

## Test the System with Hyperledger Caliper
### Run a Benchmark
- Go to caliper directory: `$ cd caliper`
- Initialize a project: `$ npm init -y`
- Install caliper: `$ npm install --only=prod @hyperledger/caliper-cli@0.4.0`
- Bind it: `$ npx caliper bind --caliper-bind-sut fabric:2.1 --caliper-bind-cwd ./`
- Run: `$ npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networkConfig.yaml --caliper-benchconfig benchmarks/CreateCampaign.yaml  --caliper-fabric-gateway-enabled --caliper-flow-only-test`

`$ npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networkConfig.yaml --caliper-benchconfig benchmarks/testConfigPeak.yaml  --caliper-fabric-gateway-enabled --caliper-flow-only-test`

- Check `caliper/report.html` for the results of the tests.

## Monitoring
Prometheus is enabled in the project as a monitoring framework. In addition, Grafana is added for better visualization. You can access Prometheus at: `http://0.0.0.0:9090` and Grafana at `http://0.0.0.0:3000`.

## Ports
There are two types of organizations: adv, bus with some additional services such as CA, couchdb, verifier,... Their ports are organized as follows:
 - ADV0:
    + Peer0: 5000
        + CouchDB: 5001
        + Verifier: 5002
    + Peer1: 5010
        + CouchDB: 5011
        + Verifier: 5012
- ADV1:
    + Peer0: 5100
        + CouchDB: 5101
        + Verifier: 5102
    + Peer1: 5110
        + CouchDB: 5111
        + Verifier: 5112
- BUS0:
    + Peer0: 6000
        + CouchDB: 6001
        + Verifier: 6002
    + Peer1: 6010
        + CouchDB: 6011
        + Verifier: 6012
- BUS1:
    + Peer0: 6100
        + CouchDB: 6101
        + Verifier: 6102
    + Peer1: 6110
        + CouchDB: 6111
        + Verifier: 6112