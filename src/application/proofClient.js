const utils = require('./utils');
const logger = require('./logger')(__filename, "debug");

exports.generatePoC = async (camId, entityId) => {
    logger.debug(`generatePoC: ${camId},${entityId}`);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        // randomly generate a user id
        if (entityId === undefined) {
            entityId = `u${utils.getId(10000)}`;
        }
        logger.info(`GeneratePoCProof: camId:${camId} - entityId:${entityId}`);
        return contract.submitTransaction("GeneratePoCProof2", camId, entityId);
    }, async response => {
        logger.debug(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof;
    });
}

exports.generateTPoCs = async (camId, cStr, rStr, numVerifiers, numTPoCs) => {
    logger.debug(`generateTPoCs: ${camId},${cStr},${rStr},${numVerifiers},${numTPoCs}`);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('poc');

        logger.info(`GenerateTPoCProofs: camId:${camId} - c:${cStr} - r:${rStr} - numTPoCs:${numTPoCs}`);
        return contract.submitTransaction("GenerateTPoCProofs", camId, cStr, rStr, numVerifiers,numTPoCs);
    }, async response => {
        logger.debug(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof;
    });
}

exports.generatePoCAndTPoCs = async (camId, entityId, numTPoCs) => {
    logger.debug(`generatePocAndTPoCs: ${camId},${entityId}`);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('poc');

        logger.info(`generatePocAndTPoCs: camId:${camId} - entityId:${entityId} - numTPoCs: ${numTPoCs}`);
        return contract.submitTransaction("GeneratePoCAndTPoCProof", camId, entityId, numTPoCs);
    }, async response => {
        logger.debug(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof;
    });
}

exports.verifyPoCProof = async (camId, comm, r) => {
    logger.debug(`verifyPoCProof: ${camId},${comm},${r}`);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');
        logger.info(`VerifyPoCProof: camId:${camId} - comm:${comm} - r: ${r}`);
        return contract.submitTransaction("VerifyPoCProof", camId, comm, r);
    }, async response => {
        logger.debug(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof;
    });
}

exports.verifyTPoCProof = async (camId, commS, rs, hs, key) => {
    logger.debug(`verifyTPoCProof: ${camId},${commS},${rs},${hs},${key}`);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');
        const rsStr = rs.join(";");
        const commStr = commS.join(";")
        const hsStr = hs.join(";")

        logger.info(`verifyTPoCProof: camId:${camId} - commStr:${commStr} - rsStr: ${rsStr} hsStr: ${hsStr} - key: ${key}`);
        return contract.submitTransaction("VerifyTPoCProof", camId, commStr, rsStr, hsStr, key);
    }, async response => {
        logger.debug(`response:${response}`);
        let resultProof = response == "true";
        logger.debug(resultProof)
        return resultProof;
    });
}


exports.addProof = async (camId, deviceId, addedTime, deviceTPoC, customerTPoC) => {
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        logger.info(`AddCampaignTokenTransaction: camid: ${camId} - deviceId: ${deviceId} - cusTPoC: ${JSON.stringify(customerTPoC)} - deviceTPoC: ${JSON.stringify(deviceTPoC)} - addedTimeStr: ${addedTime}`);

        return contract.submitTransaction("AddCampaignTokenTransaction", camId, deviceId, addedTime, deviceTPoC.tComms.join(";"), deviceTPoC.tRs.join(";"), deviceTPoC.hashes.join(";"), deviceTPoC.key, customerTPoC.tComms.join(";"), customerTPoC.tRs.join(";"), customerTPoC.hashes.join(";"), customerTPoC.key);
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