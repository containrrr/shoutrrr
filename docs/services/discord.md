# Discord

## URL Format

Your Discord Webhook-URL will look like this:

!!! info ""
    https://discord.com/api/webhooks/__`webhookid`__/__`token`__  

The shoutrrr service URL should look like this:  

!!! info ""
    discord://__`token`__@__`webhookid`__

--8<-- "docs/services/discord/config.md"

## Optional: Sending messages to a specific thread

Discord supports sending messages via webhook to threads. You can target a specific thread by appending
`?thread_id=<id>` to the end of the url.

!!! info ""
### Example

```
discord://<token>@<webhook_id>?thread_id=1234567890123456789
```

To get the thread ID:
 - Open a Discord Thread
 - Right click -> copy link
 - OR Right click -> copy thread ID (Developer mode turned on in discord settings)

!!! warning ""
A valid thread_id must is 19 digits long. If extracting from a link don't confuse the channel ID and thread ID.

## Creating a webhook in Discord

1. Open your channel settings by first clicking on the gear icon next to the name of the channel.
![Screenshot 1](discord/sc-1.png)

2. In the menu on the left, click on *Integrations*.
![Screenshot 2](discord/sc-2.png)

3. In the menu on the right, click on *Create Webhook*.
![Screenshot 3](discord/sc-3.png)

4. Set the name, channel and icon to your liking and click the *Copy Webhook URL* button.
![Screenshot 4](discord/sc-4.png)

5. Press the *Save Changes* button.
![Screenshot 5](discord/sc-5.png)

6. Format the service URL:
```
https://discord.com/api/webhooks/693853386302554172/W3dE2OZz4C13_4z_uHfDOoC7BqTW288s-z1ykqI0iJnY_HjRqMGO8Sc7YDqvf_KVKjhJ
                                 └────────────────┘ └──────────────────────────────────────────────────────────────────┘
                                     webhook id                                    token

discord://W3dE2OZz4C13_4z_uHfDOoC7BqTW288s-z1ykqI0iJnY_HjRqMGO8Sc7YDqvf_KVKjhJ@693853386302554172
          └──────────────────────────────────────────────────────────────────┘ └────────────────┘
                                          token                                    webhook id
```
