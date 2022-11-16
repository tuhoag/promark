const { start } = require('repl');
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

const CreateCampaignUnequalVerifiersArgs = (numOrgsPerType, numPeersPerOrgs, numVerifiers, numDevices, startDate, endDate) => {
    const camId = "c" + Math.floor(Math.random()*10000);
    const name = "Campaign " + camId;
    const advName = "adv"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const pubName = "pub"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const startTimeStr = Math.floor(new Date(startDate).getTime() / 1000);
    const endTimeStr = Math.floor(new Date(endDate).getTime() / 1000);

    var allVerifiersUrls = [];

    logger.debug(`orgs: ${numOrgsPerType} - peers: ${numPeersPerOrgs}`)

    for (let peerId = 0; peerId < numPeersPerOrgs; peerId++) {
        const advPeerURL = `peer${peerId}.${advName}.promark.com:5000`;
        const pubPeerURL = `peer${peerId}.${pubName}.promark.com:5000`;

        allVerifiersUrls.push(advPeerURL);
        allVerifiersUrls.push(pubPeerURL);
    }

    logger.debug(`all verifiers: ${JSON.stringify(allVerifiersUrls)}`);

    var verifierURLs = [];
    for (let i = 0; i < numVerifiers * 2; i++) {
        // randomly select a peer to be verifier
        let verifierUrlIdx = Math.floor(Math.random()*10000) % allVerifiersUrls.length;
        let verifierUrl = allVerifiersUrls[verifierUrlIdx];

        allVerifiersUrls.pop(verifierUrlIdx);
        verifierURLs.push(verifierUrl);

        logger.debug(`selected verifiers: ${JSON.stringify(verifierURLs)}`);
    }

    const verifierURLsStr = verifierURLs.join(";");

    var deviceIds = [];
    for (let i = 0; i < numDevices; i ++) {
        const deviceId = "w" + Math.floor(Math.random()*10000) % numPeersPerOrgs + "." + pubName;
        deviceIds.push(deviceId);
    }

    const deviceIdsStr = deviceIds.join(";");

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

exports.createRandomCampaign = async (numVerifiers, deviceIdsStr) => {
    return utils.callChaincodeFn(async network => {

        const contract = await network.getContract('campaign');
        startDate = "2022-11-13"
        endDate = "2022-11-14"
        const campaign = CreateCampaignUnequalVerifiersArgs(global.numOrgsPerType, global.numPeersPerOrg, numVerifiers, deviceIdsStr, startDate, endDate);
        logger.info(`Create Campaign: ${JSON.stringify(campaign)}`);

        const verifierAddressesStr = campaign.verifierURLs.join(";")

        logger.info(`Create Campaign: ${JSON.stringify(campaign)} and ${verifierAddressesStr} - devices: ${campaign.deviceIds.join(";")}`);

        return contract.submitTransaction("CreateCampaign", campaign.id, campaign.name, campaign.advertiser, campaign.publisher, campaign.startTimeStr, campaign.endTimeStr, verifierAddressesStr, campaign.deviceIds.join(";"));
    }, async response => {
        const resultCampaign = JSON.parse(response);
        logger.debug(`raw response: ${response}`);
        return resultCampaign;
    });
}


exports.getCampaignById = async (camId) => {
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


exports.getAllCampaigns = async () => {
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


exports.deleteCampaignById = async (camId) => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        logger.info("DeleteCampaignById");
        return contract.submitTransaction("DeleteCampaignById", camId);
    }, async (response) => {
        logger.debug(`response: ${response}`);
        return response;
    });
}

exports.deleteAllCampaigns = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        logger.info("DeleteAllCampaigns");
        return contract.submitTransaction("DeleteAllCampaigns");
    }, async (response) => {
        logger.debug(`response: ${response}`);
        return response;
    });
}

exports.getChaincodeData = async () => {
    return utils.callChaincodeFn(async (network) => {
        const contract = await network.getContract("campaign");
        logger.info("GetChaincodeData");
        return contract.submitTransaction("GetChaincodeData");
    }, async (response) => {
        logger.debug(`response: ${response}`);
        return response;
    });
}

exports.CreateCampaignUnequalVerifiersArgs = (numPeersPerOrgs, numOrgsPerType, numVerifiers, numDevices, startDate, endDate) => {
    const camId = "c" + Math.floor(Math.random()*10000);
    const name = "Campaign " + camId;
    const advertiser = "adv"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const publisher = "pub"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const startTimeStr = Math.floor(new Date(startDate).getTime() / 1000);
    const endTimeStr = Math.floor(new Date(endDate).getTime() / 1000);

    var verifierURLs = [];

    var allVerifiersUrls = new Set();

    for (let orgId = 0; orgId < numOrgsPerType; orgId++) {
        for (let peerId = 0; peerId < numPeersPerOrgs; peerId++) {
            const advPeerURL = `peer${peerId}.adv${orgId}.promark.com:5000`;
            const pubPeerURL = `peer${peerId}.pub${orgId}.promark.com:5000`;

            allVerifiersUrls.add(advPeerURL);
            allVerifiersUrls.add(pubPeerURL);
        }
    }

    for (let i = 0; i < numVerifiers; i++) {
        // randomly select a peer to be verifier
        let verifierUrl = Math.floor(Math.random()*10000) % allVerifiersUrls.length;

        allVerifiersUrls.delete(verifierUrl);
        verifierURLs.push(verifierUrl);
    }

    const verifierURLsStr = verifierURLs.join(";");

    var deviceIds = [];
    for (let i = 0; i < numDevices; i ++) {
        const deviceId = "w" + Math.floor(Math.random()*10000) % numPeersPerOrgs + "." + publisher;
        deviceIds.push(deviceId);
    }

    const deviceIdsStr = deviceIds.join(";");

    return {
        camId, name, advertiser, publisher, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr
    };
}
