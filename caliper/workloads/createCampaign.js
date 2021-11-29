'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

// const id = ['id1', 'id2', 'id3', 'id4', 'id5', 'id6', 'id7', 'id8', 'id9', 'id10'];
// const names = ['campaign1', 'campaign2', 'campaign3', 'campaign4', 'campaign5', 'campaign6', 'campaign7', 'campaign8', 'campaign9', 'campaign10'];
const advs = ['adv0',
// 'adv1',
// 'adv2'
];
const buss = ['bus0',
// 'bus1', 'bus2'
];
const vers = ['http://peer0.bus0.promark.com:9000',
             'http://peer0.adv0.promark.com:8500',
            //  'http://peer0.adv1.promark.com:8600',
            //  'http://peer0.adv2.promark.com:8700',
            //  'http://peer0.bus1.promark.com:9100',
            //  'http://peer0.bus2.promark.com:9200'
            ];

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
        const id = 'ID' + this.workerIndex + '_' + this.txIndex.toString();

        let campaignName = 'campaign' + id;
        // let campaignName = names[Math.floor(Math.random() * names.length)];
        let campaignAdv = advs[Math.floor(Math.random() * advs.length)].toString();
        let campaignBus = buss[Math.floor(Math.random() * buss.length)].toString();
        let ver1 = vers.find(a =>a.includes(campaignAdv));
        let ver2 = vers.find(a =>a.includes(campaignBus));
        console.log(`ver1: ${campaignAdv}:${ver1}`);
        console.log(`ver1: ${campaignBus}:${ver2}`);

        let args = {
            contractId: 'campaign',
            contractVersion: 'v1',
            contractFunction: 'CreateCampaign',
            // invokerIdentity: 'peer0.adv0.promark.com',
            contractArguments: [id, campaignName, campaignAdv, campaignBus, ver1, ver2],
            timeout: 30
        };
        console.log(`submitTransaction: ${args}`);
        await this.sutAdapter.sendRequests(args);

    }

    async cleanupWorkloadModule() {
        this.txIndex++;

        const id = 'ID' + this.workerIndex + '_' + this.txIndex.toString();
        const request = {
            contractId: 'campaign',
            contractFunction: 'DeleteCampaignByID',
            // invokerIdentity: 'peer0.adv0.promark.com',
            contractArguments: [id],
            readOnly: false
        };

        console.log(`cleanupWorkloadModule: ${id}`);
        await this.sutAdapter.sendRequests(request);
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
