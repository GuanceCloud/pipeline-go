---
name: pipeline-script
description: Create, refine, and validate GuanceCloud/DataKit pipeline scripts from raw message samples. Use when Codex needs to author pipeline .p scripts, choose json/load_json versus grok extraction for text or JSON messages, run pipeline-go validation, or debug parsing/runtime/extraction errors.
---

# Pipeline Script

## Workflow

1. Treat the sample input as a single `message` value unless the user provides more point context. Pipeline's original input `_` maps to the `message` field.
2. Inspect the message shape before writing code:
   - Valid JSON object/array: prefer `json(_, path, target_key)` for a small fixed set of fields.
   - JSON that needs repeated access, mutation, or conditionals: use `data = load_json(_)`, then read `data["field"]`.
   - Plain text or log lines: use `grok(_, "...")`; add `:int`, `:float`, or `:bool` in grok captures when the type is known.
   - Text with an embedded JSON fragment: grok the fragment into a temporary key, then use `json(temp_key, ...)` or `load_json(temp_key)`.
   - Text containing `key=value` segments: grok the surrounding structure, then use `kv_split()` on the segment so key order can vary.
3. If function syntax is uncertain, query embedded docs before guessing:

```sh
./bin/pipeline-check --search-functions json --function-lang all
./bin/pipeline-check --function-doc grok
```

4. Write the script with stable field names. Keep high-cardinality values as fields unless the user asks for tags. Use `default_time(time)` after extracting a timestamp intended to become point time.
5. For non-trivial grok, split long regex fragments into several short `add_pattern()` definitions and validate from a `.p` file with `--script`. Keep every quoted pipeline string on one physical line; never wrap inside the quotes. Use inline `--cmd` only for short scripts or quick experiments.
6. Validate with the skill-local checker, always requiring the fields that prove extraction worked:

```sh
./bin/pipeline-check --script ./example.p --message 'raw message here' --require-key service
```

For inline drafts:

```sh
./bin/pipeline-check --cmd 'json(_, service)' --message '{"service":"api"}' --require-key service
```

7. If validation returns `ok: false`, fix the script and rerun. A grok mismatch can otherwise look like a successful run with missing fields, so use `--require-key` for every expected output key and `--expect key=value` for important sample values.
8. Report the final `.p` script content and the checker execution result. Prefer `result.extracted_fields`, `result.extracted_tags`, `result.time`, and `errors`; the original input is available as `result.extracted_fields.message`. Inspect `output.fields`/`output.tags` only when the full transformed point is needed. Do not include shell line-continuation backslashes in the script itself. Do not claim the script is validated unless `./bin/pipeline-check` exits successfully.

## References

Read [pipeline-patterns.md](references/pipeline-patterns.md) when you need function syntax, extraction examples, or the full checker option list.

Read [production-guidelines.md](references/production-guidelines.md) when the user asks for production-ready scripts, robust parsing, multiple samples, or avoiding overfitting.

Read [common-log-examples.md](references/common-log-examples.md) when you need validated starter scripts for common JSON app logs, Nginx access logs, SSH/syslog lines, Java application logs, or logfmt messages.

Read [troubleshooting.md](references/troubleshooting.md) when parser errors mention `unknown escape sequence`, grok brackets, copied shell continuations, missing `time`, or required keys not found.

Read [datakit-pipeline-examples.md](references/datakit-pipeline-examples.md) when `/usr/local/datakit/pipeline` is available and you need the built-in script inventory, shared patterns, caveats, or routing to narrower DataKit references.

For DataKit built-in examples, prefer the narrowest reference:
- [datakit-web-examples.md](references/datakit-web-examples.md): `nginx.p`, `apache.p`, `tomcat.p`.
- [datakit-database-examples.md](references/datakit-database-examples.md): `mysql.p`, `mongodb.p`, `postgresql.p`, `dameng.p`, `kingbase.p`, `sqlserver.p`.
- [datakit-middleware-examples.md](references/datakit-middleware-examples.md): `redis.p`, `elasticsearch.p`, `kafka.p`, `rabbitmq.p`, `solr.p`, `consul.p`, `jenkins.p`, `tdengine.p`.
