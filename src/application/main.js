'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class

const fs = require('fs');
const path = require('path');
const utils = require('./utils');
const camClient = require('./campaignClient');
const proofClient = require('./proofClient');
const dataClient = require('./dataClient');
const setting = require('./setting');
const { createLogger } = require('./logger');
logger = createLogger(__filename);

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
        logger.error(`Unsupported campaign command ${command}`);
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
    } else if (command == "allids") {
        return proofClient.getAllProofIds();
    } else if (command == "verify") {
        camId = argv[1];
        let proofId = argv[2];

        return proofClient.verifyProof(camId, proofId);
    } else if (command == "query-time") {
            let startTime = argv[1];
            let endTime = argv[2];

            return proofClient.queryByTimestamps(startTime, endTime);
    } else {
        logger.error(`Unsupported campaign command ${command}`);
        throw `Unsupported proof command ${command}`;
    }
}

// wrap settimeout in a promise to create a wait
const sleep = (ms) => {
    return new Promise((resolve) => setTimeout(resolve, ms));
};


const testCommandHandler = async argv => {

    try {
        const numVerifiersPerOrg = 1;
        const cusId = "u1";
        const deviceId = "d1";
        const numTPoCs = [3, 2];
        const numCampaigns = 2;
        var numTrans = 0

        for (let i = 0; i < numTPoCs.length; i ++) {
            numTrans += numTPoCs[i];
        }

        let campaigns = [];

        for (let i = 0; i < numCampaigns; i++) {
            let campaign = await camClient.createRandomCampaign(numVerifiersPerOrg, "d1,d2");

            logger.debug(`campaign ${i}: ${JSON.stringify(campaign)}`);
            campaigns.push(campaign);
        }

        logger.info(`added ${campaigns.length} campaigns`);

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

        logger.info(`generated ${customersPocs.length} tokens`);


        for (let i = 0; i < numCampaigns; i++) {
            let count1 = await proofClient.addTokenTransactions(campaigns[i], deviceId, devicesPocs[i], customersPocs[i], numTPoCs[i]);

            logger.info(`cam ${campaigns[i].id} - added ${count1} token transactions`);
        }

        for (let i = 0; i < numCampaigns; i++) {
            let addedTokenTransactions = await proofClient.getTokenTransactionsByCampaignId(campaigns[i].id, "device")
            logger.info(`camId ${campaigns[i].id} - numTPoC: ${customersPocs[i].tpocs.length} - ${addedTokenTransactions.length} token transactions - valid: ${customersPocs[i].length == addedTokenTransactions.length}: ${JSON.stringify(addedTokenTransactions)}`);

            addedTokenTransactions = await proofClient.getTokenTransactionsByCampaignId(campaigns[i].id, "device")
            logger.info(`camId ${campaigns[i].id} - numTPoC: ${customersPocs[i].tpocs.length} - ${addedTokenTransactions.length} token transactions - valid: ${customersPocs[i].length == addedTokenTransactions.length}: ${JSON.stringify(addedTokenTransactions)}`);

            let counts = await proofClient.simulateGetTokenTransactionsByCampaignId(campaigns[i].id, "device", numTrans)
            logger.info(`camId ${campaigns[i].id} - numTPoC: ${customersPocs[i].tpocs.length} - ${addedTokenTransactions.length} token transactions - counts: ${counts}`);

            counts = await proofClient.simulateGetTokenTransactionsByCampaignId(campaigns[i].id, "device", 20)
            logger.info(`camId ${campaigns[i].id} - numTPoC: ${customersPocs[i].tpocs.length} - ${addedTokenTransactions.length} token transactions - counts: ${counts}`);
        }
    } catch (error) {
        logger.debug(`error: ${error.stack}`);
    }
}


const dataCommandHandler = async argv => {
    logger.debug(`data handler args: ${argv}`);
    const numOrgsPerType = argv[0];
    const numPeersPerOrg = argv[1];
    const numVerifiers = argv[4];
    const numTrans = argv[5];
    const command = argv[3].trim();

    if (command == "verify-cam-tpocs") {
        return await dataClient.generateDataForBatchVerificationEvaluation(numOrgsPerType, numPeersPerOrg, numVerifiers, 2, numTrans);
    } else {
        logger.error(`Unsupported data command ${command}`);
        throw `Unsupported data command ${command}`;
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
            return dataCommandHandler(argv);
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