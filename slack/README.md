# slack

Sends slack notifications. 

## Basic Usage

Export the slack webhook url

```
export SLACK_WEBHOOK=https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
```

Submit a normal message

```
slack --message="hello world?"
```

You can use Slack's [formatting](https://api.slack.com/docs/message-formatting) on messages, for example

```
slack --message="hello, *world*! _hi_"
```

Submit a message as an ERROR

```
slack --error --message="Something went wrong!"
```

You can also alert `@here` on messages by passing `--alert`

```
slack --error --alert --message="Something went wrong!"
```
- This slack message will `cc @here`

Message with a Field

```
slack --error --message="Oops!" --field='{"title":"Environment","value":"UAT","short":true}"
```

Message with multiple Fields

```
slack --error --message="Oops!" --fields='[{"title":"Environment","value":"UAT","short":true},{"title":"Application","value":"My New App *Abc*","short":true}]"
```

That's all.

Optionally: specify webhook directly on command

```
slack --webhook="https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXX" --message="Something went wrong!"
```


### Tip

You can also pass arguments like this:
```
slack -message "Hello World!" -error -alert
```
