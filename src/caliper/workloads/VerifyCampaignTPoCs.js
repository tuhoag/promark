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



        // // generate campaigns
        // this.campaigns = utils.CreateCampaignsWithEqualVerifiersArgs(numOrgsPerType, numPeersPerOrgs, numVerifiers, numDevices);

        // // throw Error(`${numPeersPerOrgs} * ${numOrgsPerType * 2} / ${numVerifiers} = ${this.campaigns.length}`);
        // // throw Error(`${JSON.stringify(this.campaigns)}`);
        // for (let i = 0; i < this.campaigns.length; i++) {
        //     const {camId, name, advName, pubName, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr} = this.campaigns[i];

        //     let transArgs = {
        //         contractId: "campaign",
        //         contractFunction: 'CreateCampaign',
        //         contractArguments: [camId, name, advName, pubName, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr],
        //         readOnly: false
        //     };

        //     // throw Error(`${JSON.stringify(transArgs.contractArguments)}`);

        //     await this.sutAdapter.sendRequests(transArgs);
        // }

        // this.PoCs = [];
        // this.TPoCs = [];
        // this.tokenTrans = [];
        // // let numTrans = Math.floor(Math.random()*10000 % maxNumTrans + minNumTrans);
        // let numAddedTrans = 0;
        // let numTransPerCampaigns = Math.floor(numTrans / this.campaigns.length);

        // for (let i = 0; i < this.campaigns.length; i++) {
        //     let userId = Math.floor(Math.random()*10000);
        //     let deviceId = Math.floor(Math.random()*10000 % numDevices);

        //     let pTransArgs = {
        //         contractId: "proof",
        //         contractFunction: "GeneratePoCProof2",
        //         contractArguments: [this.campaigns[i].camId, userId],
        //         readOnly: true
        //     };

        //     let result = await this.sutAdapter.sendRequests(pTransArgs);
        //     let customerPoC = JSON.parse(result["result"]);

        //     pTransArgs = {
        //         contractId: "proof",
        //         contractFunction: "GeneratePoCProof2",
        //         contractArguments: [this.campaigns[i].camId, deviceId],
        //         readOnly: true
        //     };
        //     result = await this.sutAdapter.sendRequests(pTransArgs);
        //     let devicePoC = JSON.parse(result["result"]);

        //     // generate number of tpocs
        //     let curNumTransPerCampaigns = numTransPerCampaigns;
        //     if (i == numCampaigns - 1) {
        //         curNumTransPerCampaigns = numTrans - numAddedTrans;
        //     }

        //     numAddedTrans = numAddedTrans + curNumTransPerCampaigns;
        //     let diff = this.campaigns[i].endTimeStr - this.campaigns[i].startTimeStr;
        //     for (let j = 0; j < curNumTransPerCampaigns; j++) {
        //         // let addedTime = campaign.startTime + endTime;
        //         let tpTransArgs = {
        //             contractId: "poc",
        //             contractFunction: "GenerateTPoCProofs",
        //             contractArguments: [this.campaigns[i].camId, customerPoC.comm, customerPoC.r, customerPoC.numVerifiers, 1],
        //             readOnly: true
        //         };

        //         result = await this.sutAdapter.sendRequests(tpTransArgs);
        //         let customerTPoC = JSON.parse(result["result"]).tpocs[0];

        //         tpTransArgs = {
        //             contractId: "poc",
        //             contractFunction: "GenerateTPoCProofs",
        //             contractArguments: [this.campaigns[i].camId, devicePoC.comm, devicePoC.r, devicePoC.numVerifiers, 1],
        //             readOnly: true
        //         };

        //         result = await this.sutAdapter.sendRequests(tpTransArgs);
        //         let deviceTPoC = JSON.parse(result["result"]).tpocs[0];
        //         let addingTime = Math.floor(Math.random() * 10000 % diff + this.campaigns[i].startTimeStr);

        //         let transArgs = {
        //             contractId: "proof",
        //             contractFunction: "AddCampaignTokenTransaction",
        //             contractArguments: [this.campaigns[i].camId, deviceId, addingTime, deviceTPoC.tComms.join(";"), deviceTPoC.tRs.join(";"), deviceTPoC.hashes.join(";"), deviceTPoC.key, customerTPoC.tComms.join(";"), customerTPoC.tRs.join(";"), customerTPoC.hashes.join(";"), customerTPoC.key],
        //             readOnly: false
        //         };

        //         result = await this.sutAdapter.sendRequests(transArgs);
        //         let tokenTran = JSON.parse(result["result"]);
        //         this.tokenTrans.push(tokenTran);
        //     }

        // }
        const args = this.roundArguments;
        // this.contractId = args.contractId;
        // this.contractVersion = args.contractVersion;
        const {numPeersPerOrgs, numOrgsPerType, numVerifiers, numDevices, numTrans, numCampaigns} = this.roundArguments;

        try {
            let path = `../caliper/data/cams-${numOrgsPerType}-${numPeersPerOrgs}-${numVerifiers}.txt`;
            const data = fs.readFileSync(path, 'utf8');
            this.camIds = data.split(",");
            // console.log(data);
        } catch (err) {
            console.error(err.stack);
        }


        // throw Error(JSON.stringify(this.tokenTrans))
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        const {mode} = this.roundArguments;
        const camIdx = Math.floor(Math.random()*10000) % this.campaigns.length;
        const camId = this.camIds[camIdx];

        // camId string, csStr string, rsStr string, hashesStr string, keyStr string
        const transArgs = {
            contractId: "proof",
            contractFunction: "FindTokenTransactionsByCampaignId",
            contractArguments: [camId, mode],
            readOnly: true
        };

        return this.sutAdapter.sendRequests(transArgs);
    }

    async cleanupWorkloadModule() {
        // let transArgs = {
        //     contractId: "campaign",
        //     contractFunction: 'DeleteAllCampaigns',
        //     contractArguments: [],
        //     readOnly: false
        // };

        // await this.sutAdapter.sendRequests(transArgs);

        // // for (let i = 0; i < this.tokenTrans.length; i++) {
        // //     transArgs = {
        // //         contractId: "proof",
        // //         contractFunction: 'DeleteProofById',
        // //         contractArguments: [this.tokenTrans[i]["id"]],
        // //         readOnly: false
        // //     };
        // //     await this.sutAdapter.sendRequests(transArgs);
        // // }

        // transArgs = {
        //     contractId: "proof",
        //     contractFunction: 'DeleteAllProofs',
        //     contractArguments: [],
        //     readOnly: false
        // };
        // await this.sutAdapter.sendRequests(transArgs);

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
