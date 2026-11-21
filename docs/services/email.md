# Email

## URL Format

!!! info ""
    smtp://__`username`__:__`password`__@__`host`__:__`port`__/?from=__`fromAddress`__&to=__`recipient1`__[,__`recipient2`__,...]

--8<-- "docs/services/smtp/config.md"

Recipients can include a display name, for example:

```text
smtp://host:25/?from=sender%40example.com&to=John%20Smith%20%3Cjohn.smith%2Btag%40example.com%3E
```

This produces `To: John Smith <john.smith+tag@example.com>`. In query values,
`%20` or `+` represents a space in the display name; use `%2B` for a literal
plus in a name or mailbox. Legacy bare-plus mailbox tags remain supported.
