'use strict';

const fs = require("fs");
const process = require("process");

const logger = require('@hyperledger/caliper-core').CaliperUtils.getLogger('promark');

exports.CreateCampaignUnequalVerifiersArgs = (numOrgsPerType, numPeersPerOrgs, numVerifiers, numDevices) => {
    const camId = "c" + Math.floor(Math.random()*10000);
    const name = "Campaign " + camId;
    const advName = "adv"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const pubName = "pub"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const startTimeStr = Math.floor(new Date("2022-05-01").getTime() / 1000);
    const endTimeStr = Math.floor(new Date("2022-07-01").getTime() / 1000);

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
    for (let i = 0; i < numVerifiers; i++) {
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

    return {
        camId, name, advName, pubName, startTimeStr, endTimeStr, verifierURLsStr, deviceIdsStr
    };
}

exports.CreateCampaignArgs = (numPeersPerOrgs, numOrgsPerType, numVerifiersPerType, numDevices) => {
    const camId = "c" + Math.floor(Math.random()*10000);
    const name = "Campaign " + camId;
    const advertiser = "adv"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const publisher = "pub"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const startTimeStr = Math.floor(new Date("2022-05-01").getTime() / 1000);
    const endTimeStr = Math.floor(new Date("2022-07-01").getTime() / 1000);

    var verifierURLs = [];

    for (let i = 0; i < numVerifiersPerType; i++) {
        const advertierPeerName = "peer" + Math.floor(Math.random()*10000) % numPeersPerOrgs;
        const publisherPeerName = "peer" + Math.floor(Math.random()*10000) % numPeersPerOrgs;

        const advPeerURL = advertierPeerName + "."+advertiser + ".promark.com:5000";
        const pubPeerURL = publisherPeerName + "."+publisher + ".promark.com:5000";

        verifierURLs.push(advPeerURL);
        verifierURLs.push(pubPeerURL);
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

exports.loadInitData = (numCampaigns, numProofs, numVerifiersPerType) => {
    const path = `./data/initData-${numCampaigns}-${numProofs}-${numVerifiersPerType}.json`;
    // throw new Error(process.cwd());
    try {
        const data = fs.readFileSync(path, "utf8");
        logger.info("raw init data: ", data);

        return JSON.parse(data);
    } catch (err) {
        throw err;
        logger.error(err);
    }
}

// module.exports = { CreateCampaignArgs, loadInitData };