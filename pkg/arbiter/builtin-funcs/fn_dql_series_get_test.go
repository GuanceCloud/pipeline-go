package funcs

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/dql"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

func TestFuncDQLSeries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockOpenAPIMetric))
	defer server.Close()
	cases := []ProgCase{}
	cases = append(cases, cDQLSeriesGet.Progs...)
	cases = append(cases, cDQLSeriesFirst.Progs...)
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runCase(t, tc, map[runtimev2.TaskP]any{
				PDQLCli: dql.NewDQLOpenAPI(
					server.URL,
					dql.OpenAPIPath,
					"abc", nil,
				),
			})
		})
	}
}

func TestDQLSeriesFirstValue(t *testing.T) {
	cases := []struct {
		name   string
		result map[string]any
		field  string
		want   any
	}{
		{
			name: "column wins over tag on first point",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"host": "column-host"},
							"tags":    map[string]any{"host": "tag-host"},
						},
					},
				},
			},
			field: "host",
			want:  "column-host",
		},
		{
			name: "tag fallback on first point",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"time": int64(1)},
							"tags":    map[string]any{"host": "tag-host"},
						},
					},
				},
			},
			field: "host",
			want:  "tag-host",
		},
		{
			name: "missing first point returns nil",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"time": int64(1)},
						},
						map[string]any{
							"columns": map[string]any{"host": "later-host"},
						},
					},
				},
			},
			field: "host",
			want:  nil,
		},
		{
			name: "missing first point tag returns nil",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"time": int64(1)},
						},
						map[string]any{
							"tags": map[string]any{"host": "later-host"},
						},
					},
				},
			},
			field: "host",
			want:  nil,
		},
		{
			name: "skip empty series",
			result: map[string]any{
				"series": []any{
					[]any{},
					[]any{
						map[string]any{
							"columns": map[string]any{"host": "second-series-host"},
						},
					},
				},
			},
			field: "host",
			want:  "second-series-host",
		},
		{
			name: "empty result returns nil",
			result: map[string]any{
				"series": []any{},
			},
			field: "host",
			want:  nil,
		},
		{
			name:   "missing series returns nil",
			result: map[string]any{},
			field:  "host",
			want:   nil,
		},
		{
			name: "malformed series returns nil",
			result: map[string]any{
				"series": []any{
					[]any{"bad-point"},
					[]any{
						map[string]any{
							"columns": map[string]any{"host": "later-series-host"},
						},
					},
				},
			},
			field: "host",
			want:  nil,
		},
		{
			name: "bool value",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"ok": true},
						},
					},
				},
			},
			field: "ok",
			want:  true,
		},
		{
			name: "list value",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"values": []any{int64(1), "a"}},
						},
					},
				},
			},
			field: "values",
			want:  []any{int64(1), "a"},
		},
		{
			name: "map value",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"meta": map[string]any{"host": "h1"}},
						},
					},
				},
			},
			field: "meta",
			want:  map[string]any{"host": "h1"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := dqlSeriesFirstValue(tc.result, tc.field)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestDQLSeriesFirstValueMatchesSeriesGetFirst(t *testing.T) {
	cases := []struct {
		name   string
		result map[string]any
		field  string
	}{
		{
			name: "normal column",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"time": int64(1), "host": "h1"},
							"tags":    map[string]any{"host": "tag-h1"},
						},
						map[string]any{
							"columns": map[string]any{"time": int64(2), "host": "h2"},
							"tags":    map[string]any{"host": "tag-h2"},
						},
					},
					[]any{
						map[string]any{
							"columns": map[string]any{"time": int64(3), "host": "h3"},
							"tags":    map[string]any{"host": "tag-h3"},
						},
					},
				},
			},
			field: "host",
		},
		{
			name: "tag fallback",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"time": int64(1)},
							"tags":    map[string]any{"host": "tag-h1"},
						},
					},
					[]any{
						map[string]any{
							"columns": map[string]any{"host": "column-h2"},
							"tags":    map[string]any{"host": "tag-h2"},
						},
					},
				},
			},
			field: "host",
		},
		{
			name: "first non-empty series",
			result: map[string]any{
				"series": []any{
					[]any{},
					[]any{
						map[string]any{
							"columns": map[string]any{"host": "h1"},
						},
					},
				},
			},
			field: "host",
		},
		{
			name: "missing field on first point",
			result: map[string]any{
				"series": []any{
					[]any{
						map[string]any{
							"columns": map[string]any{"time": int64(1)},
						},
						map[string]any{
							"columns": map[string]any{"host": "later-host"},
						},
					},
				},
			},
			field: "host",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := dqlSeriesFirstValue(tc.result, tc.field)
			want := firstValueFromSeriesGet(tc.result, tc.field)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("expected %#v from dqlSeriesValues, got %#v", want, got)
			}
		})
	}
}

