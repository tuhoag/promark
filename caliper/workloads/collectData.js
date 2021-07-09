'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

// const id = ['id1', 'id2', 'id3', 'id4', 'id5', 'id6', 'id7', 'id8', 'id9', 'id10'];
// const names = ['campaign1', 'campaign2', 'campaign3', 'campaign4', 'campaign5', 'campaign6', 'campaign7', 'campaign8', 'campaign9', 'campaign10'];
const totalComm = 'dsUDvULKSaMk6/eaWBmThy7vqd4HSszBlv2mA+MDr1s=';
const r1 = 'sLc2AFBxAOEhGqpiOSnVJmEX8/fKPM//62XykBc26wM=';
const r2 = 'Zd7/Jz2CguWf+sR6dDwhyqFwiS1TsnCRUj+AKnIDtAY=';

/**
 * Workload module for the benchmark round.
 */
class CreateCarWorkload extends WorkloadModuleBase {
    /**
     * Initializes the workload module instance.
     */
    constructor() {
        super();
        this.txIndex = 0;
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        this.txIndex++;
        let id = 'id5'
        let index = this.workerIndex + '_' + this.txIndex.toString();
        let userName = 'username' + index;
        // let campaignName = names[Math.floor(Math.random() * names.length)];
        // let campaignAdv = advs[Math.floor(Math.random() * advs.length)].toString();
        // let campaignBus = buss[Math.floor(Math.random() * buss.length)].toString();
        // let carOwner = owners[Math.floor(Math.random() * owners.length)];

        let args = {
            contractId: 'campaign',
            contractVersion: 'v1',
            contractFunction: 'AddCollectedData',
            contractArguments: [id, userName, totalComm, r1, r2],
            timeout: 30
        };

        await this.sutAdapter.sendRequests(args);
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new CreateCarWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
