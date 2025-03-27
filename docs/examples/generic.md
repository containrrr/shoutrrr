# Examples

Examples of service URLs that can be used with [the generic service](../../services/generic) together with common service providers.

## Home Assistant

The service URL needs to be:
```
generic://HAIPAddress:HAPort/api/webhook/WebhookIDFromHA?template=json
```

And, if you need http://
```
generic://HAIPAddress:HAPort/api/webhook/WebhookIDFromHA?template=json&disabletls=yes
```

Then, in HA, use `{{ trigger.json.message }}` to get the message sent from the JSON.

_Credit [@JeffCrum1](https://github.com/JeffCrum1), source: [https://github.com/nicholas-fedor/shoutrrr/issues/325#issuecomment-1460105065]_
