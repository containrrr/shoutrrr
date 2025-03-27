# Telegram

## URL Format

!!! info ""
    telegram://__`token`__@telegram?chats=__`channel-1`__[,__`chat-id-1`__,...]
    
--8<-- "docs/services/telegram/config.md"

## Getting a token for Telegram

Talk to [the botfather](https://core.telegram.org/bots#6-botfather).

## Identifying the target chats/channels

The `chats` param consists of one or more `Chat ID`s or `channel name`s. 

### Public Channels
The channel names can be retrieved in the telegram client in the `Channel info` section for public channels. 
Replace the `t.me/` prefix from the link with a `@`.

!!! note
    Channels names need to be prefixed by `@` to identify them as such.

!!! note
    If your channel only has an invite link (starting with `t.me/+`), you have to use it's Chat ID (see below)

!!! note
    A `message_thread_id` param ([reference](https://core.telegram.org/bots/api#sendmessage)) can be added, with the format of `$chat_id:$message_thread_id`. [More info](https://stackoverflow.com/questions/74773675/how-to-get-topic-id-for-telegram-group-chat/75178418#75178418) on how to obtain the `message_thread_id`.

### Chats
Private channels, Group chats and private chats are identified by `Chat ID`s. Unfortunatly, they are generally not visible in the
telegram clients.
The easiest way to retrieve them is by using the `shoutrrr generate telegram` command which will guide you through
creating a URL with your target chats.

!!! tip
    You can use the `nickfedor/shoutrrr` image in docker to run it without download/installing the `shoutrrr` CLI using:
    ```
    docker run --rm -it nickfedor/shoutrrr generate telegram
    ```

### Asking @shoutrrrbot
Another way of retrieving the Chat IDs, is by forwarding a message from the target chat to the [@shoutrrrbot](https://t.me/shoutrrrbot).
It will reply with the Chat ID for the chat where the forwarded message was originally posted.
Note that it will not work correctly for Group chats, as those messages are just seen as being posted by a user, not in a specific chat.
Instead you can use the second method, which is to invite the @shoutrrrbot into your group chat and address a message to it (start the message with @shoutrrrbot). You can then safely kick the bot from the group. 

The bot should be constantly online, unless it's usage exceeds the free tier on GCP. It's source is available at [github.com/nicholas-fedor/shoutrrrbot](https://github.com/nicholas-fedor/shoutrrrbot).



## Optional parameters

You can optionally specify the __`notification`__, __`parseMode`__ and __`preview`__ parameters in the URL:  

!!! info ""
    <pre>telegram://__`token`__@__`telegram`__/?channels=__`channel`__&notification=no&preview=false&parseMode=html</pre>

See [the telegram documentation](https://core.telegram.org/bots/api#sendmessage) for more information.

!!! note
    `preview` and `notification` are inverted in regards to their API counterparts (`disable_web_page_preview` and `disable_notification`)

### Parse Mode and Title

If a parse mode is specified, the message needs to be escaped as per the corresponding sections in
[Formatting options](https://core.telegram.org/bots/api#formatting-options).

When a title has been specified, it will be prepended to the message, but this is only supported for
the `HTML` parse mode. Note that, if no parse mode is specified, the message will be escaped and sent using `HTML`.

Since the markdown modes are really hard to escape correctly, it's recommended to stick to `HTML` parse mode.
