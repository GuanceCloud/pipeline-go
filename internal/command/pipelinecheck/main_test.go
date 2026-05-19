// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2026-present Guance, Inc.

package pipelinecheck

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRunJSONMessage(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--cmd", "json(_, service)\njson(_, status, status_code)",
		"--message", `{"service":"api","status":200}`,
		"--require-key", "service",
		"--expect", "status_code=200",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res checkResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if !res.OK {
		t.Fatalf("expected ok result: %+v", res)
	}
	if res.Input == nil || !res.Input.MessageIsJSON {
		t.Fatalf("expected JSON message inspection: %+v", res.Input)
	}
	if got := res.Output.Fields["service"]; got != "api" {
		t.Fatalf("service = %#v", got)
	}
	if res.Result == nil {
		t.Fatalf("expected execution result summary")
	}
	if got := res.Result.ExtractedFields["message"]; got != `{"service":"api","status":200}` {
		t.Fatalf("result message = %#v", got)
	}
	if got := res.Result.ExtractedFields["service"]; got != "api" {
		t.Fatalf("result service = %#v", got)
	}
	if res.Result.FieldCount != len(res.Result.ExtractedFields) {
		t.Fatalf("unexpected field count: %+v", res.Result)
	}
}

func TestRunMissingRequiredKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--cmd", "json(_, service)",
		"--message", `{"service":"api"}`,
		"--require-key", "missing",
	}, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("expected non-zero exit code, stdout = %s", stdout.String())
	}

	var res checkResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.OK {
		t.Fatalf("expected failed result: %+v", res)
	}
	if len(res.Errors) == 0 || !strings.Contains(res.Errors[0], "required key") {
		t.Fatalf("expected required key error: %+v", res.Errors)
	}
}

func TestRunResultIncludesChangedInputField(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--cmd", `add_key(foo, "new")`,
		"--message", "x",
		"--field", "foo=old",
		"--expect", "foo=new",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res checkResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if got := res.Result.ExtractedFields["foo"]; got != "new" {
		t.Fatalf("result foo = %#v", got)
	}
}

func TestFunctionDocSearch(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--search-functions", "json",
		"--function-lang", "all",
		"--function-limit", "3",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res functionDocResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if !res.OK || res.Mode != "search" {
		t.Fatalf("unexpected search result: %+v", res)
	}
	if len(res.Functions) == 0 || res.Functions[0].Name != "json" {
		t.Fatalf("expected json to be first search result: %+v", res.Functions)
	}
	if res.Functions[0].Snippet == "" || res.Functions[0].Markdown != "" {
		t.Fatalf("expected compact search result: %+v", res.Functions[0])
	}
}

func TestFunctionDocShow(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{
		"--function-doc", "grok",
		"--function-lang", "zh",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s, stdout = %s", code, stderr.String(), stdout.String())
	}

	var res functionDocResult
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if !res.OK || len(res.Functions) != 1 {
		t.Fatalf("unexpected doc result: %+v", res)
	}
	doc := res.Functions[0]
	if doc.Name != "grok" || doc.Signature == "" || doc.Markdown == "" {
		t.Fatalf("unexpected grok doc: %+v", doc)
	}
}
