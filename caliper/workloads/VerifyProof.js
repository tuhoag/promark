'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const utils = require('./utils');

const logger = require('@hyperledger/caliper-core').CaliperUtils.getLogger('promark');

let count = 0;
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
        const {numPeersPerOrgs, numOrgsPerType, numVerifiersPerType, numCampaigns, numProofs} = this.roundArguments;

        this.addedProofIds = [];
        count += 1;

        this.initData = utils.loadInitData(numCampaigns, numProofs, numVerifiersPerType);
        logger.info("this.initData: ", this.initData);
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // proofId string, comm string, rsStr string
        const proofId = `p${Math.floor(Math.random()*100000)}`;
        const proofIdx = Math.floor(Math.random()*10000) % this.initData.proofs.length;
        const proof = this.initData.proofs[proofIdx];
        // const proofId = `p:${proof.deviceId}:${proof.camId}:${proof.customerId}:${proof.addedTimeStr}`;
        // camId string, deviceId string, cusId string, cusComm string, cusRsStr string, addedTimeStr string
        // "camId": "c5842",
        // "deviceId": "d1",
        // "customerId": "u239",
        // "addedTimeStr": "2022-01-30T22:29:27.890Z",
        // "cusComm": "hCMWcCIZ0n1oY7moL+lzepz1XcPDo8UjnnM5hCvpdgE=",
        // "cusRsStr": "d/ZBfcHKWCNXjw5VsCGsxr9vQHf1rYrOUlAb1bJNqQc=;qQ9w12OSOJsjcc5C2eAfx0TGDSiUbSJjGdnQTIpEuwo="
        // const transArgs = {
        //     contractId: this.roundArguments.contractId,
        //     contractFunction: "AddCustomerProofCampaign2",
        //     contractArguments: [proof.camId, proof.deviceId, proof.customerId, proof.cusComm, proof.cusRsStr, proof.addedTimeStr],
        //     readOnly: false
        // };

        const camId = "c7679";
        const proofId = "p9999";

        const transArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: "VerifyCampaignProof",
            contractArguments: [camId, proofId],
            readOnly: true
        };
        // proofId := fmt.Sprintf("p:%s:%s:%s:%s", deviceId, camId, cusId, addedTimeStr)

        this.addedProofIds.push(proofId);
        logger.debug(`submitTransaction count: ${count}`);
        // throw new Error(`submitTransaction count: ${count}`);

        return this.sutAdapter.sendRequests(transArgs);
    }

    async cleanupWorkloadModule() {
        logger.info("addedProofIds.length:", this.addedProofIds.length);

        // logger.info(this.addedProofIds);
        // // throw new Error();
        // for (let proofId of this.addedProofIds) {
        //     const transArgs = {
        //         contractId: "proof",
        //         contractFunction: 'DeleteProofById',
        //         contractArguments: [proofId],
        //         readOnly: false
        //     };

        //     await this.sutAdapter.sendRequests(transArgs);
        // }
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
