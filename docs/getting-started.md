# Getting started

## As a package

Using shoutrrr is easy! There is currently two ways of using it as a package.

### Using the direct send command

```go
  url := "slack://token-a/token-b/token-c"
  err := shoutrrr.Send(url, "Hello world (or slack channel) !")
```

### Using a sender

```go
  url := "slack://token-a/token-b/token-c"
  sender, err := shoutrrr.CreateSender(url)
  sender.Send("Hello world (or slack channel) !", map[string]string { /* ... */ })
```

## Through the CLI

Start by running the `build.sh` script.
You may then run the shoutrrr executable:

```shell
$ ./shoutrrr

Usage:
./shoutrrr <ActionVerb> [...]
Possible actions: send, verify, generate
```

### Commands

#### Send

Send a notification using the supplied notification service url.

```bash
$ shoutrrr send \
    --url "<SERVICE_URL>" \
    --message "<MESSAGE BODY>"
```

#### Verify

Verify the validity of a notification service url.

```bash
$ shoutrrr verify \
    --url "<SERVICE_URL>"
```

#### Generate

Generate and display the configuration for a notification service url.

```bash
$ shoutrrr generate \
    --url "<SERVICE_URL>"
```

### Options

#### Debug

Enables debug output from the CLI.

| Flags           | Env.             | Default | Required |
| --------------- | ---------------- | ------- | -------- |
| `--debug`, `-d` | `SHOUTRRR_DEBUG` | `false` |          |

#### URL

The target url for the notifications generated, see [overview](./services/overview).

| Flags         | Env.           | Default | Required |
| ------------- | -------------- | ------- | -------- |
| `--url`, `-u` | `SHOUTRRR_URL` | N/A     | ✅       |

### Action details

```shell
$ ./shoutrrr generate
Usage:
./shoutrrr generate [OPTIONS] <service>
```
