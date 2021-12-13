'use strict';


function CreateCampaignArgs(numPeersPerOrgs, numOrgsPerType, numVerifiersPerType) {
    const camId = "c" + Math.floor(Math.random()*10000);
    const name = "Campaign " + camId;
    const advertiser = "adv"+Math.floor(Math.random()*10000) % numOrgsPerType;
    const business = "bus"+Math.floor(Math.random()*10000) % numOrgsPerType;

    var verifierURLs = [];

    for (let i = 0; i < numVerifiersPerType; i++) {
        const advertierPeerName = "peer" + Math.floor(Math.random()*10000) % numPeersPerOrgs;
        const businessPeerName = "peer" + Math.floor(Math.random()*10000) % numPeersPerOrgs;

        const advPeerURL = "http://" + advertierPeerName + "."+advertiser + ".promark.com:5000";
        const busPeerURL = "http://" + businessPeerName + "."+business + ".promark.com:5000";

        verifierURLs.push(advPeerURL);
        verifierURLs.push(busPeerURL);
    }

    const verifierURLsStr = verifierURLs.join(";");
    return {
        camId, name, advertiser, business, verifierURLsStr
    };
}

module.exports = { CreateCampaignArgs };