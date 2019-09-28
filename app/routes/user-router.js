const express = require('express'),
    user = require('./../src/user.js'),
    responseMessage = require('./../utils/createResponseJson.js'),
    routes = express.Router();

// Register a new identity
routes.post('/register', async function (req, res) {
    var userName = req.body.userID,
        userRole = req.body.role;

    if (!userName) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'userID'));
        return;
    }
    if (!userRole) {
        res.status(200).json(responseMessage.getFieldErrorResponse(1001, 'role'));
        return;
    }

    let response = await user.registerIdentity(userName, userRole);
    res.status(200).json(response)
})

module.exports = routes
