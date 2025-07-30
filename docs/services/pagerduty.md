# OpsGenie

## URL Format

--8<-- "docs/services/pagerduty/config.md"

## Create a Service Integration in PagerDuty

Follow to official [PagerDuty documentation](https://support.pagerduty.com/main/docs/services-and-integrations) to

1. create a service, and then
2. add an 'Events API V2' integration to the service. Note teh value of the `Integration Key`

The host is always `events.pagerduty.com`, so you do not need to explicitly specify it, however you can provide one if
you prefer, of if there is a need to override the default.

```
pagerduty:///eb243592-faa2-4ba2-a551q-1afdf565c889
             └───────────────────────────────────┘
                       integration key
                       

pagerduty://events.pagerduty.com/eb243592-faa2-4ba2-a551q-1afdf565c889
                                 └───────────────────────────────────┘
                                           integration key
                                              
```

## Passing parameters via code

If you want to, you can pass additional parameters to the `send` function.
<br/>
The following example contains all parameters that are currently supported.

```go
service.Send("An example alert message", &types.Params{
"severity": "critical",
"source":   "The source of the alert",
"action":   "trigger",

})
```

See the [PagerDuty documentation](https://developer.pagerduty.com/docs/send-alert-event) for details on which fields
are required and what values are permitted for each field

## Passing parameters via URL

You can optionally specify the parameters in the URL:

!!! info ""
pagerduty://events.pagerduty.com/145d44a18bb44a0bc06161d5f541a90a?severity=critical&source=beszel&action=trigger
!!!

Example using the command line:

```shell
shoutrrr send -u 'pagerduty://events.pagerduty.com/145d44a18bb44a0bc06161d5f541a90a?severity=critical&source=beszel&action=trigger' -m 'This is a test'
```


