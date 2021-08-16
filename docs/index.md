# Shoutrrr

<div align="center">
<img src="https://raw.githubusercontent.com/containrrr/shoutrrr/main/docs/shoutrrr-logotype.png" height="450" width="450" />
</div>

<p align="center">
Notification library for gophers and their furry friends.<br />
Heavily inspired by <a href="https://github.com/caronc/apprise">caronc/apprise</a>.
</p>

<p align="center" class="badges">
    <a target="_blank" rel="noopener noreferrer" href="https://github.com/containrrr/shoutrrr/workflows/Main%20Workflow/badge.svg">
        <img src="https://github.com/containrrr/shoutrrr/workflows/Main%20Workflow/badge.svg" alt="github actions workflow status">
    </a>
    <a href="https://codecov.io/gh/containrrr/shoutrrr" rel="nofollow">
        <img alt="codecov" src="https://codecov.io/gh/containrrr/shoutrrr/branch/main/graph/badge.svg">
    </a>
    <a href="https://www.codacy.com/gh/containrrr/shoutrrr/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=containrrr/shoutrrr&amp;utm_campaign=Badge_Grade" rel="nofollow">
        <img alt="Codacy Badge" src="https://app.codacy.com/project/badge/Grade/47eed72de79448e2a6e297d770355544">
    </a>
    <a href="https://goreportcard.com/badge/github.com/containrrr/shoutrrr" rel="nofollow">
        <img alt="report card" src="https://goreportcard.com/badge/github.com/containrrr/shoutrrr">
    </a>
    <a href="https://pkg.go.dev/github.com/containrrr/shoutrrr" rel="nofollow">
        <img alt="go.dev reference" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&amp;logoColor=white&amp;style=flat-square">
    </a>
    <a href="https://github.com/containrrr/shoutrrr">
        <img alt="github code size in bytes" src="https://img.shields.io/github/languages/code-size/containrrr/shoutrrr.svg?style=flat-square">
    </a>
    <a href="https://github.com/containrrr/shoutrrr/blob/main/LICENSE">
        <img alt="license" src="https://img.shields.io/github/license/containrrr/shoutrrr.svg?style=flat-square">
    </a>
</p>



To make it easy and streamlined to consume shoutrrr regardless of the notification service you want to use,
we've implemented a notification service url schema. To send notifications, instantiate the `ShoutrrrClient` using one of
the service urls from the [overview](services/overview.md).
