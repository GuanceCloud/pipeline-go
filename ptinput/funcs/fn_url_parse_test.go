// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"
	"time"

	"github.com/GuanceCloud/cliutils/point"
	tu "github.com/GuanceCloud/cliutils/testutil"
	"github.com/GuanceCloud/pipeline-go/ptinput"
	"github.com/GuanceCloud/platypus/pkg/ast"
)

func TestURLParse(t *testing.T) {
	cases := []struct {
		name     string
		pl, in   string
		outKey   string
		expected interface{}
		absent   []string
		fail     bool
	}{
		{
			name: "port",
			pl: `
json(_, url)
m = url_parse(url)
add_key(scheme, m["scheme"])
`,
			in:       `{"url": "https://www.baidu.com"}`,
			outKey:   "scheme",
			expected: "https",
			fail:     false,
		},
		{
			name: "host",
			pl: `
json(_, url)
m = url_parse(url)
add_key(host, m["host"])
`,
			in:       `{"url": "http://127.0.0.1:9529"}`,
			outKey:   "host",
			expected: "127.0.0.1:9529",
			fail:     false,
		},
		{
			name: "port",
			pl: `
json(_, url)
m = url_parse(url)
add_key(port, m["port"])
`,
			in:       `{"url": "http://127.0.0.1:9529"}`,
			outKey:   "port",
			expected: "9529",
			fail:     false,
		},
		{
			name: "path",
			pl: `
json(_, url)
m = url_parse(url)
add_key(path, m["path"])
`,
			in:       `{"url": "http://127.0.0.1:9529/v1/metrics"}`,
			outKey:   "path",
			expected: "/v1/metrics",
			fail:     false,
		},
		{
			name: "arg1",
			pl: `
json(_, url)
m = url_parse(url)
add_key(a, m["params"]["arg1"])
`,
			in:       `{"url": "http://127.0.0.1:9529/v1/metrics?arg1=v1&arg2=v2"}`,
			outKey:   "a",
			expected: "v1",
			fail:     false,
		},
		{
			name: "arg2",
			pl: `
json(_, url)
m = url_parse(url)
add_key(a, m["params"]["arg2"])
`,
			in:       `{"url": "http://127.0.0.1:9529/v1/metrics?arg1=v1&arg2=v2&arg2=v3"}`,
			outKey:   "a",
			expected: "v2,v3",
			fail:     false,
		},
		{
			name: "with prefix",
			pl: `
json(_, url)
m = url_parse(url, "up_")
add_key(scheme, m["up_scheme"])
if m["scheme"] != nil {
    add_key(unexpected_unprefixed_key, true)
}
`,
			in:       `{"url": "https://www.baidu.com"}`,
			outKey:   "scheme",
			expected: "https",
			absent:   []string{"unexpected_unprefixed_key"},
			fail:     false,
		},
		{
			name: "with prefix params",
			pl: `
json(_, url)
m = url_parse(url, "up_")
add_key(a, m["up_params"]["arg1"])
`,
			in:       `{"url": "http://127.0.0.1:9529/v1/metrics?arg1=v1&arg2=v2"}`,
			outKey:   "a",
			expected: "v1",
			fail:     false,
		},
		{
			name: "with prefix variable",
			pl: `
json(_, url)
p = "up_"
m = url_parse(url, p)
add_key(scheme, m["up_scheme"])
`,
			in:       `{"url": "https://www.baidu.com"}`,
			outKey:   "scheme",
			expected: "https",
			fail:     false,
		},
		{
			name: "with named prefix",
			pl: `
json(_, url)
m = url_parse(url, prefix="up_")
add_key(scheme, m["up_scheme"])
`,
			in:       `{"url": "https://www.baidu.com"}`,
			outKey:   "scheme",
			expected: "https",
			fail:     false,
		},
		{
			name: "empty prefix",
			pl: `
json(_, url)
m = url_parse(url, "")
add_key(scheme, m["scheme"])
`,
			in:       `{"url": "https://www.baidu.com"}`,
			outKey:   "scheme",
			expected: "https",
			fail:     false,
		},
		{
			name: "relative path",
			pl: `
json(_, url)
m = url_parse(url)
add_key(p, m["path"])
`,
			in:       `{"url": "/var/log/datakit/log"}`,
			outKey:   "p",
			expected: "/var/log/datakit/log",
		},
		{
			name: "too many args",
			pl: `
json(_, url)
m = url_parse(url, "up_", 2)
`,
			in:   `{"url": "http://127.0.0.1:9529/v1/metrics?arg1=v1&arg2=v2"}`,
			fail: true,
		},
		{
			name: "invalid prefix type",
			pl: `
json(_, url)
m = url_parse(url, 2)
`,
			in:   `{"url": "http://127.0.0.1:9529/v1/metrics?arg1=v1&arg2=v2"}`,
			fail: true,
		},
		{
			name: "invalid prefix type variable",
			pl: `
json(_, url)
p = 123
m = url_parse(url, p)
`,
			in:   `{"url": "http://127.0.0.1:9529/v1/metrics?arg1=v1&arg2=v2"}`,
			fail: true,
		},
	}

	for idx, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner, err := NewTestingRunner(tc.pl)
			if err != nil {
				if !tc.fail {
					t.Fatalf("[%d] unexpected compile error: %s", idx, err)
				}
				return
			}

			pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())
			errR := runScript(runner, pt)
			if errR != nil {
				if !tc.fail {
					t.Fatalf("[%d] unexpected runtime error: %s", idx, errR)
				}
				return
			}
			if tc.fail {
				t.Fatal("expected an error, got nil")
			}

			if v, istag, err := pt.Get(tc.outKey); err != nil {
				t.Errorf("[%d]key %s, error: %s", idx, tc.outKey, err)
			} else {
				if istag != ast.String {
					t.Errorf("key %s should be a field", tc.outKey)
				} else {
					tu.Equals(t, tc.expected, v)
					t.Logf("[%d] PASS", idx)
				}
			}
			for _, k := range tc.absent {
				if _, _, err := pt.Get(k); err == nil {
					t.Errorf("key %q indicates an unprefixed result entry was generated", k)
				}
			}
		})
	}
}
