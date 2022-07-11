
## Service status

[ Ported |  Passes all tests ]

- `[x] [x] bark`
- `[x] [x] discord`
- `[x] [x] generic`
- `[x] [x] googlechat`
- `[x] [x] gotify`
- `[x] [x] ifttt`
- `[x] [x] join`
- `[x] [x] logger`
- `[x] [x] matrix`
- `[x] [x] mattermost`
- `[x] [x] opsgenie`
- `[x] [x] pushbullet`
- `[x] [x] pushover`
- `[x] [x] rocketchat`
- `[x] [x] slack`
- `[x] [x] smtp`
- `[x] [x] teams`
- `[x] [x] telegram`
- `[x] [x] zulip`

### Generic
~~Cannot be ported to `GeneratedConfig` right now due to how it handles the URL parsing very differently.
Needs special handling of escaped query vars. Perhaps this could be implemented by putting them in a `CustomVars` map?~~ **Done!**

### ifttt
~~Value validation is not checked, fails 1/14 specs:~~
 - ~~the ifttt package when creating a config when given an url [It] should return an error if message value is above 3~~ **Done!**

### gotify
~~Uses the _last_ path item as the token, which is not supported. Fails 3/10 specs:~~
 - ~~the Gotify plugin URL building and token validation functions creating a config when parsing the configuration URL [It] should be identical after de-/serialization (without path)~~
- ~~the Gotify plugin URL building and token validation functions sending the payload [It] should not report an error if the server accepts the payload~~
- ~~the Gotify plugin URL building and token validation functions sending the payload [It] should not panic if an error occurs when sending the payload~~ **Done!**
