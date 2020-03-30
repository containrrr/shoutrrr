# Services overview

Currently shoutrrr supports the following services:

| Service                                      | URL format   |
| -------------------------------------------- | ------------ |
| [Discord](/services/discord)                 | `discord://<token>@<channel>`    |
| [Telegram](/services/telegram)               | `telegram://<token>@telegram?channels=<channel-1>[,<channel-2>,...]`   |
| [Pushover](/services/pushover)               | `pushover://<token>/<user>/<device>`   |
| [Slack](/services/slack)                     | `slack://[<botname>@]<token-a>/<token-b>/<token-c>`      |
| [Email](/services/email)                     | `smtp://<username>:<password>@<host>:<port>/?fromAddress=<fromAddress>&toAddresses=<recipient1>[,<recipient2>,...]`       |
| [Microsoft Teams](/services/microsoft-teams) | `teams://<token-a>/<token-b>/<token-c>`      |
| [Gotify](/services/gotify)                   | `gotify://<gotify-host>/<token>`     |
| [Pushbullet](/services/pushbullet)           | `pushbullet://api-token[/<device>/#<channel>/<email>]` |
| [IFTTT](/services/IFTTT)                     | `ifttt://<key>/?events=<event1>[,<event2>,...]&value1=<value1>&value2=<value2>&value3=<value3>`      |
| [Mattermost](/services/mattermost)           | `mattermost://<mattermost-host>/<token>[/<username>/<channel>]` |