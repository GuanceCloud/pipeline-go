# Pipeline Troubleshooting

Use this reference when `./bin/pipeline-check` fails or the generated script is copied into DataKit and reports parser/runtime errors.

## Parser Errors

### `unknown escape sequence U+000A`

Cause: a backslash immediately before a newline was copied into the `.p` script, usually from a shell command like:

```sh
./bin/pipeline-check \
  --script example.p
```

Fix:

- Do not put shell continuation backslashes in pipeline scripts.
- Keep a grok string on one line, or split reusable pieces with `add_pattern()`.
- Prefer saving the script as `example.p` and validating with `--script example.p`.

### `unknown escape sequence U+005B '['`

Cause: the script contains `\[` inside a double-quoted pipeline string. Pipeline strings need `\\[` to represent a regex literal `[`.

Correct script content:

```python
grok(_, "\\[%{DATA:value}\\]")
```

When using shell inline `--cmd`, an extra shell escaping layer is often needed. Avoid that by writing the `.p` file and running `--script`.

## Grok Did Not Match

Symptoms:

- `./bin/pipeline-check` exits `ok: false` only because `--require-key` failed.
- Script parsed and ran, but expected fields are missing.

Fix:

- Require at least one field from every important pattern branch with `--require-key`.
- Add typed captures for important numbers: `%{INT:status_code:int}`.
- Use `%{SPACE}` for tabs/spaces unless the format requires literal `\t`.
- Split long patterns into short fragments. Keep every quoted pipeline string on one physical line:

```python
add_pattern("prefix", "%{TIMESTAMP_ISO8601:time}%{SPACE}%{LOGLEVEL:status}")
add_pattern("source", "%{NOTSPACE:component}%{SPACE}%{DATA:file_path}:%{INT:line:int}")
add_pattern("ids", "\\[%{DATA:request_id}\\]")
grok(_, "%{prefix}%{SPACE}%{source}%{SPACE}%{ids}%{SPACE}%{GREEDYDATA:msg}")
```

## Time Field Is Missing

`default_time(time)` converts the extracted `time` field into point time and deletes the `time` field from output fields.

Do not validate with:

```sh
--require-key time
```

Instead, inspect `output.time` / `output.time_unix_nano` in the JSON result and require other extracted fields.

## Inline Validation Guidance

Use `--cmd` only for small scripts:

```sh
./bin/pipeline-check --cmd 'json(_, service)' --message '{"service":"api"}' --require-key service
```

Use a `.p` file for grok scripts with brackets, quotes, tabs, or multiline samples:

```sh
./bin/pipeline-check --script ./example.p --message-file ./sample.log --require-key service
```
