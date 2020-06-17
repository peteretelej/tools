# active_conn

Prints the number of active MySQL connections to a database at a desired interval (default: 5seconds)

Needs DB credentials exported, for example
```
export DB_HOST=localhost
export DB_USER=demo
export DB_PASS=demo123
```

## Usage:

```
./active_conn -interval 10s
```
- Print number of active connections every 10 seconds
