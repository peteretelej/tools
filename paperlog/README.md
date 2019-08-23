## paperlog

Logs to papertrail

- Uses new connection for each log entry
- Useful for bash scripts, or apps that for some reason can't use remote_syslog2 or papertrail integration via syslog.

### Basic Usage

Export the Papertrail address to log to and your app name. You can also pass this values as arguments to the command instead of exporting.

```
export PAPERLOG_ADDR=logsN.papertrailapp.com:XXXXX
export PAPERLOG_APPNAME=myapp
```

Submit a normal log (INFO)

```
paperlog --message="Hello"
```

Submit error log (ERROR)

```
paperlog --error --message="Error here"
```

Submit a warning

```
paperlog --warn --message="Warning here"
```

That's all

The papertrail configs can also be specified as part of the CLI

```
paperlog --addr="logsN.papertrailapp.com:XXXXX" --appname="myapp" --message="hello"
```
