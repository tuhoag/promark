'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');
const { Wallets, Gateway } = require('fabric-network');

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

    return gateway
}

const createCampaign = async (network, camId, advName, busName, verifierAddresses) => {
    // Get addressability to commercial paper contract
    console.log('Use campaign.promark smart contract.');
    const contract = await network.getContract('campaign');

    // issue commercial paper
    console.log('Submit campaign issue transaction.');
    const verifierAddressesStr = verifierAddresses.join(";")
    const issueResponse = await contract.submitTransaction("CreateCampaign", camId, camId, advName, busName, verifierAddressesStr);

    // process response
    console.log("response: " + issueResponse);
    console.log('Process issue transaction response.' + issueResponse);

    return issueResponse
}

const getId = (maxNum) => {
    return Math.floor(Math.random() * 100) % maxNum;
}

const findAdvertiser = (numAdvs) => {
    return "adv" + getId(numAdvs);
}

const findBusiness = (numBuses) => {
    return "bus" + getId(numBuses);
}

const findVerifierAddresses = (orgName, numPeers, numVerifiers) => {
    // select peer
    // peer0.adv0.promark.com:5000
    const peerName = "peer" + getId(numPeers);

    if (numVerifiers > numPeers) {
        throw `numVerifiers (${numVerifiers}$) is higher than numPeers (${numPeers}).`;
    }

    let verifierAddresses = [];

    for (let i = 0; i < numVerifiers; i++) {
        const address = `${peerName}.${orgName}.promark.com:5000`;
        verifierAddresses.push(address);
    }

    return verifierAddresses
}

// Main program function
const main = async (argv) => {
    const subCommand = argv[0];
    const numOrgsPerType = argv[1];
    const numPeersPerOrg = argv[2];
    const numVerifiers = argv[3];

    console.log(argv);
    console.log(subCommand);
    console.log(numOrgsPerType);
    console.log(numPeersPerOrg);
    console.log(numVerifiers);

    const userName = "ken";
    const orgName = "adv0";
    const orgUserName = "User1";
    const channelName = "mychannel";

    const camId = "c" + getId(10000);

    const advName = findAdvertiser(numOrgsPerType);
    const busName = findBusiness(numOrgsPerType);
    const verifierAddresses = findVerifierAddresses(advName, numPeersPerOrg, numVerifiers);
    verifierAddresses.push(...findVerifierAddresses(busName, numPeersPerOrg, numVerifiers));

    let gateway;
    try {
        gateway = await connectToGateway(userName, orgName, orgUserName);
        // Access PaperNet network
        console.log('Use network channel: ' + channelName);
        const network = await gateway.getNetwork(channelName);

        const campaignJSON = await createCampaign(network, camId, advName, busName, verifierAddresses);

        console.log("campaignJSON: " + campaignJSON);

        const campaign = JSON.parse(campaignJSON);

        console.log(`campaign.Id: ${campaign.Id} - Name: ${campaign.Name} - Adv: ${campaign.Advertiser} - Bus: ${campaign.Business} - Verifiers: ${campaign.VerifierURLs}`)
        console.log('Transaction complete.');
    } catch (error) {
        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);
    } finally {
        // Disconnect from the gateway
        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();
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