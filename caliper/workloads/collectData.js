'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

// const id = ['id1', 'id2', 'id3', 'id4', 'id5', 'id6', 'id7', 'id8', 'id9', 'id10'];
// const names = ['campaign1', 'campaign2', 'campaign3', 'campaign4', 'campaign5', 'campaign6', 'campaign7', 'campaign8', 'campaign9', 'campaign10'];
const totalComm = 'FqGivC5kKeJDcshwer6Mjru6JbM3yRNxL+FKdvVng34=';
const r1 = 'tuHEjEtLQ6hY/JMhsnUEOB/sthSe3mjHofNdzpFR/AA=';
const r2 = 'kB/AV8tfG3D771odJ/NC5UEkAr7swRJrglc6tAoEtwk=';
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
        let id = 'id4'
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
            contractArguments: [id, userName, totalComm, r1, r2, ver1, ver2],
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
