'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
// const {}
const utils = require('./utils');
/**
 * Workload module for the benchmark round.
 */
class AddProofWorkload extends WorkloadModuleBase {

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
        const {numPeersPerOrgs, numOrgsPerType, numVerifiersPerType} = this.roundArguments;

        this.campaignIds = []

        // generate campaigns
        for (let i = 0; i < args.numCampaigns; i++) {
            const {camId, name, advertiser, business, verifierURLsStr} = utils.CreateCampaignArgs(numPeersPerOrgs, numOrgsPerType, numVerifiersPerType)
            const newCampaignId = "c" + i
            const newCampaignName = "campaign " + i
            const transArgs = {
                contractId: "campaign",
                contractFunction: 'CreateCampaign',
                contractArguments: [camId, name, advertiser, business, verifierURLsStr],
                readOnly: false
            };

            let raw_result = await this.sutAdapter.sendRequests(transArgs);
            let returnedCampaign = JSON.parse(raw_result.result.toString());
            this.campaignIds.push(returnedCampaign.id);
            // throw new Error(returnedCampaign.id);
        }

        // generate proofs
        this.proofs = []
        for (let i = 0; i < args.numProofs; i++) {
            const camIdx = Math.floor(Math.random()*10000) % this.campaignIds.length;
            const userId = Math.floor(Math.random()*10000);
            const camId = this.campaignIds[camIdx]

            const transArgs = {
                contractId: this.roundArguments.contractId,
                contractFunction: "GenerateCustomerCampaignProof",
                contractArguments: [camId, userId],
                readOnly: true
            };

            let raw_result = await this.sutAdapter.sendRequests(transArgs);
            let returnedProof = JSON.parse(raw_result.result.toString());
            this.proofs.push(returnedProof)
            // throw new Error(returnedProof.Rs.join(";"));
        }
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // proofId string, comm string, rsStr string
        // const {numCampaigns, numPeersPerOrgs, numOrgsPerType, numVerifiersPerType} = this.roundArguments;

        // const camIdx = Math.floor(Math.random()*10000) % numCampaigns;
        // const userId = Math.floor(Math.random()*10000);
        const proofId = Math.floor(Math.random()*10000);
        const proofIdx = Math.floor(Math.random()*10000) % this.proofs.length;
        const proof = this.proofs[proofIdx]

        const transArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: "AddCustomerProofCampaign",
            contractArguments: [proofId, proof.Comm, proof.Rs.join(";")],
            readOnly: false
        };

        throw new Error(JSON.stringify(transArgs));

        return this.sutAdapter.sendRequests(transArgs);
    }

    async cleanupWorkloadModule() {

    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new AddProofWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
