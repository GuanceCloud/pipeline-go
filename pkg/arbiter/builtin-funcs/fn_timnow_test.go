// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	originTimeNow := timeNow
	timeNow = func() time.Time {
		return time.Unix(1672531500, 123456789)
	}
	defer func() {
		timeNow = originTimeNow
	}()

	cases := []ProgCase{
		{
			Name: "time_now_returns_current_time",
			Script: `printf("%v,%v,%v,%v", time_now("s"), time_now("ms"), time_now("us"), time_now("ns"))
`,
			Stdout: "1672531500,1672531500123,1672531500123456,1672531500123456789",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			runCase(t, c)
		})
	}
}
