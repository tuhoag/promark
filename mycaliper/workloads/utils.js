'use strict';

const fs = require("fs");
const process = require("process");

const logger = require('@hyperledger/caliper-core').CaliperUtils.getLogger('promark');

const CreateCampaignArgs = (numPeersPerOrgs, numOrgsPerType, numVerifiersPerType) => {
    const camId = "c" + Math.floor(Math.random()*10000);
    const name = "Campaign " + camId;
    const advertiser = "adv"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const business = "bus"+Math.floor(Math.random()*10000) % numOrgsPerType;

    var verifierURLs = [];

    for (let i = 0; i < numVerifiersPerType; i++) {
        const advertierPeerName = "peer" + Math.floor(Math.random()*10000) % numPeersPerOrgs;
        const businessPeerName = "peer" + Math.floor(Math.random()*10000) % numPeersPerOrgs;

        const advPeerURL = advertierPeerName + "."+advertiser + ".promark.com:5000";
        const busPeerURL = businessPeerName + "."+business + ".promark.com:5000";

        verifierURLs.push(advPeerURL);
        verifierURLs.push(busPeerURL);
    }

    const verifierURLsStr = verifierURLs.join(";");
    return {
        camId, name, advertiser, business, verifierURLsStr
    };
}

const loadInitData = (numCampaigns, numProofs, numVerifiersPerType) => {
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

module.exports = { CreateCampaignArgs, loadInitData };