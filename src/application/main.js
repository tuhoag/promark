'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class

const fs = require('fs');
const path = require('path');
const utils = require('./utils');
const camClient = require('./campaignClient');
const proofClient = require('./proofClient');
const setting = require('./setting');
const logger = require('./logger')(__filename, "debug");

const campaignCommandHandler = async (argv) => {
    logger.debug(argv);
    const command = argv[0];

    if (command == "add") {
        const numVerifiers = argv[1];
        const deviceIdsStr = argv[2];

        return camClient.createRandomCampaign(numVerifiers, deviceIdsStr);
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
        logger.info(`Unsupported campaign command ${command}`);
        throw `Unsupported campaign command ${command}`;
    }
}

const proofCommandHandler = async argv => {
    logger.debug(`proof handler args: ${argv}`);
    const command = argv[0].trim();
    let camId;

    if (command == "gen") {
        camId = argv[1];
        let userId = argv[2];

        logger.debug(`generatePoC: ${camId}, ${userId}`);
        return proofClient.generatePoC(camId, userId);
    } else if (command == "gentpocs") {
        camId = argv[1];
        let userId = argv[2];
        let numTPoCs = argv[3];

        // logger.debug(`generating poc: ${camId} - ${userId}`);
        let poc = await proofClient.generatePoC(camId, userId);

        logger.debug(`generated poc: ${poc.comm} - ${poc.r}`);
        return proofClient.generateTPoCs(camId, poc.comm, poc.r, poc.numVerifiers, numTPoCs);
    } else if (command == "verifypoc") {
        camId = argv[1];
        const comm = argv[2];
        const r = argv[3];

        logger.debug(`verifyPoCProof: ${camId}, ${comm}, ${r}`);
        return proofClient.verifyPoCProof(camId, comm, r);
    } else if (command == "verifytpoc") {
        camId = argv[1];
        const csStr = argv[2];
        const rsStr = argv[3];

        const cs = argv[2].split(';');
        const rs = argv[3].split(';');
        const hs = ["h1", "h2"];
        const key = "";

        logger.debug(`verifytPoCProof: ${camId}, ${cs}, ${rs}`);
        return proofClient.verifyTPoCProof(camId, cs, rs, hs, key);
    } else if (command == "add") {
        const camId = argv[1];
        const deviceId = argv[2];
        const cusId = argv[3];
        const cusComm = argv[4];
        const cusRsStr = argv[5];
        const addedTime = Date.now();

        return proofClient.addProof(camId, deviceId, cusId, cusComm, cusRsStr, addedTime);
    } else if (command == "all") {
        return proofClient.getAllProofs();
    } else if (command == "delall") {
        return proofClient.deleteAllProofs();
    } else if (command == "verify") {
        camId = argv[1];
        let proofId = argv[2];

        return proofClient.verifyProof(camId, proofId);
    } else if (command == "query-time") {
            let startTime = argv[1];
            let endTime = argv[2];

            return proofClient.queryByTimestamps(startTime, endTime);
    } else {
        logger.info(`Unsupported campaign command ${command}`);
        throw `Unsupported proof command ${command}`;
    }
}

// wrap settimeout in a promise to create a wait
const sleep = (ms) => {
    return new Promise((resolve) => setTimeout(resolve, ms));
};

const addTokenTransactions = async (campaign, deviceId, devicePocAndTPoCs, customerPocAndTPoCs, numAdditions) => {
    let numTPoCs = customerPocAndTPoCs.tpocs.length;

    if (numAdditions > numTPoCs) {
        throw `numAdditions > numTPoCs: ${numAdditions} > ${numTPoCs}`;
    }

    const diff = campaign.endTime - campaign.startTime;
    for (let i = 0; i  < numAdditions; i++) {
        logger.debug(JSON.stringify(customerPocAndTPoCs));
        logger.debug(JSON.stringify(customerPocAndTPoCs.tpocs[i]));
        let customerTPoC = customerPocAndTPoCs.tpocs[i];
        let deviceTPoC = devicePocAndTPoCs.tpocs[i];
        let addingTime = Math.floor(Math.random() * diff) % diff;
        let addedTime = campaign.startTime + addingTime;

        logger.debug(`diff: ${diff} - adding time: ${addingTime} - valid: ${(campaign.startTime < addedTime) && (addedTime < campaign.endTime)}`);

        let transaction = await proofClient.addProof(campaign.id, deviceId, addedTime, deviceTPoC, customerTPoC);

        logger.debug(`added transaction: ${JSON.stringify(transaction)}`);

        await sleep(1000);
    }

    return numAdditions;
}

const testCommandHandler = async argv => {

    try {
        const numVerifiersPerOrg = 1;
        const cusId = "u1";
        const deviceId = "d1";
        const numTPoCs = [3, 2];
        const numCampaigns = 2;


        let campaigns = [];

        for (let i = 0; i < numCampaigns; i++) {
            let campaign = await camClient.createRandomCampaign(numVerifiersPerOrg, "d1,d2");

            logger.info(`campaign ${i}: ${JSON.stringify(campaign)}`);
            campaigns.push(campaign);
        }

        let customersPocs = [];
        let devicesPocs = [];

        for (let i = 0; i < numCampaigns; i++) {
            let customerPoc = await proofClient.generatePoC(campaigns[i].id, cusId);

            let customerPocAndTPoCs = await proofClient.generateTPoCs(campaigns[i].id, customerPoc.comm, customerPoc.r, numVerifiersPerOrg * 2, numTPoCs[i]);

            let devicePoC = await proofClient.generatePoC(campaigns[i].id, deviceId);
            let devicePocAndTPoCs = await proofClient.generateTPoCs(campaigns[i].id, devicePoC.comm, devicePoC.r, numVerifiersPerOrg * 2, numTPoCs[i]);

            customersPocs.push(customerPocAndTPoCs);
            devicesPocs.push(devicePocAndTPoCs);
        }


        for (let i = 0; i < numCampaigns; i++) {
            let count1 = await addTokenTransactions(campaigns[i], deviceId, devicesPocs[i], customersPocs[i], numTPoCs[i]);
        }


        for (let i = 0; i < numCampaigns; i++) {
            let addedTokenTransactions = await proofClient.getTokenTransactionsByCampaignId(campaigns[i].id, "device")
            logger.info(`camId ${campaigns[i].id} - numTPoC: ${customersPocs[i].tpocs.length} - ${addedTokenTransactions.length} token transactions - valid: ${customersPocs[i].length == addedTokenTransactions.length}: ${JSON.stringify(addedTokenTransactions)}`);
        }
    } catch (error) {
        logger.debug(`error: ${error}`);
    }
}




// Main program function
const main = async (argv) => {
    logger.info(`main argv: ${argv}`);
    const numOrgsPerType = argv[0];
    const numPeersPerOrg = argv[1];
    const command = argv[2];
    const subArgs = argv.slice(3);

    global.numOrgsPerType = numOrgsPerType;
    global.numPeersPerOrg = numPeersPerOrg;

    switch (command) {
        case "campaign":
            return campaignCommandHandler(subArgs);
        case "proof":
            return proofCommandHandler(subArgs);
        case "test":
            return testCommandHandler(subArgs);
        case "data":
            return dataCommandHandler(subArgs);
        case "chaincode":
            return chaincodeCommandHandler(subArgs);
        default:
            throw `Unsupported command ${command}`;
    }
}

main(process.argv.slice(2)).then(() => {
    logger.info("Program complete.")
}).catch((err) => {
    if (!err) {
        logger.error("Error: " + err.stack);
    }
    process.exit(-1);
});