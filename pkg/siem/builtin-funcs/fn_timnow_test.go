// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/GuanceCloud/platypus/pkg/engine"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

func TestTimestamp(t *testing.T) {
	cases := []ProgCase{
		{
			Name: "timenow",
			Script: `printf("%v", time_now("s"))
`,
			Stdout: fmt.Sprintf("%d", time.Now().Unix()),
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			s, err := engine.ParseV2(c.Name, c.Script, Funcs)
			if err != nil {
				t.Error(err)
				return
			}

			var privateMap map[runtimev2.TaskP]any

			stdout := bytes.NewBuffer([]byte{})
			if err := s.Run(nil, runtimev2.WithPrivate(privateMap)); err != nil {
				t.Error(err.Error())
			}
			aV, _ := strconv.ParseInt(stdout.String(), 10, 64)
			eV, _ := strconv.ParseInt(c.Stdout, 10, 64)
			if (aV - eV) > 2 {
				t.Error("not equal")
			}
		})
	}
}
