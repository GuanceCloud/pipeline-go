package funcs

import (
	"github.com/GuanceCloud/pipeline-go/pkg/siem/trigger"
)

type ProgCase struct {
	Name          string
	Script        string
	Stdout        string
	jsonout       bool
	TriggerResult []trigger.Data
}

type FuncExample struct {
	FnName string
	Progs  []ProgCase
}

var FnExps = []*FuncExample{
	cCast,
	cCIDR,
	cDelete,
	cDQL,
	cExit,
	cTrigger,
}

var cCast = &FuncExample{
	FnName: FnCastDesc.Name,
	Progs: []ProgCase{
		{
			Script: `v1 = "1.1"
v2 = "1"
v2_1 = "-1"
v3 = "true"

printf("%v; %v; %v; %v; %v; %v; %v; %v\n",
	cast(v1, "float") + 1,
	cast(v2, "int") + 1,
	cast(v2_1, "int"),
	cast(v3, "bool") + 1,

	cast(cast(v3, "bool") - 1, "bool"),
	cast(1.1, "str"),
	cast(1.1, "int"),
	cast(1.1, "bool")
)
`,
			Stdout: "2.1; 2; -1; 2; false; 1.1; 1; true\n",
		},
	},
}

var cCIDR = &FuncExample{
	FnName: FnCIDRDesc.Name,
	Progs: []ProgCase{
		{
			Name: "ipv4_contains",
			Script: `ip = "192.0.2.233"
if cidr(ip, "192.0.2.1/24") {
	printf("%s", ip)
}`,
			Stdout: "192.0.2.233",
		},
		{
			Name: "ipv4_not_contains",
			Script: `ip = "192.0.2.233"
if cidr(mask="192.0.1.1/24", ip=ip) {
	printf("%s", ip)
}`,
			Stdout: "",
		},
	},
}

var cDelete = &FuncExample{
	FnName: FnDeleteDesc.Name,
	Progs: []ProgCase{
		{
			Name: "delete_map",
			Script: `v = {
    "k1": 123,
    "k2": {
        "a": 1,
        "b": 2,
    },
    "k3": [{
        "c": 1.1, 
        "d":"2.1",
    }]
}
delete(v["k2"], "a")
delete(v["k3"][0], "d")
printf("result group 1: %v; %v\n", v["k2"], v["k3"])

v1 = {"a":1}
v2 = {"b":1}
delete(key="a", m=v1)
delete(m=v2, key="b")
printf("result group 2: %v; %v\n", v1, v2)
`,
			Stdout: "result group 1: {\"b\":2}; [{\"c\":1.1}]\nresult group 2: {}; {}\n",
		},
	},
}

var cExit = &FuncExample{
	FnName: FnExitDesc.Name,
	Progs: []ProgCase{
		{
			Name: "cast int",
			Script: `printf("1\n")
printf("2\n")
exit()
printf("3\n")
	`,
			Stdout: "1\n2\n",
		},
	},
}

var cDQL = &FuncExample{
	FnName: FnDQLDesc.Name,
	Progs: []ProgCase{
		{
			Name: "dql",
			Script: `v, ok = dql("M::cpu limit 3 slimit 3")
if ok {
	v, ok = dump_json(v, "    ")
	if ok {
		printf("%v", v)
	}
}
`,
			jsonout: true,
			Stdout: `{
    "series": [
        [
            {
                "columns": {
                    "time": 1744866108991,
                    "total": 7.18078381,
                    "user": 4.77876106
                },
                "tags": {
                    "cpu": "cpu-total",
                    "guance_site": "testing",
                    "host": "172.16.241.111",
                    "host_ip": "172.16.241.111",
                    "name": "cpu",
                    "project": "cloudcare-testing"
                }
            },
            {
                "columns": {
                    "time": 1744866103991,
                    "total": 10.37376049,
                    "user": 7.17009916
                },
                "tags": {
                    "cpu": "cpu-total",
                    "guance_site": "testing",
                    "host": "172.16.241.111",
                    "host_ip": "172.16.241.111",
                    "name": "cpu",
                    "project": "cloudcare-testing"
                }
            }
        ],
        [
            {
                "columns": {
                    "time": 1744866107975,
                    "total": 21.75562864,
                    "user": 5.69187959
                },
                "tags": {
                    "cpu": "cpu-total",
                    "guance_site": "testing",
                    "host": "172.16.242.112",
                    "host_ip": "172.16.242.112",
                    "name": "cpu",
                    "project": "cloudcare-testing"
                }
            },
            {
                "columns": {
                    "time": 1744866102975,
                    "total": 16.59466328,
                    "user": 5.28589581
                },
                "tags": {
                    "cpu": "cpu-total",
                    "guance_site": "testing",
                    "host": "172.16.242.112",
                    "host_ip": "172.16.242.112",
                    "name": "cpu",
                    "project": "cloudcare-testing"
                }
            }
        ]
    ],
    "status_code": 200
}`},
	},
}

var cTrigger = &FuncExample{
	FnName: FnTriggerDesc.Name,
	Progs: []ProgCase{
		{
			Name: "trigger",
			Script: `trigger(1, "critical", {"tag_abc": "1"}, {"a": "1", "a1": 2.1})

trigger(2, dim_tags={"a": "1", "b": "2"}, related_data={"b": {}})

trigger(false, related_data={"a": 1, "b": 2}, level="critical")

trigger("hello",  dim_tags={}, related_data={"a": 1, "b": [1]}, level="critical")
`,
			TriggerResult: []trigger.Data{
				{
					Result:      int64(1),
					Level:       "critical",
					DimTags:     map[string]string{"tag_abc": "1"},
					RelatedData: map[string]any{"a": "1", "a1": float64(2.1)},
				},
				{
					Result:      int64(2),
					Level:       "",
					DimTags:     map[string]string{"a": "1", "b": "2"},
					RelatedData: map[string]any{"b": map[string]any{}},
				},
				{
					Result:      false,
					Level:       "critical",
					DimTags:     map[string]string{},
					RelatedData: map[string]any{"a": int64(1), "b": int64(2)},
				},
				{
					Result:      "hello",
					Level:       "critical",
					DimTags:     map[string]string{},
					RelatedData: map[string]any{"a": int64(1), "b": []any{int64(1)}},
				},
			},
		},
	},
}

// cGeoIP = &FuncExample{
// 	FnName: FnGeoIPDesc.Name,
// 	Progs: []ProgCase{
// 		{
// 			Name: "geoip",
// 			Script: `v, ok = geoip("8.8.8.8")
// if ok {
// 	printf("%v", v)
// }
// `,
// 			Stdout: "US\n",
// 		},
// 	},
// }
