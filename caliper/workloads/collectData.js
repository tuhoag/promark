'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

// const ids = ['id10', 'id11', 'id12', 'id13'];
const ids = ['id10', 'id11', 'id12', 'id13'];
// const names = ['campaign1', 'campaign2', 'campaign3', 'campaign4', 'campaign5', 'campaign6', 'campaign7', 'campaign8', 'campaign9', 'campaign10'];
const totalComms = ['iiIhxLouay6sTS2uhDpr9hTTPVg2I8k+AtlG869m+FI=',
                    '8m91gDWPbfv0sSOHuM5aq6ZcgyQAwz2g7xmT//pODHs=',
                    '9ngd/x4NmukQh61x8wNFGw82xWcjy9x7RSP9QMkXABc=',
                    'dGVL9Lu2yPiePmUTjqmIh3LOhVvRHswrR5H+QN/gVBY=',
                    ];
const r1s = ['v7g8Lm0VjieFg+7fBULDIMOu9gBrvhY54a+gU2eLnwc=',
            'plB/YF932SjCOP37BN/8lGWYMUjthyrzB4g+tbU72Qk=',
            'ppGBvNBUpN9MgnlaS5kvtDFh8DdiU0Kk+g0QIr/3/AI=',
            'JwLzkwM1BC0U+mcTn98U845aPlCGrwwoQmIO3YLpbgM=',
            ];
const r2s = ['VEgJObUFPo88iwqzfwJN3HD1dJthqf007DPe35JusgA=',
            'qjJPpPxfRxu0zyBbbhNf0zrfg4GCex88xdVgOonc1wc=',
            'M1tL2lfgu6Wu9yiqo3w1IwRDvohLgu8q+vpDoB/bkQk=',
            '3eE5k9hCDc6k8hlPDoIGJfplMahFeMzqo+MiYp9GBAM=',
            ];

const ver1s = ['http://peer0.adv0.promark.com:8500',
                'http://peer0.adv1.promark.com:8510',
                // 'http://peer0.adv0.promark.com:8500',
                // 'http://peer0.adv0.promark.com:8500',
                'http://peer0.adv2.promark.com:8520',
                'http://peer0.adv3.promark.com:8530',
              ];
const ver2s = ['http://peer0.bus0.promark.com:9000',
                'http://peer0.bus1.promark.com:9010',
                // 'http://peer0.bus0.promark.com:9000',
                // 'http://peer0.bus0.promark.com:9000',
                'http://peer0.bus2.promark.com:9020',
                'http://peer0.bus3.promark.com:9030',               
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
        let id = ids[Math.floor(Math.random() * ids.length)].toString();
        let index = this.workerIndex + '_' + this.txIndex.toString();

        let userName = 'user' + id.toString() + this.roundArguments.testRound + '_' + index;
        let totalComm = totalComms[ids.indexOf(id)].toString();
        let r1 = r1s[ids.indexOf(id)].toString();
        let r2 = r2s[ids.indexOf(id)].toString();
        let ver1 = ver1s[ids.indexOf(id)].toString();
        let ver2 = ver2s[ids.indexOf(id)].toString();

        let args = {
            contractId: 'campaign',
            contractVersion: 'v1',
            contractFunction: 'AddCollectedData',
            // invokerIdentity: 'peer0.bus1.promark.com',
            contractArguments: [id, userName, totalComm, r1, r2, ver1, ver2],
            readOnly: false,
            timeout: 30
        };

        await this.sutAdapter.sendRequests(args);
        // }
    }

    // async cleanupWorkloadModule() {
    //     // this.txIndex++;
    //     // let id = ids[Math.floor(Math.random() * ids.length)].toString();
    //     let index = this.workerIndex + '_' + this.txIndex.toString();
    //     let userName = 'user' + this.roundArguments.testRound + '_' + index;

    //     const request = {
    //         contractId: 'campaign',
    //         contractVersion: 'v1',
    //         contractFunction: 'DeleteDataByUserId',
    //         invokerIdentity: 'peer0.bus1.promark.com',
    //         contractArguments: [userName],
    //         readOnly: false
    //     };

    //     // console.log(`cleanupWorkloadModule: ${userName}`);
    //     await this.sutAdapter.sendRequests(request);
    //     // }
    // }

}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new CreateCarWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
