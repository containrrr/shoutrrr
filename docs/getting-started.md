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

### Action details

```shell
$ ./shoutrrr send
Usage:
./shoutrrr send [OPTIONS] <URL> <Message [...]>

OPTIONS:
  -verbose
        display additional output
```

```shell
$ ./shoutrrr verify
Usage:
./shoutrrr send [OPTIONS] <URL> <Message [...]>
```

```shell
$ ./shoutrrr generate
Usage:
./shoutrrr generate [OPTIONS] <service>
```