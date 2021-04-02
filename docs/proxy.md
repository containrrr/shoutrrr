To use a proxy with shoutrrr, you could either set the proxy URL in the environment variable `HTTP_PROXY` or override the default HTTP client like this:

```go
proxyurl, err := url.Parse("socks5://localhost:1337")
if err != nil {
	log.Fatalf("Error parsing proxy URL: %q", err)
}

http.DefaultClient.Transport = &http.Transport{
	Proxy: http.ProxyURL(proxyurl),
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
```
