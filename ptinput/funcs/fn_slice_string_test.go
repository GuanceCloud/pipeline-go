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

func TestSliceString(t *testing.T) {

	cases := []struct {
		name, pl, in string
		keyName      string
		expect       any
		fail         bool
	}{
		{
			name: "normal1",
			pl: `
			substring = slice_string("█汉字15384073392",0,5)
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "█汉字15",
			fail:    false,
		},
		{
			name: "normal2",
			pl: `
			substring = slice_string("15384073392",5,10)
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "07339",
			fail:    false,
		},
		{
			name: "normal3",
			pl: `
			substring = slice_string("abcdefghijklmnop",0,10)
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "abcdefghij",
			fail:    false,
		},
		{
			name: "out of range1",
			pl: `
			substring = slice_string("abcdefghijklmnop",-1,10)
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "",
			fail:    false,
		},
		{
			name: "out of range2",
			pl: `
			substring = slice_string("abcdefghijklmnop",0,100)
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "abcdefghijklmnop",
			fail:    false,
		},
		{
			name: "not integer1",
			pl: `
			substring = slice_string("abcdefghijklmnop","a","b")
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "",
			fail:    true,
		},
		{
			name: "not integer2",
			pl: `
			substring = slice_string("abcdefghijklmnop","abc","def")
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "",
			fail:    true,
		},
		{
			name: "not string",
			pl: `
			substring = slice_string(12345,0,3)
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "",
			fail:    true,
		},
		{
			name: "not correct args",
			pl: `
			substring = slice_string("abcdefghijklmnop",0)
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "",
			fail:    true,
		},
		{
			name: "not correct args",
			pl: `
			substring = slice_string("abcdefghijklmnop",0,1,2)
			pt_kvs_set("result", substring)
			`,
			keyName: "result",
			expect:  "",
			fail:    true,
		},
		{
			name: "panic",
			pl: `
			val = "123你好123123123123123123123123123"
			## len 32, cap 32
			#
			add_key("result", slice_string(val, 0, len(val)))
			`,
			keyName: "result",
			expect:  "123你好123123123123123123123123123",
		},
	}

	for idx, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			script, err := NewTestingRunner(tc.pl)
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
			errR := script.Run(pt, nil)
			if errR != nil {
				t.Fatal(errR.Error())
			}

			v, _, _ := pt.Get(tc.keyName)
			assert.Equal(t, tc.expect, v)
			t.Logf("[%d] PASS", idx)

		})
	}
}
