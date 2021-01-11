# Services overview

Click on the service for a more thorough explanation.

| Service                           | URL format                                                                                                                                      |
| --------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| [Discord](./discord.md)           | *discord://__`token`__
@__`id`__*                                                                                                                |
| [Telegram](./telegram.md)         | *telegram://__`token`__
@telegram?channels=__`channel-1`__[,__`channel-2`__,...]*                                                                |
| [Pushover](./pushover.md)         | *pushover://shoutrrr:__`apiToken`__@__`userKey`__/?devices=__`device1`__[,__`device2`__, ...]*                                                  |
| [Slack](./not-documented.md)      | *slack://[__`botname`__@]__`token-a`__/__`token-b`__/__`token-c`__*                                                                             |
| [Email](./not-documented.md)      | *smtp://__`username`__:__`password`__@__`host`__:__`port`__/?fromAddress=__`fromAddress`__&toAddresses=__`recipient1`__[,__`recipient2`__,...]* |
| [Microsoft Teams](./teams.md)     | *teams://__`token-a`__/__`token-b`__/__`token-c`__*                                                                                             |
| [Gotify](./not-documented.md)     | *gotify://__`gotify-host`__/__`token`__*                                                                                                        |
| [Pushbullet](./not-documented.md) | *pushbullet://__`api-token`__[/__`device`__/#__`channel`__/__`email`__]*                                                                        |
| [IFTTT](./not-documented.md)      | *ifttt://__`key`__/?events=__`event1`__[,__`event2`__,...]&value1=__`value1`__&value2=__`value2`__&value3=__`value3`__*                         |
| [Mattermost](./not-documented.md) | *mattermost://[__`username`__@]__`mattermost-host`__
/__`token`__[/__`channel`__]*                                                               |
| [Matrix](./matrix.md)             | *matrix://__`username`__:__`password`__@__`host`__:__`port`__/[
?rooms=__`!roomID1`__[,__`roomAlias2`__]]*                                        |
| [Hangouts Chat](./hangouts.md)    | *hangouts:
//chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz*                                                                       |
| [Zulip Chat](./zulip.md)          | *zulip://__`bot-mail`__:__`bot-key`__@__`zulip-domain`__/?stream=__`name-or-id`__&topic=__`name`__*                                             |
| [Join](./not-documented.md)       | *join://shoutrrr:__`api-key`__
@join/?devices=__`device1`__[,__`device2`__, ...][&icon=__`icon`__][&title=__`title`__]*                          |
| [Rocketchat](./rocketchat.md)     | *rocketchat://[__`username`__@]__`rocketchat-host`__/__`token`__[/__`channel`
&#124;`@recipient`__]*                                             |
