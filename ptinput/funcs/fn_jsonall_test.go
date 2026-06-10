// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"
	"time"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/GuanceCloud/pipeline-go/ptinput"
	"github.com/stretchr/testify/assert"
)

func TestJSONAllIncludeKeys(t *testing.T) {
	runner, err := NewTestingRunner(`json_all(_, include_keys=[
		"age",
		"active",
		"service",
		"name.first",
	])`)
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": `{
			"service": "api",
			"name": {"first": "Tom", "last": "Anderson"},
			"age": 37,
			"active": true,
			"children": ["Sara", "Alex", "Jack"],
			"friends": [
				{"first": "Dale", "age": 44},
				{"first": "Roger", "age": 68}
			]
		}`,
	}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	cases := map[string]any{
		"service": "api",
		"age":     float64(37),
		"active":  true,
	}
	for key, want := range cases {
		got, _, err := pt.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	}

	_, _, err = pt.Get("name")
	assert.Error(t, err)
	_, _, err = pt.Get("name.first")
	assert.Error(t, err)
	_, _, err = pt.Get("children")
	assert.Error(t, err)
	_, _, err = pt.Get("children[0]")
	assert.Error(t, err)
	_, _, err = pt.Get("friends[0].first")
	assert.Error(t, err)
}

func TestJSONAllIncludeKeysPositional(t *testing.T) {
	runner, err := NewTestingRunner(`json_all(_, ["age"])`)
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": `{"name": "Tom", "age": 37}`,
	}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	got, _, err := pt.Get("age")
	assert.NoError(t, err)
	assert.Equal(t, float64(37), got)

	_, _, err = pt.Get("name")
	assert.Error(t, err)
}

func TestJSONAllKeyPatterns(t *testing.T) {
	runner, err := NewTestingRunner(`json_all(
		_,
		include_keys=["trace_*"],
		key_patterns=["trace_?d"],
	)`)
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": `{
			"trace_*": "literal",
			"trace_id": "abc",
			"trace_span": "def",
			"service": "api"
		}`,
	}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	got, _, err := pt.Get("trace_*")
	assert.NoError(t, err)
	assert.Equal(t, "literal", got)

	got, _, err = pt.Get("trace_id")
	assert.NoError(t, err)
	assert.Equal(t, "abc", got)

	_, _, err = pt.Get("trace_span")
	assert.Error(t, err)
	_, _, err = pt.Get("service")
	assert.Error(t, err)
}

func TestJSONAllWithoutKeysDoesNothing(t *testing.T) {
	runner, err := NewTestingRunner(`json_all(_)`)
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": `{"name": "Tom", "age": 37}`,
	}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	for _, key := range []string{"name", "age"} {
		_, _, err := pt.Get(key)
		assert.Error(t, err)
	}
}

func TestJSONAllDynamicIncludeKeys(t *testing.T) {
	runner, err := NewTestingRunner(`
keys = ["service", "status", "name.first"]
json_all(_, include_keys=keys)
`)
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": `{"service": "api", "status": "ok", "name": {"first": "Tom", "last": "Anderson"}, "age": 37}`,
	}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	got, _, err := pt.Get("service")
	assert.NoError(t, err)
	assert.Equal(t, "api", got)

	got, _, err = pt.Get("status")
	assert.NoError(t, err)
	assert.Equal(t, "ok", got)

	_, _, err = pt.Get("name.first")
	assert.Error(t, err)
	_, _, err = pt.Get("age")
	assert.Error(t, err)
}

func TestJSONAllTopLevelArray(t *testing.T) {
	runner, err := NewTestingRunner(`json_all(_, include_keys=["[0]", "[1]", "[2].name"])`)
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": `["first", 2, {"name": "nested"}]`,
	}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	got, _, err := pt.Get("[0]")
	assert.NoError(t, err)
	assert.Equal(t, "first", got)

	got, _, err = pt.Get("[1]")
	assert.NoError(t, err)
	assert.Equal(t, float64(2), got)

	_, _, err = pt.Get("[2]")
	assert.Error(t, err)
	_, _, err = pt.Get("[2].name")
	assert.Error(t, err)
}

func TestJSONAllEmptyIncludeKeys(t *testing.T) {
	runner, err := NewTestingRunner(`json_all(_, include_keys=[])`)
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": `{"name": "Tom", "age": 37}`,
	}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	_, _, err = pt.Get("name")
	assert.Error(t, err)
	_, _, err = pt.Get("age")
	assert.Error(t, err)
}

func TestJSONAllChecking(t *testing.T) {
	_, err := NewTestingRunner(`json_all()`)
	assert.Error(t, err)

	_, err = NewTestingRunner(`json_all(include_keys=["age"])`)
	assert.Error(t, err)

	_, err = NewTestingRunner(`json_all(["age"])`)
	assert.Error(t, err)

	_, err = NewTestingRunner(`json_all(_, limit=3)`)
	assert.Error(t, err)

	_, err = NewTestingRunner(`json_all(_, include_keys=[1])`)
	assert.Error(t, err)

	_, err = NewTestingRunner(`json_all(_, key_patterns=[1])`)
	assert.Error(t, err)

	_, err = NewTestingRunner(`json_all(_, include_keys="age")`)
	assert.Error(t, err)

	_, err = NewTestingRunner(`json_all(_, key_patterns="trace_*")`)
	assert.Error(t, err)
}
