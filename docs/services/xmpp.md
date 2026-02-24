# XMPP (Jabber)

## URL Format

!!! info ""
    xmpp://__`username`__:__`password`__@__`host`__:__`port`__/?receiver=__`recipient1`__[,__`recipient2`__,...]

--8<-- "docs/services/xmpp/config.md"

## Examples

!!! example "Send to a single recipient"

    ```uri
    xmpp://user:pass@jabber.example.com?receiver=recipient@jabber.example.com
    ```

!!! example "Send to multiple recipients"

    ```uri
    xmpp://user:pass@jabber.example.com?receiver=alice@example.com&receiver=bob@example.com
    ```

!!! example "Using a non-standard port with TLS"

    ```uri
    xmpp://user:pass@jabber.example.com:5223?receiver=recipient@example.com&tls=Yes
    ```
