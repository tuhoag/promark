'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const utils = require('./utils');
const camClient = require('./campaignClient');
const proofClient = require('./proofClient');
const fs = require('fs');
const path = require('path');
const setting = require('./setting');

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
    } else if (command == "delall") {
        return camClient.deleteAllCampaigns();
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

        case "all":
            return proofClient.getAllProofs();
        case "delall":
            return proofClient.deleteAllProofs();

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

const dataCommandHandler = async argv => {
    const command = argv[0];

    switch(command) {
        case "init":
            const numCampaigns = argv[1];
            const numProofs = argv[2];
            const numVerifiers = argv[3];
            return initData(numCampaigns, numProofs, numVerifiers);

        case "delall":
            return deleteAllData();

        default:
            throw `Unsupported proof command ${command}`;
    }
}

const initData = async (numCampaigns, numProofs, numVerifiers) => {
    let campaigns = [];

    for (let i = 0; i < numCampaigns; i ++) {
        let campaign = await camClient.createRandomCampaign(numVerifiers);

        campaigns.push(campaign);
    }

    let proofs = [];

    for (let i = 0; i < numProofs; i ++) {
        const camIdx = utils.getId(campaigns.length);
        console.log(JSON.stringify(campaigns[camIdx]));
        let proof = await proofClient.generateProofForRandomUser(campaigns[camIdx].id);
        console.log(`generated proof: ${JSON.stringify(proof)}`);
        proofs.push(proof);
    }

    console.log(JSON.stringify(proofs));
    console.log("printing proofs");

    let initData = {
        "campaigns": campaigns,
        "proofs": [],
    }

    for (let proof of proofs) {
        console.log(`comm:${proof.comm}`);
        console.log(`rsStr:${proof.rs.join(";")}`);

        initData["proofs"].push({
            "comm": proof.comm,
            "rsStr": proof.rs.join(";"),
        });
    }

    const initDataPath = path.join(setting.initDataDirPath, `initData-${numCampaigns}-${numProofs}-${numVerifiers}.json`);
    fs.writeFileSync(initDataPath, JSON.stringify(initData));
    console.log(`saved init data to: ${initDataPath}`);
    return initData;
}

const deleteAllData = async () => {
    // remove all init data
    fs.readdirSync(setting.initDataDirPath).forEach(file => {
        const filePath = path.join(setting.initDataDirPath, file);
        console.log(`deleting file: ${filePath}`);

        fs.unlink(filePath, err => {
            if (err) throw err;
        });
    });

    return camClient.deleteAllCampaigns().then(proofClient.deleteAllProofs());
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
        case "data":
            return dataCommandHandler(subArgs);
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