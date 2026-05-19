# Arbiter Patterns

Use this reference when authoring or validating Arbiter scripts.

## Core Functions

- `dql(query, qtype="dql", limit=10000, offset=0, slimit=0, time_range=[], workspace_uuid="") -> map`: query GuanceCloud data with DQL or PromQL.
- `dql_series_get(series, name) -> list`: extract a column or tag from every point in a DQL result. It returns a nested list grouped by series.
- `dql_timerange_get() -> list`: get the caller-provided query time range in milliseconds, or the default last 15 minutes.
- `trigger(result, status="", dimension_tags={}, related_data={}, check_workspace_uuid="")`: emit a security event.
- `printf(format, ...)`: debug output captured in `./bin/arbiter-check` JSON field `stdout`.
- `dump_json(value, indent="")`: render maps/lists for debug output.

Function docs:

```sh
./bin/arbiter-check --list-functions --function-limit 20
./bin/arbiter-check --function-doc dql
./bin/arbiter-check --function-doc trigger
```

## DQL Result Shape

`dql()` returns a map commonly shaped like:

```json
{
  "series": [
    [
      {
        "columns": {"time": 1779172110000, "count": 3},
        "tags": {"host": "host-a"}
      }
    ]
  ],
  "status_code": 200
}
```

`dql_series_get(data, "host")` returns nested lists, for example `[["host-a"]]`.

Use `dql_series_get()` when only one field is needed. Use explicit loops when you need to preserve row context, combine columns/tags, filter nil values, or construct sample evidence.

## DQL Alias And Series Safety

- Alias DQL expressions to stable simple names, then read those exact names with `dql_series_get()`.
- Avoid reading raw dotted or `@` field names from scripts unless the DQL result is known to expose them as tags.
- When grouping by a nested field, also project it with an alias if the script needs it in `related_data`.
- `dql_series_get()` returns nested lists. Check list lengths before reading sibling fields with the same indexes.
- Prefer C-style indexed loops when reading multiple nested series in parallel. Platypus supports `for item in list`, but not `for i, item in list`.
- Test every severity branch with mock DQL results. Repeated `if`/`elif` conditions often hide a branch that can never run.

## Count/Threshold Rule

```python
data = dql("L::`default`:(count(*) as cnt) BY `source`")
counts = dql_series_get(data, "cnt")

total = 0
for group in counts {
    for cnt in group {
        if cnt != nil {
            total = total + cnt
        }
    }
}

if total > 100 {
    trigger(total, "high", {"rule_type": "log_volume"}, {
        "count": total,
        "threshold": 100,
    })
}
```

## Entity Evidence Rule

```python
data = dql("L::`default`:(count(*) as cnt) BY `host`, `source`")
hosts = dql_series_get(data, "host")
counts = dql_series_get(data, "cnt")

evidence = []
for i=0; i<len(hosts); i+=1 {
    group = hosts[i]
    for j=0; j<len(group); j+=1 {
        host = group[j]
        cnt = counts[i][j]
        if cnt != nil && cnt > 10 {
            evidence = append(evidence, {"host": host, "count": cnt})
        }
    }
}

if len(evidence) > 0 {
    trigger(len(evidence), "medium", {"rule_type": "entity_anomaly"}, {
        "matched": evidence,
    })
}
```

## CSPM Posture Rule

```python
data = dql("O::`cloud_bucket`:(count(*) as cnt) BY `provider`, `account_id`, `region`")
providers = dql_series_get(data, "provider")
accounts = dql_series_get(data, "account_id")
counts = dql_series_get(data, "cnt")

for i=0; i<len(counts); i+=1 {
    group = counts[i]
    for j=0; j<len(group); j+=1 {
        cnt = group[j]
        if cnt != nil && cnt > 0 {
            trigger(cnt, "high", {
                "rule_type": "cspm",
                "provider": providers[i][j],
                "account_id": accounts[i][j],
            }, {
                "count": cnt,
                "reason": "cloud resource violates posture policy",
            })
        }
    }
}
```

## Offline Validation

Use a mock DQL result:

```sh
./bin/arbiter-check \
  --script ./rule.p \
  --dql-result '{"series":[[{"columns":{"cnt":3},"tags":{"host":"h1"}}]],"status_code":200}' \
  --check-dql \
  --require-trigger \
  --expect-status high
```

## Live DQL Validation

Use live mode only when credentials are available and the DQL itself must be checked against real backend data:

```sh
./bin/arbiter-check \
  --script ./rule.p \
  --live-dql \
  --guance https://openapi.guance.com \
  --guance-key "$GUANCE_API_KEY" \
  --check-dql \
  --duration 15m \
  --require-trigger
```

Use `--time-range start,end` with millisecond timestamps for deterministic live checks:

```sh
./bin/arbiter-check --script ./rule.p --live-dql --guance-key "$GUANCE_API_KEY" --check-dql --time-range 1779171210000,1779172110000
```

`--guance-key` defaults to `GUANCE_API_KEY` or `DF_API_KEY`. `--guance` defaults to `https://openapi.guance.com` and can be overridden with `GUANCE_OPENAPI_ENDPOINT`.

## Standalone DQL Validation

Use the local `dql` skill when producing or repairing any DQL. Do not route DQL through a generic SQL skill; DQL is only accepted after `dqlcheck` passes.

```sh
../dql/bin/dqldocs
../dql/bin/dqlcheck -q 'M::cpu:(usage_total) BY host'
../dql/bin/dqlcheck --file ./query.dql --format=json --pretty
```

For complex quoting, prefer `--file` or `--stdin`.

Important output fields:

- `ok`: overall parse/run/assertion result.
- `dql_mode`: `mock` or `live`.
- `stdout`: text produced by `printf()`.
- `triggers`: events emitted by `trigger()`.
- `dql_queries`: every `dql()` call with query, qtype, limits, explicit time range, and workspace UUIDs.
- `dql_checks`: `dqlcheck` result for each recorded DQL when `--check-dql` is used.
- `mock_dql_result`: mock response used by `dql()` when `dql_mode` is `mock`.
- `errors`: parse, type-check, runtime, missing trigger, or status assertion failures.

Use `--parse-only` for scripts that require live DQL or HTTP and cannot be safely run offline.
