const utils = require('./utils');
const logger = require('./logger')(__filename, "debug");


const findAdvertiser = (numAdvs) => {
    return "adv" + utils.getId(numAdvs);
}

const findPublisher = (numPubs) => {
    return "pub" + utils.getId(numPubs);
}

const findVerifierURLs = (orgName, numPeers, numVerifiers) => {
    // select peer
    // peer0.adv0.promark.com:5000
    const peerName = "peer" + utils.getId(numPeers);

    if (numVerifiers > numPeers) {
        throw `numVerifiers (${numVerifiers}$) is higher than numPeers (${numPeers}).`;
    }

    let verifierAddresses = [];

    for (let i = 0; i < numVerifiers; i++) {
        const address = `${peerName}.${orgName}.promark.com:5000`;
        verifierAddresses.push(address);
    }

    return verifierAddresses
}

const generateCampaignArgs = (numOrgsPerType, numPeersPerOrg, numVerifiers, deviceIdsStr) => {
    const camId = "c" + utils.getId(10000);
    const advName = findAdvertiser(numOrgsPerType);
    const pubName = findPublisher(numOrgsPerType);
    const verifierURLs = findVerifierURLs(advName, numPeersPerOrg, numVerifiers);
    verifierURLs.push(...findVerifierURLs(pubName, numPeersPerOrg, numVerifiers));
    const deviceIds = deviceIdsStr.split(",");
    const startTimeStr = Math.floor(new Date("2022-05-01").getTime() / 1000);
    const endTimeStr = Math.floor(new Date("2022-06-01").getTime() / 1000);

    const result = {
        id: camId,
        name: `Campaign ${camId}`,
        advertiser: advName,
        publisher: pubName,
        startTimeStr,
        endTimeStr,
        verifierURLs,
        deviceIds,
    }

    logger.debug(result)

    return result
}

const createRandomCampaign = async (numVerifiers, deviceIdsStr) => {
    return utils.callChaincodeFn(async network => {

        const contract = await network.getContract('campaign');
        const campaign = generateCampaignArgs(global.numOrgsPerType, global.numPeersPerOrg, numVerifiers, deviceIdsStr);
        const verifierAddressesStr = campaign.verifierURLs.join(";")

        logger.info(`Create Campaign: ${JSON.stringify(campaign)} and ${verifierAddressesStr} - devices: ${campaign.deviceIds.join(";")}`);

        return contract.submitTransaction("CreateCampaign", campaign.id, campaign.name, campaign.advertiser, campaign.publisher, campaign.startTimeStr, campaign.endTimeStr, verifierAddressesStr, campaign.deviceIds.join(";"));
    }, async response => {
        const resultCampaign = JSON.parse(response);
        logger.debug(`raw response: ${response}`);
        return resultCampaign;
    });
}


const getCampaignById = async (camId) => {
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('campaign');
        logger.info(`Get Campaign: ${camId}`);
        return contract.submitTransaction("GetCampaignById", camId);
    }, async response => {
        const resultCampaign = JSON.parse(response);
        logger.debug(`raw response: ${response}`);
        return resultCampaign;
    });
}


const getAllCampaigns = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        logger.info("GetAllCampaigns");
        return contract.submitTransaction("GetAllCampaigns");
    }, async (response) => {
        if (response.length == 0) {
            logger.info("No campaigns");
            return [];
        }

        const campaigns = JSON.parse(response);
        logger.debug(`got ${campaigns.length} campaigns`);
        logger.debug(`campaigns: ${response}`);
        return campaigns;
    });
}


const deleteCampaignById = async (camId) => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        logger.info("DeleteCampaignById");
        return contract.submitTransaction("DeleteCampaignById", camId);
    }, async (response) => {
        logger.debug(`response: ${response}`);
        return response;
    });
}

const deleteAllCampaigns = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        logger.info("DeleteAllCampaigns");
        return contract.submitTransaction("DeleteAllCampaigns");
    }, async (response) => {
        logger.debug(`response: ${response}`);
        return response;
    });
}

const getChaincodeData = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        logger.info("GetChaincodeData");
        return contract.submitTransaction("GetChaincodeData");
    }, async (response) => {
        logger.debug(`response: ${response}`);
        return response;
    });
}

module.exports = {
    createRandomCampaign,
    getCampaignById,
    getAllCampaigns,
    deleteCampaignById,
    deleteAllCampaigns,
    getChaincodeData,
}