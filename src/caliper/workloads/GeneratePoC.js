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
        const {numPeersPerOrgs, numOrgsPerType, numVerifiers, numDevices} = this.roundArguments;
        this.campaignIds = [];
        this.txIndex = 0;

        for (let i = 0; i < args.numCampaigns; i++) {

            const {camId, name, advName, pubName, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr} = utils.CreateCampaignUnequalVerifiersArgs(numOrgsPerType, numPeersPerOrgs, numVerifiers, numDevices)

            // throw Error(advertiser)
            const newCampaignId = "c" + i
            const newCampaignName = "campaign " + i

            const transArgs = {
                contractId: "campaign",
                contractFunction: 'CreateCampaign',
                contractArguments: [newCampaignId, newCampaignName, advName, pubName, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr],
                readOnly: false
            };

            this.campaignIds.push(newCampaignId)
            await this.sutAdapter.sendRequests(transArgs);
        }
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        const {numCampaigns, numPeersPerOrgs, numOrgsPerType, numVerifiersPerType} = this.roundArguments;

        const camIdx = Math.floor(Math.random()*10000) % this.campaignIds.length;
        const userId = Math.floor(Math.random()*10000);
        const camId = this.campaignIds[camIdx]

        const transArgs = {
            contractId: "poc",
            contractFunction: "GeneratePoCProof",
            contractArguments: [camId, userId],
            invokerIdentity: 'peer0.adv0.promark.com',
            readOnly: true
        };

        return this.sutAdapter.sendRequests(transArgs);
    }

    async cleanupWorkloadModule() {
        const transArgs = {
            contractId: "campaign",
            contractFunction: 'DeleteAllCampaigns',
            contractArguments: [],
            readOnly: false
        };
        await this.sutAdapter.sendRequests(transArgs);
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
