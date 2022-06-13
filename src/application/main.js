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
        const userId = argv[2];

        logger.debug(`generatePoCForRandomUser: ${camId}, ${userId}`);
        return proofClient.generatePoCForRandomUser(camId, userId);
    } else if (command == "gentpoc") {
        camId = argv[1];
        const userId = argv[2];
        const numTPoCs = argv[3];

        logger.debug(`generatePoCAndTPoCs: ${camId}, ${userId}`);
        return proofClient.generatePoCAndTPoCs(camId, userId, numTPoCs);
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
        const comm = argv[1];
        const rsStr = argv[2];
        return proofClient.addProof(comm, rsStr);
    } else if (command == "add2") {
        const camId = argv[1];
        const deviceId = argv[2];
        const cusId = argv[3];
        const cusComm = argv[4];
        const cusRsStr = argv[5];
        const addedTime = Date.now();

        return proofClient.addProof2(camId, deviceId, cusId, cusComm, cusRsStr, addedTime);
    } else if (command == "all") {
        return proofClient.getAllProofs();
    } else if (command == "delall") {
        return proofClient.deleteAllProofs();
    } else if (command == "verify") {
        camId = argv[1];
        let proofId = argv[2];

        return proofClient.verifyProof(camId, proofId);
    } else {
        logger.info(`Unsupported campaign command ${command}`);
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
    let campaign = await camClient.createRandomCampaign(numVerifiers, "d1,d2");
    logger.info("campaign:" + JSON.stringify(campaign));
    // console.log(campaign);

    // sleep(2000);
    // generate its proof
    const cusId = "u1";
    const numTPoCs = 3;
    let customerPocAndTPoCs = await proofClient.generatePoCAndTPoCs(campaign.id, cusId, numTPoCs);
    logger.info("customer PoC & TPoCs:" + JSON.stringify(customerPocAndTPoCs));

    // sleep(2000);

    const deviceId = "d1";
    let devicePocAndTPoCs = await proofClient.generatePoCAndTPoCs(campaign.id, deviceId, numTPoCs);
    logger.info("device PoC & TPoCs:" + JSON.stringify(devicePocAndTPoCs));

    // sleep(2000);

    // verify pocs
    let result = await proofClient.verifyPoCProof(campaign.id, customerPocAndTPoCs.poc.comm, customerPocAndTPoCs.poc.r);

    logger.info(`verification customer poc result-true: ${result}`);

    result = await proofClient.verifyPoCProof(campaign.id, devicePocAndTPoCs.poc.comm, devicePocAndTPoCs.poc.r);

    logger.info(`verification device poc result-true: ${result}`);

    for (let i = 0; i  < numTPoCs; i++) {
        let customerTPoC = customerPocAndTPoCs.tpocs[i];
        let deviceTPoC = devicePocAndTPoCs.tpocs[i];

        logger.debug(`Got customer tpoc ${i}: ${JSON.stringify(customerTPoC)}`);
        logger.debug(`Got device tpoc ${i}: ${JSON.stringify(deviceTPoC)}`);

        result = await proofClient.verifyTPoCProof(campaign.id, customerTPoC.tComms, customerTPoC.tRs, customerTPoC.hashes, customerTPoC.key);
        logger.info(`verification customer tpoc ${i} result-true: ${result}`);

        result = await proofClient.verifyTPoCProof(campaign.id, deviceTPoC.tComms, deviceTPoC.tRs, deviceTPoC.hashes, deviceTPoC.key);
        logger.info(`verification device tpoc ${i} result-true: ${result}`);
    }

    // add transaction
    const diff = campaign.endTime - campaign.startTime;

    for (let i = 0; i  < numTPoCs; i++) {
        let customerTPoC = customerPocAndTPoCs.tpocs[i];
        let deviceTPoC = devicePocAndTPoCs.tpocs[i];
        let addedTime = campaign.startTime + Math.floor(Math.random()) % diff;

        let transaction = await proofClient.addProof(campaign.id, deviceId, addedTime, deviceTPoC, customerTPoC);

        logger.info(`added transaction: ${JSON.stringify(transaction)}`);
    }

    const allAddedTransactions = await proofClient.getAllProofs();
    logger.info(`all ${len(allAddedTransactions)} added transactions: ${JSON.stringify(allAddedTransactions)}`);
}

const dataCommandHandler = async argv => {
    const command = argv[0];

    switch (command) {
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
    const deviceIdsStr = ["d1", "d2"].join(",");

    for (let i = 0; i < numCampaigns; i++) {
        let campaign = await camClient.createRandomCampaign(numVerifiers, deviceIdsStr);

        campaigns.push(campaign);
    }

    let proofs = [];

    for (let i = 0; i < numProofs; i++) {
        const camIdx = utils.getId(campaigns.length);
        logger.debug(JSON.stringify(campaigns[camIdx]));
        const userId = `u${utils.getId(10000)}`;
        let proof = await proofClient.generatePoCForRandomUser(campaigns[camIdx].id, userId);

        // camId string, deviceId string, cusId string, cusComm string, cusRsStr string, addedTimeStr string
        proof["camId"] = campaigns[camIdx].id;
        proof["deviceId"] = campaigns[camIdx].deviceIds[0];
        proof["customerId"] = userId;
        proof["addedTimeStr"] = utils.randomDate(new Date(2022, 0, 1), new Date()).toISOString();

        logger.debug(`generated proof: ${JSON.stringify(proof)}`);
        proofs.push(proof);

        sleep(1000);
    }

    logger.debug(JSON.stringify(proofs));
    logger.debug("printing proofs");

    let initData = {
        "campaigns": campaigns,
        "proofs": [],
    }

    for (let proof of proofs) {
        logger.debug(`comm:${proof.comm}`);
        logger.debug(`rsStr:${proof.rs.join(";")}`);

        initData["proofs"].push({
            "camId": proof["camId"],
            "deviceId": proof["deviceId"],
            "customerId": proof["customerId"],
            "addedTimeStr": proof["addedTimeStr"],
            "cusComm": proof.comm,
            "cusRsStr": proof.rs.join(";"),
        });
    }

    const initDataPath = path.join(setting.initDataDirPath, `initData-${numCampaigns}-${numProofs}-${numVerifiers}.json`);
    fs.writeFileSync(initDataPath, JSON.stringify(initData));
    logger.info(`saved init data to: ${initDataPath}`);
    return initData;
}

const deleteAllData = async () => {
    // remove all init data
    fs.readdirSync(setting.initDataDirPath).forEach(file => {
        const filePath = path.join(setting.initDataDirPath, file);
        logger.info(`deleting file: ${filePath}`);

        fs.unlink(filePath, err => {
            if (err) throw err;
        });
    });

    return camClient.deleteAllCampaigns().then(proofClient.deleteAllProofs());
}

const chaincodeCommandHandler = async argv => {
    // return utils.callChaincodeFn(async (network) => {
    //     const contract = await network.getContract("campaign");
    //     logger.info("DeleteCampaignById");
    //     return contract.submitTransaction("DeleteCampaignById", camId);
    // }, async (response) => {
    //     logger.debug(`response: ${response}`);
    //     return response;
    // });

    const response = await camClient.getChaincodeData();
    logger.info(response);
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