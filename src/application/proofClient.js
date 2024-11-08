const utils = require('./utils');
const { createLogger } = require('./logger');
logger = createLogger(__filename);

exports.generatePoC = async (camId, entityId) => {
    logger.debug(`generatePoC: ${camId},${entityId}`);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');

        // randomly generate a user id
        if (entityId === undefined) {
            entityId = `u${utils.getId(10000)}`;
        }
        logger.debug(`GeneratePoCProof: camId:${camId} - entityId:${entityId}`);
        return contract.submitTransaction("GeneratePoCProof2", camId, entityId);
    }, async response => {
        logger.debug(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof;
    });
}

exports.generateTPoCs = async (camId, cStr, rStr, numVerifiersPerOrg, numTPoCs) => {
    logger.debug(`generateTPoCs: ${camId},${cStr},${rStr},${numVerifiersPerOrg},${numTPoCs}`);

    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('poc');

        logger.debug(`GenerateTPoCProofs: camId:${camId} - c:${cStr} - r:${rStr} - numTPoCs:${numTPoCs}`);
        return contract.submitTransaction("GenerateTPoCProofs", camId, cStr, rStr, numVerifiersPerOrg,numTPoCs);
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

        logger.debug(`generatePocAndTPoCs: camId:${camId} - entityId:${entityId} - numTPoCs: ${numTPoCs}`);
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
        logger.debug(`VerifyPoCProof: camId:${camId} - comm:${comm} - r: ${r}`);
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

        logger.debug(`verifyTPoCProof: camId:${camId} - commStr:${commStr} - rsStr: ${rsStr} hsStr: ${hsStr} - key: ${key}`);
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

        logger.debug(`AddCampaignTokenTransaction: camid: ${camId} - deviceId: ${deviceId} - cusTPoC: ${JSON.stringify(customerTPoC)} - deviceTPoC: ${JSON.stringify(deviceTPoC)} - addedTimeStr: ${addedTime}`);

        return contract.submitTransaction("AddCampaignTokenTransaction", camId, deviceId, addedTime, deviceTPoC.tComms.join(";"), deviceTPoC.tRs.join(";"), deviceTPoC.hashes.join(";"), deviceTPoC.key, customerTPoC.tComms.join(";"), customerTPoC.tRs.join(";"), customerTPoC.hashes.join(";"), customerTPoC.key);
    }, async response => {
        logger.debug(`response:${response}`);
        const resultProof = JSON.parse(response);
        return resultProof;
    });

    // c29,0,0,1D3R3+IW+Y7yQDnPXI0iPm9T7yn9QN2XKa1xB70ueRM=,4S2tcsDuFIW8OVMQwMactH3l/CqOyVtdNoVK4SA7VQ8=,,,yl/lbEIkMRZKzkO6XYMftxHMzuUpsTON/jgsC3hbpCA=,V5Xzqld9AHeUF1/t640Gp7La1p0TRYu9pMNfCWBLZAc=,,
}

exports.deleteAllProofs = async () => {
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('proof');
        logger.debug("DeleteAllProofs");
        return contract.submitTransaction("DeleteAllProofs");
    }, async response => {
        logger.debug(`response:${response}`);
        return response;
    });
}

exports.getAllProofs = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("proof");
        logger.debug("GetAllProofs");
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

exports.getTokenTransactionsByCampaignId = async (camId, mode) => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("proof");
        logger.debug("getTokenTransactionsByCampaignId");
        return contract.submitTransaction("FindTokenTransactionsByCampaignId", camId, mode);
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

exports.simulateGetTokenTransactionsByCampaignId = async (camId, mode, limit) => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("poc");
        logger.debug(`simulateGetTokenTransactionsByCampaignId: ${camId} - ${mode} - ${limit}`);
        return contract.submitTransaction("SimulateFindTokenTransactionsByCampaignId", camId, mode, limit);
    }, async (response) => {
        if (response.length == 0) {
            logger.debug("No proofs");
            return [];
        }

        const proofs = JSON.parse(response);
        logger.debug(`got: ${response}`);
        return proofs;
    });
}


exports.getTokenTransactionsByTimestamps = async (startTime, endTime) => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("proof");
        logger.debug("getTokenTransactionsByTimestamps");
        return contract.submitTransaction("FindTokenTransactionsByTimestamps", startTime, endTime);
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
        logger.debug(`VerifyCampaignProof: camId: ${camId} - proofId: ${proofId}`);
        return contract.submitTransaction("VerifyCampaignProof", camId, proofId);
    }, async (response) => {
        logger.debug(`response: ${response}`);
        return response;
    });
}

exports.queryByTimestamps = async (startTime, endTime) => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("proof");
        logger.debug("queryByTimestamps");
        return contract.submitTransaction("FindTokenTransactionsByTimestamps", startTime, endTime);
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

exports.getAllProofIds = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("proof");
        logger.debug("GetAllProofIds");
        return contract.submitTransaction("GetAllProofIds");
    }, async (response) => {
        if (response.length == 0) {
            logger.debug("No proofs");
            return [];
        }

        const proofs = JSON.parse(response);
        logger.debug(`got ${proofs.length} proofIds: ${response}`);
        return proofs;
    });
}

exports.addTokenTransactions = async (campaign, deviceId, devicePocAndTPoCs, customerPocAndTPoCs, numAdditions) => {
    let numTPoCs = customerPocAndTPoCs.tpocs.length;

    if (numAdditions > numTPoCs) {
        throw `numAdditions > numTPoCs: ${numAdditions} > ${numTPoCs}`;
    }

    const diff = campaign.endTime - campaign.startTime;
    for (let i = 0; i  < numAdditions; i++) {
        logger.debug(JSON.stringify(customerPocAndTPoCs));
        logger.debug(JSON.stringify(customerPocAndTPoCs.tpocs[i]));
        let customerTPoC = customerPocAndTPoCs.tpocs[i];
        let deviceTPoC = devicePocAndTPoCs.tpocs[i];
        let addingTime = Math.floor(Math.random() * diff + campaign.startTime);
        // let addedTime = campaign.startTime + relativeAddingTime;

        logger.debug(`start: ${campaign.startTime} - end: ${campaign.endTime} - diff: ${diff} - adding time: ${addingTime} - valid: ${(campaign.startTime < addingTime) && (addingTime < campaign.endTime)}`);

        let transaction = await exports.addProof(campaign.id, deviceId, addingTime, deviceTPoC, customerTPoC);

        logger.debug(`added transaction: ${JSON.stringify(transaction)}`);

        // await sleep(1000);
    }

    return numAdditions;
}