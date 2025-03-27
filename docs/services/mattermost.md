# MatterMost

## URL Format

!!! info ""
    mattermost://[__`username`__@]__`mattermost-host`__/__`token`__[/__`channel`__][?icon=__`smiley`__&disabletls=__`yes`__]

--8<-- "docs/services/mattermost/config.md"



## Creating a Webhook in MatterMost

1. Open up the Integrations page by clicking on *Integrations* within the menu
![Screenshot 1](mattermost/1.PNG)

2. Click *Incoming Webhooks*
![Screenshot 2](mattermost/2.PNG)

3. Click *Add Incoming Webhook*
![Screenshot 3](mattermost/3.PNG)

4. Fill in the information for the webhook and click *Save*
![Screenshot 4](mattermost/4.PNG)

5. If you did everything correctly, MatterMost will give you the *URL* to your newly created webhook
![Screenshot 5](mattermost/5.PNG)

6. Format the service URL
```
https://your-domain.com/hooks/bywsw8zt5jgpte3nm65qjiru6h
                              └────────────────────────┘
                                        token
mattermost://your-domain.com/bywsw8zt5jgpte3nm65qjiru6h
                             └────────────────────────┘
                                       token
```

## Additional URL configuration

Mattermost provides functionality to post as another user or to another channel, compared to the webhook configuration.
<br/>
To do this, you can add a *user* and/or *channel* to the service URL.

```
mattermost://shoutrrrUser@your-domain.com/bywsw8zt5jgpte3nm65qjiru6h/shoutrrrChannel
             └──────────┘                 └────────────────────────┘ └─────────────┘
                 user                               token                channel
```

## Passing parameters via code

If you want to, you also have the possibility to pass parameters to the `send` function.
<br/>
The following example contains all parameters that are currently supported.

```gotemplate
params := (*types.Params)(
	&map[string]string{
		"username": "overwriteUserName",
		"channel": "overwriteChannel",
        "icon": "overwriteIcon",
	},
)

service.Send("this is a message", params)
```

This will overwrite any options, that you passed via URL.