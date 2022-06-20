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
        const {numPeersPerOrgs, numOrgsPerType, numVerifiersPerType, numDevices} = this.roundArguments

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
            contractFunction: "GeneratePoCAndTPoCProof",
            contractArguments: [this.campaignId, userId, 1],
            readOnly: true
        };

        let result = await this.sutAdapter.sendRequests(pTransArgs);
        this.PoC = JSON.parse(result["result"]);
        this.TPoC = this.PoC.tpocs[0];

    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // logger.info(`submit: ${JSON.stringify(this.TPoC)} - ${this.campaignId}`);

        // camId string, csStr string, rsStr string, hashesStr string, keyStr string
        const transArgs = {
            contractId: "proof",
            contractFunction: "VerifyTPoCProof",
            contractArguments: [this.campaignId, this.TPoC.tComms.join(";"), this.TPoC.tRs.join(";"), this.TPoC.hashes.join(";"), this.TPoC.key],
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
