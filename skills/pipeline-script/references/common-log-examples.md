# Common Log Examples

These examples were validated with `./bin/pipeline-check`. When testing grok patterns that contain brackets from a shell command, prefer saving the script to a `.p` file and using `--script`; inline `--cmd` requires an extra shell escaping layer.

## JSON Application Log

Sample:

```json
{"ts":"2026-05-19T03:22:01Z","level":"error","service":"checkout","trace_id":"abc123","latency_ms":87,"msg":"payment declined","user":{"id":"u-42"}}
```

Script:

```python
json(_, ts, time)
json(_, level)
json(_, service)
json(_, trace_id)
json(_, latency_ms)
json(_, msg)
json(_, user.id, user_id)
default_time(time)
```

Validation keys: `level`, `service`, `trace_id`, `user_id`; expect `latency_ms=87`.

## Nginx Combined Access Log

Sample:

```text
192.168.1.10 - john [10/Oct/2023:13:55:36 +0800] "GET /api/v1/users?id=42 HTTP/1.1" 200 1234 "https://example.com/" "Mozilla/5.0"
```

Script:

```python
grok(_, "%{IPORHOST:client_ip} %{NOTSPACE:http_ident} %{NOTSPACE:http_auth} \\[%{HTTPDATE:time}\\] \"%{WORD:http_method} %{GREEDYDATA:http_url} HTTP/%{NUMBER:http_version}\" %{INT:status_code:int} %{INT:bytes:int} \"%{DATA:http_referer}\" \"%{GREEDYDATA:user_agent}\"")
default_time(time)
```

Validation keys: `client_ip`, `http_method`, `http_url`; expect `status_code=200`, `bytes=1234`.

## SSH Auth Syslog

Sample:

```text
Jan 12 06:25:34 web01 sshd[12345]: Failed password for invalid user admin from 10.0.0.5 port 54321 ssh2
```

Script:

```python
grok(_, "%{NOTSPACE:month} +%{INT:day:int} %{TIME:clock} %{HOSTNAME:host} %{DATA:process}\\[%{INT:pid:int}\\]: Failed password for invalid user %{NOTSPACE:user} from %{IPORHOST:source_ip} port %{INT:source_port:int} %{WORD:protocol}")
```

Validation keys: `host`, `process`, `source_ip`; expect `pid=12345`, `source_port=54321`.

## Java Application Log

Sample:

```text
2026-05-19 11:42:31,123 ERROR [checkout] c.g.PaymentService - payment failed order_id=ord-9 latency_ms=531
```

Script:

```python
grok(_, "%{TIMESTAMP_ISO8601:time} %{LOGLEVEL:level} \\[%{DATA:service}\\] %{NOTSPACE:logger} - %{DATA:msg} order_id=%{NOTSPACE:order_id} latency_ms=%{INT:latency_ms:int}")
default_time(time)
```

Validation keys: `level`, `service`, `order_id`; expect `latency_ms=531`.

## Logfmt Message

Sample:

```text
level=info ts=2026-05-19T03:25:12Z service=worker msg="job completed" job_id=abc123 duration_ms=42
```

Script:

```python
grok(_, "level=%{WORD:level} ts=%{NOTSPACE:time} service=%{NOTSPACE:service} msg=\"%{DATA:msg}\" job_id=%{NOTSPACE:job_id} duration_ms=%{INT:duration_ms:int}")
default_time(time)
```

Validation keys: `level`, `service`, `job_id`; expect `duration_ms=42`.

## Tab-Separated Go/Zap-Style Log

Sample:

```text
2026-05-19T13:43:04.269+0800	INFO	arbiter-worker	wkr/arbiter.go:261	[workspace_uuid=wksp_bb331414b5d2497f96593551529cd73a monitor_checker_id=srul_62dcb0d154694f0191440754b7282ca9 task_id=task-jaOV8ZIeKKtH] query time range=2026/05/19 13:42:00 ~ 2026/05/19 13:43:00
```

Script:

```python
add_pattern("arbiter_prefix", "%{TIMESTAMP_ISO8601:time}%{SPACE}%{LOGLEVEL:status}%{SPACE}%{NOTSPACE:service}%{SPACE}%{DATA:source_file}:%{INT:source_line}")
grok(_, "%{arbiter_prefix}%{SPACE}\\[%{DATA:attrs}\\] query time range=%{DATA:query_start} ~ %{GREEDYDATA:query_end}")
kv_split(attrs, include_keys=["workspace_uuid", "monitor_checker_id", "task_id"])
cast(source_line, "int")
drop_key(attrs)
default_time(time)
```

Validation keys: `status`, `service`, `source_file`, `source_line`, `workspace_uuid`, `monitor_checker_id`, `task_id`, `query_start`, `query_end`; expect `source_line=261`.

Use `--script` or `--message-file` when validating this style because the sample contains real tab characters and bracket escaping. `kv_split()` makes the bracketed key-value fields order-insensitive.
