# Generic
The Generic service can be used for any target that is not explicitly supported by Shoutrrr, as long as it
supports recieving the message via a POST request.
Usually, this requires customization on the recieving end to interpret the payload that it recives, and might
not be a viable approach.

## Shortcut URL
You can just add `generic+` as a prefix to your target URL to use it with the generic service, so
```
https://example.com/api/v1/postStuff
```
would become
```
generic+https://example.com/api/v1/postStuff
```

## Forwarded query variables
All query variables that are not listed in the [Query/Param Props](#queryparam_props) section will be
forwarded to the target endpoint.
If you need to pass a query variable that _is_ reserved, you can prefix it with an underscore (`_`).

!!! example
    The URL `generic+https://example.com/api/v1/postStuff?contenttype=text/plain` would send a POST message
    to `https://example.com/api/v1/postStuff` using the `Content-Type: text/plain` header.

    If instead escaped, `generic+https://example.com/api/v1/postStuff?_contenttype=text/plain` would send a POST message
    to `https://example.com/api/v1/postStuff?contenttype=text/plain` using the `Content-Type: application/json` header (as it's the default).


## URL Format

--8<-- "docs/services/generic/config.md"