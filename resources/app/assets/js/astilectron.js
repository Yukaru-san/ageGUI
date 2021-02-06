// Using "path"
const { time } = require('console');
var path = require('path');
var basePath = "";

// Wait for astilectron to be ready
document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {

        if (message.name == "generatedKeyPair") {
            console.log("received key pair: "+message.payload);
        }

    });

    // Initial data request
    var json = {
        name: "getBaseDirectory",
        payload: ""
    }

    astilectron.sendMessage(json, function(message) {
        basePath = message.payload;
    });
})
