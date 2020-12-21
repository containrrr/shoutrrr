# Basic generator

The basic generator looks at the `key:""`, `desc:""` and `default:""` tags on service configuration structs and uses them to ask the user to fill in their corresponding values.

Example:
```shell
$ shoutrrr generate telegram
```
```yaml
Generating URL for telegram using basic generator
Enter the configuration values as prompted

Token: 110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw
Preview[Yes]: No
Notification[Yes]:
ParseMode[None]:
Channels: @mychannel

URL: telegram://110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw@telegram?channels=@mychannel&notification=Yes&parsemode=None&preview=No
```
