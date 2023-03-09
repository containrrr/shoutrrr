# Generic
The Generic service can be used for any target that is not explicitly supported by Shoutrrr, as long as it
supports receiving the message via a POST request.
Usually, this requires customization on the receiving end to interpret the payload that it receives, and might
not be a viable approach.

Common examples for use with service providers can be found under [examples](../examples/generic.md).

## JSON template
By using the built in `JSON` template (`template=json`) you can create a generic JSON payload. The keys used for `title` and `message` can be overriden
by supplying the params/query values `titleKey` and `messageKey`.

!!! example
    ```json
    {
        "title": "Oh no!",
        "message": "The thing happened and now there is stuff all over the area!"
    }
    ```

## Shortcut URL
You can just add `generic+` as a prefix to your target URL to use it with the generic service, so
```
https://example.com/api/v1/postStuff
```
would become
```
generic+https://example.com/api/v1/postStuff
```

!!! note
    Any query variables added to the URL will be escaped so that they can be forwarded to the remote server. That means that you cannot use `?format=json` with the  `generic+https://`, just use `generic://` instead!

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
