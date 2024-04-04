const fs = require('fs');

const utils = require('./utils');
const camClient = require('./campaignClient');
const proofClient = require('./proofClient');

const { createLogger } = require('./logger');
const { createImportSpecifier } = require('typescript');
logger = createLogger(__filename);

exports.generateDataForBatchVerificationEvaluation = async (numOrgsPerType, numPeersPerOrgs, numVerifiers, numDevices, numTrans) => {
    logger.debug(`generateDataForBatchVerificationEvaluation: ${numOrgsPerType}, ${numPeersPerOrgs}, ${numVerifiers}, ${numDevices}, ${numTrans}`);
    try {
        let campaigns = utils.CreateCampaignsWithEqualVerifiersArgs(numOrgsPerType, numPeersPerOrgs, numVerifiers, numDevices);

        logger.debug(`${JSON.stringify(campaigns)}`);
        let camIds = [];

        for (let i = 0; i < campaigns.length; i++) {
            const {id, name, advName, pubName, startTime, endTime, verifierURLsStr, deviceIdsStr} = campaigns[i];

            await camClient.createCampaign(id, name, advName, pubName, startTime, endTime, verifierURLsStr, deviceIdsStr);

            camIds.push(id);
        }

        let numAddedTrans = 0;
        let numTransPerCampaigns;

        if (numTrans == 0) {
            numTrans = campaigns.length;
        }

        numTransPerCampaigns = Math.floor(numTrans / campaigns.length);


        logger.debug(`numTransPerCampaigns: ${numTransPerCampaigns}`);

        for (let i = 0; i < campaigns.length; i++) {
            let userId = Math.floor(Math.random()*10000);
            let deviceId = Math.floor(Math.random()*10000 % numDevices);

            let customerPoc = await proofClient.generatePoC(campaigns[i].id, userId);
            let devicePoC = await proofClient.generatePoC(campaigns[i].id, deviceId);

            // generate number of tpocs
            let curNumTransPerCampaigns = numTransPerCampaigns;
            if (i == campaigns.length - 1) {
                curNumTransPerCampaigns = numTrans - numAddedTrans;
            }

            numAddedTrans = numAddedTrans + curNumTransPerCampaigns;

            let customerPocAndTPoCs = await proofClient.generateTPoCs(campaigns[i].id, customerPoc.comm, customerPoc.r, numVerifiers, curNumTransPerCampaigns);
            let devicePocAndTPoCs = await proofClient.generateTPoCs(campaigns[i].id, devicePoC.comm, devicePoC.r, numVerifiers, curNumTransPerCampaigns);

            await proofClient.addTokenTransactions(campaigns[i], deviceId, devicePocAndTPoCs, customerPocAndTPoCs, curNumTransPerCampaigns);

            logger.info(`numAddedTrans: ${numAddedTrans}/${numTrans}`);
        }


        let content = camIds.join(",");
        fs.writeFileSync(`../caliper/data/cams-${numOrgsPerType}-${numPeersPerOrgs}-${numVerifiers}.txt`, content);
    } catch (err) {
        logger.error(err.stack);
    }
}

exports.evaluateFindTPoCs = async (mode, limit) => {
    // find the first campaign id
    const campaigns = await camClient.getAllCampaigns();
    logger.debug(campaigns[0]);
    logger.debug(campaigns[0].verifierURLs.length);

    // loop through all proofs and evaluate
    console.time("verify");
    await proofClient.simulateGetTokenTransactionsByCampaignId(campaigns[0].id, mode, limit);
    console.timeEnd("verify");
}

exports.evaluateFindTPoCs2 = async (mode, limit) => {
    // find the first campaign id
    const campaigns = await camClient.getAllCampaigns();
    logger.debug(campaigns[0]);
    logger.debug(campaigns[0].verifierURLs.length);

    // loop through all proofs and evaluate
    console.time("verify");
    const trans = await proofClient.getAllProofs();

    for (const tran of trans) {
        const splits = tran.id.split(":");
        const camId = splits[1];

        let tpocs = [tran.deviceTPoC];

        if (mode == "all") {
            tpocs.push(tran.customerTPoC);
        }

        let flag = true;
        for (const tpoc of tpocs) {
            const curFlag = await proofClient.verifyTPoCProof(camId, tpoc.tComms, tpoc.tRs, tpoc.hashes, tpoc.key);
            if (curFlag == false) {
                flag = false;
            }
        }
    }
    console.dir(trans, {depth: null});
    console.timeEnd("verify");
}