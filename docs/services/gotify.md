# Gotify

## URL Format

--8<-- "docs/services/gotify/config.md"

## Examples

!!! example "Common usage"
    
    ```uri
    gotify://gotify.example.com:443/AzyoeNS.D4iJLVa/?title=Great+News&priority=1
    ```

!!! example "With subpath"
    ```uri
    gotify://example.com:443/path/to/gotify/AzyoeNS.D4iJLVa/?title=Great+News&priority=1
    ```

!!! example "Gotify 3 application token"
    ```uri
    gotify://gotify.example.com:443/gtfya.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrs?title=Great+News&priority=1
    ```

Legacy application tokens remain supported. Copy the complete private application token from Gotify; the token above is a synthetic example.
