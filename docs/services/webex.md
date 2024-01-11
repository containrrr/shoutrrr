# Webex

Adds Webex service functionality

!!! info "Your bot must be in the room"
    Your bot must be invited to the room before a message can be sent by it.

## URL Format

The URL format requires the `@webex?` to be present. This is to handle the URL
being parsed properly and is ignored.

!!! info ""
    ```uri
   webex://__`token`__@webex?rooms=__`room-1`__[,__`room-2`__,...]

--8<-- "docs/services/webex/config.md"