'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

const peers = ['peer0.adv0.promark.com', 'peer0.pub0.promark.com'];

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        for (let i=0; i<this.roundArguments.test; i++) {
            const campaignID = `CAMPAIGN_${this.workerIndex}_${i}`;
            const peerId = Math.floor(Math.random() * (peers.length - 1));
            console.log(`Worker ${this.workerIndex}: Creating backup ${campaignID} for peer ${peers[peerId]}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'InitLedger',
                invokerIdentity: 'peer0.pub0.promark.com',
                contractArguments: [],
                contractArguments: [campaignID, 'campaign3','Adv0','Pub0'],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
        }
    }

    async submitTransaction() {
        const randomId = Math.floor(Math.random()*this.roundArguments.backups);
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'InitLedger',
            invokerIdentity: 'peer0.pub0.promark.com',
            contractArguments: [`BACKUP_${this.workerIndex}_${randomId}`],
            readOnly: true
        };

        await this.sutAdapter.sendRequests(myArgs);
    }

    async cleanupWorkloadModule() {
        for (let i=0; i<this.roundArguments.test; i++) {
            const campaignID = `CAMPAIGN_${this.workerIndex}_${i}`;
            console.log(`Worker ${this.workerIndex}: Deleting backup ${campaignID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'InitLedger',
                invokerIdentity: 'peer0.pub0.promark.com',
                contractArguments: [campaignID],
                readOnly: false
            };

            console.log(`cleanupWorkloadModule: ${request}`);
            await this.sutAdapter.sendRequests(request);
        }
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;