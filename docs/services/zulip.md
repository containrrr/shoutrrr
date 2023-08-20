# Zulip Chat

## URL Format

The shoutrrr service URL should look like this:
!!! info ""
    zulip://__`botmail`__:__`botkey`__@__`host`__/?stream=__`stream`__&topic=__`topic`__

--8<-- "docs/services/zulip/config.md"

!!! note
    Since __`botmail`__  is a mail address you need to URL escape the `@` in it to `%40`.

### Examples

Stream and topic are both optional and can be given as parameters to the Send method:

```go
  sender, _ := shoutrrr.CreateSender(url)

  params := make(types.Params)
  params["stream"] = "mystream"
  params["topic"] = "This is my topic"

  sender.Send(message, &params)
```

!!! example "Example service URL"
    zulip://my-bot%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar
