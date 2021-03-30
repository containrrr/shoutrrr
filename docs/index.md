# Shoutrrr

<div align="center">
<img src="https://raw.githubusercontent.com/containrrr/shoutrrr/master/docs/shoutrrr-logotype.png" width="450" />


Notification library for gophers and their furry friends.
Heavily inspired by <a href="https://github.com/caronc/apprise">caronc/apprise</a>.

![github actions workflow status](https://github.com/containrrr/shoutrrr/workflows/Main%20Workflow/badge.svg)
[![codecov](https://codecov.io/gh/containrrr/shoutrrr/branch/master/graph/badge.svg)](https://codecov.io/gh/containrrr/shoutrrr)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/47eed72de79448e2a6e297d770355544)](https://www.codacy.com/gh/containrrr/shoutrrr/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=containrrr/shoutrrr&amp;utm_campaign=Badge_Grade)
[![report card](https://goreportcard.com/badge/github.com/containrrr/shoutrrr)](https://goreportcard.com/badge/github.com/containrrr/shoutrrr)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/containrrr/shoutrrr)
[![github code size in bytes](https://img.shields.io/github/languages/code-size/containrrr/shoutrrr.svg?style=flat-square)](https://github.com/containrrr/shoutrrr)
[![license](https://img.shields.io/github/license/containrrr/shoutrrr.svg?style=flat-square)](https://github.com/containrrr/shoutrrr/blob/master/LICENSE)

</div>

To make it easy and streamlined to consume shoutrrr regardless of the notification service you want to use,
we've implemented a notification service url schema. To send notifications, instantiate the `ShoutrrrClient` using one of
the service urls from the [overview](/shoutrrr/services/overview).
