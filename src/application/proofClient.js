const utils = require('./utils');
const logger = require('./logger')(__filename, "debug");

exports.generateProofForRandomUser = async (camId, entityId) => {
    logger.debug(`generateProofForRandomUser: ${camId},${entityId}`);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('poc');

        // randomly generate a user id
        if (entityId === undefined) {
            entityId = `u${utils.getId(10000)}`;
        }
        logger.info(`GeneratePoCProof: camId:${camId} - entityId:${entityId}`);
        return contract.submitTransaction("GeneratePoCProof", camId, entityId);
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

exports.addProof2 = async (camId, deviceId, cusId, cusComm, cusRsStr, addedTime) => {
    logger.info(`addProof2: ${camId},${deviceId},${cusId},${cusComm},${cusRsStr},${addedTime}`);
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        // camId string, deviceId string, cusId string, cusComm string, cusRsStr string, addedTimeStr string
        logger.info(`AddCustomerProofCampaign2: camid: ${camId} - deviceId: ${deviceId} - cusId: ${cusId} - cusComm: ${cusComm} - cusRsStr: ${cusRsStr} - addedTimeStr: ${addedTime}`)
        return contract.submitTransaction("AddCustomerProofCampaign2", camId, deviceId, cusId, cusComm, cusRsStr, addedTime);
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

exports.verifyProof = async (camId, proofId) => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("proof");
        logger.info(`VerifyCampaignProof: camId: ${camId} - proofId: ${proofId}`);
        return contract.submitTransaction("VerifyCampaignProof", camId, proofId);
    }, async (response) => {
        logger.debug(`response: ${response}`);
        return response;
    });
}