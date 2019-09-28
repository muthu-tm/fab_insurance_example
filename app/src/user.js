const FabricCAServices = require('fabric-ca-client'),
    { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network'),
    path = require('path'),
    fs = require('fs'),
    constants = require('./../utils/constants.js'),
    responseMessage = require('./../utils/createResponseJson.js');

const ccpPath = path.resolve(process.cwd(), '../', constants.NETWORK_CONFIG_FILE),
    ccpJSON = fs.readFileSync(ccpPath, 'utf8'),
    ccp = JSON.parse(ccpJSON);

async function registerIdentity(username, userRole) {
    try {
        // Create a new file system based wallet for managing identities.
        var walletPath = path.join(process.cwd(), '../', constants.WALLET_PATH);
        const wallet = new FileSystemWallet(walletPath);

        // Check to see if we've already enrolled the user.
        var userExists = await wallet.exists(username);
        if (userExists) {
            console.error(`An identity for the user ${username} already exists in the wallet`);
            return responseMessage.getApiResponse(1007);
        }

        var adminExists = await wallet.exists(constants.ORG_ADMINID);
        if (!adminExists) {
            console.log(`An identity for the admin user ${constants.ORG_ADMINID} does not exist in the wallet`);
            // Create a new CA client for interacting with the CA.
            var caURL = ccp.certificateAuthorities[constants.CA_NAME].url;
            const ca = new FabricCAServices(caURL);

            // Enroll the admin user, and import the new identity into the wallet.
            var enrollment = await ca.enroll({ enrollmentID: constants.ORG_ADMINID, enrollmentSecret: constants.ORG_ADMIN_SECRET });
            var identity = X509WalletMixin.createIdentity(constants.ORG_MSP, enrollment.certificate, enrollment.key.toBytes());
            await wallet.import(constants.ORG_ADMINID, identity);
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: constants.ORG_ADMINID });

        // Get the CA client object from the gateway for interacting with the CA.
        var ca = gateway.getClient().getCertificateAuthority();
        var adminIdentity = gateway.getCurrentIdentity();

        var keyValueAttribute = {
            name: "role",
            value: userRole,
            ecert: true
        }

        // Register the user, enroll the user, and import the new identity into the wallet.
        var secret = await ca.register({ enrollmentID: username, attrs: [keyValueAttribute] }, adminIdentity);
        var userEnrollment = await ca.enroll({ enrollmentID: username, enrollmentSecret: secret });
        var userIdentity = X509WalletMixin.createIdentity(constants.ORG_MSP, userEnrollment.certificate, userEnrollment.key.toBytes());
        await wallet.import(username, userIdentity);
        // await gateway.getClient().setUserContext({username: username, password: secret});

        // Disconnect from the gateway.
        gateway.disconnect();

        return responseMessage.getApiResponse(0000)
    } catch (error) {
        console.error('Failed to register user; Error: %s', error.toString());
        return responseMessage.getApiResByMessage(1006, error.toString())
    }
}

exports.registerIdentity = registerIdentity;