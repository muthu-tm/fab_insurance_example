const HOSTNAME = "localhost",
    PORT = 3000

const CHANNEL_NAME = "insuranceChannel",
CHAINCODE_NAME = "insurance-network"

// ROLE_ATTRIBUTE - identity attribute name
const ROLE_ATTRIBUTE = "role"
// CUSTOMER_ROLE - identity attribute value
const CUSTOMER_ROLE = "Customer"
// INSURER_ROLE - identity attribute value
const INSURER_ROLE = "Insurer"
// CREATE_NEW_POLICY - Creates new insurance policy
const CREATE_NEW_POLICY = "createNewPolicy"
// QUERY_POLICY_BY_ID - query a specific policy from ledger
const QUERY_POLICY_BY_ID = "getPolicyByID"
// QUERY_ALL_POLICIES - query all insurance policies
const QUERY_ALL_POLICIES = "getAllPolicies"
// GET_CUSTOMER - query and retriev customer information
const GET_CUSTOMER = "getCustomer"
// TRANSFER_POLICY - tranfer a policy to customer
const TRANSFER_POLICY = "transferPolicy"
// GET_POLICY_PAYMENTS - get payments of a insurance policy
const GET_POLICY_PAYMENTS = "getPolicyPayments"
// RENEW_POLICY - renew the insurance policy
const RENEW_POLICY = "renewPolicy"

const NETWORK_CONFIG_FILE= "network/connection.json",
WALLET_PATH="wallet",
ORG_ADMINID="admin",
ORG_ADMIN_SECRET="adminpw",
CA_NAME="ca.insurer.example.com",
ORG_MSP="InsurerMSP"


module.exports = {
    HOSTNAME,
    PORT,
    CHANNEL_NAME,
    CHAINCODE_NAME,
    ROLE_ATTRIBUTE,
    CUSTOMER_ROLE,
    INSURER_ROLE,
    CREATE_NEW_POLICY,
    QUERY_POLICY_BY_ID,
    QUERY_ALL_POLICIES,
    GET_CUSTOMER,
    TRANSFER_POLICY,
    GET_POLICY_PAYMENTS,
    RENEW_POLICY,
    NETWORK_CONFIG_FILE,
    WALLET_PATH,
    ORG_ADMINID,
    ORG_ADMIN_SECRET,
    CA_NAME,
    ORG_MSP
}