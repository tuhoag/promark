'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const utils = require('./utils');
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

        this.campaignIds = []
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // camId string, name string, advertiser string, publisher string, verifierURLStr string
        const {numPeersPerOrgs, numOrgsPerType, numVerifiersPerType, numDevices} = this.roundArguments
        const {camId, name, advertiser, publisher, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr} = utils.CreateCampaignArgs(numPeersPerOrgs, numOrgsPerType, numVerifiersPerType, numDevices)

        const transArgs = {
            contractId: "campaign",
            contractFunction: 'CreateCampaign',
            contractArguments: [camId, name, advertiser, publisher, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr],
            readOnly: false
        };

        this.campaignIds.push(camId);

        return this.sutAdapter.sendRequests(transArgs);
    }

    async cleanupWorkloadModule() {
        // const args = this.roundArguments;

        // for (let i = 0; i < this.campaignIds.length; i++) {
        //     const transArgs = {
        //         contractId: "campaign",
        //         contractFunction: 'DeleteCampaignById',
        //         contractArguments: [this.campaignIds[i]],
        //         readOnly: false
        //     };

        //     // this.campaignIds.push(newCampaignId)
        //     await this.sutAdapter.sendRequests(transArgs);
        // }
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
