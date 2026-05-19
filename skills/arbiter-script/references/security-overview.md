# Security Overview

Use this reference to align Arbiter scripts with GuanceCloud security monitoring semantics.

Sources:

- https://docs.guance.com/security/
- https://docs.guance.com/security/create-detection-rules/
- https://docs.guance.com/security/signals/
- https://docs.guance.com/security/arbiter-index/

## Product Model

Security monitoring combines two capabilities:

- SIEM focuses on runtime activity security: logs, events, network activity, cloud audit events, and behavioral indicators that may reveal active threats.
- CSPM focuses on static cloud infrastructure configuration security: public buckets, security group exposure, IAM policy drift, service configuration, and compliance deviations.

Arbiter is the SIEM data analysis engine. In security detection rules, Arbiter runs user-authored scripts, queries data through `dql()`, processes results, and emits security events through `trigger()`.

## Rule Lifecycle

A security detection rule typically has:

- Detection type: SIEM or CSPM.
- Detection frequency: scheduled interval such as 1m, 5m, 15m, 30m, 1h, 6h, 12h, 24h, or crontab.
- Detection interval: query time range for each run, constrained by frequency.
- Detection logic: Arbiter script with one or more DQL/PromQL queries and trigger conditions.
- Rule description: event title and description shown to users.
- Alert policy: whether a generated event also sends notifications.

Events become signals for triage and analysis. The script should therefore preserve enough evidence in `related_data` for investigation.

## Severity Guidance

Use statuses consistently:

- `critical`: confirmed or highly likely active compromise, destructive action, exposed critical asset, or privilege escalation with strong evidence.
- `high`: serious threat or misconfiguration requiring prompt action, such as suspicious cross-region login, public sensitive storage, or broad admin policy exposure.
- `medium`: suspicious pattern or important control drift that needs investigation but has lower confidence or limited blast radius.
- `low`: weak signal, hygiene issue, or limited exposure.
- `info`: informational detection, audit trail, or non-actionable enrichment.

## Event Payload Rules

Prefer:

- `dimension_tags`: low-cardinality grouping fields such as cloud provider, account, region, rule_id, resource_type, host, service, or namespace.
- `related_data`: high-cardinality and evidence fields such as IP lists, user IDs, resource IDs, request IDs, matched rows, full DQL, counts, thresholds, and sample records.
- `check_workspace_uuid`: set only when the triggered event should be associated with a specific workspace different from the rule runtime.

Avoid:

- Triggering on every raw row unless the rule explicitly requires per-entity events.
- Putting unbounded identifiers, raw queries, or large arrays into `dimension_tags`.
- Emitting a trigger without enough related evidence to explain why the signal fired.
