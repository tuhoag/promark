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
        const {numPeersPerOrgs, numOrgsPerType, numVerifiers, numDevices, numTrans} = this.roundArguments

        this.data = [];
        this.txIndex = 0;
        this.campaigns = utils.CreateCampaignsWithEqualVerifiersArgs(numOrgsPerType, numPeersPerOrgs, numVerifiers, numDevices);
        // this.startTime = 0;
        // logger.debug(`${JSON.stringify(campaigns)}`);
        // console.log(this.campaigns);
        // throw Error();
        let camIds = [];

        for (let i = 0; i < this.campaigns.length; i++) {
            const {id, name, advName, pubName, startTime, endTime, verifierURLsStr, deviceIdsStr} = this.campaigns[i];
            this.startTime = startTime;

            const cTransArgs = {
                contractId: "campaign",
                contractFunction: 'CreateCampaign',
                contractArguments: [id, name, advName, pubName, startTime, endTime, verifierURLsStr, deviceIdsStr],
                readOnly: false
            };

            camIds.push(id);

            await this.sutAdapter.sendRequests(cTransArgs);
        }

        let numAddedTrans = 0;
        let numTransPerCampaigns;

        if (numTrans == 0) {
            numTrans = campaigns.length;
        }

        numTransPerCampaigns = Math.floor(numTrans / this.campaigns.length);

        for (let i =0; i < numTrans; i++) {
            let userId = Math.floor(Math.random()*10000);
            let deviceId = Math.floor(Math.random()*10000 % numDevices);

            let iCam = i % this.campaigns.length;

            let pTransArgs = {
                contractId: "poc",
                contractFunction: "GeneratePoCAndTPoCProof",
                contractArguments: [this.campaigns[iCam].id, userId, 1],
                readOnly: true
            };

            let result = await this.sutAdapter.sendRequests(pTransArgs);
            let customerPoCAndTPoCs = JSON.parse(result["result"]);

            pTransArgs = {
                contractId: "poc",
                contractFunction: "GeneratePoCAndTPoCProof",
                contractArguments: [this.campaigns[iCam].id, deviceId, 1],
                readOnly: true
            };

            result = await this.sutAdapter.sendRequests(pTransArgs);
            let devicePoCAndTPoCs = JSON.parse(result["result"]);
            // let deviceTPoC = this.devicePoC.tpocs[0];

            // for (let j = 0; j < curNumTransPerCampaigns; j ++) {
                this.data.push({
                    camId: this.campaigns[iCam].id,
                    deviceId: deviceId,
                    deviceTPoC: devicePoCAndTPoCs.tpocs[0],
                    customerTPoC: customerPoCAndTPoCs.tpocs[0],
                })
            // }
        }

        // for (let i = 0; i < this.campaigns.length; i++) {
        //     let userId = Math.floor(Math.random()*10000);
        //     let deviceId = Math.floor(Math.random()*10000 % numDevices);

        //     // generate number of tpocs
        //     let curNumTransPerCampaigns = numTransPerCampaigns;
        //     if (i == this.campaigns.length - 1) {
        //         curNumTransPerCampaigns = numTrans - numAddedTrans;
        //     }

        //     let pTransArgs = {
        //         contractId: "poc",
        //         contractFunction: "GeneratePoCAndTPoCProof",
        //         contractArguments: [this.campaigns[i].id, userId, curNumTransPerCampaigns],
        //         readOnly: true
        //     };

        //     let result = await this.sutAdapter.sendRequests(pTransArgs);
        //     let customerPoCAndTPoCs = JSON.parse(result["result"]);

        //     pTransArgs = {
        //         contractId: "poc",
        //         contractFunction: "GeneratePoCAndTPoCProof",
        //         contractArguments: [this.campaigns[i].id, deviceId, curNumTransPerCampaigns],
        //         readOnly: true
        //     };

        //     result = await this.sutAdapter.sendRequests(pTransArgs);
        //     let devicePoCAndTPoCs = JSON.parse(result["result"]);
        //     // let deviceTPoC = this.devicePoC.tpocs[0];

        //     for (let j = 0; j < curNumTransPerCampaigns; j ++) {
        //         this.data.push({
        //             camId: this.campaigns[i].id,
        //             deviceId: deviceId,
        //             deviceTPoC: devicePoCAndTPoCs.tpocs[j],
        //             customerTPoC: customerPoCAndTPoCs.tpocs[j],
        //         })
        //     }
        // }

        // console.dir(this.data, {depth: null});
        // throw Error();
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // logger.info(`submit: ${JSON.stringify(this.TPoC)} - ${this.campaignId}`);

        // return contract.submitTransaction("AddCampaignTokenTransaction", camId, deviceId, addedTime, deviceTPoC.tComms.join(";"), deviceTPoC.tRs.join(";"), deviceTPoC.hashes.join(";"), deviceTPoC.key, customerTPoC.tComms.join(";"), customerTPoC.tRs.join(";"), customerTPoC.hashes.join(";"), customerTPoC.key);
        // if (this.txIndex == 1) {
        //     console.log()
        // }

        const addedTime = 0
        const {camId, deviceId, deviceTPoC, customerTPoC} = this.data[this.txIndex];

        // console.log(this.txIndex);
        // console.log(camId, deviceId, deviceTPoC, customerTPoC);

        const transArgs = {
            contractId: "proof",
            contractFunction: "AddCampaignTokenTransaction",
            contractArguments: [camId, deviceId, addedTime, deviceTPoC.tComms.join(";"), deviceTPoC.tRs.join(";"), deviceTPoC.hashes.join(";"), deviceTPoC.key, customerTPoC.tComms.join(";"), customerTPoC.tRs.join(";"), customerTPoC.hashes.join(";"), customerTPoC.key],
            readOnly: false
        };
        this.txIndex += 1;
        return this.sutAdapter.sendRequests(transArgs);
    }

    async cleanupWorkloadModule() {
        let transArgs = {
            contractId: "campaign",
            contractFunction: 'DeleteAllCampaigns',
            contractArguments: [],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(transArgs);

        transArgs = {
            contractId: "proof",
            contractFunction: 'DeleteAllProofs',
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
