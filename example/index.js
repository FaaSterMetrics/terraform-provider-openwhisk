/*
 * http://localhost:3233/api/v1/web/guest/faastermetrics/ENDPOINT_NAME
*/

const fs = require('fs'); 


function main(args) {
    process.env = {...process.env, ...args}; 
    var blubb = fs.readFileSync("file.txt", "utf8");
    return {body: blubb + process.env.testvar}

}

exports.main = main;
