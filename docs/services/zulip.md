# Zulip Chat

## URL Format

The shoutrrr service URL should look like this:  
> zulip://__`bot-mail`__:__`bot-key`__@__`zulip-domain`__/?stream=__`name-or-id`__&topic=__`name`__

Stream and topic are both optional and can be given as parameters to the Send method:

```go
  sender, __ := shoutrrr.CreateSender(url)

  params := make(types.Params)
  params["stream"] = "mystream"
  params["topic"] = "This is my topic"

  sender.Send(message, &params)
```

Since __`bot-mail`__  is a mail address you need to URL escape the `@` in it to `%40`.

An example service URL would look like:
> zulip://my-bot%40zulipchat.com:correcthorsebatterystable@example.zulipchat.com?stream=foo&topic=bar
