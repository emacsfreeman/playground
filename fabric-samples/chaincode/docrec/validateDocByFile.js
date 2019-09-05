'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const { KJUR, KEYUTIL, X509 } = require('jsrsasign');
const CryptoJS = require('crypto-js');

const ccpPath = path.resolve(__dirname, '..', 'basic-network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

const caCertPath = path.resolve(__dirname, '..', 'basic-network', 'crypto-config', 'peerOrganizations', 'org1.example.com', 'ca', 'ca.org1.example.com-cert.pem');
const caCert = fs.readFileSync(caCertPath, 'utf8');


async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);

        // Collect input parameters
        // user: who initiates this query, can be anyone in the wallet
        // filename: the file to be validated
        // certfile: the cert file owner who signed the document
        const user = process.argv[2];
        const filename = process.argv[3];
        const certfile = process.argv[4];

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(user);
        if (!userExists) {
            console.log('An identity for the user ' + user + ' does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // calculate Hash from the file
        const fileLoaded = fs.readFileSync(filename, 'utf8');
        var hashToAction = CryptoJS.SHA256(fileLoaded).toString();
        console.log("Hash of the file: " + hashToAction);

        // get certificate from the certfile
        const certLoaded = fs.readFileSync(certfile, 'utf8');

        // retrieve record from ledger

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: false } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('docrec');

        // Submit the specified transaction.
        const result = await contract.evaluateTransaction('queryDocRecord', hashToAction);
        console.log("Transaction has been evaluated");
        var resultJSON = JSON.parse(result);
        console.log("Doc record found, created by " + resultJSON.time);
        console.log("");

        // Show info about certificate provided
        const certObj = new X509();
        certObj.readCertPEM(certLoaded);
        console.log("Detail of certificate provided")
        console.log("Subject: " + certObj.getSubjectString());
        console.log("Issuer (CA) Subject: " + certObj.getIssuerString());
        console.log("Valid period: " + certObj.getNotBefore() + " to " + certObj.getNotAfter());
        console.log("CA Signature validation: " + certObj.verifySignature(KEYUTIL.getKey(caCert)));
        console.log("");

        // perform signature checking
        var userPublicKey = KEYUTIL.getKey(certLoaded);
        var recover = new KJUR.crypto.Signature({"alg": "SHA256withECDSA"});
        recover.init(userPublicKey);
        recover.updateHex(hashToAction);
        var getBackSigValueHex = new Buffer(resultJSON.signature, 'base64').toString('hex');
        console.log("Signature verified with certificate provided: " + recover.verify(getBackSigValueHex));

        // perform certificate validation
        // var caPublicKey = KEYUTIL.getKey(caCert);
        // var certValidate = new KJUR.crypto.Signature({"alg": "SHA256withECDSA"});
        // certValidate.init(caPublicKey);
        // certValidate.update

        // Disconnect from the gateway.
        await gateway.disconnect();


    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

main();
