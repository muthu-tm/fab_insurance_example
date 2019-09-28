const resCodes = require('./responseCodeConstants.js');

function getApiResponse(resCode, ccOperation, payload, type, txId) {
    var apiResponse = {
        "timeStamp": getUTCTimestampISOFormat(),
        "response": {
            "responsecode": resCode,
            "responsemessage": getResMessage(resCode),
            "chaincodeoperation": checkField(ccOperation),
            "payload": checkField(payload)
        },
        "type": checkField(type),
        "fabricTransactionID": checkField(txId)
    }

    return apiResponse;
}

function getApiResByMessage(resCode, desc, ccOperation, payload, type, txId) {
    var apiResponse = {
        "timeStamp": getUTCTimestampISOFormat(),
        "response": {
            "responsecode": resCode,
            "responsemessage": desc,
            "chaincodeoperation": checkField(ccOperation),
            "payload": checkField(payload)
        },
        "type": checkField(type),
        "fabricTransactionID": checkField(txId)
    }

    return apiResponse;
}

function getFieldErrorResponse(resCode, field) {
    var resMsg = getResMessage(resCode);
    
    var apiResponse = {
        "timeStamp": getUTCTimestampISOFormat(),
        "response": {
            "responsecode": resCode,
            "responsemessage": field + " " + resMsg,
            "chaincodeoperation": "NA",
            "payload": {}
        },
        "type": "NA",
        "fabricTransactionID": "NA"
    }

    return apiResponse;
}

//returns local timestamp in ISO format (ex, 2019-01-10T12:34:26.414Z)
function getLocalTimestampISOFormat() {
    var timeZoneOffset = (new Date()).getTimezoneOffset() * 60000, //offset in milliseconds
        localISOTime = (new Date(Date.now() - timeZoneOffset)).toISOString();

    return localISOTime;
}

//returns UTC timestamp in ISO format (ex, 2019-01-10T07:02:51.162Z)
function getUTCTimestampISOFormat() {
    return new Date(Date.now()).toISOString();
}

function checkField(field) {
    if (!field) {
        field = "NA"
    }
    return field
}

function getResMessage(code) {
    var message = resCodes.responseCodesMap.get(code);
    if (!message) {
        return ""
    }
    return message;
}

module.exports = {
    getApiResponse,
    getApiResByMessage,
    getFieldErrorResponse
}
