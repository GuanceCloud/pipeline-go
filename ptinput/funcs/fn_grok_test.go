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

func TestGrok(t *testing.T) {
	cases := []struct {
		name, pl, in string
		expected     interface{}
		fail         bool
		outkey       string
	}{
		{
			name: "normal_return_t",
			pl: `
add_pattern("_second", "(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)")
add_pattern("_minute", "(?:[0-5][0-9])")
add_pattern("_hour", "(?:2[0123]|[01]?[0-9])")
add_pattern("time", "([^0-9]?)%{_hour:hour}:%{_minute:minute}(?::%{_second:second})([^0-9]?)")
add_key(grok_match_ok, grok(_, "%{time}"))`,
			in:       "12:13:14.123",
			expected: true,
			outkey:   "grok_match_ok",
		},
		{
			name: "normal_return_f",
			pl: `
add_pattern("_second", "(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)")
add_pattern("_minute", "(?:[0-5][0-9])")
add_pattern("_hour", "(?:2[0123]|[01]?[0-9])")
add_pattern("time", "([^0-9]?)%{_hour:hour}:%{_minute:minute}(?::%{_second:second})([^0-9]?)")
add_key(grok_match_ok, grok(_, "%{time}"))`,
			in:       "12 :13:14.123",
			expected: false,
			outkey:   "grok_match_ok",
		},
		{
			name: "normal_return_sample_t",
			pl: `
add_key(grok_match_ok, grok(_, "12 :13:14.123"))`,
			in:       "12 :13:14.123",
			expected: true,
			outkey:   "grok_match_ok",
		},
		{
			name: "normal_return_sample_f",
			pl: `
add_key(grok_match_ok, grok(_, "12 :13:14.123"))`,
			in:       "12:13:14.123",
			expected: false,
			outkey:   "grok_match_ok",
		},
		{
			name: "normal",
			pl: `
add_pattern("_second", "(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)")
add_pattern("_minute", "(?:[0-5][0-9])")
add_pattern("_hour", "(?:2[0123]|[01]?[0-9])")
add_pattern("time", "([^0-9]?)%{_hour:hour}:%{_minute:minute}(?::%{_second:second})([^0-9]?)")
grok(_, "%{time}")`,
			in:       "12:13:14.123",
			expected: "14.123",
			outkey:   "second",
		},
		{
			name: "normal",
			pl: `
add_pattern("_second", "(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)")
add_pattern("_minute", "(?:[0-5][0-9])")
add_pattern("_hour", "(?:2[0123]|[01]?[0-9])")
add_pattern("time", "([^0-9]?)%{_hour:hour}:%{_minute:minute}(?::%{_second:second})([^0-9]?)")
grok(_, "%{time}")`,
			in:       "12:13:14",
			expected: "13",
			outkey:   "minute",
		},
		{
			name: "normal",
			pl: `
add_pattern("_second", "(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)")
add_pattern("_minute", "(?:[0-5][0-9])")
add_pattern("_hour", "(?:2[0123]|[01]?[0-9])")
add_pattern("time", "([^0-9]?)%{_hour:hour}:%{_minute:minute}(?::%{_second:second})([^0-9]?)")
grok(_, "%{time}")`,
			in:       "12:13:14",
			expected: "12",
			outkey:   "hour",
		},
		{
			name: "normal",
			pl: `
add_pattern("_second", "(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)")
add_pattern("_minute", "(?:[0-5][0-9])")
add_pattern("_hour", "(?:2[0123]|[01]?[0-9])")
add_pattern("time", "([^0-9]?)%{_hour:hour}:%{_minute:minute}(?::%{_second:second})([^0-9]?)")
grok(_, "%{time}")`,
			in:       "12:13:14",
			expected: "14",
			outkey:   "second",
		},
		{
			name: "normal",
			pl: `
add_pattern("time", "%{NUMBER:time:float}")
grok(_, '''%{time}
%{WORD:word:string}
	%{WORD:code:int}
%{WORD:w1}''')`,
			in: `1.1
s
	123cvf
aa222`,
			expected: int64(0),
			outkey:   "code",
		},
		{
			name: "normal",
			pl: `
add_pattern("time", "%{NUMBER:time:float}")
grok(_, '''%{time}
%{WORD:word:string}
	%{WORD:code:int}
%{WORD:w1}''')`,
			in: `1.1
s
	123
aa222`,
			expected: int64(123),
			outkey:   "code",
		},
		{
			name: "normal",
			pl: `
add_pattern("time", "%{NUMBER:time:float}")
grok(_, '''%{time}
%{WORD:word:str}
	%{WORD:code:int}
%{WORD:w1}''')`,
			in: `1.1
s
	123
aa222`,
			expected: int64(123),
			outkey:   "code",
		},
		{
			name: "normal",
			pl: `
add_pattern("_second", "(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)")
add_pattern("_minute", "(?:[0-5][0-9])")
add_pattern("_hour", "(?:2[0123]|[01]?[0-9])")
add_pattern("time", "([^0-9]?)%{_hour:hour:string}:%{_minute:minute:int}(?::%{_second:second:float})([^0-9]?)")
grok(_, "%{WORD:date} %{time}")`,
			in:       "2021/1/11 2:13:14.123",
			expected: float64(14.123),
			outkey:   "second",
		},
		{
			name: "trim_space",
			in:   " not_space ",
			pl: `add_pattern("d", "[\\s\\S]*")
			grok(_, "%{d:item}")`,
			expected: "not_space",
			outkey:   "item",
		},
		{
			name: "trim_space, enable",
			in:   " not_space ",
			pl: `add_pattern("d", "[\\s\\S]*")
			grok(_, "%{d:item}", true)`,
			expected: "not_space",
			outkey:   "item",
		},
		{
			name: "trim_space, disable",
			in:   " not_space ",
			pl: `add_pattern("d", "[\\s\\S]*")
			grok(_, "%{d:item}", false)`,
			expected: " not_space ",
			outkey:   "item",
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

			pt := ptinput.NewPlPt(
				point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())
			errR := runScript(runner, pt)

			if errR != nil {
				t.Fatal(errR.Error())
			}

			if tc.fail {
				t.Logf("[%d]expect error: %s", idx, errR.Error())
			}
			v, _, _ := pt.Get(tc.outkey)
			tu.Equals(t, tc.expected, v)
			t.Logf("[%d] PASS", idx)
		})
	}
}

