# Tools

A suite CLI tools & utilities

- [paperlog](./paperlog) - logs messages to papertrail
- [slack](./slack) - send slack notifications
- [get_secret](./get_secret) - fetch AWS Secrets Manager Secrets

## Basic Usage

Download a binary for your OS from the [Releases page](https://github.com/peteretelej/tools/releases)

For Example, for the slack cli tool, download an approriate version from Releases then:

```
# rename binary
mv slack*amd64 slack

# Add execute permissions to binary (for *nix Oss)
chmod +x slack

# just run it
./slack

# getting command usage help
./slack --help
```

- See usage for each binary in their respective folders

## Building from source

Get the repo

```
git clone https://github.com/peteretelej/tools.git
cd tools
```

To build a specific tool, cd into it's directory and build it:

```
cd slack
go build
```

- This builds a binary (`slack` or `slack.exe` depending on OS) that you can use directly. It has no dependencies, so you can copy it elsewhere and use as desired (eg /usr/local/bin)
