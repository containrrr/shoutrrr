# Services overview

Click on the service for a more thorough explanation. <!-- @formatter:off -->

| Service                           | URL format                                                                                                                                      |
| --------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| [Bark](./bark.md)                 | *bark://__`devicekey`__@__`host`__*                                                                                                             |
| [Discord](./discord.md)           | *discord://__`token`__@__`id`__*                                                                                                                |
| [Email](./email.md)               | *smtp://__`username`__:__`password`__@__`host`__:__`port`__/?from=__`fromAddress`__&to=__`recipient1`__[,__`recipient2`__,...]*                 |
| [Gotify](./gotify.md)             | *gotify://__`gotify-host`__/__`token`__*                                                                                                        |
| [Google Chat](./googlechat.md)    | *googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz*                                                                     |
| [IFTTT](./ifttt.md)               | *ifttt://__`key`__/?events=__`event1`__[,__`event2`__,...]&value1=__`value1`__&value2=__`value2`__&value3=__`value3`__*                         |
| [Join](./join.md)                 | *join://shoutrrr:__`api-key`__@join/?devices=__`device1`__[,__`device2`__, ...][&icon=__`icon`__][&title=__`title`__]*                          |
| [Mattermost](./mattermost.md)     | *mattermost://[__`username`__@]__`mattermost-host`__/__`token`__[/__`channel`__]*                                                               |
| [Matrix](./matrix.md)             | *matrix://__`username`__:__`password`__@__`host`__:__`port`__/[?rooms=__`!roomID1`__[,__`roomAlias2`__]]*                                       |
| [Ntfy](./ntfy.md)                 | *ntfy://__`username`__:__`password`__@ntfy.sh/__`topic`__*                                                                                      |
| [OpsGenie](./opsgenie.md)         | *opsgenie://__`host`__/token?responders=__`responder1`__[,__`responder2`__]*                                                                    |
| [Pushbullet](./pushbullet.md)     | *pushbullet://__`api-token`__[/__`device`__/#__`channel`__/__`email`__]*                                                                        |
| [Pushover](./pushover.md)         | *pushover://shoutrrr:__`apiToken`__@__`userKey`__/?devices=__`device1`__[,__`device2`__, ...]*                                                  |
| [Rocketchat](./rocketchat.md)     | *rocketchat://[__`username`__@]__`rocketchat-host`__/__`token`__[/__`channel`&#124;`@recipient`__]*                                             |
| [Slack](./slack.md)               | *slack://[__`botname`__@]__`token-a`__/__`token-b`__/__`token-c`__*                                                                             |
| [Teams](./teams.md)               | *teams://__`group`__@__`tenant`__/__`altId`__/__`groupOwner`__?host=__`organization`__.webhook.office.com*                                      |
| [Telegram](./telegram.md)         | *telegram://__`token`__@telegram?chats=__`@channel-1`__[,__`chat-id-1`__,...]*                                                                  |
| [Zulip Chat](./zulip.md)          | *zulip://__`bot-mail`__:__`bot-key`__@__`zulip-domain`__/?stream=__`name-or-id`__&topic=__`name`__*                                             |

## Specialized services

| Service                           | Description                                                                                                                                     |
| --------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| [Logger](./logger.md)             | Writes notification to a configured go `log.Logger`                                                                                             |
| [Generic Webhook](./generic.md)   | *generic+https://`url`* to send notifications directly to a webhook                                                                                                   |

