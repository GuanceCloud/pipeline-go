# DataKit Middleware Pipeline Examples

Examples from `/usr/local/datakit/pipeline` for middleware, search, queue, and service logs.

## Redis

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/redis.p \
  --message '1234:M 19 May 2026 11:22:33.456 # Server initialized' \
  --require-key pid \
  --require-key role \
  --require-key status \
  --expect status=warning
```

Key behavior: maps Redis severity symbols to status: `.` debug, `-` verbose, `*` notice, `#` warning.

## Elasticsearch Slow Query

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/elasticsearch.p \
  --message '[2026-05-19T11:42:31,123][WARN ][i.s.s.query   ] [node-1] [orders][0] took[12ms], took_millis[12], total_hits[1 hits], types[]' \
  --require-key status \
  --require-key index \
  --require-key duration \
  --expect duration=12000000
```

Key behavior: `duration` is captured as milliseconds, cast to int, then converted from ms to ns with `duration_precision(duration, "ms", "ns")`.

## Kafka

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/kafka.p \
  --message '2026-05-19 11:42:31 INFO kafka.server.KafkaServer:123 - started Kafka server' \
  --require-key status \
  --require-key name \
  --require-key line \
  --expect line=123
```

Key behavior: multiple grok fallbacks cover thread-format, duration-format, logger line-format, and bracketed timestamp-format logs. `line` remains a string.

## RabbitMQ

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/rabbitmq.p \
  --message '2026-05-19 11:42:31.123 [error] connection refused' \
  --require-key status \
  --require-key msg \
  --expect status=error
```

Key behavior: supports both `==== time ===` and bracketed status log variants.

## Solr

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/solr.p \
  --message '2026-05-19 11:42:31.123 INFO  (qtp123-456) [ x:collection1] o.a.s.c.S.Request [collection1] webapp=/solr path=/select params={q=*:*} hits=10 status=0 QTime=12' \
  --require-key status \
  --require-key thread \
  --require-key reporter \
  --require-key qtime
```

Key behavior: `hits`, `qstatus`, and `qtime` remain strings because the template does not cast them.

## Consul

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/consul.p \
  --message 'May 19 11:42:31 host01 consul[1234]: 2026-05-19T11:42:31123 [INFO] agent.server: cluster leadership acquired' \
  --require-key level \
  --require-key character \
  --require-key msg \
  --expect level=INFO
```

Key behavior: calls `drop_origin_data()`, so `message` is removed from output after extraction.

## Jenkins

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/jenkins.p \
  --message $'2026-05-19 11:42:31.123+0000 [id=42]\tSEVERE\t' \
  --require-key id \
  --require-key status \
  --expect status=error
```

Key behavior: the template expects literal tab separators and maps `SEVERE`/`ERROR` to `error`.

## TDengine

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/tdengine.p \
  --message '2026-05-19 11:42:31 TAOS_HTTP info "| 200 | 12ms | 10.0.0.5 | GET | /rest/sql"' \
  --require-key module \
  --require-key level \
  --require-key status \
  --require-key code \
  --require-key cost_time \
  --expect status=info
```

Key behavior: uses `parse_duration(cost_time)` and then `duration_precision(cost_time, "ns", "ms")`; the final `cost_time` is milliseconds.
