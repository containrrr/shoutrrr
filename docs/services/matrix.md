# Matrix

## URL Format

*matrix://__`user`__:__`password`__@__`host`__:__`port`__/[?rooms=__`!roomID1`__[,__`roomAlias2`__]][&disableTLS=yes]*

## Authentication

If no `user` is specified, the `password` is treated as the authentication token. This means that no matter what login
flow your server uses, if you can manually retrieve a token, then Shoutrrr can use it.

### Password Login Flow

If a `user` and `password` is supplied, the `m.login.password` login flow is attempted if the server supports it.

## Rooms

If `rooms` are *not* specified, the service will send the message to all the rooms that the user has currently joined.

Otherwise, the service will only send the message to the specified rooms. If the user is *not* in any of those rooms,
but have been invited to it, it will automatically accept that invite.

**Note**: The service will **not** join any rooms unless they are explicitly specified in `rooms`. If you need the user
to join those rooms, you can send a notification with `rooms` explicitly set once.

### Room Lookup

Rooms specified in `rooms` will be treated as room IDs if the start with a `!` and used directly to identify rooms. If
they have no such prefix (or use a *correctly escaped* `#`) they will instead be treated as aliases, and a directory
lookup will be used to resolve their corresponding IDs.

**Note**: Don't use unescaped `#` for the channel aliases as that will be treated as the `fragment` part of the URL.
Either omit them or URL encode them, I.E. `rooms=%23alias:server` or `rooms=alias:server`

### TLS

If you do not have TLS enabled on the server you can disable it by providing `disableTLS=yes`. This will effectively
use `http` intead of `https` for the API calls.