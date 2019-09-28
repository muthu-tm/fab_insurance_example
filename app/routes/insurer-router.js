const express = require('express'),
    routes = express.Router(),
    query = require('../src/query-chaincode.js'),
    invoke = require('../src/invoke-chaincode.js'),
    constants = require('./../utils/constants.js'),
    responseMessage = require('../utils/createResponseJson.js');

var channelName = constants.CHANNEL_NAME,
    chaincodeId = constants.CHAINCODE_NAME;

// Create new insurance policy for a customer
routes.post('/insurancePolicy', async function (req, res) {
    var userID = req.body.userID,
        customerID = req.body.customerID,
        customerName = req.body.customerName,
        insurerID = req.body.insurerID,
        insurerName = req.body.insurerName,
        policyID = req.body.policyID,
        appliedDate = req.body.appliedDate,
        expiryDate = req.body.expiryDate,
        policyType = req.body.policyType,
        amountVal = req.body.amountVal,
        payType = req.body.payType,
        datePaid = req.body.datePaid;

    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!customerID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'customerID'));
        return;
    }
    if (!customerName) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'customerName'));
        return;
    }
    if (!insurerID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'insurerID'));
        return;
    }
    if (!insurerName) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'insurerName'));
        return;
    }
    if (!policyID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'policyID'));
        return;
    }
    if (!appliedDate) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'appliedDate'));
        return;
    }
    if (!expiryDate) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'expiryDate'));
        return;
    }
    if (!policyType) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'policyType'));
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

    if (amountVal == undefined || Math.sign(amountVal) === 0) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1003));
        return;
    }

    if (Math.sign(amountVal) < 0) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1015));
        return;
    }

    let args = [customerID, customerName, insurerID, insurerName, policyID, appliedDate, expiryDate, policyType, amountVal.toString(), payType, datePaid]
    let response = await invoke.invokeChaincode(channelName, chaincodeId,
        constants.CREATE_NEW_POLICY, args, userID);

    res.status(200).json(response)
});

// Get Insurance policy by unique policyID
routes.get('/policy', async function (req, res) {
    var userID = req.query.userID,
        insurerID = req.query.insurerID,
        policyID = req.query.policyID;


    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!insurerID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'insurerID'));
        return;
    }
    if (!policyID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'policyID'));
        return;
    }

    var args = [insurerID, policyID]
    var response = await query.queryChaincode(channelName, chaincodeId,
        constants.QUERY_POLICY_BY_ID, args, userID);

    res.status(200).json(response)
});

// Get all policies for the given insureID
routes.get('/policies', async function (req, res) {
    var userID = req.query.userID,
        insurerID = req.query.insurerID;

    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!insurerID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'insurerID'));
        return;
    }

    var args = [insurerID]
    var response = await query.queryChaincode(channelName, chaincodeId,
        constants.QUERY_ALL_POLICIES, args, userID);

    res.status(200).json(response)
});


// Retrieve a particular customer's details using customerID
routes.get('/customer', async function (req, res) {
    var userID = req.body.userID,
    customerID = req.body.customerID;

    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!customerID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'customerID'));
        return;
    }

    var args = [customerID];
    var response = await query.queryChaincode(channelName, chaincodeId,
        constants.GET_CUSTOMER, args, userID);

    res.status(200).json(response)
});

// Transfer the policy ownership to customer
routes.put('/transferPolicy', async function (req, res) {
    var userID = req.body.userID,
    customerID = req.body.customerID,
    insurerID = req.body.insurerID,
    policyID = req.body.policyID,
    transferOwnership = req.body.transferOwnership;

    if (!userID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!customerID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'customerID'));
        return;
    }
    if (!insurerID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'insurerID'));
        return;
    }
    if (!policyID) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'policyID'));
        return;
    }
    if (!transferOwnership) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'transferOwnership'));
        return;
    }

    var args = [customerID, insurerID, policyID, transferOwnership.toString()];
    var response = await invoke.invokeChaincode(channelName, chaincodeId,
        constants.TRANSFER_POLICY, args, userID);

    res.status(200).json(response)
});


module.exports = routes