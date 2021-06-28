'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

const peers = ['peer0.adv0.promark.com', 'peer0.bus0.promark.com'];

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }
    
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        for (let i=0; i<this.roundArguments.backups; i++) {
            const backupID = `BACKUP_${this.workerIndex}_${i}`;
            const peerId = Math.floor(Math.random() * (peers.length - 1));
            console.log(`Worker ${this.workerIndex}: Creating backup ${backupID} for peer ${peers[peerId]}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'InitLedger',
                invokerIdentity: 'peer0.bus0.promark.com',
                contractArguments: [],
                // contractArguments: ['id8', 'campaign3','Adv0','Bus0'],
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
            invokerIdentity: 'peer0.bus0.promark.com',
            contractArguments: [`BACKUP_${this.workerIndex}_${randomId}`],
            readOnly: true
        };

        await this.sutAdapter.sendRequests(myArgs);
    }
    
    async cleanupWorkloadModule() {
        for (let i=0; i<this.roundArguments.backups; i++) {
            const backupID = `BACKUP_${this.workerIndex}_${i}`;
            console.log(`Worker ${this.workerIndex}: Deleting backup ${backupID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'InitLedger',
                invokerIdentity: 'peer0.bus0.promark.com',
                contractArguments: [backupID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
        }
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;