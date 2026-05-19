---
name: arbiter-script
description: Create, review, and validate GuanceCloud Arbiter scripts for security detection rules. Use when Codex needs to author SIEM or CSPM detection logic, query data with dql or PromQL, process DQL series, call trigger to generate security events, inspect Arbiter built-in functions, or debug Arbiter syntax/runtime/trigger-output problems.
---

# Arbiter Script

## Workflow

1. Classify the rule before writing code:
   - SIEM: runtime activity detection from logs, events, metrics, traces, network, or cloud audit data.
   - CSPM: static cloud/resource configuration posture checks, baseline drift, and compliance gaps.
2. Define the detection contract: data source and DQL/PromQL query, time range, grouping dimensions, trigger condition, severity, dimension tags, related data, and optional target workspace.
   - For `dql()` query authoring, repairing, or final delivery, also use the local `dql` skill workflow. Do not use a generic SQL skill to generate DQL; DQL has different grammar and must pass `dqlcheck`.
3. Query function docs before guessing syntax:

```sh
./bin/arbiter-check --search-functions dql --function-lang all
./bin/arbiter-check --function-doc trigger
```

4. For any final executable DQL, follow the `dql` skill rules. Run `../dql/bin/dqldocs` once before reading DQL references, and validate each DQL with `../dql/bin/dqlcheck` or `./bin/arbiter-check --check-dql`.
5. Write scripts in Platypus syntax. Use `dql()` to fetch data, `dql_series_get()` for common column/tag extraction, explicit loops when row-level context matters, and `trigger()` only when the condition is met.
   - Alias every DQL expression to the exact simple field name read by `dql_series_get()`. Do not read `accountId` if the DQL aliases the field as `userId`.
   - Validate both positive and negative branches with mock DQL results when status/severity depends on an extracted value.
6. Keep `dimension_tags` small and stable. Put high-cardinality evidence such as IPs, user IDs, resource IDs, request IDs, matched rows, and raw query details into `related_data`.
7. Validate with the skill-local checker. Use mock DQL for script logic and trigger payload tests:

```sh
./bin/arbiter-check \
  --script ./rule.p \
  --dql-result-file ./sample-dql-result.json \
  --check-dql \
  --require-trigger \
  --expect-status high
```

When OpenAPI credentials are available and the task requires real DQL verification, run live DQL:

```sh
./bin/arbiter-check \
  --script ./rule.p \
  --live-dql \
  --guance https://openapi.guance.com \
  --guance-key "$GUANCE_API_KEY" \
  --check-dql \
  --duration 15m
```

For syntax-only checks:

```sh
./bin/arbiter-check --script ./rule.p --parse-only
```

8. If validation fails, fix DQL check failures first, then parse/type/runtime errors, then trigger assertions. Verify `dql_checks`, `dql_mode`, `dql_queries`, `stdout`, and `triggers` in the JSON output. Do not claim DQL is executable unless every `dql_checks[].ok` is true or each query passed `dqlcheck` individually. Do not claim DQL is verified against real data unless `dql_mode` is `live` and `./bin/arbiter-check` exits successfully.
9. Report the final script, the validation command, and the important execution output: `dql_checks`, `triggers`, `dql_queries`, `stdout`, and any `errors`.

## References

Read [security-overview.md](references/security-overview.md) when deciding SIEM versus CSPM behavior, event lifecycle, severity, or signal/event handling.

Read [arbiter-patterns.md](references/arbiter-patterns.md) when writing detection scripts, choosing `dql_series_get()` versus loops, building trigger payloads, or preparing mock DQL results.

Read [cloud-security-patterns.md](references/cloud-security-patterns.md) when writing cloud audit rules such as AWS CloudTrail, root account activity, IAM changes, console login, cloud posture, or account/resource evidence.

Read [troubleshooting.md](references/troubleshooting.md) when the checker reports parse errors, missing DQL context, no trigger output, unexpected status, or DQL result shape issues.
