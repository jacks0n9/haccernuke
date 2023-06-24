# Haccernuke nuker
## Welcome to haccernuke, a nuker written in golang.

Haccernuke can be self botted or used as a normal bot (more information in config).

What does haccernuke do exactly?
A nuker is something that is meant to completely obliterate a server and different nukers do it in different ways
Haccernuke does this with the following features:
- Ban all members
- Delete all roles
- Spam create channels
- Spam messages in these channels
- Grant admin to other users. (see config for limitations)
- Delete emojis (discordgo doesn't have sticker support yet sorry)
- Whatever you want it to do. Put in a feature request in the issues tab and if it's good, I'll add it
## Setup for haccernuke:
- Rename nukeconf_example.toml to nukeconf.toml
- Open nukeconf.toml
- Put your token into the token field. If it is a user token just put it into the quotes. If it is a bot token, prepend "Bot " before it like this: `Bot OTQ0NTM5NDk4OTc3MjUwODQ...`
- Go into config, disable or enable features. Some features are in their own category because they are more configurable. You can also configure parameters of these functions.
- Run the bot by doing `go run *.go` in the main directory


