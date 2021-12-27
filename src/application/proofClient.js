const utils = require('./utils');


exports.generateProofForRandomUser = (camId, userId) => {
    utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        console.log('Submit campaign transaction.');
        // randomly generate a user id
        if (userId === undefined) {
            userId = `u${utils.getId(10000)}`;
        }
        // const userId = `u${utils.getId(10000)}`;
        console.log(`userId:${userId}`);
        // userId = "u51"
        return contract.submitTransaction("GenerateCustomerCampaignProof", camId, userId);
    }, async response => {
        console.log(`response:${response}`);
    });

    console.log('Transaction complete.');
}

// module.exports = {
//     // generateProofForRandomUser
// };