var irc = require("irc"),
      fs = require("fs");

var settings = {
    userName: 'OkmaBot',
    realName: 'OkmaBot',
    password: 'oauth:fixmizy88ga63drb2yzo41g1gs7z93',
    port: 6667,
    localAddress: null,
    debug: false,
    showErrors: false,
    autoRejoin: false,
    autoConnect: true,
    channels: [],
    secure: false,
    selfSigned: false,
    certExpired: false,
    floodProtection: false,
    floodProtectionDelay: 1000,
    sasl: false,
    retryCount: 0,
    retryDelay: 2000,
    stripColors: false,
    channelPrefixes: "&#",
    messageSplit: 512,
    encoding: ''
};

var bot = new irc.Client("irc.twitch.tv", settings.userName, {
    channels: ["#OkmaLol"],
    password: settings.password,
    username: settings.userName,
    autoConnect: false,
});


/***** Begin Events ******/

// Listen for joins
bot.addListener("join", function(channel, nick, message) {

	// Welcome them in!
	bot.say(channel, "Welcome " + nick + "!");

});

// catch error event
bot.addListener("error", function(message) {

  // error logging
  console.log(message);
});

// catch user quitting
bot.addListener("quit", function (nick, reason, channels, message) {

    //
    console.log();
});


/***** End Events ******/


/*** Connect the bot ***/
bot.connect();
