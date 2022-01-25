const utils = require('./utils');
const logger = require('./logger')(__filename, "debug");


const findAdvertiser = (numAdvs) => {
    return "adv" + utils.getId(numAdvs);
}

const findBusiness = (numBuses) => {
    return "bus" + utils.getId(numBuses);
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

const generateCampaignArgs = (numOrgsPerType, numPeersPerOrg, numVerifiers) => {
    const camId = "c" + utils.getId(10000);
    const advName = findAdvertiser(numOrgsPerType);
    const busName = findBusiness(numOrgsPerType);
    const verifierURLs = findVerifierURLs(advName, numPeersPerOrg, numVerifiers);
    verifierURLs.push(...findVerifierURLs(busName, numPeersPerOrg, numVerifiers));

    return {
        id: camId,
        name: `Campaign ${camId}`,
        advertiser: advName,
        business: busName,
        verifierURLs,
    }
}

const createRandomCampaign = async (numVerifiers) => {
    return utils.callChaincodeFn(async network => {

        const contract = await network.getContract('campaign');
        const campaign = generateCampaignArgs(global.numOrgsPerType, global.numPeersPerOrg, numVerifiers);
        const verifierAddressesStr = campaign.verifierURLs.join(";")
        logger.info(`Create Campaign: ${JSON.stringify(campaign)} and ${verifierAddressesStr}`);

        return contract.submitTransaction("CreateCampaign", campaign.id, campaign.name, campaign.advertiser, campaign.business, verifierAddressesStr);
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


module.exports = {
    createRandomCampaign,
    getCampaignById,
    getAllCampaigns,
    deleteCampaignById,
    deleteAllCampaigns,
}