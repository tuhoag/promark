const utils = require('./utils');
const logger = require('./logger')(__filename, "debug");

exports.generateProofForRandomUser = async (camId, userId) => {
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        // randomly generate a user id
        if (userId === undefined) {
            userId = `u${utils.getId(10000)}`;
        }
        logger.info(`GenerateCustomerCampaignProof: camId:${$camId} - userId:${userId}`);
        return contract.submitTransaction("GenerateCustomerCampaignProof", camId, userId);
    }, async response => {
        logger.debug(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof;
    });
}

exports.addProof = async (comm, rsStr) => {
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        // randomly generate a user id
        const proofId = `p${utils.getId(10000)}`;
        logger.info(`AddCustomerProofCampaign: proofId: ${proofId} - comm: ${comm} - rsStr: ${rsStr}`)
        return contract.submitTransaction("AddCustomerProofCampaign", proofId, comm, rsStr);
    }, async response => {
        logger.debug(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof;
    });
}

exports.deleteAllProofs = async () => {
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');
        logger.info("DeleteAllProofs");
        return contract.submitTransaction("DeleteAllProofs");
    }, async response => {
        logger.debug(`response:${response}`);
        return response;
    });
}

exports.getAllProofs = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("proof");
        logger.info("GetAllProofs");
        return contract.submitTransaction("GetAllProofs");
    }, async (response) => {
        if (response.length == 0) {
            logger.debug("No proofs");
            return [];
        }

        const proofs = JSON.parse(response);
        logger.debug(`got ${proofs.length} proofs: ${response}`);
        return proofs;
    });
}