'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

const ids = ['id10', 'id11', 'id12', 'id13'];
// const names = ['campaign1', 'campaign2', 'campaign3', 'campaign4', 'campaign5', 'campaign6', 'campaign7', 'campaign8', 'campaign9', 'campaign10'];
const totalComms = ['v3jlBUHAExMNsEe9Q/WYZNKjKoFQDqMbsinsOyi32s=',
                    'oBNJqHMQWP3o7m/ltchj/Vq4emeY3H5OLgryHDVd/Rw=',
                    'FClP6BQxXhxyymrXQDbwJXwki+xl/PukIvWOR31aJ04=',
                    '6oZXEY2la7BH0oLQM4Q/KHSAntknB4cO77ZuwhFm52k=',
                    ];
const r1s = ['MJopGWsp5MSVnHen2bp3hcAxNNFUi48Ra2525hTikg4=',
            'WvdLSYK+HnD+hm/ozrxzFwDPGqWOzJrURzMpy8GIDAs=',
            'V/7IGFlDV+nuBA0ACkgWB07QElhldA8jdxJHsWCLdQo=',
            'PMegpMnPZ/FJjmLr7Pzy8DU/sA1UnBz8tJc/WuBStwo=',
            ];
const r2s = ['Xjpa8hlGtDRaQN1RyN0MTdToK2ZqJq+dzXMx4sPnxAM=',
            '6rcouDaC3naJBoRQ08adCTIydO7TzChfF8ax0Om1/AQ=',
            'MOZxS9LxYJ9TOyHReLgAjfPKAnKBmkkmwGhAoT9PEgo=',
            'nsjRqUikxqvf6I+ZpK5InHYfefqMmZiDFK+aqm4VNAg=',
            ];

const ver1s = ['http://peer0.adv0.promark.com:8500',
                'http://peer0.adv1.promark.com:8510',
                'http://peer0.adv2.promark.com:8520',
                'http://peer0.adv3.promark.com:8530',
              ];
const ver2s = ['http://peer0.bus0.promark.com:9000',
                'http://peer0.bus1.promark.com:9010',
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

        let userName = 'user' + this.roundArguments.testRound + '_' + index;
        let totalComm = totalComms[ids.indexOf(id)].toString();
        let r1 = r1s[ids.indexOf(id)].toString();
        let r2 = r2s[ids.indexOf(id)].toString();
        let ver1 = ver1s[ids.indexOf(id)].toString();
        let ver2 = ver2s[ids.indexOf(id)].toString();

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
