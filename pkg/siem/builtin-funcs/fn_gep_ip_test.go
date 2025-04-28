// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"testing"

	"github.com/GuanceCloud/pipeline-go/ptinput/ipdb"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

type mockGEO struct{}

func (m *mockGEO) Init(dataDir string, config map[string]string) {}
func (m *mockGEO) SearchIsp(ip string) string                    { return "" }

func (m *mockGEO) Geo(ip string) (*ipdb.IPdbRecord, error) {
	switch ip {
	case "114.114.114.114":
		return &ipdb.IPdbRecord{
			City:    "Ji'an",
			Country: "CN",
			Region:  "Jiangxi",
			Isp:     "chinanet",
		}, nil
	default:
		return nil, nil
	}
}

func TestFnGepIp(t *testing.T) {
	cases := []ProgCase{}
	cases = append(cases, cGeoIP.Progs...)
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runCase(t, tc, map[runtimev2.TaskP]any{
				PGeoIPDB: &mockGEO{},
			})
		})
	}
}
