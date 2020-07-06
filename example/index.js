const lib = require('@faastermetrics/lib')
const path = require('path')
const fs = require('fs')

module.exports = lib.serverless.rpcHandler( event => ({
  ok: true,
  file: fs.readFileSync(path.join(__dirname, 'file.txt'), 'utf8'),
  event
}))
