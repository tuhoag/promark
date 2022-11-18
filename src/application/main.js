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
        const numVerifiers = 1;
        // add a campaign
        let campaign1 = await camClient.createRandomCampaign(numVerifiers, "d1,d2");
        logger.info("campaign:" + JSON.stringify(campaign1));
        let campaign2 = await camClient.createRandomCampaign(numVerifiers, "d1,d2");
        logger.info("campaign:" + JSON.stringify(campaign2));
        // console.log(campaign);

        // sleep(2000);
        // generate its proof
        const cusId = "u1";
        const deviceId = "d1";
        const numTPoCs1 = 3;
        const numTPoCs2 = 2;
        // let customerPoC = await proofClient.generatePoC(campaign.id, cusId)
        let customerPocAndTPoCs1 = await proofClient.generatePoCAndTPoCs(campaign1.id, cusId, numTPoCs1);
        let devicePocAndTPoCs1 = await proofClient.generatePoCAndTPoCs(campaign1.id, deviceId, numTPoCs1);

        let customerPocAndTPoCs2 = await proofClient.generatePoCAndTPoCs(campaign2.id, cusId, numTPoCs2);
        let devicePocAndTPoCs2 = await proofClient.generatePoCAndTPoCs(campaign2.id, deviceId, numTPoCs2);

        // add transaction
        let count1 = await addTokenTransactions(campaign1, deviceId, devicePocAndTPoCs1, customerPocAndTPoCs1, numTPoCs1);
        // let count2 = 0
        let count2 = await addTokenTransactions(campaign2, deviceId, devicePocAndTPoCs2, customerPocAndTPoCs2, numTPoCs2);

        const addedTokenTransactions1 = await proofClient.getTokenTransactionsByCampaignId(campaign1.id, "device")
        logger.info(`camId ${campaign1.id} - numTPoC: ${customerPocAndTPoCs1.tpocs.length} - ${addedTokenTransactions1.length} token transactions - valid: ${customerPocAndTPoCs1.length == addedTokenTransactions1.length}: ${JSON.stringify(addedTokenTransactions1)}`);

        const addedTokenTransactions2 = await proofClient.getTokenTransactionsByCampaignId(campaign2.id, "device")
        logger.info(`camId ${campaign2.id} - numTPoC: ${customerPocAndTPoCs2.tpocs.length} - ${addedTokenTransactions2.length} token transactions - valid: ${customerPocAndTPoCs2.length == addedTokenTransactions2.length}: ${JSON.stringify(addedTokenTransactions2)}`);

        // const addedTokenTransactions2 = await proofClient.getTokenTransactionsByTimestamps(campaign1.startTime, campaign1.endTime)
        // logger.info(`camId ${campaign1.id} - numTPoC: ${customerPocAndTPoCs2.length} - ${addedTokenTransactions2.length} token transactions - valid: ${customerPocAndTPoCs2.length == addedTokenTransactions2.length}: ${JSON.stringify(addedTokenTransactions2)}`);

        // const allAddedTransactions = await proofClient.getAllProofs();
        // logger.info(`all ${allAddedTransactions.length} vs count1: ${count1} count2: ${count2}`);
        // logger.info(`all ${allAddedTransactions.length} added transactions: ${JSON.stringify(allAddedTransactions)}`);
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