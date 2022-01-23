const utils = require('./utils');


exports.generateProofForRandomUser = (camId, userId) => {
    console.log(camId)
    console.log(userId)
    return utils.callChaincodeFn(async network => {
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
        const resultProof = JSON.parse(response);
        return resultProof
    });
}

exports.addProof = (comm, rsStr) => {
    // console.log(camId);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        console.log('Submit campaign transaction.');
        // randomly generate a user id
        const proofId = `p${utils.getId(10000)}`;
        console.log(`proofId:${proofId}`);
        console.log(`Comm:${comm}`);
        console.log(`rsStr:${rsStr}`);
        // proofId string, comm string, rsStr string
        return contract.submitTransaction("AddCustomerProofCampaign", proofId, comm, rsStr);
    }, async response => {
        console.log(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof
    });
}