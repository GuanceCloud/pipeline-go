# Pipeline Authoring Patterns

## JSON Message

Use `json()` when the message is JSON and fields are known:

```python
json(_, service)
json(_, status, status_code)
json(_, trace.id, trace_id)
json(_, items[0].name, first_item_name)
```

Use `load_json()` when the script benefits from a parsed object:

```python
data = load_json(_)
add_key(service, data["service"])
add_key(status_code, data["status"])
```

Validate:

```sh
./bin/pipeline-check \
  --cmd 'json(_, service)
json(_, status, status_code)' \
  --message '{"service":"api","status":200}' \
  --require-key service \
  --expect status_code=200
```

## Text Message With Grok

Prefer typed captures for numeric values:

```python
grok(_, "%{IPORHOST:client_ip} %{NOTSPACE:ident} %{NOTSPACE:auth} \\[%{HTTPDATE:time}\\] \"%{WORD:method} %{NOTSPACE:path} HTTP/%{NUMBER:http_version}\" %{INT:status_code:int} %{INT:bytes:int}")
default_time(time)
```

For long grok patterns, use several short `add_pattern()` fragments and validate from a script file. Keep each quoted string on one physical line; do not insert line breaks inside `"..."`. This avoids shell escaping issues with `\\[`, quotes, tabs, and multiline samples.

Validate with required keys:

```sh
./bin/pipeline-check \
  --script ./nginx.p \
  --message '192.168.1.10 - - [19/Jun/2021:04:04:58 +0000] "GET / HTTP/1.1" 200 612' \
  --require-key client_ip \
  --require-key status_code \
  --expect status_code=200
```

## Embedded JSON

When a line wraps JSON inside text, extract the JSON substring first:

```python
grok(_, "%{TIMESTAMP_ISO8601:time} %{LOGLEVEL:level} payload=%{GREEDYDATA:payload_json}")
json(payload_json, user.id, user_id)
json(payload_json, action)
default_time(time)
```

## Key-Value Segment Inside Text

When a text log contains a delimited `key=value` segment, capture the segment and parse it with `kv_split()` instead of hardcoding key order into grok:

```python
grok(_, "%{TIMESTAMP_ISO8601:time}%{SPACE}\\[%{DATA:attrs}\\] %{GREEDYDATA:msg}")
kv_split(attrs, include_keys=["workspace_uuid", "task_id"])
drop_key(attrs)
default_time(time)
```

## Common Functions

- `json(input, path, target_key)`: Extract one JSON path from `_` or a field.
- `load_json(value)`: Parse JSON into map/list for expression access.
- `grok(input, pattern)`: Extract fields from text. Use typed captures like `%{INT:code:int}`.
- `kv_split(key, include_keys=[...])`: Extract selected key-value pairs from a string; use it for order-insensitive key-value segments.
- `cast(key, "int"|"float"|"bool"|"str")`: Convert after extraction when typed grok capture is not used.
- `default_time(key)`: Convert an extracted time into point time.
- `set_tag(key)`: Promote a field to a tag only when the value is low-cardinality and intended for indexing.
- `drop_key(key)` / `delete(key)`: Remove temporary extraction fields after use.

## Pipeline Check

Use this skill's `./bin/pipeline-check` wrapper. It runs a skill-local `bin/pipeline-go` binary when present, otherwise falls back to the local `pipeline-go` repository or a `pipeline-go` binary in `PATH`.

Important options:

- `--script FILE` or `-P FILE`: Read a script file.
- `--cmd TEXT` or `-c TEXT`: Validate an inline script.
- `--message TEXT` or `-T TEXT`: Run with a single `message` field.
- `--message-file FILE`: Read the message from a file.
- `--category logging`: Point category; defaults to `logging`.
- `--name NAME`: Point name; defaults to `pipeline_check`.
- `--tag key=value`: Add an input tag; repeatable.
- `--field key=value`: Add an input field; JSON-looking values are decoded.
- `--require-key key`: Fail if the output tag/field is missing.
- `--expect key=value`: Fail if the output value differs; JSON-looking expected values are decoded.
- `--add-status`: Enable logging status normalization during the run.

Practical validation rules:

- Prefer `--script FILE` for grok-heavy scripts.
- Use `--cmd TEXT` only for short scripts or quick JSON extraction checks.
- Use `--message-file FILE` for multiline messages or messages containing real tab characters.
- Do not copy shell line-continuation backslashes into `.p` scripts.

Function doc options:

- `--list-functions`: List embedded pipeline function docs with signatures and summaries.
- `--search-functions QUERY` or `--fn-search QUERY`: Search function names, signatures, summaries, and markdown content.
- `--function-doc NAME` or `--fn-doc NAME`: Show the full markdown doc for one function.
- `--function-lang zh|en|all`: Choose function doc language; defaults to `zh`.
- `--function-limit N`: Limit list/search result count; `0` means no limit.

Examples:

```sh
./bin/pipeline-check --search-functions json --function-lang all
./bin/pipeline-check --function-doc grok
```

Output JSON contains:

- `ok`: Overall validation result.
- `input.message_is_json`, `input.message_json_type`, `input.message_json_keys`: Message inspection hints.
- `output.fields` and `output.tags`: The transformed point after script execution.
- `result.extracted_fields`: Fields produced or changed by the script, plus the original input as `message`.
- `result.extracted_tags`: Tags produced or changed by the script, excluding unchanged input tags.
- `result.time`, `result.time_unix_nano`, `result.dropped`: Execution summary values from the transformed point.
- `errors`: Parse, runtime, missing-key, or expectation failures.
