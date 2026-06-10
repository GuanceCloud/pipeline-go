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
)

func BenchmarkJSONAllCompare(b *testing.B) {
	const in = `{"service":"api","status":"ok","host":"web-01","region":"us-west","trace_id":"abc123","span_id":"def456","code":200,"duration":12.34,"success":true,"nested":{"ignored":"x"},"items":[1,2,3]}`

	benchmarks := []struct {
		name   string
		script string
	}{
		{
			name: "json_all",
			script: `json_all(_, include_keys=[
	"service",
	"status",
	"host",
	"region",
	"trace_id",
	"span_id",
	"code",
	"duration",
	"success",
])`,
		},
		{
			name: "json_repeated",
			script: `json(_, service)
json(_, status)
json(_, host)
json(_, region)
json(_, trace_id)
json(_, span_id)
json(_, code)
json(_, duration)
json(_, success)`,
		},
		{
			name: "load_json_pt_kvs_set",
			script: `data = load_json(_)
pt_kvs_set("service", data["service"])
pt_kvs_set("status", data["status"])
pt_kvs_set("host", data["host"])
pt_kvs_set("region", data["region"])
pt_kvs_set("trace_id", data["trace_id"])
pt_kvs_set("span_id", data["span_id"])
pt_kvs_set("code", data["code"])
pt_kvs_set("duration", data["duration"])
pt_kvs_set("success", data["success"])`,
		},
		{
			name: "load_json_pt_kvs_set_map",
			script: `data = load_json(_)
pt_kvs_set_map(data, include_keys=[
	"service",
	"status",
	"host",
	"region",
	"trace_id",
	"span_id",
	"code",
	"duration",
	"success",
])`,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			runner, err := NewTestingRunner(bm.script)
			if err != nil {
				b.Fatal(err)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				pt := ptinput.NewPlPt(
					point.Logging, "test", nil, map[string]any{"message": in}, time.Now())
				if err := runScript(runner, pt); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
