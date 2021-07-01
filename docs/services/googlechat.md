# Google Chat

## URL Format

Your Google Chat Incoming Webhook URL will look like this:

!!! info ""
    https://chat.googleapis.com/v1/spaces/__`FOO`__/messages?key=__`bar`__&token=__`baz`__

The shoutrrr service URL should look like this:

!!! info ""
    googlechat://chat.googleapis.com/v1/spaces/__`FOO`__/messages?key=__`bar`__&token=__`baz`__

In other words the incoming webhook URL with `https` replaced by `googlechat`.

Google Chat was previously known as Hangouts Chat. Using `hangouts` in
the service URL instead `googlechat` is still supported, although
deprecated.

## Creating an incoming webhook in Google Chat

1. Open the room you would like to add Shoutrrr to and open the chat
room menu.
![Screenshot 1](googlechat/hangouts-1.png)

2. Then click on *Configure webhooks*.
![Screenshot 2](googlechat/hangouts-2.png)

3. Name the webhook and save.
![Screenshot 3](googlechat/hangouts-3.png)

4. Copy the URL.
![Screenshot 4](googlechat/hangouts-4.png)


5. Format the service URL by replacing `https` with `googlechat`.
