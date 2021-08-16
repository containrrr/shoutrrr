# Telegram

## URL Format

!!! info ""
    telegram://__`token`__@telegram?channels=__`channel-1`__[,__`channel-2`__,...]
    
--8<-- "docs/services/telegram/config.md"

## Getting a token for Telegram

Talk to [the botfather](https://core.telegram.org/bots#6-botfather).

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