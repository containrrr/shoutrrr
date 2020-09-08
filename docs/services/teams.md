# Teams

## URL Format

*teams://__`token-a`__/__`token-b`__/__`token-c`__*

## Setting up a webhook

To be able to use the Microsoft Teams notification service, you first need to set up a custom webhook.
Instructions on how to do this can be found in [this guide](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using#setting-up-a-custom-incoming-webhook)

## Extracting the token

The token is extracted from your webhook URL:

```
  https://outlook.office.com/webhook/{tokenA}/IncomingWebhook/{tokenB}/{tokenC}
```