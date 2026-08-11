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
	"github.com/GuanceCloud/pipeline-go/ptinput/ipdb"

	"github.com/stretchr/testify/assert"
)

type mockGEO struct{}

func (m *mockGEO) Init(dataDir string, config map[string]string) {}
func (m *mockGEO) SearchIsp(ip string) string                    { return geoDefaultVal }
func (m *mockGEO) GeoWithChecker(ip string, check ipdb.CheckData) (*ipdb.IPdbRecord, error) {
	return m.Geo(ip)
}

func (m *mockGEO) Geo(ip string) (*ipdb.IPdbRecord, error) {
	return &ipdb.IPdbRecord{
		City: func() string {
			if ip == "unknown-city" {
				return geoDefaultVal
			} else {
				return "Shanghai"
			}
		}(),
		Region: func() string {
			if ip == "unknown-region" {
				return geoDefaultVal
			} else {
				return "Shanghai"
			}
		}(),
		Country: func() string {
			if ip == "unknown-country-short" {
				return geoDefaultVal
			} else {
				return "CN"
			}
		}(),
		Isp: m.SearchIsp(ip),
	}, nil
}

func TestGeoIpFunc(t *testing.T) {
	cases := []struct {
		in     string
		script string

		expected map[string]string
		absent   []string

		fail bool
	}{
		{
			in: `{"ip":"1.2.3.4-something", "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, ip)
				geoip(ip)`,
			expected: map[string]string{
				"city":     "Shanghai",
				"country":  "CN",
				"province": "Shanghai",
				"isp":      geoDefaultVal,
			},
		},

		{
			in: `{"ip":"unknown-city", "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, ip)
				geoip(ip)`,
			expected: map[string]string{
				"city":     geoDefaultVal,
				"country":  "CN",
				"province": "Shanghai",
				"isp":      geoDefaultVal,
			},
		},

		{
			in: `{"aa": {"ip":"116.228.89.xxx"}, "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, aa.ip)
				geoip(aa.ip)`,
			expected: map[string]string{
				"city":     "Shanghai",
				"country":  "CN",
				"province": "Shanghai",
				"isp":      geoDefaultVal,
			},
		},

		{
			in: `{"aa": {"ip":"unknown-region"}, "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, aa.ip)
				geoip(aa.ip)`,
			expected: map[string]string{
				"city":     "Shanghai",
				"country":  "CN",
				"province": geoDefaultVal,
				"isp":      geoDefaultVal,
			},
		},

		{
			in: `{"aa": {"ip":"unknown-country-short"}, "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, aa.ip)
				geoip(aa.ip)`,
			expected: map[string]string{
				"city":     "Shanghai",
				"country":  geoDefaultVal,
				"province": "Shanghai",
				"isp":      geoDefaultVal,
			},
		},

		{
			in: `{"ip":"1.2.3.4-something", "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, ip)
				add_key(city, "existing")
				geoip(ip, "geo_")`,
			expected: map[string]string{
				"city":         "existing",
				"geo_city":     "Shanghai",
				"geo_country":  "CN",
				"geo_province": "Shanghai",
				"geo_isp":      geoDefaultVal,
			},
			absent: []string{"country", "province", "isp"},
		},
		{
			in: `{"ip":"1.2.3.4-something", "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, ip)
				p = "geo_"
				geoip(ip, p)`,
			expected: map[string]string{
				"geo_city":     "Shanghai",
				"geo_country":  "CN",
				"geo_province": "Shanghai",
				"geo_isp":      geoDefaultVal,
			},
		},
		{
			in: `{"ip":"1.2.3.4-something", "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, ip)
				geoip(ip, prefix="geo_")`,
			expected: map[string]string{
				"geo_city":     "Shanghai",
				"geo_country":  "CN",
				"geo_province": "Shanghai",
				"geo_isp":      geoDefaultVal,
			},
		},
		{
			in: `{"ip":"1.2.3.4-something", "second":2,"third":"abc","forth":true}`,
			script: `
				json(_, ip)
				geoip(ip, "")`,
			expected: map[string]string{
				"city":     "Shanghai",
				"country":  "CN",
				"province": "Shanghai",
				"isp":      geoDefaultVal,
			},
		},
		{
			in: `{"second":2,"third":"abc","forth":true}`,
			script: `
				geoip(ip, "geo_")`,
			expected: map[string]string{},
		},
		{
			in: `{"ip":"1.2.3.4-something"}`,
			script: `
				json(_, ip)
				geoip(ip, 2)`,
			fail: true,
		},
		{
			in: `{"ip":"1.2.3.4-something"}`,
			script: `
				json(_, ip)
				p = 123
				geoip(ip, p)`,
			fail: true,
		},
		{
			in: `{"ip":"1.2.3.4-something"}`,
			script: `
				json(_, ip)
				geoip(ip, "a", "b")`,
			fail: true,
		},
	}

	for idx, tc := range cases {
		t.Logf("case %d...", idx)

		runner, err := NewTestingRunner(tc.script)
		if err != nil {
			if !tc.fail {
				t.Fatalf("[%d] unexpected compile error: %s", idx, err)
			}
			continue
		}

		pt := ptinput.NewPlPt(
			point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())
		pt.SetIPDB(&mockGEO{})
		errR := runScript(runner, pt)

		if errR != nil {
			if !tc.fail {
				t.Fatalf("[%d] unexpected runtime error: %s", idx, errR)
			}
			continue
		}
		if tc.fail {
			t.Errorf("[%d] expected an error, got nil", idx)
			continue
		}

		for k, v := range tc.expected {
			r, _, e := pt.Get(k)
			assert.NoError(t, e)
			assert.Equal(t, v, r, "`%s` != `%s`, key: `%s`", r, v, k)
		}
		for _, k := range tc.absent {
			_, _, err := pt.Get(k)
			assert.Error(t, err, "key %q should not be generated without the prefix", k)
		}
	}
}