func TestGrokRunObserver(t *testing.T) {
	events := make(chan GrokRunInfo, 1)
	SetGrokRunObserver(func(info GrokRunInfo) {
		events <- info
	})
	defer SetGrokRunObserver(nil)

	runner, err := NewTestingRunner(`
grok(_, "%{NOTSPACE:client_ip} %{NOTSPACE:http_ident} %{NOTSPACE:http_auth} \\[%{HTTPDATE:time}\\] \"%{DATA:http_method} %{GREEDYDATA:http_url} HTTP/%{NUMBER:http_version}\" %{INT:status_code} %{INT:bytes}")
`)
	if err != nil {
		t.Fatal(err)
	}

	pt := ptinput.NewPlPt(
		point.Logging, "test", nil,
		map[string]any{"message": `127.0.0.1 - - [21/Jul/2021:14:14:38 +0800] "GET /?1 HTTP/1.1" 200 2178`},
		time.Now())
	if errR := runScript(runner, pt); errR != nil {
		t.Fatal(errR)
	}

	select {
	case info := <-events:
		if info.ScriptName != "default.p" {
			t.Fatalf("script name = %q", info.ScriptName)
		}
		if info.Line <= 0 || info.Column <= 0 {
			t.Fatalf("expected positive source position, got line=%d column=%d", info.Line, info.Column)
		}
		if info.PatternHash == 0 {
			t.Fatal("expected pattern hash")
		}
		if info.Path == "" {
			t.Fatal("expected path")
		}
		if info.WorkUnits <= 0 {
			t.Fatalf("expected positive work units, got %d", info.WorkUnits)
		}
		if info.Cost <= 0 {
			t.Fatalf("expected positive cost, got %s", info.Cost)
		}
	default:
		t.Fatal("expected grok observer event")
	}
}

