/**
 * responseCodesMap is used to store all the response codes with their string descriptions
 * To add new response code, update the map with responseCodesMap.set(key, value)
 */
var responseCodesMap = new Map();
responseCodesMap.set(1001, "Field is mandatory");
responseCodesMap.set(1002, "Field with wrong type");
responseCodesMap.set(1003, "Insurance Policy amount cannot be zero");
responseCodesMap.set(1004, "Query response payload was null");
responseCodesMap.set(1005, "Failed to QUERY chaincode. Error!!");
responseCodesMap.set(1006, "Failed to register the user. Error!!");
responseCodesMap.set(1007, "User was already registered in the network");
responseCodesMap.set(1008, "Failed to invoke chaincode. Proposal response null or status is not 200");
responseCodesMap.set(1009, "Query chaincode failure.");
responseCodesMap.set(1010, "Authentication failure!! check user details")
responseCodesMap.set(1011, "Error: Amount must be a positive float value")

/**
 * these codes are represents the general error/success codes with their string descriptions
 * To add new general codes, update the map in incremental way
 */
responseCodesMap.set(0000, "SUCCESS");
responseCodesMap.set(9999, "Connection timed out. Please try again later");
responseCodesMap.set(9998, "System Error. Please contact your network administrator");
responseCodesMap.set(9997, "Server Error occurred!!");

/**
 * 3000 series codes are used to represent the chaincode error codes with their string descriptions
 * To add new chaincode error code, update the map in incremental way
 */
responseCodesMap.set(3001, "required chaincode parameters missing");

exports.responseCodesMap = responseCodesMap;