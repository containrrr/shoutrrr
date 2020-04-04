# Hangouts Chat

## URL format

Your Hangouts Chat Incoming Webhook URL will look like this:

> https://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz

The shoutrrr service URL should look like this:

> hangouts://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz

In other words the incoming webhook URL with `https` replaced by `hangouts`.

## Creating an incoming webhook in Hangouts Chat

1. Open the room you would like to add Shoutrrr to and open the chat
room menu.
![Screenshot 1](hangouts/hangouts-1.png)

2. Then click on *Configure webhooks*.
![Screenshot 2](hangouts/hangouts-2.png)

3. Name the webhook and save.
![Screenshot 3](hangouts/hangouts-3.png)

4. Copy the URL.
![Screenshot 4](hangouts/hangouts-4.png)


5. Format the service URL by replacing `https` with `hangouts`.
