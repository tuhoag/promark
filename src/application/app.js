'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const utils = require('./utils');
const camClient = require('./campaignClient');
const proofClient = require('./proofClient');
const { syncBuiltinESMExports } = require('module');

const campaignCommandHandler = async (argv) => {
    console.log(argv)
    const command = argv[0];

    if (command == "create") {
        const numVerifiers = argv[1];

        return camClient.createRandomCampaign(numVerifiers);
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
    let camId;

    switch(command) {
        case "gen":
            camId = argv[1];
            const userId = argv[2];
            return proofClient.generateProofForRandomUser(camId, userId);

        case "add":
            // camId = argv[1];
            const comm = argv[1];
            const rsStr = argv[2];
            return proofClient.addProof(comm, rsStr);

        default:
            throw `Unsupported proof command ${command}`;
    }
}

// wrap settimeout in a promise to create a wait
const sleep = (ms) => {
    return new Promise((resolve) => setTimeout(resolve, ms));
};

const testCommandHandler = async argv => {
    const numVerifiers = 1;
    // add a campaign
    let campaign = await camClient.createRandomCampaign(numOrgsPerType, numPeersPerOrg, numVerifiers);
    console.log("outside:" + JSON.stringify(campaign));
    // console.log(campaign);

    sleep(2000);
    // generate its proof
    let proof = await proofClient.generateProofForRandomUser(campaign.id);
    console.log("outside:" + JSON.stringify(proof));

    sleep(2000);
    // add the generated proof
    let addedProof = await proofClient.addProof(campaign.id, proof.Comm, proof.Rs.join(";"));
    console.log("outside:" + JSON.stringify(addedProof));
    // show all campaigns

    // show all proofs
}

// Main program function
const main = async (argv) => {
    console.log(argv)
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
        case "test":
            return testCommandHandler(subArgs);
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