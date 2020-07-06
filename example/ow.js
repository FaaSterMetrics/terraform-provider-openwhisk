/*
 * http://localhost:3233/api/v1/web/guest/faastermetrics/ENDPOINT_NAME
*/

module.exports.main = require('@faastermetrics/lib/openwhisk')(() => require('./index.js'))
