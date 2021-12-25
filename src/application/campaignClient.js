const utils = require('./utils');
const setting = require('./setting');

const getId = (maxNum) => {
    return Math.floor(Math.random() * 100) % maxNum;
}

const findAdvertiser = (numAdvs) => {
    return "adv" + getId(numAdvs);
}

const findBusiness = (numBuses) => {
    return "bus" + getId(numBuses);
}

const findVerifierURLs = (orgName, numPeers, numVerifiers) => {
    // select peer
    // peer0.adv0.promark.com:5000
    const peerName = "peer" + getId(numPeers);

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
    const camId = "c" + getId(10000);
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

const createRandomCampaign = async (numOrgsPerType, numPeersPerOrg, numVerifiers) => {
    utils.callChaincodeFn(async network => {
        const contract = await network.getContract('campaign');

        console.log('Submit campaign transaction.');
        const campaign = generateCampaignArgs(numOrgsPerType, numPeersPerOrg, numVerifiers);
        const verifierAddressesStr = campaign.verifierURLs.join(";")
        return contract.submitTransaction("CreateCampaign", campaign.id, campaign.name, campaign.advertiser, campaign.business, verifierAddressesStr);
    }, async response => {
        const resultCampaign = JSON.parse(response);
        console.log(`result campaign.Id: ${resultCampaign.id} - Name: ${resultCampaign.name} - Adv: ${resultCampaign.advertiser} - Bus: ${resultCampaign.business} - Verifiers: ${resultCampaign.verifierURLs}`);
    });

    console.log('Transaction complete.');
}

const callGetCampaignById = async (network, camId) => {
    // Get addressability to commercial paper contract
    console.log("camId: " + camId)

    console.log('Use campaign.promark smart contract.');
    const contract = await network.getContract('campaign');

    // issue commercial paper
    console.log('Submit campaign issue transaction.');
    const response = await contract.submitTransaction("GetCampaignById", camId);

    // process response
    console.log("response: " + response);

    return response
}

const getCampaignById = async (camId) => {
    utils.callChaincodeFn(async network => {
        const contract = await network.getContract('campaign');

        console.log('Submit campaign transaction.');
        return contract.submitTransaction("GetCampaignById", camId);
    }, async response => {
        const resultCampaign = JSON.parse(response);
        console.log(`result campaign.Id: ${resultCampaign.id} - Name: ${resultCampaign.name} - Adv: ${resultCampaign.advertiser} - Bus: ${resultCampaign.business} - Verifiers: ${resultCampaign.verifierURLs}`);
    });

    console.log('Transaction complete.');
}


const getAllCampaigns = async () => {
    utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        console.log('Submit transaction.');
        return contract.submitTransaction("GetAllCampaigns");
    }, async (response) => {
        console.log("response: " + response);
    });
}


const deleteCampaignById = async (camId) => {
    utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        console.log('Submit transaction.');
        return contract.submitTransaction("DeleteCampaignById", camId);
    }, async (response) => {
        console.log("response: " + response);
    });
}


module.exports = {
    createRandomCampaign,
    getCampaignById,
    getAllCampaigns,
    deleteCampaignById
}