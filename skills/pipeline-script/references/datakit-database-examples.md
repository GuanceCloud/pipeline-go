# DataKit Database Pipeline Examples

Examples from `/usr/local/datakit/pipeline` for database logs.

## MySQL Error

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/mysql.p \
  --message '2026-05-19T03:22:01.123456Z 12 [ERROR] Access denied for user root@localhost' \
  --require-key thread_id \
  --require-key status \
  --require-key msg \
  --expect thread_id=12
```

Key behavior: multiple grok branches cover standard error logs and slow-query logs; optional missing slow-query fields may produce debug output.

## MySQL Slow Query

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/mysql.p \
  --message $'# Time: 2026-05-19T03:22:01.123456Z\n# User@Host: app_user @ app_host [10.0.0.8] Id: 77\n# Query_time: 1.234 Lock_time: 0.001 Rows_sent: 10 Rows_examined: 200\n# Thread_id: 88 Killed: 0 Errno: 0\n# Bytes_sent: 4096 Bytes_received: 128\nselect * from orders where id = 42;' \
  --require-key db_user \
  --require-key db_ip \
  --require-key query_time \
  --require-key rows_examined \
  --require-key db_slow_statement \
  --expect query_time=1.234 \
  --expect rows_examined=200
```

Key behavior: uses multiline `add_pattern("sqls", "(?s)(.*)")`; cast query timings to float and row/byte counters to int.

## MongoDB JSON

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/mongodb.p \
  --message '{"t":{"$date":"2026-05-19T03:22:01.123Z"},"s":"I","c":"NETWORK","ctx":"conn1","msg":"connection accepted"}' \
  --require-key status \
  --require-key component \
  --require-key context \
  --expect component=NETWORK
```

Key behavior: extracts nested `$date` by first extracting `t` into a temporary field, then `json(tmp, `$date`, "time")`, then `drop_key(tmp)`.

## PostgreSQL

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/postgresql.p \
  --message '2026-05-19 11:42:31 CST [12345] LOG:  database system is ready to accept connections' \
  --require-key process_id \
  --require-key status \
  --expect status=LOG
```

Key behavior: uses custom date/status/session patterns; `process_id` is not cast in the template and may remain a string.

## SQL Server

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/sqlserver.p \
  --message '2026-05-19 11:42:31.123 spid51      Error: 18456, Severity: 14, State: 1.' \
  --require-key origin \
  --require-key msg \
  --expect origin=spid51
```

Key behavior: calls `default_time(time, "+0")`, so timezone behavior differs from templates that call `default_time(time)`.

## Dameng

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/dameng.p \
  --message '2026-05-19 11:42:31 [ERROR] database startup failed' \
  --require-key status \
  --require-key msg \
  --expect status=ERROR
```

## Kingbase

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/kingbase.p \
  --message '2026-05-19 11:42:31 CST [9876] ERROR: duplicate key value violates unique constraint' \
  --require-key process_id \
  --require-key status \
  --require-key msg \
  --expect status=ERROR
```
