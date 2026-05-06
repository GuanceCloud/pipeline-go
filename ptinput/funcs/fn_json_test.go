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

func TestJSON(t *testing.T) {
	testCase := []*funcCase{
		{
			in: `{
			  "name": {"first": "Tom", "last": "Anderson"},
			  "age":37,
			  "children": ["Sara","Alex","Jack"],
			  "fav.movie": "Deer Hunter",
			  "friends": [
			    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
			    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
			    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
			  ]
			}`,
			script: `json(_, name)
			json(name, first)`,
			expected: "Tom",
			key:      "first",
		},
		{
			in: `{
			  "name": {"first": "Tom", "last": "Anderson"},
			  "age":37,
			  "children": ["Sara","Alex","Jack"],
			  "fav.movie": "Deer Hunter",
			  "friends": [
			    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
			    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
			    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
			  ]
			}`,
			script: `json(_, friends)
			json(friends, .[1].first, f_first)`,
			expected: "Roger",
			key:      "f_first",
		},
		{
			in: `[
				    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
				    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
				    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
				]`,
			script:   `json(_, .[0].nets[-1])`,
			expected: "tw",
			key:      "[0].nets[-1]",
		},
		{
			in: `[
				    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
				    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
				    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
				]`,
			script:   `json(_, .[1].age)`,
			expected: float64(68),
			key:      "[1].age",
		},
		{
			name:     "trim_space auto",
			in:       `{"item": " not_space "}`,
			script:   `json(_, item, item)`,
			key:      "item",
			expected: "not_space",
		},
		{
			name:     "trim_space disable",
			in:       `{"item": " not_space "}`,
			script:   `json(_, item, item, false)`,
			key:      "item",
			expected: " not_space ",
		},
		{
			name:     "trim_space enable",
			in:       `{"item": " not_space "}`,
			script:   `json(_, item, item, true)`,
			key:      "item",
			expected: "not_space",
		},
		{
			name:     "path_with_dot_in_key",
			in:       `{"a.b": 123}`,
			script:   "json(_, `a.b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_numeric_key",
			in:       `{"0": 123}`,
			script:   "json(_, `0`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_wildcard_in_key",
			in:       `{"a*b": 123}`,
			script:   "json(_, `a*b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_question_in_key",
			in:       `{"a?b": 123}`,
			script:   "json(_, `a?b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_array_count_in_key",
			in:       `{"a#b": 123}`,
			script:   "json(_, `a#b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_pipe_in_key",
			in:       `{"a|b": 123}`,
			script:   "json(_, `a|b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_backslash_in_key",
			in:       `{"a\\b": 123}`,
			script:   "json(_, `a\\b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_modifier_name_in_key",
			in:       `{"@this": 123}`,
			script:   "json(_, `@this`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_static_prefix_in_key",
			in:       `{"!foo": 123}`,
			script:   "json(_, `!foo`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_multipath_prefix_in_key",
			in:       `{"[key": 123}`,
			script:   "json(_, `[key`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "path_with_object_multipath_prefix_in_key",
			in:       `{"{key": 123}`,
			script:   "json(_, `{key`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_multipath_prefix_in_key",
			in:       `{"root": {"[key": 123}}`,
			script:   "json(_, root.`[key`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_object_multipath_prefix_in_key",
			in:       `{"root": {"{key": 123}}`,
			script:   "json(_, root.`{key`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_dot_in_key",
			in:       `{"root": {"a.b": 123}}`,
			script:   "json(_, root.`a.b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_wildcard_in_key",
			in:       `{"root": {"a*b": 123}}`,
			script:   "json(_, root.`a*b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_question_in_key",
			in:       `{"root": {"a?b": 123}}`,
			script:   "json(_, root.`a?b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_array_count_in_key",
			in:       `{"root": {"a#b": 123}}`,
			script:   "json(_, root.`a#b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_pipe_in_key",
			in:       `{"root": {"a|b": 123}}`,
			script:   "json(_, root.`a|b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_backslash_in_key",
			in:       `{"root": {"a\\b": 123}}`,
			script:   "json(_, root.`a\\b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_modifier_name_in_key",
			in:       `{"root": {"@this": 123}}`,
			script:   "json(_, root.`@this`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "nested_path_with_static_prefix_in_key",
			in:       `{"root": {"!foo": 123}}`,
			script:   "json(_, root.`!foo`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "special_key_map_value",
			in:       `{"a|b": {"inner": 123}}`,
			script:   "json(_, `a|b`, out)",
			key:      "out",
			expected: `{"inner":123}`,
		},
		{
			name:     "special_key_list_value",
			in:       `{"a|b": [1, 2, 3]}`,
			script:   "json(_, `a|b`, out)",
			key:      "out",
			expected: `[1,2,3]`,
		},
		{
			name:     "special_key_then_index",
			in:       `{"a|b": [1, 2, 3]}`,
			script:   "json(_, `a|b`[1], out)",
			key:      "out",
			expected: float64(2),
		},
		{
			name:     "index_then_special_key",
			in:       `[{"a|b": 123}]`,
			script:   "json(_, .[0].`a|b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name:     "negative_index_then_special_key",
			in:       `[{"a|b": 123}]`,
			script:   "json(_, .[-1].`a|b`, out)",
			key:      "out",
			expected: float64(123),
		},
		{
			name: "source_replaced_between_json_calls",
			in:   `{"a":{"first":1}}`,
			script: `json(_, a.first, first)
			add_key("message", "{\"a\":{\"second\":2}}")
			json(_, a.second, second)`,
			key:      "second",
			expected: float64(2),
		},
		{
			name:     "map_delete_after",
			in:       `{"item": " not_space "}`,
			script:   `json(_, item, item, true, true)`,
			key:      "message",
			expected: "{}",
			fail:     false,
		},
		{
			name:     "map_delete_after1",
			in:       `{"item": " not_space ", "item2":{"item3": [123]}}`,
			script:   `json(_, item2.item3, item, delete_after_extract = true)`,
			key:      "message",
			expected: `{"item":" not_space ","item2":{}}`,
		},
		{
			name:     "list_delete_after1",
			in:       `{"item": " not_space ", "item2": [[1,2,3,4,5],[6]]}`,
			script:   `json(_, .[0].item2[0][2].a[0], item, true, true)`,
			key:      "item",
			expected: "1",
			fail:     true,
		},
		{
			name:     "list_delete_after2",
			in:       `{"item": " not_space ", "item2": [[1,2,3,4,5],[6]]}`,
			script:   `json(_, .[0], item, true, true)`,
			key:      "item",
			expected: "1",
			fail:     true,
		},
		{
			name:     "list_delete_after3",
			in:       `{"item": " not_space ", "item2": [[1,2,3,4,5],[6]]}`,
			script:   `json(_, a[0][1], item, true, true)`,
			key:      "item",
			expected: "1",
			fail:     true,
		},
	}

	for idx, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			runner, err := NewTestingRunner(tc.script)

			if err != nil && tc.fail {
				return
			} else if err != nil || tc.fail {
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.fail, err != nil)
			}

			pt := ptinput.NewPlPt(
				point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())
			errR := runScript(runner, pt)
			if errR != nil {
				t.Fatal(errR.Error())
			}

			r, _, e := pt.Get(tc.key)
			assert.NoError(t, e)
			if tc.key == "[2].age" {
				t.Log(1)
			}
			assert.Equal(t, tc.expected, r)

			t.Logf("[%d] PASS", idx)
		})
	}
}