func TestGrokRunObserverPanicDoesNotBreakGrok(t *testing.T) {
	SetGrokRunObserver(func(info GrokRunInfo) {
		panic("observer panic")
	})
	defer SetGrokRunObserver(nil)

	runner, err := NewTestingRunner(`
grok(_, "%{NOTSPACE:client_ip} %{NOTSPACE:http_ident} %{NOTSPACE:http_auth} \\[%{HTTPDATE:time}\\] \"%{DATA:http_method} %{GREEDYDATA:http_url} HTTP/%{NUMBER:http_version}\" %{INT:status_code} %{INT:bytes}")
`)
	if err != nil {
		t.Fatal(err)
	}

	pt := ptinput.NewPlPt(
		point.Logging, "test", nil,
		map[string]any{"message": `127.0.0.1 - - [21/Jul/2021:14:14:38 +0800] "GET /?1 HTTP/1.1" 200 2178`},
		time.Now())
	if errR := runScript(runner, pt); errR != nil {
		t.Fatal(errR)
	}

	v, _, err := pt.Get("client_ip")
	if err != nil || v != "127.0.0.1" {
		t.Fatalf("expected grok output after observer panic, got value=%#v err=%v", v, err)
	}
}

func TestGrokFastPathCompatibility(t *testing.T) {
	cases := []struct {
		name, pl, in string
		outkey       string
		expected     interface{}
	}{
		{
			name: "inverse_space_class_captures_whole_line_without_trim",
			in:   " not_space ",
			pl: `add_pattern("ANY", "[\\s\\S]*")
grok(_, "%{ANY:item}", false)`,
			outkey:   "item",
			expected: " not_space ",
		},
		{
			name: "inverse_space_class_captures_multiline",
			in:   "line1\nline2",
			pl: `add_pattern("ANY", "[\\s\\S]*")
grok(_, "%{ANY:item}", false)`,
			outkey:   "item",
			expected: "line1\nline2",
		},
		{
			name: "inverse_digit_class_preserves_punctuation_and_space",
			in:   "abc- ",
			pl: `add_pattern("NONDIGITS", "[\\D]*")
grok(_, "%{NONDIGITS:item}", false)`,
			outkey:   "item",
			expected: "abc- ",
		},
		{
			name: "inverse_word_class_before_literal",
			in:   " -\tword",
			pl: `add_pattern("NONWORDS", "[\\W]+")
grok(_, "^%{NONWORDS:item}word$", false)`,
			outkey:   "item",
			expected: " -\t",
		},
		{
			name: "dot_all_style_class_between_literals",
			in:   "prefix: a:b :suffix",
			pl: `add_pattern("ANY", "[\\s\\S]*")
grok(_, "^prefix:%{ANY:item}:suffix$")`,
			outkey:   "item",
			expected: "a:b",
		},
		{
			name:     "greedydata_keeps_text_before_following_literal",
			in:       "prefix=alpha beta suffix=ok",
			pl:       `grok(_, "^prefix=%{GREEDYDATA:item} suffix=%{WORD:tail}$")`,
			outkey:   "item",
			expected: "alpha beta",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner, err := NewTestingRunner(tc.pl)
			if err != nil {
				t.Fatal(err)
			}

			pt := ptinput.NewPlPt(
				point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())
			if errR := runScript(runner, pt); errR != nil {
				t.Fatal(errR.Error())
			}

			v, _, err := pt.Get(tc.outkey)
			if err != nil {
				t.Fatal(err)
			}
			tu.Equals(t, tc.expected, v)
		})
	}
}
