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
)

func TestUserAgent(t *testing.T) {
	cases := []struct {
		name     string
		pl, in   string
		expected map[string]interface{}
		absent   []string
		fail     bool
	}{
		{
			name: "normal",
			pl: `json(_, userAgent)
			user_agent(userAgent)`,
			in: `
{
   "userAgent" : "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36",
   "second"    : 2,
   "third"     : "abc",
   "forth"     : true
}
`,
			expected: map[string]interface{}{
				"isMobile":   false,
				"isBot":      false,
				"os":         "Windows 7",
				"browser":    "Chrome",
				"browserVer": "36.0.1985.125",
				"engine":     "AppleWebKit",
				"engineVer":  "537.36",
				"ua":         "Windows",
			},
			fail: false,
		},
		{
			name: "normal",
			pl: `json(_, userAgent)
			user_agent(userAgent)`,
			in: `
{
    "userAgent" : "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15"
}
`,
			expected: map[string]interface{}{
				"isMobile":   false,
				"isBot":      false,
				"os":         "Intel Mac OS X 10_15_7",
				"browser":    "Safari",
				"browserVer": "15.1",
				"engine":     "AppleWebKit",
				"engineVer":  "605.1.15",
				"ua":         "Macintosh",
			},
			fail: false,
		},

		{
			name: "normal",
			pl: `json(_, userAgent)
			user_agent(agent)`,
			in: `
{
    "userAgent" : "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15"
}
`,
			expected: map[string]interface{}{},
			fail:     false,
		},

		{
			name: "invalid arg type",
			pl: `json(_, userAgent)
			user_agent("userAgent")`,
			in: `
		{
		   "userAgent" : "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15"
		}
		`,
			expected: map[string]interface{}{},
			fail:     false,
		},

		{
			name: "with prefix",
			pl: `json(_, userAgent)
			add_key(os, "existing")
			user_agent(userAgent, "ua_")`,
			in: `
{
   "userAgent" : "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36",
   "second"    : 2,
   "third"     : "abc",
   "forth"     : true
}
`,
			expected: map[string]interface{}{
				"os":            "existing",
				"ua_isMobile":   false,
				"ua_isBot":      false,
				"ua_os":         "Windows 7",
				"ua_browser":    "Chrome",
				"ua_browserVer": "36.0.1985.125",
				"ua_engine":     "AppleWebKit",
				"ua_engineVer":  "537.36",
				"ua_ua":         "Windows",
			},
			absent: []string{"isMobile", "isBot", "browser", "browserVer", "engine", "engineVer", "ua"},
			fail:   false,
		},

		{
			name: "with prefix variable",
			pl: `json(_, userAgent)
			p = "ua_"
			user_agent(userAgent, p)`,
			in: `
{
   "userAgent" : "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36",
   "second"    : 2,
   "third"     : "abc",
   "forth"     : true
}
`,
			expected: map[string]interface{}{
				"ua_isMobile":   false,
				"ua_isBot":      false,
				"ua_os":         "Windows 7",
				"ua_browser":    "Chrome",
				"ua_browserVer": "36.0.1985.125",
				"ua_engine":     "AppleWebKit",
				"ua_engineVer":  "537.36",
				"ua_ua":         "Windows",
			},
			fail: false,
		},

		{
			name: "with named prefix",
			pl: `json(_, userAgent)
			user_agent(userAgent, prefix="ua_")`,
			in: `
{
   "userAgent" : "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36",
   "second"    : 2,
   "third"     : "abc",
   "forth"     : true
}
`,
			expected: map[string]interface{}{
				"ua_isMobile":   false,
				"ua_isBot":      false,
				"ua_os":         "Windows 7",
				"ua_browser":    "Chrome",
				"ua_browserVer": "36.0.1985.125",
				"ua_engine":     "AppleWebKit",
				"ua_engineVer":  "537.36",
				"ua_ua":         "Windows",
			},
			fail: false,
		},

		{
			name: "empty prefix",
			pl: `json(_, userAgent)
			user_agent(userAgent, "")`,
			in: `
{
   "userAgent" : "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36",
   "second"    : 2,
   "third"     : "abc",
   "forth"     : true
}
`,
			expected: map[string]interface{}{
				"isMobile":   false,
				"isBot":      false,
				"os":         "Windows 7",
				"browser":    "Chrome",
				"browserVer": "36.0.1985.125",
				"engine":     "AppleWebKit",
				"engineVer":  "537.36",
				"ua":         "Windows",
			},
			fail: false,
		},

		{
			name:     "missing key with prefix",
			pl:       `user_agent(agent, "ua_")`,
			in:       `{}`,
			expected: map[string]interface{}{},
			fail:     false,
		},

		{
			name: "invalid prefix type",
			pl: `json(_, userAgent)
			user_agent(userAgent, 123)`,
			in: `
{
   "userAgent" : "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36"
}
`,
			fail: true,
		},

		{
			name: "invalid prefix type variable",
			pl: `json(_, userAgent)
			p = 123
			user_agent(userAgent, p)`,
			in: `
{
   "userAgent" : "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36"
}
`,
			fail: true,
		},

		{
			name: "too many args",
			pl: `json(_, userAgent)
			user_agent(userAgent, someArg, anotherArg)`,
			in: `
		{
		   "userAgent" : "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15"
		}
		`,
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

			pt := ptinput.NewPlPt(
				point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())
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

			fieldsToCompare := make(map[string]interface{})
			for k := range tc.expected {
				fieldsToCompare[k], _, _ = pt.Get(k)
			}
			tu.Equals(t, tc.expected, fieldsToCompare)
			for _, k := range tc.absent {
				if _, _, err := pt.Get(k); err == nil {
					t.Errorf("key %q should not be generated without the prefix", k)
				}
			}
			t.Logf("[%d] PASS", idx)
		})
	}
}
