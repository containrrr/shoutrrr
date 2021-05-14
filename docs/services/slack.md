# Slack

The slack notification service uses [Slack Webhook](https://api.slack.com/messaging/webhooks)s to send messages.  
Follow the [Getting started with Incoming Webhooks](https://api.slack.com/messaging/webhooks#getting_started) guide and
replace the initial `https://hooks.slack.com/services/` part of the webhook URL with `slack://` to get your Shoutrrr URL.

!!! example
    Slack Webhook URL:
    ```uri
    https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
    ```
    Shoutrrr URL:
    ```uri
    slack://T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
    ```
    
--8<-- "docs/services/slack/config.md"

!!! note "Color format"
    The format for the `Color` prop follows the [slack docs](https://api.slack.com/reference/messaging/attachments#fields)
    but `#` needs to be escaped as `%23` when passed in a URL.  
    So <span style="background:#ff8000;width:.9em;height:.9em;display:inline-block;vertical-align:middle"></span><code>#ff8000</code> would be `%23ff8000` etc.

## Examples

!!! example
    All fields set:
    ```uri
    slack://ShoutrrrBot@T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX?color=good&title=Great+News
    ```