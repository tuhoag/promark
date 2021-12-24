const { Wallets, Gateway } = require('fabric-network')
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');

const buildCPP = async () => {
    return yaml.safeLoad(fs.readFileSync('connection-profile.yaml', 'utf8'));
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

    const cryptoPath = path.resolve(__dirname, '..', '..', 'organizations', 'peerOrganizations', orgFullName, "users", orgUserFullName, "msp");
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
    let connectionProfile = await buildCPP();

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


module.exports = {
    connectToGateway: connectToGateway,
}