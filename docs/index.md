<p align="center">
    <a href="https://github.com/containrrr/shoutrrr"><img src="https://github.com/containrrr/shoutrrr/raw/gh-pages/shoutrrr.jpg" width="450" /></a>
</p>
<h1 align="center">
    Shoutrrr
</h1>
<p align="center">
    Notification library for gophers and their furry friends.
    Heavily inspired by <a href="https://github.com/caronc/apprise">caronc/apprise</a>.
</p>
<p align="center">
    <a href="https://app.codacy.com/app/containrrr/shoutrrr?utm_source=github.com&utm_medium=referral&utm_content=containrrr/shoutrrr&utm_campaign=Badge_Grade_Dashboard"><img
          alt="codacy coverage"
          src="https://img.shields.io/codacy/coverage/30ce077eecde418ca328f4f7868f70c8.svg?style=flat-square"
        /></a>
    <a href="https://circleci.com/gh/containrrr/shoutrrr"><img
          alt="circleci"
          src="https://img.shields.io/circleci/project/github/containrrr/shoutrrr/master.svg?style=flat-square"
        /></a>
    <a href="https://app.codacy.com/app/containrrr/shoutrrr?utm_source=github.com&utm_medium=referral&utm_content=containrrr/shoutrrr&utm_campaign=Badge_Grade_Dashboard"><img
 alt="codacy grade" src="https://img.shields.io/codacy/grade/30ce077eecde418ca328f4f7868f70c8/master.svg?style=flat-square" /></a>
    <a href="https://github.com/containrrr/shoutrrr"><img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/containrrr/shoutrrr.svg?style=flat-square" /></a>
    <a href="https://github.com/containrrr/shoutrrr/blob/master/LICENSE"><img alt="license" src="https://img.shields.io/github/license/containrrr/shoutrrr.svg?style=flat-square" /></a>
    <a href="https://godoc.org/github.com/containrrr/shoutrrr"><img           src="https://godoc.org/github.com/containrrr/shoutrrr?status.svg" alt="GoDoc" /></a>
</p>

To make it easy and streamlined to consume shoutrrr regardless of the notification service you want to use,
we've implemented a notification service url schema. To send notifications, instantiate the `ShoutrrrClient` using one of
the service urls below.