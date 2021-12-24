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

const callCreateCampaign = async (network, camId, advName, busName, verifierAddresses) => {
    // Get addressability to commercial paper contract
    console.log("camId: " + camId)
    console.log("advName: " + advName)
    console.log("busName: " + busName)
    console.log("verifierAddresses: " + verifierAddresses)

    console.log('Use campaign.promark smart contract.');
    const contract = await network.getContract('campaign');

    // issue commercial paper
    console.log('Submit campaign issue transaction.');
    const verifierAddressesStr = verifierAddresses.join(";")
    const response = await contract.submitTransaction("CreateCampaign", camId, camId, advName, busName, verifierAddressesStr);

    // process response
    console.log("response: " + response);

    return response
}

const generateCampaignArgs = (numOrgsPerType, numPeersPerOrg, numVerifiers) => {
    const camId = "c" + getId(10000);
    const advName = findAdvertiser(numOrgsPerType);
    const busName = findBusiness(numOrgsPerType);
    const verifierURLs = findVerifierURLs(advName, numPeersPerOrg, numVerifiers);
    verifierURLs.push(...findVerifierURLs(busName, numPeersPerOrg, numVerifiers));

    return {
        id: camId,
        advertiser: advName,
        business: busName,
        verifierURLs,
    }
}

const createRandomCampaign = async (numOrgsPerType, numPeersPerOrg, numVerifiers) => {
    const campaign = generateCampaignArgs(numOrgsPerType, numPeersPerOrg, numVerifiers);

    let gateway;
    try {
        gateway = await utils.connectToGateway(setting.userName, setting.orgName, setting.orgUserName);
        const network = await gateway.getNetwork(setting.channelName);

        const response = await callCreateCampaign(network, campaign.id, campaign.advertiser, campaign.business, campaign.verifierURLs);
        console.log("response: " + response);

        const resultCampaign = JSON.parse(response);
        console.log(`result campaign.Id: ${resultCampaign.id} - Name: ${resultCampaign.name} - Adv: ${resultCampaign.advertiser} - Bus: ${resultCampaign.business} - Verifiers: ${resultCampaign.verifierURLs}`)
        console.log('Transaction complete.');
    } catch (error) {
        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);
    } finally {
        // Disconnect from the gateway
        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();
    }
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
    let gateway;
    try {
        gateway = await utils.connectToGateway(setting.userName, setting.orgName, setting.orgUserName);
        const network = await gateway.getNetwork(setting.channelName);

        const response = await callGetCampaignById(network, camId);
        console.log("response: " + response);

        const resultCampaign = JSON.parse(response);
        console.log(`result campaign.Id: ${resultCampaign.id} - Name: ${resultCampaign.name} - Adv: ${resultCampaign.advertiser} - Bus: ${resultCampaign.business} - Verifiers: ${resultCampaign.verifierURLs}`)
        console.log('Transaction complete.');
    } catch (error) {
        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);
    } finally {
        // Disconnect from the gateway
        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();
    }
}

class ChaincodeCaller {
    constructor() {
        this.gateway = null;
        this.network = null;
    }

    static async init() {
        let caller = new ChaincodeCaller();
        caller.gateway = await utils.connectToGateway(setting.userName, setting.orgName, setting.orgUserName);
        caller.network = await caller.gateway.getNetwork(setting.channelName);

        return caller
    }

    async call(fn) {
        return fn(this.network);
    }

    async disconnect() {
        this.gateway.disconnect();
    }
}

const callChaincodeFn = async (requestFn, responseFn) => {
    let gateway;
    try {
        gateway = await utils.connectToGateway(setting.userName, setting.orgName, setting.orgUserName);
        const network = await gateway.getNetwork(setting.channelName);
        console.log('Use campaign.promark smart contract.');
        const response = await requestFn(network);
        await responseFn(response);
    } catch (error) {
        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);
    } finally {
        // Disconnect from the gateway
        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();
    }
}

const getAllCampaigns = async () => {
    caller = await ChaincodeCaller.init();
    caller.call(async (network) => {
        const contract = await network.getContract("campaign");
        console.log('Submit transaction.');
        return contract.submitTransaction("GetAllCampaigns");
    }).then(response => {
        console.log("response: " + response);
    }).finally(() => caller.disconnect());
}

const callGetAllCampaigns = async (network) => {
    console.log('Use campaign.promark smart contract.');
    const contract = await network.getContract('campaign');

    // issue commercial paper
    console.log('Submit transaction.');
    const response = await contract.submitTransaction("GetAllCampaigns");

    // process response
    console.log("response: " + response);

    return response
}

module.exports = {
    // findAdvertiser,
    // findBusiness,
    // findVerifierAddresses,
    // getId,
    // generateCampaignArgs,
    // createCampaign,
    createRandomCampaign,
    getCampaignById,
    getAllCampaigns,
}