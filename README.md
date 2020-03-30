<p align="center">
    <a href="https://github.com/containrrr/shoutrrr">
        <img src="https://raw.githubusercontent.com/containrrr/shoutrrr/gh-pages/shoutrrr.jpg" width="450" /></a>
</p>
<h1 align="center">
    Shoutrrr
</h1>
<p align="center">
    Notification library for gophers and their furry friends.
    Heavily inspired by <a href="https://github.com/caronc/apprise">caronc/apprise</a>.
</p>
<p align="center">
    <img src="https://github.com/containrrr/shoutrrr/workflows/Main%20Workflow/badge.svg" />
    <a href="https://app.codacy.com/app/containrrr/shoutrrr?utm_source=github.com&utm_medium=referral&utm_content=containrrr/shoutrrr&utm_campaign=Badge_Grade_Dashboard"><img
      alt="codacy coverage"
      src="https://img.shields.io/codacy/coverage/30ce077eecde418ca328f4f7868f70c8.svg?style=flat-square"
    /></a>
    <a href="https://app.codacy.com/app/containrrr/shoutrrr?utm_source=github.com&utm_medium=referral&utm_content=containrrr/shoutrrr&utm_campaign=Badge_Grade_Dashboard"><img
 alt="codacy grade" src="https://img.shields.io/codacy/grade/30ce077eecde418ca328f4f7868f70c8/master.svg?style=flat-square" /></a>
    <a href="https://github.com/containrrr/shoutrrr"><img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/containrrr/shoutrrr.svg?style=flat-square" /></a>
    <a href="https://github.com/containrrr/shoutrrr/blob/master/LICENSE"><img alt="license" src="https://img.shields.io/github/license/containrrr/shoutrrr.svg?style=flat-square" /></a>
    <a href="https://godoc.org/github.com/containrrr/shoutrrr"><img           src="https://godoc.org/github.com/containrrr/shoutrrr?status.svg" alt="GoDoc" /></a>
</p>

## Quick Start

### As a package

Using shoutrrr is easy! There is currently two ways of using it as a package.

#### Using the direct send command

```go
  url := "slack://token-a/token-b/token-c"
  err := shoutrrr.Send(url, "Hello world (or slack channel) !")

```

#### Using a sender
```go
  url := "slack://token-a/token-b/token-c"
  sender := shoutrrr.CreateSender(url)
  sender.Send("Hello world (or slack channel) !", map[string]string { /* ... */ })
```

### Through the CLI

Start by running the `build.sh` script.
You may then run send notifications using the shoutrrr executable:

```shell
$ shoutrrr send [OPTIONS] <URL> <Message [...]>
```

## Documentation
For additional details, visit the [full documentation](https://containrrr.github.io/shoutrrr). 
