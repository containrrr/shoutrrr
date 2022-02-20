# Telegram

## URL Format

!!! info ""
    telegram://__`token`__@telegram?chats=__`channel-1`__[,__`chat-id-1`__,...]
    
--8<-- "docs/services/telegram/config.md"

## Getting a token for Telegram

Talk to [the botfather](https://core.telegram.org/bots#6-botfather).

## Identifying the target chats/channels

The `chats` param consists of one or more `Chat ID`s or `channel name`s. 

### Channels
The channel names can be retrieved in the telegram client in the `Channel info` section. 
Replace the `t.me/` prefix from the link with a `@`.

!!! note
    Channels names need to be prefixed by `@` to identify them as such.

### Chats
Group chats and private chats are identifieds by `Chat ID`s. Unfortunatly, they are generally not visible in the
telegram clients.
The easiest way to retrieve them is by using the `shoutrrr generate telegram` command which will guide you through
creating a URL with your target chats.

!!! tip
    You can use the `containrrr/shoutrrr` image in docker to run it without download/installing the `shoutrrr` CLI using:
    ```
    docker run --rm -it containrrr/shoutrrr generate telegram
    ```

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