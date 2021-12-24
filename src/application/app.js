'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const utils = require('./utils');
const camClient = require('./campaignClient');

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
    } else {
        throw `Unsupported campaign command ${command}`;
    }
}

// Main program function
const main = async (argv) => {
    const command = argv[0];

    if (command == "campaign") {
        return campaignCommandHandler(argv.slice(1));
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