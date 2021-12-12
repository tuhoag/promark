'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

/**
 * Workload module for the benchmark round.
 */
class CreateCampaignWorkload extends WorkloadModuleBase {

    /**
     * Initializes the workload module instance.
     */
    constructor() {
        super();
        this.contractId = '';
        this.contractVersion = '';
    }

    /**
     * Initialize the workload module with the given parameters.
     * @param {number} workerIndex The 0-based index of the worker instantiating the workload module.
     * @param {number} totalWorkers The total number of workers participating in the round.
     * @param {number} roundIndex The 0-based index of the currently executing round.
     * @param {Object} roundArguments The user-provided arguments for the round from the benchmark configuration file.
     * @param {ConnectorBase} sutAdapter The adapter of the underlying SUT.
     * @param {Object} sutContext The custom context object provided by the SUT adapter.
     * @async
     */
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        const args = this.roundArguments;
        this.contractId = args.contractId;
        this.contractVersion = args.contractVersion;
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // camId string, name string, advertiser string, business string, verifierURLStr string
        const camId = "c" + Math.floor(Math.random()*10000)
        const name = "Campaign " + camId
        const advertiser = "adv"+Math.floor(Math.random()*10000) % this.roundArguments.numOrgsPerType
        const business = "bus"+Math.floor(Math.random()*10000) % this.roundArguments.numOrgsPerType

        var verifierURLs = []

        // http://peer0.adv0.promark.com:5000
        // add adv verifiers
        for (let i = 0; i < this.roundArguments.numVerifiersPerType; i++) {
            const advertierPeerName = "peer" + Math.floor(Math.random()*10000) % this.roundArguments.numPeersPerOrgs
            const businessPeerName = "peer" + Math.floor(Math.random()*10000) % this.roundArguments.numPeersPerOrgs

            const advPeerURL = "http://" + advertierPeerName + "."+advertiser + ".promark.com:5000"
            const busPeerURL = "http://" + businessPeerName + "."+business + ".promark.com:5000"

            verifierURLs.push(advPeerURL)
            verifierURLs.push(busPeerURL)
        }

        var verifierURLsStr = verifierURLs.join(";")
        // add bus verifiers

        // const myArgs = {
        //     contractId: this.roundArguments.contractId,
        //     contractFunction: 'QueryBackup',
        //     invokerIdentity: 'peer0.org1.example.com',
        //     contractArguments: [`BACKUP_${this.workerIndex}_${randomId}`],
        //     readOnly: true
        // };

        const myArgs = {
            contractId: this.contractId,
            contractFunction: 'CreateCampaign',
            contractArguments: [camId, name, advertiser, business, verifierURLsStr],
            readOnly: true
        };
        return this.sutAdapter.sendRequests(myArgs);
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new CreateCampaignWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
