// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"
)

func TestExit(t *testing.T) {
	cases := []ProgCase{}
	cases = append(cases, cExit.Progs...)
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runCase(t, tc)
		})
	}
}
