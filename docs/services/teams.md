# Teams

!!! attention Webhook URL scheme changed
    Microsoft has changed the URL scheme for Teams webhooks. You will now have to specify the hostname using:
    ```text
    ?host=example.webhook.office.com
    ```
    Where `example` is your organization short name

## URL Format

!!! info ""
    teams://__`group`__@__`tenant`__/__`altId`__/__`groupOwner`__?host=__`organization`__.webhook.office.com

--8<-- "docs/services/teams/config.md"

## Setting up a webhook

To be able to use the Microsoft Teams notification service, you first need to set up a custom webhook.
Instructions on how to do this can be found in [this guide](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using#setting-up-a-custom-incoming-webhook)

## Extracting the token

The token is extracted from your webhook URL:

<pre><code>https://<b>&lt;organization&gt;</b>.webhook.office.com/webhookb2/<b>&lt;group&gt;</b>@<b>&lt;tenant&gt;</b>/IncomingWebhook/<b>&lt;altId&gt;</b>/<b>&lt;groupOwner&gt;</b></code></pre>
