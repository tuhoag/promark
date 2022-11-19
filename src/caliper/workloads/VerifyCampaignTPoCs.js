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
        const {numPeersPerOrgs, numOrgsPerType, numVerifiers, numDevices} = this.roundArguments

        // generate campaigns
        this.campaignIds = [];

        let numCampaigns = Math.floor(numPeersPerOrgs * numOrgsPerType / numVerifiers);
        if (numCampaigns < 1) {
            numCampaigns = 1;
        }

        for (let i = 0; i < numCampaigns; i++) {

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

        this.PoCs = [];
        this.TPoCs = [];
        for (let i = 0; i < numCampaigns; i++) {
            let userId = Math.floor(Math.random()*10000);

            let pTransArgs = {
                contractId: "proof",
                contractFunction: "GeneratePoCProof2",
                contractArguments: [this.campaignIds[i], userId],
                readOnly: true
            };

            let result = await this.sutAdapter.sendRequests(pTransArgs);
            let PoC = JSON.parse(result["result"]);

            let tpTransArgs = {
                contractId: "poc",
                contractFunction: "GenerateTPoCProofs",
                contractArguments: [this.campaignIds[i], PoC.comm, PoC.r, PoC.numVerifiers, 1],
                readOnly: true
            };

            result = await this.sutAdapter.sendRequests(tpTransArgs);
            let TPoC = JSON.parse(result["result"]).tpocs[0];
            this.PoCs.push(PoC);
            this.TPoCs.push(TPoC);
        }

        // throw Error(JSON.stringify(this.TPoCs))
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // logger.info(`submit: ${JSON.stringify(this.TPoC)} - ${this.campaignId}`);
        const tpocIdx = Math.floor(Math.random()*10000) % this.TPoCs.length;
        const tpoc = this.TPoCs[tpocIdx];

        // camId string, csStr string, rsStr string, hashesStr string, keyStr string
        const transArgs = {
            contractId: "proof",
            contractFunction: "VerifyTPoCProof",
            contractArguments: [this.campaignIds[tpocIdx], tpoc.tComms.join(";"), tpoc.tRs.join(";"), tpoc.hashes.join(";"), tpoc.key],
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
