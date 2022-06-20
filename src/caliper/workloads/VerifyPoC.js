'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const utils = require('./utils');

const logger = require('@hyperledger/caliper-core').CaliperUtils.getLogger('my-module');

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
        const {numPeersPerOrgs, numOrgsPerType, numVerifiersPerType, numDevices, numPoCs} = this.roundArguments


        const {camId, name, advertiser, publisher, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr} = utils.CreateCampaignArgs(numPeersPerOrgs, numOrgsPerType, numVerifiersPerType, numDevices)
        const cTransArgs = {
            contractId: "campaign",
            contractFunction: 'CreateCampaign',
            contractArguments: [camId, name, advertiser, publisher, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr],
            readOnly: false
        };

        this.campaignId = camId;
        await this.sutAdapter.sendRequests(cTransArgs);

        const userId = Math.floor(Math.random()*10000);

        const pTransArgs = {
            contractId: "poc",
            contractFunction: "GeneratePoCProof",
            contractArguments: [this.campaignId, userId],
            readOnly: true
        };

        let result = await this.sutAdapter.sendRequests(pTransArgs);
        // logger.info(`after create poc:${JSON.stringify(result)} - data: ${result["result"]}`);
        this.PoC = JSON.parse(result["result"]);
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // logger.info(`submit: ${JSON.stringify(this.PoC)} - ${this.campaignId} - ${this.PoC.comm} - ${this.PoC.rs}}`);

        // camId string, cStr string, rsStr string
        const transArgs = {
            contractId: "proof",
            contractFunction: "VerifyPoCProof",
            contractArguments: [this.campaignId, this.PoC.comm, this.PoC.rs.join(";")],
            readOnly: true
        };

        return this.sutAdapter.sendRequests(transArgs);
    }

    async cleanupWorkloadModule() {
        const transArgs = {
            contractId: "campaign",
            contractFunction: 'DeleteCampaignById',
            contractArguments: [this.campaignId],
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
