# OMyBot

**Abstract**

A Bot handles requests from Discord and replies seamlessly.

**Prerequisite**

Setup should be ready by these steps

* Navigate to [Discord Application](https://discordapp.com/developers/applications/me/create), create an App

* Under Bot Section, create a Bot User

* Make it a public Bot

* Generate OAuth2 URL and choose "Send Message", then select a channel, click Authorize

**Usage**

In project root directory, build a docker image <code>docker build -t omybot .</code>

Open Discord with the selected channel, send "!quote 0005.HK" and see the magic
