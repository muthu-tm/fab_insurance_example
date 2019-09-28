const express = require('express'),
    bodyParser = require('body-parser'),

    insurerRoutes = require('./routes/insurer-router.js'),
    customerRoutes = require('./routes/customer-router.js'),
    constants = require('./utils/constants.js'),
    userRoutes = require('./routes/user-router.js');

var host = constants.HOSTNAME,
    port = constants.PORT;

// Instantiate the app
var app = express();

// Support parsing of application/json type POST data
app.use(bodyParser.json());
// Support parsing of application/x-www-form-urlencoded POST data
app.use(bodyParser.urlencoded({
    extended: false
}));

app.use('/insurer', insurerRoutes);
app.use('/customer', customerRoutes);
app.use('/user', userRoutes)

// Start the server and listen on port 
app.listen(port, function () {
    console.log('*************** SERVER STARTED: http://%s:%s ******************', host, port);
});