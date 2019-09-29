const express = require('express'),
    routes = express.Router(),
    query = require('../src/query-chaincode.js'),
    invoke = require('../src/invoke-chaincode.js'),
    constants = require('./../utils/constants.js'),
    responseMessage = require('../utils/createResponseJson.js');

var channelName = constants.CHANNEL_NAME,
    chaincodeId = constants.CHAINCODE_NAME;

// Get all policies for the given customerID
routes.get('/policies', async function (req, res) {
    var userID = req.query.userID,
        customerID = req.query.customerID;

    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!customerID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'customerID'));
        return;
    }

    var args = [customerID]
    var response = await query.queryChaincode(channelName, chaincodeId,
        constants.QUERY_CUSTOMER_POLICY, args, userID);

    res.status(200).json(response)
});

// Get all payments for the given customerID
routes.get('/payments', async function (req, res) {
    var userID = req.query.userID,
        policyID = req.query.policyID;

    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!policyID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'policyID'));
        return;
    }

    var args = [policyID]
    var response = await query.queryChaincode(channelName, chaincodeId,
        constants.GET_POLICY_PAYMENTS, args, userID);

    res.status(200).json(response)
});

// Renew the current policy by payment
routes.post('/renewPolicy', async function (req, res) {
    var userID = req.query.userID,
        policyID = req.query.policyID,
        amountVal = req.query.amountVal,
        payType = req.query.payType,
        datePaid = req.query.datePaid,
        expiryDate = req.query.expiryDate;

    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!policyID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'policyID'));
        return;
    }
    if (!payType) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'payType'));
        return;
    }
    if (!datePaid) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'datePaid'));
        return;
    }
    if (!expiryDate) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'expiryDate'));
        return;
    }

    if (amountVal == undefined || Math.sign(amountVal) === 0) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1003));
        return;
    }

    if (Math.sign(amountVal) < 0) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1015));
        return;
    }

    var args = [policyID, amountVal.toString(), payType, datePaid, expiryDate]
    var response = await invoke.invokeChaincode(channelName, chaincodeId,
        constants.RENEW_POLICY, args, userID);

    res.status(200).json(response)
});

routes.get('/policyHistory', async function (req, res) {
    var userID = req.query.userID,
        policyID = req.query.policyID;

    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!policyID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'policyID'));
        return;
    }

    var args = [policyID]
    var response = await query.queryChaincode(channelName, chaincodeId,
        constants.QUERY_POLICY_HISTORY, args, userID);

    res.status(200).json(response)
});