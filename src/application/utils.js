const { Wallets, Gateway } = require('fabric-network');
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');

const setting = require('./setting');

const buildCPP = async (numOrgsPerType, numPeersPerOrg) => {
    const connectionProfileName = `connectionProfile-${numOrgsPerType}-${numPeersPerOrg}.yaml`;
    return yaml.safeLoad(fs.readFileSync(connectionProfileName, 'utf8'));
}

const buildWallet = async (userName) => {
    // Create a new  wallet : Note that wallet is for managing identities.
    const walletPath = './wallet/' + userName
    console.log(`Built a file system wallet at ${walletPath}`);
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

const connectToGateway = async (userName, orgName, orgUserName) => {
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
    console.log('Connect to Fabric gateway.');

    await gateway.connect(connectionProfile, connectionOptions);

    return gateway;
}

const connectToNetwork = async (userName, orgName, orgUserName) => {
    const gateway = connectToGateway(userName, orgName, orgUserName);
    console.log('Use network channel: ' + setting.channelName);
    return gateway.getNetwork(setting.channelName);
}


class ChaincodeCaller {
    constructor() {
        this.gateway = null;
        this.network = null;
    }

    static async init() {
        let caller = new ChaincodeCaller();
        caller.gateway = await connectToGateway(setting.userName, setting.orgName, setting.orgUserName);
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

const callChaincodeFn = async (requestFn, responseFn) => {
    let gateway;
    try {
        gateway = await connectToGateway(setting.userName, setting.orgName, setting.orgUserName);
        const network = await gateway.getNetwork(setting.channelName);
        console.log('Use campaign.promark smart contract.');
        const response = await requestFn(network);
        return responseFn(response);
    } catch (error) {
        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);
    } finally {
        // Disconnect from the gateway
        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();
    }
}

const getId = (maxNum) => {
    return Math.floor(Math.random() * 100) % maxNum;
}

module.exports = {
    connectToGateway,
    callChaincodeFn,
    ChaincodeCaller,
    getId,
}