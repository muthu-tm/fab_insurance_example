const FabricCAServices = require('fabric-ca-client'),
	{ FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network'),
	fs = require('fs'),
	path = require('path'),
	constants = require('./constants.js');

const ccpPath = path.resolve(process.cwd(), '../', constants.NETWORK_CONFIG_FILE),
	ccpJSON = fs.readFileSync(ccpPath, 'utf8'),
	ccp = JSON.parse(ccpJSON);

async function getGateway(userID) {
	try {
		// Create a new file system based wallet for managing identities.
		var walletPath = path.join(process.cwd(), '../', constants.WALLET_PATH);
		const wallet = new FileSystemWallet(walletPath);
		console.log(`Wallet path: ${walletPath}`);

		const gateway = new Gateway();
		if (userID) {
			var userExists = await wallet.exists(userID);
			if (!userExists) {
				console.error(`An identity for the user "${userID}" does not exist in the wallet`);
				throw new Error(`User "${userID}" not found in the network`)
			}
			// Create a new gateway for connecting to our peer node.
			await gateway.connect(ccp, { wallet, identity: userID, discovery: { enabled: false, asLocalhost: false } });
		} else {
			var adminExist = await wallet.exists(constants.ORG_ADMINID);
			if (!adminExist) {
				console.error(`Admin identity "${constants.ORG_ADMINID}" does not exist in the wallet`);
				var caURL = ccp.certificateAuthorities[constants.CA_NAME].url;
				const ca = new FabricCAServices(caURL);

				// Enroll the admin user, and import the new identity into the wallet.
				var enrollment = await ca.enroll({ enrollmentID: constants.ORG_ADMINID, enrollmentSecret: constants.ORG_ADMIN_SECRET });
				var identity = X509WalletMixin.createIdentity(constants.ORG_MSP, enrollment.certificate, enrollment.key.toBytes());
				await wallet.import(constants.ORG_ADMINID, identity);
			}

			await gateway.connect(ccp, { wallet, identity: constants.ORG_ADMINID, discovery: { enabled: false, asLocalhost: false } });
		}

		return gateway;
	} catch (error) {
		console.error(`Failed to submit transaction: ${error}`);
		throw new Error(error.message);
	}
}

exports.getGateway = getGateway;
