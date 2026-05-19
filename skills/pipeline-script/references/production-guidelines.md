# Production Pipeline Guidelines

Use this when the target script should be production-ready or when the user explicitly asks to avoid overfitting.

## Acceptance Criteria

A production candidate should have:

- An output contract: required fields, optional fields, expected types, and which fields become tags.
- At least two positive samples when available: normal case plus one variant such as reordered fields, missing optional data, alternate status, or different path.
- A validation command for each sample using `./bin/pipeline-check`.
- No hardcoded sample IDs, timestamps, hosts, UUIDs, or task names unless they are literal protocol tokens.
- No leftover temporary extraction fields unless they are intentionally part of output.

## Parsing Strategy

Choose the least fragile parser:

- JSON payload: use `json()` for fixed fields; use `load_json()` for repeated reads or logic.
- Stable text envelope with variable key-value interior: use `grok()` for the envelope and `kv_split()` for the key-value segment.
- Pure key-value/logfmt text: prefer `kv_split()` with `include_keys`.
- Free-form text: use `grok()` with delimiters and required anchors.

Avoid these overfitting patterns:

- Capturing exact sample values such as `wksp_...`, `task-...`, or concrete dates.
- One long grok with many ordered key-value captures when the key order may change.
- `GREEDYDATA` before fields that must be captured later, unless it is bounded by a clear delimiter.
- Promoting high-cardinality IDs to tags by default.
- Requiring `time` after `default_time(time)`; that field is removed.

## Field And Tag Rules

- Keep UUIDs, trace IDs, task IDs, query strings, request bodies, and user IDs as fields unless the user explicitly wants them as tags.
- Tags should be low-cardinality identifiers such as service, environment, source, host, component, or status.
- Cast numeric fields that downstream users will aggregate or compare.
- Normalize status only when the desired mapping is known; otherwise preserve the original level/status.

## Validation Matrix

For each candidate script, run at least:

```sh
./bin/pipeline-check --script ./candidate.p --message-file ./sample-1.log --require-key service
./bin/pipeline-check --script ./candidate.p --message-file ./sample-2.log --require-key service
```

Use `--expect key=value` for representative values, not every field. Use `--require-key` for the contract.

When only one sample is available, synthesize one minimal structural variant without changing semantics. Example: reorder key-value pairs or change an ID/date/path. State that this is synthetic validation.

## Review Checklist

Before presenting the script:

- Verify all required fields exist.
- Verify numeric casts produce numeric JSON values.
- Verify `result.extracted_fields` / `result.extracted_tags` contain the intended parsed values; inspect `output.fields` / `output.tags` for full point state.
- Verify `result.time` changed when time parsing is expected.
- Verify optional missing fields do not fail the script unless they are required.
- Verify temporary fields such as `attrs`, `payload_json`, or `tmp` are dropped when not needed.
- Report validation commands and whether samples were real or synthetic.
