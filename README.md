### Test the System with Hyperledger Caliper
#### Run a Benchmark
- Go to caliper directory: `$ cd caliper`
- Initialize a project: `$ npm init -y`
- Install caliper: `$ npm install --only=prod @hyperledger/caliper-cli@0.4.0`
- Bind it: `$ npx caliper bind --caliper-bind-sut fabric:2.1 --caliper-bind-cwd ./`
- Run: `$ npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networkConfig.yaml --caliper-benchconfig benchmarks/queryCampaign.yaml  --caliper-fabric-gateway-enabled --caliper-flow-only-test`
- Check `caliper/report.html` for the results of the tests. 

## Monitoring
Prometheus is enabled in the project as a monitoring framework. In addition, Grafana is added for better visualization. You can access Prometheus at: `http://0.0.0.0:9090` and Grafana at `http://0.0.0.0:3000`.