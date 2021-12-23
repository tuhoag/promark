const { Wallets, Gateway } = require('fabric-network')
const fs = require('fs');
const yaml = require('js-yaml');

exports.createWallet = async function(userName, orgsMemberName, orgsName) {
    const walletPath = './wallet/user/' + userName
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    const orgsFullName = orgsName + ".promark.com"
    const orgsMemberFullName = orgsMemberName + "@" + orgsFullName

    const credPath = "../../organizations/peerOrganizations/" + orgsFullName

    const certPath = credPath + "/users/" + orgsMemberFullName + "/msp/signcerts/" + orgsMemberFullName + "-cert.pem"
    const keyPath = credPath + "/users/" + orgsMemberFullName + "/msp/keystore/priv_sk"

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

    await wallet.put(orgsMemberFullName, identity);

    return wallet
}

exports.connect = async function() {
    const wallet = exports.createWallet("ken", "Admin", "adv0")

    // const wallet = await Wallets.newFileSystemWallet('./identity/user/ken/wallet');

    // A gateway defines the peers used to access Fabric networks
    const gateway = new Gateway();

    let connectionProfile = yaml.safeLoad(fs.readFileSync('connection-profile.yaml', 'utf8'));

    // Set connection options; identity and wallet
    let connectionOptions = {
        identity: "Admin@adv0.promark.com",
        wallet: wallet,
        discovery: { enabled: false, asLocalhost: false }
    };

    // Connect to gateway using application specified parameters
    console.log('Connect to Fabric gateway.');

    await gateway.connect(connectionProfile, connectionOptions);

    return gateway
}