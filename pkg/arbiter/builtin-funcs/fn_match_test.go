// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"
)

func TestMatch(t *testing.T) {
	cases := []ProgCase{
		{
			Name: "match",
			Script: `
v, ok =match("hello", "hello", 1)
printf("%v,%v", v[0][0], v[0])
`,

			Stdout: "hello,[\"hello\"]",
		},
	}
	cases = append(cases, cMatch.Progs...)
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runCase(t, tc)
		})
	}
}
