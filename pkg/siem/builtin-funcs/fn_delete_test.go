// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"
)

func TestDelete(t *testing.T) {
	cases := []ProgCase{}
	cases = append(cases, cDelete.Progs...)

	for _, tc := range cases {
		runCase(t, tc)
	}
}
