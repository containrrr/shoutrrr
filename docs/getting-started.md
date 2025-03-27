# Getting started

## As a package

Using shoutrrr is easy! There is currently two ways of using it as a package.

### Using the direct send command
Easiest to use, but very limited.

```go
url := "slack://token-a/token-b/token-c"
err := shoutrrr.Send(url, "Hello world (or slack channel) !")
```

### Using a sender
Using a sender gives you the ability to preconfigure multiple notification services and send to all of them with the same `Send(message, params)` method.

```go
urlA := "slack://token-a/token-b/token-c"
urlB := "telegram://110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw@telegram?channels=@mychannel"
sender, err := shoutrrr.CreateSender(urlA, urlB)

// Send notifications instantly to all services
sender.Send("Hello world (or slack/telegram channel)!", map[string]string { "title": "He-hey~!"  })

// ...or bundle notifications... 
func doWork() error {
    // ...and send them when leaving the scope
    defer sender.Flush(map[string]string { "title": "Work Result" })
    
    sender.Enqueue("Started doing %v", stuff)
    
    // Maybe get creative...?
    defer func(start time.Time) { 
    	sender.Enqueue("Elapsed: %v", time.Now().Sub(start)) 
    }(time.Now())
    
    if err := doMoreWork(); err != nil {
        sender.Enqueue("Oh no! %v", err)
    	
        // This will send the currently queued up messages...
        return
    }   
    
    sender.Enqueue("Everything went very well!")
    
    // ...or this:
}

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

On a system with Go installed you can install the latest Shoutrrr CLI
command with:

```shell
go install github.com/nicholas-fedor/shoutrrr/shoutrrr@latest
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
$ shoutrrr generate [OPTIONS] <SERVICE>
```

| Flags                        | Description                                     |
| ---------------------------- | ------------------------------------------------|
| `-g, --generator string`     |  The generator to use (default "basic")         |
| `-p, --property stringArray` |  Configuration property in key=value format     |
| `-s, --service string`       |  The notification service to generate a URL for |

**Note**: Service can either be supplied as the first argument or using the `-s` flag.

For more information on generators, see [Generators](./generators/overview.md).

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

## From a GitHub Actions workflow

You can also use Shoutrrr from a GitHub Actions workflow.

See this example and the [action on GitHub
Marketplace](https://github.com/marketplace/actions/shoutrrr-action):

```yaml
name: Deploy
on:
  push:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Some other steps needed for deploying
        run: ...
      - name: Shoutrrr
        uses: nicholas-fedor/shoutrrr-action@v1
        with:
          url: ${{ secrets.SHOUTRRR_URL }}
          title: Deployed ${{ github.sha }}
          message: See changes at ${{ github.event.compare }}.
```
