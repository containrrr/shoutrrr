## Prerequisites
To contribute code changes to this project you will need to install the go distribution.
 * [Go](https://golang.org/doc/install)
 
Also, as shoutrrr utilizes go modules for vendor locking, you'll need atleast Go 1.11.
You can check your current version of the go language as follows:
```bash
  ~ $ go version
  go version go1.12.1 darwin/amd64
```

## Checking out the code
Do not place your code in the go source path.
```bash
git clone git@github.com:<yourfork>/shoutrrr.git
cd shoutrrr
```

## Building and testing
shoutrrr is a go library and is built with go commands. The following commands assume that you are at the root level of your repo.
```bash
go build ./cmd/shoutrrr                # compiles and packages an executable stand-alone client of shoutrrr
go test ./... -v                       # runs tests with verbose output
./shoutrrr                             # runs the application
```

## Commit messages
Shoutrrr try to follow the conventional commit specification. More information is available [here](https://www.conventionalcommits.org/en/v1.0.0-beta.4/#summary)
