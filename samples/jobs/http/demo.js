#!/usr/bin/env node

'use strict';

var express = require('express');
var app = express();

// The /test endpoint populates a flow by replacing
// one of its properties (`value` in this case.)
//
// Note that only those top-level properties of the
// flow that need to be changed are output—the API
// will automatically leave the other ones alone.
app.get('/test', function(req, res) {
	res.send('{"value":123}');
});

// The /error endpoint simulates a script error by
// returning a 400 status. The Agent will capture
// this and send it on to the API, which will display
// it in your account's log.
app.get('/error', function(req, res) {
	res.status(400).send('Caterpillar drive offline.');
});

app.listen(8000);