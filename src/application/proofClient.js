const utils = require('./utils');


exports.generateProofForRandomUser = async (camId, userId) => {
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
        return resultProof;
    });
}

exports.addProof = async (comm, rsStr) => {
    // console.log(camId);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        console.log('Submit proof transaction.');
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
        return resultProof;
    });
}

exports.deleteAllProofs = async () => {
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');
        return contract.submitTransaction("DeleteAllProofs");
    }, async response => {
        console.log(`response:${response}`);
    });
}

exports.getAllProofs = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("proof");
        return contract.submitTransaction("GetAllProofs");
    }, async (response) => {
        if (response.length == 0) {
            console.log("No proofs");
            return [];
        }

        const proofs = JSON.parse(response);
        console.log(`got ${proofs.length} proofs`);
        console.log(`proofs: ${response}`);
        return proofs;
    });
}