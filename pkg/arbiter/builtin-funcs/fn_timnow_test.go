// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/dql"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

func TestTimestamp(t *testing.T) {
	cases := []ProgCase{
		{
			Name: "time_now_returns_query_start",
			Script: `printf("%v,%v,%v,%v", time_now("s"), time_now("ms"), time_now("us"), time_now("ns"))
`,
			Stdout: "1672531500,1672531500123,1672531500123000,1672531500123000000",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			runCase(t, c, map[runtimev2.TaskP]any{
				PDQLCli: dql.NewDQLOpenAPI(
					"",
					dql.OpenAPIPath,
					"abc",
					[]int64{1672531500123, 1672532100123},
				),
			})
		})
	}
}
