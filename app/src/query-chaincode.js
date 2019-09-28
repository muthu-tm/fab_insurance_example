const helper = require('../utils/helper.js'),
	constants = require('../utils/constants.js'),
	responseMessage = require('../utils/createResponseJson.js');

async function queryChaincode(channelName, chaincodeName, fcn, args, userName) {
	console.debug('\n============ Query transaction on channel %s ============\n', channelName);
	try {
		var gateway = await helper.getGateway(userName);

		// Get the network (channel) our contract is deployed to.
		var network = await gateway.getNetwork(channelName);

		// Get the contract (chaincode) from the network.
		var contract = network.getContract(chaincodeName);

		var transaction = contract.createTransaction(fcn);
		var txId = transaction.getTransactionID();

		var queryResponses;
		switch (fcn) {
			case constants.QUERY_POLICY_BY_ID:
			case constants.QUERY_ALL_POLICIES:
			case constants.GET_CUSTOMER:
			case constants.GET_POLICY_PAYMENTS:
			case constants.QUERY_CUSTOMER_POLICY:
				queryResponses = await transaction.evaluate(args[0]);
				break;
		}

		// Disconnect from the gateway.
		gateway.disconnect();

		if (queryResponses) {
			if (queryResponses[0] instanceof Error) {
				console.error('Query Failed. REASON - ', queryResponses[0].toString());
				return responseMessage.getApiResByMessage(1013, queryResponses[0].toString())
			} else {
				console.info('Payload:' + queryResponses.toString());
				return responseMessage.getApiResponse(0000, fcn, JSON.parse(queryResponses.toString()), "QUERY", txId._transaction_id)
			}
		} else {
			console.error('queryResponses is null');
			return responseMessage.getApiResponse(1008);
		}
	} catch (error) {
		console.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return responseMessage.getApiResByMessage(1005, error.toString());
	}
}

exports.queryChaincode = queryChaincode;
