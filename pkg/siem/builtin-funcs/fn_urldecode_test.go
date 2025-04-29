// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"
)

func TestURLDecode(t *testing.T) {
	cases := append([]ProgCase{}, cURLDecode.Progs...)
	for _, tc := range cases {
		runCase(t, tc)
	}
}
