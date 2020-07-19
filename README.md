# yugabyte-issue

## Run project

```
go run main.go
```

## Errors

```
pq: Operation failed. Try again.: Value write after transaction start: { physical: 1595176815752573 } >= { physical: 1595176815750985 }

pq: Query error: Restart read required at: { read: { physical: 1595176841255374 } local_limit: { physical: 1595176841303306 } global_limit: <min> in_txn_limit: <max> serial_no: 0 }

Transaction expired or aborted by a conflict: 40001
```
