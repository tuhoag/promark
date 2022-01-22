'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const utils = require('./utils');
const camClient = require('./campaignClient');
const proofClient = require('./proofClient');

const campaignCommandHandler = async (argv) => {
    const command = argv[0];

    if (command == "create") {
        const numOrgsPerType = argv[1];
        const numPeersPerOrg = argv[2];
        const numVerifiers = argv[3];

        return camClient.createRandomCampaign(numOrgsPerType, numPeersPerOrg, numVerifiers);
    } else if (command == "get") {
        const camId = argv[1];

        return camClient.getCampaignById(camId);
    } else if (command == "all") {
        return camClient.getAllCampaigns();
    } else if (command == "del") {
        const camId = argv[1];
        return camClient.deleteCampaignById(camId);
    } else {
        throw `Unsupported campaign command ${command}`;
    }
}

const proofCommandHandler = async argv => {
    const command = argv[0];

    switch(command) {
        case "gen":
            const camId = argv[1];
            const userId = argv[2];
            return proofClient.generateProofForRandomUser(camId, userId);

        default:
            throw `Unsupported proof command ${command}`;
    }
}

// Main program function
const main = async (argv) => {
    const numOrgsPerType = argv[0];
    const numPeersPerOrg = argv[1];
    const command = argv[2];
    const subArgs = argv.slice(3);

    global.numOrgsPerType = numOrgsPerType;
    global.numPeersPerOrg = numPeersPerOrg;

    switch(command) {
        case "campaign":
            return campaignCommandHandler(subArgs);
        case "proof":
            return proofCommandHandler(subArgs);
        default:
            throw `Unsupported command ${command}`;
    }
}

main(process.argv.slice(2)).then(() => {
    console.log('Issue program complete.');
}).catch((e) => {
    console.log('Issue program exception.');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);
});