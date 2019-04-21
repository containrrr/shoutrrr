<h1 align="center">
    Shoutrrr
</h1>
<p align="center">
    Notification library for gophers and their furry friends.
    Heavily inspired by <a href="https://github.com/caronc/apprise">caronc/apprise</a>.
</p>
<p align="center">
    <img
      alt="codacy coverage"
      src="https://img.shields.io/codacy/coverage/30ce077eecde418ca328f4f7868f70c8.svg?style=flat-square"
    />
    <img
      alt="circleci"
      src="https://img.shields.io/circleci/project/github/containrrr/shoutrrr/master.svg?style=flat-square"
    />
    <img
      alt="codacy grade"
      src="https://img.shields.io/codacy/grade/30ce077eecde418ca328f4f7868f70c8/master.svg?style=flat-square"
    />
    <img
      alt="GitHub code size in bytes"
      src="https://img.shields.io/github/languages/code-size/containrrr/shoutrrr.svg?style=flat-square"
    />
    <img
      alt="license"
      src="https://img.shields.io/github/license/containrrr/shoutrrr.svg?style=flat-square"
    />
    <a href="https://beerpay.io/containrrr/shoutrrr">
      <img
        alt="beerpay wish"
        src="https://beerpay.io/containrrr/shoutrrr/make-wish.svg"
      />
    </a>
      <a href="https://beerpay.io/containrrr/shoutrrr">
      <img
        alt="beerpay wish"
        src="https://beerpay.io/containrrr/shoutrrr/badge.svg?style=flat-square"
      />
    </a>
</p>


### Using Shoutrrr


Using shoutrrr is as easy as:

```
  url := "slack://token-a/token-b/token-c"
  err := shoutrrr.Send(url, "Hello world (or slack channel) !")
   
```

If you've provided the environment variable `SHOUTRRR_URL`, you may instead use


```
  err := shoutrrr.SendEnv("Hello world (or slack channel) !")
```

### Service URL:s

To make it easy and streamlined to consume shoutrrr regardless of the notification service you want to use,
we've implemented a notification service url schema. To send notifications, instantiate the ShoutrrrrClient using one of
the service urls below.

| Service   | Format                                                                                      |
| --------- | ------------------------------------------------------------------------------------------- |
| Telegram  | `telegram://api-token/channel` or `telegram://api-token/channel-a/channel-b/channel-c` etc. |
| Slack     | `slack://token-a/token-b/token-c` or `slack://botname/token-a/token-b/token-c`              |
| Discord   | `discord://channel/token`                                                                   |
| Pushover  | `pushover://token/user/device`                                                              |
