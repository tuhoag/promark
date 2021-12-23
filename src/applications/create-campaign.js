'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');
const { Wallets, Gateway } = require('fabric-network')

const mspId = 'av0MSP'
const cryptoPath = path.resolve(__dirname, '..', '..', 'organizations', 'peerOrganizations', 'adv0.promark.com');
const certPath = path.resolve(cryptoPath, 'users', 'Admin@adv0.promark.com', 'msp', 'signcerts', 'Admin@org1.example.com-cert.pem');
const keyPath = path.resolve(cryptoPath, 'users', 'Admin@adv0.promark.com', 'msp', 'keystore', 'key.pem');
const tlsCertPath = path.resolve(cryptoPath, 'peers', 'peer0.adv0.promark.com', 'tls', 'ca.crt');
const peerEndpoint = 'localhost:5000'

// Main program function
async function main() {
    // A wallet stores a collection of identities for use
    const wallet = await Wallets.newFileSystemWallet('./identity/user/ken/wallet');

    // A gateway defines the peers used to access Fabric networks
    const gateway = new Gateway();

    // Main try/catch block
    try {

    //     // Specify userName for network access
    //     // const userName = 'isabella.issuer@magnetocorp.com';
    //     const userName = 'ken';

    //     // Load connection profile; will be used to locate a gateway
        let connectionProfile = yaml.safeLoad(fs.readFileSync('connection-profile.yaml', 'utf8'));

        const certPath = "../../organizations/peerOrganizations/adv0.promark.com/users/Admin@adv0.promark.com/msp/signcerts/Admin@adv0.promark.com-cert.pem"
        const keyPath = "../../organizations/peerOrganizations/adv0.promark.com/users/Admin@adv0.promark.com/msp/keystore/priv_sk"

        const cert = fs.readFileSync(certPath).toString();
        const key = fs.readFileSync(keyPath).toString();

        const identityLabel = 'Admin@adv0.promark.com';
        const identity = {
            credentials: {
                certificate: cert,
                privateKey: key,
            },
            mspId: 'adv0MSP',
            type: 'X.509',
        };

        await wallet.put(identityLabel, identity);

    //     // Set connection options; identity and wallet
        let connectionOptions = {
            identity: "Admin@adv0.promark.com",
            wallet: wallet,
            discovery: { enabled: false, asLocalhost: false }
        };

        // Connect to gateway using application specified parameters
        console.log('Connect to Fabric gateway.');

        await gateway.connect(connectionProfile, connectionOptions);

        // Access PaperNet network
        console.log('Use network channel: mychannel.');

        const network = await gateway.getNetwork('mychannel');

        // Get addressability to commercial paper contract
        console.log('Use campaign.promark smart contract.');

        const contract = await network.getContract('campaign');

        // issue commercial paper
        console.log('Submit campaign issue transaction.');

        const issueResponse = await contract.submitTransaction("CreateCampaign", "c001121","campaign001","Adv0","Bus0","peer0.adv0.promark.com:5000;peer0.bus0.promark.com:5000");

        // process response
        console.log("response: " + issueResponse)
        console.log('Process issue transaction response.'+issueResponse);

    // //     // let paper = CommercialPaper.fromBuffer(issueResponse);

    // //     // console.log(`${paper.issuer} commercial paper : ${paper.paperNumber} successfully issued for value ${paper.faceValue}`);
    // //     console.log('Transaction complete.');

    } catch (error) {

        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);

    } finally {

        // Disconnect from the gateway
        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();

    }
}
main().then(() => {

    console.log('Issue program complete.');

}).catch((e) => {

    console.log('Issue program exception.');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);

});