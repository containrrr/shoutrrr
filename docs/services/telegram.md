# Telegram

## URL Format

*telegram://__`token`__@telegram?channels=__`channel-1`__[,__`channel-2`__,...]*

## Getting a token for Telegram

Talk to [the botfather](https://core.telegram.org/bots#6-botfather).

## Optional parameters

You can optionally specify the __`notification`__, __`parseMode`__ and __`preview`__ parameters in the URL:  
*telegram://__`token`__@__`telegram`__/?channels=__`channel`__&notification=no&preview=false&parseMode=markDownv2*

See [the telegram documentation](https://core.telegram.org/bots/api#sendmessage) for more information.

__Note:__ `preview` and `notification` are inverted in regards to their API counterparts (`disable_web_page_preview` and `disable_notification`)