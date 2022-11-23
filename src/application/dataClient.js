const fs = require('fs');

const utils = require('./utils');
const camClient = require('./campaignClient');
const proofClient = require('./proofClient');

const { createLogger } = require('./logger');
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
        let numTransPerCampaigns = Math.floor(numTrans / campaigns.length);

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