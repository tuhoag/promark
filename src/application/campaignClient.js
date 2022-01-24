const utils = require('./utils');
// const setting = require('./setting');


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
        console.log('Submit campaign transaction.');
        const campaign = generateCampaignArgs(global.numOrgsPerType, global.numPeersPerOrg, numVerifiers);
        const verifierAddressesStr = campaign.verifierURLs.join(";")
        return contract.submitTransaction("CreateCampaign", campaign.id, campaign.name, campaign.advertiser, campaign.business, verifierAddressesStr);
    }, async response => {
        const resultCampaign = JSON.parse(response);
        console.log(`result campaign.Id: ${resultCampaign.id} - Name: ${resultCampaign.name} - Adv: ${resultCampaign.advertiser} - Bus: ${resultCampaign.business} - Verifiers: ${resultCampaign.verifierURLs}`);
        return resultCampaign;
        // return resultCampaign;
    });
}


const getCampaignById = async (camId) => {
    return utils.callChaincodeFn(async network => {
        const contract = await network.getContract('campaign');

        console.log('Submit campaign transaction.');
        return contract.submitTransaction("GetCampaignById", camId);
    }, async response => {
        const resultCampaign = JSON.parse(response);
        console.log(`result campaign.Id: ${resultCampaign.id} - Name: ${resultCampaign.name} - Adv: ${resultCampaign.advertiser} - Bus: ${resultCampaign.business} - Verifiers: ${resultCampaign.verifierURLs}`);
    });
}


const getAllCampaigns = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        console.log('Submit transaction.');
        return contract.submitTransaction("GetAllCampaigns");
    }, async (response) => {
        if (response.length == 0) {
            console.log("No campaigns");
            return [];
        }

        const campaigns = JSON.parse(response);
        console.log(`got ${campaigns.length} campaigns`);
        console.log(`campaigns: ${response}`);
        return campaigns;
    });
}


const deleteCampaignById = async (camId) => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        console.log('Submit transaction.');
        return contract.submitTransaction("DeleteCampaignById", camId);
    }, async (response) => {
        console.log("response: " + response);
    });
}

const deleteAllCampaigns = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        console.log('Submit transaction.');
        return contract.submitTransaction("DeleteAllCampaigns");
    }, async (response) => {
        console.log("response: " + response);
    });
}


module.exports = {
    createRandomCampaign,
    getCampaignById,
    getAllCampaigns,
    deleteCampaignById,
    deleteAllCampaigns,
}