'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const utils = require('./utils');
/**
 * Workload module for the benchmark round.
 */
class GenerateProofWorkload extends WorkloadModuleBase {

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
        // const {numPeersPerOrgs, numOrgsPerType, numVerifiersPerType} = this.roundArguments;

        // this.campaignIds = []

        // for (let i = 0; i < args.numCampaigns; i++) {
        //     const {camId, name, advertiser, business, verifierURLsStr} = utils.CreateCampaignArgs(numPeersPerOrgs, numOrgsPerType, numVerifiersPerType)
        //     const transArgs = {
        //         contractId: "campaign",
        //         contractFunction: 'CreateCampaign',
        //         contractArguments: [camId, name, advertiser, business, verifierURLsStr],
        //         readOnly: true
        //     };

        //     this.campaignIds.push(camId)
        //     await this.sutAdapter.sendRequests(transArgs);
        // }
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // camId string, userId string
        const {numCampaigns, numPeersPerOrgs, numOrgsPerType, numVerifiersPerType} = this.roundArguments;

        const camIdx = Math.floor(Math.random()*10000) % numCampaigns;
        const userId = Math.floor(Math.random()*10000);
        // const camId = this.campaignIds[camIdx]

        const transArgs = {
            contractId: "proof",
            contractFunction: "GenerateCustomerCampaignProof",
            contractArguments: [camIdx, userId],
            readOnly: true
        };

        return this.sutAdapter.sendRequests(transArgs);
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new GenerateProofWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
