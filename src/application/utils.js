const { Wallets, Gateway } = require('fabric-network');
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');

const setting = require('./setting');
const { createLogger } = require('./logger');
logger = createLogger(__filename);


const buildCPP = async (numOrgsPerType, numPeersPerOrg) => {
    const connectionProfileName = path.resolve(__dirname, '..', '..', 'config', 'network', `networkConfig-${numOrgsPerType}-${numPeersPerOrg}.yaml`);
    return yaml.safeLoad(fs.readFileSync(connectionProfileName, 'utf8'));
}

const buildWallet = async (userName) => {
    // Create a new  wallet : Note that wallet is for managing identities.
    const walletPath = './wallet/' + userName
    logger.debug(`Built a file system wallet at ${walletPath}`);
    return Wallets.newFileSystemWallet(walletPath);
}

const getOrCreateIdentity = async (wallet, orgUserName, orgName) => {
    const orgFullName = getOrgFullName(orgName)
    const orgUserFullName = getOrgUserFullName(orgUserName, orgFullName)

    const cryptoPath = path.resolve(__dirname, '..', '..', 'credentials', 'peerOrganizations', orgFullName, "users", orgUserFullName, "msp");
    const certPath = path.resolve(cryptoPath, "signcerts", orgUserFullName + "-cert.pem");
    const keyPath = path.resolve(cryptoPath, "keystore", "priv_sk")

    const cert = fs.readFileSync(certPath).toString();
    const key = fs.readFileSync(keyPath).toString();

    const identity = {
        credentials: {
            certificate: cert,
            privateKey: key,
        },
        mspId: 'adv0MSP',
        type: 'X.509',
    };

    await wallet.put(orgUserFullName, identity);

    return orgUserFullName
}

const getOrgFullName = (orgName) => {
    return orgName + ".promark.com";
}

const getOrgUserFullName = (orgUserName, orgFullName) => {
    return orgUserName + "@" + orgFullName;
}

exports.connectToGateway = async (userName, orgName, orgUserName) => {
    const wallet = await buildWallet(userName);
    const gateway = new Gateway();

    // Load connection profile; will be used to locate a gateway
    let connectionProfile = await buildCPP(global.numOrgsPerType, global.numPeersPerOrg);

    // wallet = await buildWallet(userName)
    const orgUserFullName = await getOrCreateIdentity(wallet, orgUserName, orgName);

    // Set connection options; identity and wallet
    let connectionOptions = {
        identity: orgUserFullName,
        wallet: wallet,
        discovery: { enabled: false, asLocalhost: false }
    };

    // Connect to gateway using application specified parameters

    await gateway.connect(connectionProfile, connectionOptions);

    return gateway;
}

exports.connectToNetwork = async (userName, orgName, orgUserName) => {
    const gateway = exports.connectToGateway(userName, orgName, orgUserName);
    logger.debug('Use network channel: ' + setting.channelName);
    return gateway.getNetwork(setting.channelName);
}


exports.ChaincodeCaller = class ChaincodeCaller {
    constructor() {
        this.gateway = null;
        this.network = null;
    }

    static async init() {
        let caller = new ChaincodeCaller();
        caller.gateway = await exports.connectToGateway(setting.userName, setting.orgName, setting.orgUserName);
        caller.network = await caller.gateway.getNetwork(setting.channelName);

        return caller
    }

    async call(fn) {
        return fn(this.network);
    }

    async disconnect() {
        this.gateway.disconnect();
    }
}

// export {ChaincodeCaller};

exports.callChaincodeFn = async (requestFn, responseFn) => {
    let gateway;
    try {
        gateway = await exports.connectToGateway(setting.userName, setting.orgName, setting.orgUserName);
        const network = await gateway.getNetwork(setting.channelName);
        const response = await requestFn(network);
        return responseFn(response);
    } catch (error) {
        logger.error(error.stack);
    } finally {
        // Disconnect from the gateway
        logger.debug('Disconnect from Fabric gateway.');
        gateway.disconnect();
    }
}

exports.getId = (maxNum) => {
    return Math.floor(Math.random() * 10000) % maxNum;
}

exports.randomDate = (start, end) => {
    return new Date(start.getTime() + Math.random() * (end.getTime() - start.getTime()));
}

exports.CreateCampaignsWithEqualVerifiersArgs = (numOrgsPerType, numPeersPerOrgs, numVerifiers, numDevices) => {
    let campaigns = [];
    let verifierUrls = [];
    const startTime = Math.floor((new Date("2022.09.01").getTime() / 1000).toFixed(0));
    const endTime = Math.floor((new Date("2022.10.01").getTime() / 1000).toFixed(0));

    var deviceIds = [];
    for (let i = 0; i < numDevices; i ++) {
        const deviceId = `w${i}`;
        deviceIds.push(deviceId);
    }

    const deviceIdsStr = deviceIds.join(";");

    for (let orgId = 0; orgId < numOrgsPerType; orgId++) {
        let advName = `adv${orgId}`;
        let pubName = `pub${orgId}`;

        for (let peerId = 0; peerId < numPeersPerOrgs; peerId++) {
            const advPeerURL = `peer${peerId}.${advName}.promark.com:5000`;
            const pubPeerURL = `peer${peerId}.${pubName}.promark.com:5000`;

            verifierUrls.push(advPeerURL, pubPeerURL);
        }
    }

    let numCampaigns = Math.floor(verifierUrls.length / numVerifiers);

    for (let camIdx = 0; camIdx < numCampaigns; camIdx ++) {
        let camId = `c${camIdx}`;
        let name = `Campaign ${camId}`;

        let advName = `adv${Math.floor(Math.random()*10000) % numOrgsPerType}`;
        let pubName = `pub${Math.floor(Math.random()*10000) % numOrgsPerType}`;

        let currentVerifierUrls = [];

        for (let i = 0; i < numVerifiers; i++) {
            currentVerifierUrls.push(verifierUrls.pop());
        }

        campaigns.push({
            id: camId,
            name,
            advName,
            pubName,
            verifierURLsStr: currentVerifierUrls.join(";"),
            startTime: startTime,
            endTime: endTime,
            deviceIdsStr
        });
    }

    return campaigns;
}

// module.exports = {
//     ChaincodeCaller,
// }