func TestJSONIndexDoesNotMatchObjectNumericKey(t *testing.T) {
	runner, err := NewTestingRunner(`json(_, .[0], out)`)
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(
		point.Logging, "test", nil, map[string]any{"message": `{"0":123}`}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	_, _, err = pt.Get("out")
	assert.Error(t, err)
}

func TestJSONNumericKeyDoesNotMatchArrayIndex(t *testing.T) {
	runner, err := NewTestingRunner("json(_, `0`, out)")
	assert.NoError(t, err)

	pt := ptinput.NewPlPt(
		point.Logging, "test", nil, map[string]any{"message": `[123]`}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	_, _, err = pt.Get("out")
	assert.Error(t, err)
}

func BenchmarkJSON(b *testing.B) {
	runner, err := NewTestingRunner(`json(_, friends)
json(friends, .[1].first, f_first)`)
	if err != nil {
		b.Fatal(err)
	}

	in := `{
	  "name": {"first": "Tom", "last": "Anderson"},
	  "age":37,
	  "children": ["Sara","Alex","Jack"],
	  "fav.movie": "Deer Hunter",
	  "friends": [
	    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
	    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
	    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
	  ]
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt := ptinput.NewPlPt(
			point.Logging, "test", nil, map[string]any{"message": in}, time.Now())
		if err := runScript(runner, pt); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONRepeatedSource(b *testing.B) {
	runner, err := NewTestingRunner(`json(_, a.first, first)
json(_, a.second, second)
json(_, a.third, third)
json(_, a.forth, forth)`)
	if err != nil {
		b.Fatal(err)
	}

	in := `{"a":{"first": 2.3, "second":2,"third":"aBC","forth":true},"age":47}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt := ptinput.NewPlPt(
			point.Logging, "test", nil, map[string]any{"message": in}, time.Now())
		if err := runScript(runner, pt); err != nil {
			b.Fatal(err)
		}
	}
}