func TestFuncDQLSeriesFirstRuntimeEdgeCases(t *testing.T) {
	cases := []ProgCase{
		{
			Name: "dql_series_first_runtime_types",
			Script: `v = {"series": [[{"columns": {
    "bool_v": true,
    "int_v": 12,
    "float_v": 1.25,
    "str_v": "host-a",
    "list_v": [1, "a"],
    "map_v": {"host": "h1"}
}}]]}

printf("%v", {
    "bool_v": dql_series_first(v, "bool_v"),
    "int_v": dql_series_first(v, "int_v"),
    "float_v": dql_series_first(v, "float_v"),
    "str_v": dql_series_first(v, "str_v"),
    "list_v": dql_series_first(v, "list_v"),
    "map_v": dql_series_first(v, "map_v")
})
`,
			jsonout: true,
			Stdout:  `{"bool_v":true,"int_v":12,"float_v":1.25,"str_v":"host-a","list_v":[1,"a"],"map_v":{"host":"h1"}}`,
		},
		{
			Name: "dql_series_first_runtime_missing_first_point",
			Script: `v = {"series": [[
    {"columns": {"time": 1}},
    {"columns": {"host": "later-host"}}
]]}

printf("%v", {"host": dql_series_first(v, "host")})
`,
			jsonout: true,
			Stdout:  `{"host":null}`,
		},
		{
			Name: "dql_series_first_runtime_tag_fallback",
			Script: `v = {"series": [[{"columns": {"time": 1}, "tags": {"host": "tag-host"}}]]}

printf("%v", {"host": dql_series_first(v, "host")})
`,
			jsonout: true,
			Stdout:  `{"host":"tag-host"}`,
		},
		{
			Name: "dql_series_first_runtime_nil_value",
			Script: `v = {"series": [[{"columns": {"host": nil}}]]}

printf("%v", {"host": dql_series_first(v, "host")})
`,
			jsonout: true,
			Stdout:  `{"host":null}`,
		},
		{
			Name: "dql_series_first_runtime_malformed_first_point",
			Script: `v = {"series": [["bad-point"], [{"columns": {"host": "later-series-host"}}]]}

printf("%v", {"host": dql_series_first(v, "host")})
`,
			jsonout: true,
			Stdout:  `{"host":null}`,
		},
		{
			Name: "dql_series_first_runtime_empty_series",
			Script: `v = {"series": [[], [{"tags": {"host": "second-series-host"}}]]}

printf("%v", {"host": dql_series_first(v, "host")})
`,
			jsonout: true,
			Stdout:  `{"host":"second-series-host"}`,
		},
		{
			Name: "dql_series_first_runtime_missing_series",
			Script: `v = {}

printf("%v", {"host": dql_series_first(v, "host")})
`,
			jsonout: true,
			Stdout:  `{"host":null}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runCase(t, tc)
		})
	}
}

func firstValueFromSeriesGet(dqlResult map[string]any, name string) any {
	for _, vec := range dqlSeriesValues(dqlResult, name) {
		values := getList(vec)
		if len(values) > 0 {
			return values[0]
		}
	}
	return nil
}
