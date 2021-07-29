'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

// const id = ['id1', 'id2', 'id3', 'id4', 'id5', 'id6', 'id7', 'id8', 'id9', 'id10'];
// const names = ['campaign1', 'campaign2', 'campaign3', 'campaign4', 'campaign5', 'campaign6', 'campaign7', 'campaign8', 'campaign9', 'campaign10'];
const totalComm = 'qJ0LoddreR+pveZ0mA4buZ5Bh9gQqxUfaxJBcBVdhxs=';
const r1 = '4pJEX5dj+1UamPv9JAigzNdeokf1tUnCqsiMg0QQSQk=';
const r2 = 'R3b/jETOO7J7/IKy0jwXj9s28QyvH4bVJk4ogbLLGAI=';
const ver1 = 'http://peer0.bus0.promark.com:9000';
const ver2 = 'http://peer0.adv0.promark.com:8500';
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
        const id = 'id10'
        // for (let i=0; i<this.roundArguments.test; i++) {
            // const randomId = Math.floor(Math.random()*this.roundArguments.testRound)
        let index = this.workerIndex + '_' + this.txIndex.toString();
        // const userName = `username_${this.workerIndex}_${i}`;
        let userName = 'username' + this.roundArguments.testRound + '_' + index;
        // let campaignName = names[Math.floor(Math.random() * names.length)];
        // let campaignAdv = advs[Math.floor(Math.random() * advs.length)].toString();
        // let campaignBus = buss[Math.floor(Math.random() * buss.length)].toString();

        let args = {
            contractId: 'campaign',
            contractVersion: 'v1',
            contractFunction: 'AddCollectedData',
            contractArguments: [id, userName, totalComm, r1, r2, ver1, ver2],
            readOnly: false,
            timeout: 30
        };

        await this.sutAdapter.sendRequests(args);
        // }
    }

    async cleanupWorkloadModule() {
        const id = 'id10'

        // for (let i=0; i<this.roundArguments.test; i++) {
            // const randomId = Math.floor(Math.random()*this.roundArguments.testRound)
        let index = this.workerIndex + '_' + this.txIndex.toString();
        let userName = 'username' + this.roundArguments.testRound + '_' + index;
        // const userName = `username_${this.workerIndex}_${i}`;

        const request = {
            contractId: 'campaign',
            contractVersion: 'v1',
            contractFunction: 'DeleteDataByUserId',
            // invokerIdentity: 'peer0.bus0.promark.com',
            contractArguments: [userName],
            readOnly: false
        };

        console.log(`cleanupWorkloadModule: ${userName}`);
        await this.sutAdapter.sendRequests(request);
        // }
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
