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

        let {camId, name, advertiser, publisher, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr} = utils.CreateCampaignArgs(numPeersPerOrgs, numOrgsPerType, numVerifiersPerType, numDevices)

        camId = "c0";
        const cTransArgs = {
            contractId: "campaign",
            contractFunction: 'CreateCampaign',
            contractArguments: [camId, "campaign 0", advertiser, publisher, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr],
            readOnly: false
        };

        this.campaignId = camId;
        await this.sutAdapter.sendRequests(cTransArgs);

        this.userId = Math.floor(Math.random()*10000);

        let pTransArgs = {
            contractId: "poc",
            contractFunction: "GeneratePoCAndTPoCProof",
            contractArguments: [this.campaignId, this.userId, 1],
            readOnly: true
        };

        let result = await this.sutAdapter.sendRequests(pTransArgs);
        this.customerPoC = JSON.parse(result["result"]);
        this.customerTPoC = this.customerPoC.tpocs[0];

        this.deviceId = "d1";

        pTransArgs = {
            contractId: "poc",
            contractFunction: "GeneratePoCAndTPoCProof",
            contractArguments: [this.campaignId, this.deviceId, 1],
            readOnly: true
        };

        result = await this.sutAdapter.sendRequests(pTransArgs);
        this.devicePoC = JSON.parse(result["result"]);
        this.deviceTPoC = this.devicePoC.tpocs[0];

    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // logger.info(`submit: ${JSON.stringify(this.TPoC)} - ${this.campaignId}`);

        // camId string, deviceId string, addedTimeStr int64, dCsStr string, dRsStr string, dHashesStr string, dKeyStr string, uCsStr string, uRsStr string, uHashesStr string, uKeyStr string
        const addedTime = 0
        const transArgs = {
            contractId: "proof",
            contractFunction: "AddCampaignTokenTransaction",
            contractArguments: [this.campaignId, this.deviceId, addedTime, this.deviceTPoC.tComms.join(";"), this.deviceTPoC.tRs.join(";"), this.deviceTPoC.hashes.join(";"), this.deviceTPoC.key, this.customerTPoC.tComms.join(";"), this.customerTPoC.tRs.join(";"), this.customerTPoC.hashes.join(";"), this.customerTPoC.key],
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
