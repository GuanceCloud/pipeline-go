# Arbiter Troubleshooting

Use this when `./bin/arbiter-check` fails or a generated security rule produces no event.

## Script Does Not Parse

Symptoms:

- `ok: false`
- `script.parsed: false`
- errors mention syntax, unterminated string, unexpected token, or type-check failure

Fix:

- Remember Arbiter uses Platypus syntax, same as Pipeline.
- Keep string literals on one physical line unless the language explicitly supports the escape.
- Use backticks around identifiers with special characters.
- Run `./bin/arbiter-check --script ./rule.p --parse-only`.

## Missing DQL Context

Symptom:

```text
missing context data named dql_cli
```

Fix:

- Validate with `./bin/arbiter-check`, which injects a mock DQL client.
- If using another runner, provide Arbiter DQL context through `WithDQLClient`, `WithDQLOpenAPI`, or `WithDQLKodo`.

## No Trigger Output

Symptoms:

- `ok: false` with `required trigger output not found`
- `triggers` is empty

Fix:

- Check whether the mock DQL result actually satisfies the script condition.
- Add temporary `printf()` calls and inspect `stdout`.
- Verify `dql_series_get()` returns nested lists. Iterate group and item levels.
- Validate trigger branch with a minimal mock result before using a large real response.

## Unexpected Trigger Status

Symptoms:

- `expected trigger status "high" not found`
- `triggers` exists but status differs

Fix:

- Use statuses supported by `trigger`: `critical`, `high`, `medium`, `low`, `info`.
- Keep severity mapping explicit in the script.
- If several triggers are expected, pass repeated `--expect-status` values.

## DQL Result Shape Problems

Symptoms:

- `dql_series_get()` returns empty or nested nil values.
- Script loops produce no evidence.

Fix:

- Ensure mock DQL result has `series` as a list of series, each containing point maps.
- Put fields under `columns` and dimensions under `tags`.
- Match extracted names exactly. `dql_series_get(data, "host")` searches `columns.host` first, then `tags.host`.

Minimal valid mock:

```json
{
  "series": [
    [
      {
        "columns": {"time": 1779172110000, "cnt": 3},
        "tags": {"host": "h1"}
      }
    ]
  ],
  "status_code": 200
}
```

## Live DQL Verification Problems

Symptoms:

- DQL may be syntactically wrong, but mock validation still passes.
- `./bin/arbiter-check` returns `--live-dql requires --guance-key`.
- Live response has empty `series` or an OpenAPI error in script output.

Fix:

- Use mock validation for script branching and trigger payloads.
- Use live validation for DQL syntax, permissions, index names, workspace routing, and real field names.
- Pass `--live-dql --guance-key "$GUANCE_API_KEY"` or set `GUANCE_API_KEY` / `DF_API_KEY`.
- Use `--time-range start,end` for reproducible live validation.
- Inspect `dql_queries` to confirm the script sent the intended query, qtype, limit, slimit, time range, and workspace UUID.

## DQL Check Failures

Symptoms:

- `dql_checks[].ok` is false.
- `errors` contains `dql query N failed dqlcheck`.

Fix:

- Read the matching `dql_checks[].stdout`, `stderr`, and `error`.
- Validate the query outside the script with `../dql/bin/dqlcheck --file ./query.dql --format=json --pretty`.
- Use `--file` or `--stdin` for queries containing backticks, quotes, regex, or multi-line strings.
- Fix only the DQL reported by `dqlcheck`; keep the Arbiter trigger logic unchanged unless the query shape changes.

## Trigger Payload Problems

Rules:

- `dimension_tags` only keeps string values.
- Use `related_data` for maps, lists, numbers, booleans, raw evidence, and high-cardinality identifiers.
- Use `check_workspace_uuid` only when the signal should be associated with a specific workspace.
