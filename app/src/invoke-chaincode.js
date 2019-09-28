const helper = require('./../utils/helper.js'),
	constants = require('./../utils/constants.js'),
	responseMessage = require('./../utils/createResponseJson.js');

async function invokeChaincode(channelName, chaincodeName, fcn, args, userName) {
	console.debug('\n============ Invoke transaction on channel %s ============\n', channelName);
	try {
		var gateway = await helper.getGateway(userName);

		// Get the network (channel) our contract is deployed to.
		var network = await gateway.getNetwork(channelName);

		// Get the contract (chaincode) from the network.
		var contract = network.getContract(chaincodeName);

		var transaction = contract.createTransaction(fcn);
		var txId = transaction.getTransactionID();

		var proposalResponses;
		switch (fcn) {
			case constants.CREATE_NEW_POLICY:
				proposalResponses = await transaction.submit(args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10]);
				break;
			case constants.TRANSFER_POLICY:
				proposalResponses = await transaction.submit(args[0], args[1]);
				break;
			case constants.RENEW_POLICY:
				proposalResponses = await transaction.submit(args[0], args[1], args[2], args[3], args[4]);
				break;
		}

		console.log('Transaction has been submitted successfully', proposalResponses.toString());

		// Disconnect from the gateway.
		gateway.disconnect();

		return responseMessage.getApiResponse(0000, fcn, JSON.parse(proposalResponses.toString()), "INVOKE", txId._transaction_id)
	} catch (error) {
		console.error('Failed to invoke due to error: ' + error.stack ? error.stack : error);
		console.info('Failed to invoke due to error: ' + error.message);
		return responseMessage.getApiResByMessage(1008, error.message)
	}
}

exports.invokeChaincode = invokeChaincode;
