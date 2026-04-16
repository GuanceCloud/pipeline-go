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

func TestAddkey(t *testing.T) {
	cases := []struct {
		name, pl, in string
		expect       interface{}
		checkExport  bool
		fail         bool
	}{
		{
			name: "value type: string",
			in:   `1.2.3.4 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
add_key(add_new_key, "shanghai")
`,
			expect: "shanghai",
			fail:   false,
		},
		{
			name: "value type: number(int64)",
			in:   `1.2.3.4 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
		grok(_, "%{IPORHOST:client_ip} %{NOTSPACE} %{NOTSPACE} \\[%{HTTPDATE:time}\\] \"%{DATA} %{GREEDYDATA} HTTP/%{NUMBER}\" %{INT:status_code} %{INT:bytes}")
		add_key(add_new_key, -1)
		`,
			expect: int64(-1),
			fail:   false,
		},
		{
			name: "value type: number(float64)",
			in:   `1.2.3.4 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
		grok(_, "%{IPORHOST:client_ip} %{NOTSPACE} %{NOTSPACE} \\[%{HTTPDATE:time}\\] \"%{DATA} %{GREEDYDATA} HTTP/%{NUMBER}\" %{INT:status_code} %{INT:bytes}")
		add_key(add_new_key, 1.)
		`,
			expect: float64(1),
			fail:   false,
		},
		{
			name: "value type: number(float64)",
			in:   `1.2.3.4 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
		grok(_, "%{IPORHOST:client_ip} %{NOTSPACE} %{NOTSPACE} \\[%{HTTPDATE:time}\\] \"%{DATA} %{GREEDYDATA} HTTP/%{NUMBER}\" %{INT:status_code} %{INT:bytes}")
		add_key(add_new_key, .1)
		`,
			expect: float64(.1),
			fail:   true, // .1 not supported
		},
		{
			name: "value type: bool",
			in:   `1.2.3.4 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
		grok(_, "%{IPORHOST:client_ip} %{NOTSPACE} %{NOTSPACE} \\[%{HTTPDATE:time}\\] \"%{DATA} %{GREEDYDATA} HTTP/%{NUMBER}\" %{INT:status_code} %{INT:bytes}")
		add_key(add_new_key, true)
		`,
			expect: true,
			fail:   false,
		},
		{
			name: "value type: bool",
			in:   `1.2.3.4 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
		grok(_, "%{IPORHOST:client_ip} %{NOTSPACE} %{NOTSPACE} \\[%{HTTPDATE:time}\\] \"%{DATA} %{GREEDYDATA} HTTP/%{NUMBER}\" %{INT:status_code} %{INT:bytes}")
		add_key(add_new_key, tRue)
		`,
			expect: true,
			fail:   false,
		},
		{
			name: "value type: bool",
			in:   `1.2.3.4 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
		grok(_, "%{IPORHOST:client_ip} %{NOTSPACE} %{NOTSPACE} \\[%{HTTPDATE:time}\\] \"%{DATA} %{GREEDYDATA} HTTP/%{NUMBER}\" %{INT:status_code} %{INT:bytes}")
		add_key(add_new_key, false)
		`,
			expect: false,
			fail:   false,
		},
		{
			name: "value type: nil",
			in:   `1.2.3.4 - - [29/Nov/2021:07:30:50 +0000] "POST /?signature=b8d8ea&timestamp=1638171049 HTTP/1.1" 200 413 "-" "Mozilla/4.0"`,
			pl: `
		grok(_, "%{IPORHOST:client_ip} %{NOTSPACE} %{NOTSPACE} \\[%{HTTPDATE:time}\\] \"%{DATA} %{GREEDYDATA} HTTP/%{NUMBER}\" %{INT:status_code} %{INT:bytes}")
		add_key(add_new_key, nil)
		`,
			expect: nil,
			fail:   false,
		},
		{
			name: "value type: list compatibility",
			in:   `test`,
			pl: `
		add_key(add_new_key, [1, 2])
		`,
			expect:      "[1,2]",
			checkExport: true,
			fail:        false,
		},
		{
			name: "value type: map compatibility",
			in:   `test`,
			pl: `
		add_key(add_new_key, {"a": 1, "b": "x"})
		`,
			expect:      `{"a":1,"b":"x"}`,
			checkExport: true,
			fail:        false,
		},
	}

	for idx, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner, err := NewTestingRunner(tc.pl)
			if err != nil {
				if tc.fail {
					t.Logf("[%d]expect error: %s", idx, err)
				} else {
					t.Errorf("[%d] failed: %s", idx, err)
				}
				return
			}
			pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())

			errR := runScript(runner, pt)

			if errR == nil {
				v, _, e := pt.Get("add_new_key")
				assert.NoError(t, e)
				assert.Equal(t, tc.expect, v)
				if tc.checkExport {
					assert.Equal(t, tc.expect, pt.Fields()["add_new_key"])
					assert.Equal(t, tc.expect, pt.Point().KVs().InfluxFields()["add_new_key"])
				}
				t.Logf("[%d] PASS", idx)
			} else {
				t.Error(errR)
			}
		})
	}
}

func TestAddkeyMessageMapCompatibility(t *testing.T) {
	pl := `
result_msg = {}
all_keys = pt_kvs_keys()
for k in all_keys {
	if k != "message" {
		result_msg[k] = pt_kvs_get(k)
	}
}
add_key("message", result_msg)
`

	runner, err := NewTestingRunner(pl)
	if err != nil {
		t.Fatal(err)
	}

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": "test",
		"a":       "x",
		"b":       int64(1),
	}, time.Now())

	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	v, _, err := pt.Get("message")
	assert.NoError(t, err)
	assert.Equal(t, `{"a":"x","b":1,"status":"info"}`, v)

	assert.Equal(t, `{"a":"x","b":1,"status":"info"}`, pt.Fields()["message"])
	assert.Equal(t, `{"a":"x","b":1,"status":"info"}`, pt.Point().KVs().InfluxFields()["message"])
}
