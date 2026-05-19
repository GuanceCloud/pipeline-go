// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2026-present Guance, Inc.

package arbitercheck

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRunTrigger(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--cmd", `trigger(1, "critical", {"host":"h1"}, {"reason":"test"})`,
		"--require-trigger",
		"--expect-status", "critical",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res checkResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if !res.OK || len(res.Triggers) != 1 {
		t.Fatalf("unexpected result: %+v", res)
	}
	if got := res.Triggers[0].DimensionTags["host"]; got != "h1" {
		t.Fatalf("host tag = %#v", got)
	}
}

func TestRunMockDQL(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--cmd", `data = dql("L::` + "`default`" + `:(count())")
hosts = dql_series_get(data, "host")
trigger(len(hosts), "high", {}, {"hosts": hosts})`,
		"--dql-result", `{"series":[[{"columns":{"time":1779172110000},"tags":{"host":"h1"}}]],"status_code":200}`,
		"--require-trigger",
		"--expect-status", "high",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res checkResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if len(res.DQLQueries) != 1 {
		t.Fatalf("expected one DQL query: %+v", res.DQLQueries)
	}
	if res.DQLMode != "mock" || res.DQLQueries[0].QType != "dql" || res.Triggers[0].Status != "high" {
		t.Fatalf("unexpected execution result: %+v", res)
	}
}

func TestRunChecksDQL(t *testing.T) {
	dir := t.TempDir()
	checker := filepath.Join(dir, "dqlcheck")
	if err := os.WriteFile(checker, []byte(`#!/usr/bin/env sh
cat >/dev/null
printf '{"ok":true,"out":"check","build":"ok"}\n'
`), 0o755); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--cmd", `data = dql("L::` + "`default`" + `:(count() as cnt)")
trigger(1, "high", {}, {})`,
		"--dql-result", `{"series":[],"status_code":200}`,
		"--check-dql",
		"--dqlcheck-bin", checker,
		"--require-trigger",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res checkResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if len(res.DQLChecks) != 1 || !res.DQLChecks[0].OK {
		t.Fatalf("unexpected dql checks: %+v", res.DQLChecks)
	}
}

func TestRunLiveDQL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":200,"content":{"data":[{"series":[{"columns":["time","cnt"],"values":[[1779172110000,3]],"tags":{"host":"h1"},"name":"default"}]}]},"errorCode":"","message":"","success":true}`))
	}))
	defer server.Close()

	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--cmd", `data = dql("L::` + "`default`" + `:(count() as cnt) BY ` + "`host`" + `")
hosts = dql_series_get(data, "host")
counts = dql_series_get(data, "cnt")
trigger(counts[0][0], "high", {"host": hosts[0][0]}, {"count": counts[0][0]})`,
		"--live-dql",
		"--guance", server.URL,
		"--guance-key", "test-key",
		"--time-range", "1779171210000,1779172110000",
		"--require-trigger",
		"--expect-status", "high",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res checkResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.DQLMode != "live" || res.MockDQLResult != nil {
		t.Fatalf("unexpected dql mode/result: %+v", res)
	}
	if len(res.DQLQueries) != 1 || res.DQLQueries[0].Query == "" {
		t.Fatalf("expected recorded DQL query: %+v", res.DQLQueries)
	}
	if len(res.Triggers) != 1 || res.Triggers[0].DimensionTags["host"] != "h1" {
		t.Fatalf("unexpected trigger output: %+v", res.Triggers)
	}
}

func TestFunctionDocSearch(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--search-functions", "trigger",
		"--function-limit", "3",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res functionDocResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if !res.OK || res.Mode != "search" || len(res.Functions) == 0 {
		t.Fatalf("unexpected search result: %+v", res)
	}
	if res.Functions[0].Name != "trigger" {
		t.Fatalf("expected trigger first: %+v", res.Functions)
	}
}
