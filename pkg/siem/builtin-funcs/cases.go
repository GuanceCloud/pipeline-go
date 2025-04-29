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

var FnExps = map[string]*FuncExample{}

func AddExps(f *FuncExample) struct{} {
	FnExps[f.FnName] = f
	return struct{}{}
}

var _ = AddExps(cAppend)
var cAppend = &FuncExample{
	FnName: FnAppendDesc.Name,
	Progs: []ProgCase{
		{
			Name: "append",
			Script: `v = [1, 2, 3]
v = append(v, 4)
printf("%v", v)
`,
			Stdout: "[1,2,3,4]",
		},
		{
			Name: "append",
			Script: `v = [1, 2, 3]
v = append(v, "a", 1.1)
printf("%v", v)
`,
			Stdout: "[1,2,3,\"a\",1.1]",
		},
	},
}

var _ = AddExps(cB64Dec)
var cB64Dec = &FuncExample{
	FnName: FnB64DecDesc.Name,
	Progs: []ProgCase{
		{
			Name: "b64dec",
			Script: `v = "aGVsbG8sIHdvcmxk"
v, ok = b64dec(v)
if ok {
	printf("%v", v)
}
`,
			Stdout: "hello, world",
		},
	},
}

var _ = AddExps(cB64Enc)
var cB64Enc = &FuncExample{
	FnName: FnB64EncDesc.Name,
	Progs: []ProgCase{
		{
			Name: "b64enc",
			Script: `v = "hello, world"
v = b64enc(v)
printf("%v", v)
`,
			Stdout: "aGVsbG8sIHdvcmxk",
		},
	},
}

var _ = AddExps(cCast)
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

var _ = AddExps(cCIDR)
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

var _ = AddExps(cDelete)
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

var _ = AddExps(cExit)
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

var _ = AddExps(cDQL)
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

var _ = AddExps(cDumpJSON)
var cDumpJSON = &FuncExample{
	FnName: FnDumpJSONDesc.Name,
	Progs: []ProgCase{
		{
			Name: "dump_json",
			Script: `v = {"a": 1, "b": 2.1}
v, ok = dump_json(v)
if ok {
	printf("%v", v)
}
`,
			jsonout: true,
			Stdout:  "{\"a\":1,\"b\":2.1}\n",
		},
		{
			Name: "dump_json",
			Script: `v = {"a": 1, "b": 2.1}
v, ok = dump_json(v, "  ")
if ok {
	printf("%v", v)
}
`,
			jsonout: true,
			Stdout: `{
  "a": 1,
  "b": 2.1
}
`,
		},
	},
}

var _ = AddExps(cTrigger)
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

var _ = AddExps(cGeoIP)
var cGeoIP = &FuncExample{
	FnName: FnGeoIPDesc.Name,
	Progs: []ProgCase{
		{
			Name: "geoip",
			Script: `v = geoip("127.0.0.1")
printf("%v", v)
`,
			jsonout: true,
			Stdout:  `{"city":"","country":"","isp":"unknown","province":""}`,
		},
		{
			Name: "geoip",
			Script: `ip_addr = "114.114.114.114"
v, ok = dump_json(geoip(ip_addr), "    ");
if ok {
	printf("%v", v)
}
`,
			jsonout: true,
			Stdout: ` {
    "city": "Ji'an",
    "country": "CN",
    "isp": "chinanet",
    "province": "Jiangxi"
}`,
		},
	},
}

var _ = AddExps(cGJSON)
var cGJSON = &FuncExample{
	FnName: FnGJSONDesc.Name,
	Progs: []ProgCase{
		{
			Name: "gjson",
			Script: `v='''{
    "name": {"first": "Tom", "last": "Anderson"},
    "age": 37,
    "children": ["Sara","Alex","Jack"],
    "fav.movie": "Deer Hunter",
    "friends": [
        {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
        {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
        {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
    ]
}'''
age, ok = gjson(v, "age")
if ok {
	printf("%.0f", age)
} else {
	printf("not found")
}
`,
			Stdout: `37`,
		},
		{
			Name: "gjson",
			Script: `v='''{
    "name": {"first": "Tom", "last": "Anderson"},
    "age": 37,
    "children": ["Sara","Alex","Jack"],
    "fav.movie": "Deer Hunter",
    "friends": [
        {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
        {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
        {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
    ]
}'''
name, ok = gjson(v, "name")
printf("%v", name)
`,
			jsonout: true,
			Stdout:  `{"first": "Tom", "last": "Anderson"}`,
		},
		{
			Name: "gjson",
			Script: `v='''[
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
]'''
net, ok = gjson(v, "0.nets.2")
printf("%v", net)
`,
			Stdout: `tw`,
		},
	},
}

var _ = AddExps(cGrok)
var cGrok = &FuncExample{
	FnName: FnGrokDesc.Name,
	Progs: []ProgCase{
		{
			Name: "grok",
			Script: `app_log="2021-01-11T17:43:51.887+0800  DEBUG io  io/io.go:458  post cost 6.87021ms"

# Use built-in patterns, named capture groups, custom patterns, extract fields;
# convert the type of the extracted field by specifying the type.
v, ok = grok(
	app_log,
	"%{TIMESTAMP_ISO8601:log_time}\\s+(?P<log_level>[a-zA-Z]+)\\s+%{WORD}\\s+%{log_code_pos_pattern:log_code_pos}.*\\s%{NUMBER:log_cost:float}ms", 
	{
		"log_code_pos_pattern": "[a-zA-Z0-9/\\.]+:\\d+", 
	}
)

if ok {
	v, ok = dump_json(v, "  ")
	if ok {
		printf("%v", v)
	}
}
`,
			jsonout: true,
			Stdout: `{
  "log_code_pos": "io/io.go:458",
  "log_cost": 6.87021,
  "log_level": "DEBUG",
  "log_time": "2021-01-11T17:43:51.887+0800"
}`,
		},
	},
}

var _ = AddExps(cHash)
var cHash = &FuncExample{
	FnName: FnHashDesc.Name,
	Progs: []ProgCase{
		{
			Name: "hash",
			Script: `printf("%v", hash("abc", "md5"))
`,
			Stdout: "900150983cd24fb0d6963f7d28e17f72",
		},
		{
			Name: "hash",
			Script: `printf("%v", hash("abc", "sha1"))
`,
			Stdout: "a9993e364706816aba3e25717850c26c9cd0d89d",
		},
		{
			Name: "hash",
			Script: `printf("%v", hash("abc", "sha256"))
`,
			Stdout: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},
		{
			Name: "hash",
			Script: `printf("%v", hash("abc", "sha512"))
`,
			Stdout: "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f",
		},
		{
			Name: "hash",
			Script: `printf("%v", hash("abc", "xx"))
`,
			Stdout: "",
		},
	},
}

var _ = AddExps(cLen)
var cLen = &FuncExample{
	FnName: FnLenDesc.Name,
	Progs: []ProgCase{
		{
			Name: "len",
			Script: `printf("%v", len("abc"))
`,
			Stdout: "3",
		},
		{
			Name: "len",
			Script: `printf("%v", len([1, 2, 3]))
`,
			Stdout: "3",
		},
		{
			Name: "len",
			Script: `printf("%v", len({"a": 1, "b": 2, "c": 3}))
`,
			Stdout: "3",
		},
	},
}

var _ = AddExps(cLoadJSON)
var cLoadJSON = &FuncExample{
	FnName: FnLoadJSONDesc.Name,
	Progs: []ProgCase{
		{
			Name: "load_json",
			Script: `jstr = '{"a": 1, "b": 2, "c": 3}'
v, ok = load_json(jstr)
if ok {
	printf("%v", v["b"])
}
`,
			Stdout: `2`,
		},
	},
}

var _ = AddExps(cLowercase)
var cLowercase = &FuncExample{
	FnName: FnLowercaseDesc.Name,
	Progs: []ProgCase{
		{
			Name: "lower_case",
			Script: `printf("%s", lowercase("ABC"))
`,
			Stdout: "abc",
		},
	},
}

var _ = AddExps(cMatch)
var cMatch = &FuncExample{
	FnName: FnMatchDesc.Name,
	Progs: []ProgCase{
		{
			Name: "match",
			Script: `text="abc def 123 abc def 123"
v, ok = match(text, "(abc) (?:def) (?P<named_group>123)")
if ok {
	printf("%v", v)
}
`,
			jsonout: true,
			Stdout:  `[["abc def 123","abc","123"]]`,
		},
		{
			Name: "match",
			Script: `text="abc def 123 abc def 123"
v, ok = match(text, "(abc) (?:def) (?P<named_group>123)", -1)
if ok {
	printf("%v", v)
}
`,
			jsonout: true,
			Stdout:  `[["abc def 123","abc","123"],["abc def 123","abc","123"]]`,
		},
	},
}

var _ = AddExps(cParseDate)
var cParseDate = &FuncExample{
	FnName: FnParseDateDesc.Name,
	Progs: []ProgCase{
		{
			Name: "parse_date",
			Script: `v, ok = parse_date("2021-12-2T11:55:43.123+0800")
if ok {
	printf("%v", v)
}
`,
			Stdout: "1638417343123000000",
		},
		{
			Name: "parse_date",
			Script: `v, ok = parse_date("2021-12-2T11:55:43.123", "+8")
if ok {
	printf("%v", v)
}
`,
			Stdout: "1638417343123000000",
		},
		{
			Name: "parse_date",
			Script: `v, ok = parse_date("2021-12-2T11:55:43.123", "Asia/Shanghai")
if ok {
	printf("%v", v)
}
`,
			Stdout: "1638417343123000000",
		},
	},
}

var _ = AddExps(cParseDuration)
var cParseDuration = &FuncExample{
	FnName: FnParseDurationDesc.Name,
	Progs: []ProgCase{
		{
			Name: "parse_duration",
			Script: `v, ok = parse_duration("1s")
if ok {
	printf("%v", v)
}
`,
			Stdout: "1000000000",
		},
		{
			Name: "parse_duration",
			Script: `v, ok = parse_duration("100ns")
if ok {
	printf("%v", v)
}
`,
			Stdout: "100",
		},
	},
}

var _ = AddExps(cParseInt)
var cParseInt = &FuncExample{
	FnName: FnParseIntDesc.Name,
	Progs: []ProgCase{
		{
			Name: "parse_int",
			Script: `v, ok = parse_int("123", 10)
if ok {
	printf("%v", v)
}
`,
			Stdout: "123",
		},
		{
			Name: "parse_int",
			Script: `v, ok = parse_int("123", 16)	
if ok {
	printf("%v", v)
}
`,
			Stdout: "291",
		},
	},
}

var _ = AddExps(cFormatInt)
var cFormatInt = &FuncExample{
	FnName: FnFormatIntDesc.Name,
	Progs: []ProgCase{
		{
			Name: "format_int",
			Script: `v = format_int(16, 16)
printf("%s", v)
`,
			Stdout: "10",
		},
	},
}

var _ = AddExps(cPrintf)
var cPrintf = &FuncExample{
	FnName: FnPrintfDesc.Name,
	Progs: []ProgCase{
		{
			Name: "printf",
			Script: `printf("hello, %s", "world")
`,
			Stdout: "hello, world",
		},
	},
}

var _ = AddExps(cReplace)
var cReplace = &FuncExample{
	FnName: FnReplaceDesc.Name,
	Progs: []ProgCase{
		{
			Name: "replace",
			Script: `v, ok = replace("abcdef", "bc", "123")
printf("%s", v)
`,
			Stdout: "a123def",
		},
		{
			Name: "replace",
			Script: `v, ok = replace("bonjour; 你好", "[\u4e00-\u9fa5]+", "hello")
printf("%s", v)
`,
			Stdout: "bonjour; hello",
		},
	},
}

var _ = AddExps(cSQLCover)
var cSQLCover = &FuncExample{
	FnName: FnSQLCoverDesc.Name,
	Progs: []ProgCase{
		{
			Name: "sql_cover",
			Script: `v, ok = sql_cover("select abc from def where x > 3 and y < 5")
if ok {
	printf("%s",v)
}
`,
			Stdout: "select abc from def where x > ? and y < ?",
		},
		{
			Name: "sql_cover",
			Script: `v, ok = sql_cover("SELECT $func$INSERT INTO table VALUES ('a', 1, 2)$func$ FROM users")
if ok {
	printf("%s",v)
}
`,
			Stdout: "SELECT ? FROM users",
		},
		{
			Name: "sql_cover",
			Script: `v, ok = sql_cover("SELECT ('/uffd')")
if ok {
	printf("%s",v)
}
`,
			Stdout: "SELECT ( ? )",
		},
	},
}

var _ = AddExps(cStrJoin)
var cStrJoin = &FuncExample{
	FnName: FnStrJoinDesc.Name,
	Progs: []ProgCase{
		{
			Name: "strjoin",
			Script: `v = str_join(["a", "b", "c"], "##")
printf("%s", v)
`,
			Stdout: "a##b##c",
		},
	},
}

var _ = AddExps(cStrFmt)
var cStrFmt = &FuncExample{
	FnName: FnStrFmtDesc.Name,
	Progs: []ProgCase{
		{
			Name: "strfmt",
			Script: `v = strfmt("abc %s def %d", "123", 456)
printf("%s", v)
`,
			Stdout: "abc 123 def 456",
		},
	},
}

var _ = AddExps(cTimeNow)
var cTimeNow = &FuncExample{
	FnName: FnTimeNowDesc.Name,
	Progs: []ProgCase{
		{
			Name: "timenow",
			Script: `printf("%v", time_now("s"))
`,
			Stdout: "1745823860",
		},
	},
}

var _ = AddExps(cTrim)
var cTrim = &FuncExample{
	FnName: FnTrimDesc.Name,
	Progs: []ProgCase{
		{
			Name: "trim",
			Script: `printf("%s", trim(" abcdef "))
`,
			Stdout: "abcdef",
		},
		{
			Name: "trim",
			Script: `printf("%s", trim("#-abcdef-#", "-#", 2))
`,
			Stdout: "#-abcdef",
		},
		{
			Name: "trim",
			Script: `printf("%s", trim("#-abcdef-#", "-#", 1))
`,
			Stdout: "abcdef-#",
		},
		{
			Name: "trim",
			Script: `printf("%s", trim("#-abcdef-#", side=0, cutset="-#"))
`,
			Stdout: "abcdef",
		},
	},
}

var _ = AddExps(cUppercase)
var cUppercase = &FuncExample{
	FnName: FnUppercaseDesc.Name,
	Progs: []ProgCase{
		{
			Name: "upper_case",
			Script: `printf("%s", uppercase("abc"))
`,
			Stdout: "ABC",
		},
	},
}

var _ = AddExps(cURLParse)
var cURLParse = &FuncExample{
	FnName: FnURLParseDesc.Name,
	Progs: []ProgCase{
		{
			Name: "url_parse",
			Script: `v, ok = url_parse("http://www.example.com:8080/path/to/file?query=abc")
if ok {
	v, ok = dump_json(v, "  ")
	if ok {
		printf("%v", v)
	}
}
`,
			jsonout: true,
			Stdout: `{
  "host": "www.example.com:8080",
  "params": {
    "query": [
      "abc"
    ]
  },
  "path": "/path/to/file",
  "port": "8080",
  "scheme": "http"
}`,
		},
	},
}

var _ = AddExps(cURLDecode)
var cURLDecode = &FuncExample{
	FnName: FnURLDecodeDesc.Name,
	Progs: []ProgCase{
		{
			Name: "url_decode",
			Script: `v, ok = url_decode("https:%2F%2Fkubernetes.io%2Fdocs%2Freference%2Faccess-authn-authz%2Fbootstrap-tokens%2F")
if ok {
	printf("%s", v)
}
`,
			Stdout: `https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/`,
		},
	},
}

var _ = AddExps(cUserAgent)
var cUserAgent = &FuncExample{
	FnName: FnUserAgentDesc.Name,
	Progs: []ProgCase{
		{
			Name: "user_agent",
			Script: `v = user_agent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
printf("%s", v)
`,
			jsonout: true,
			Stdout:  `{"browser":"Chrome","browserVer":"96.0.4664.110","engine":"AppleWebKit","engineVer":"537.36","isBot":false,"isMobile":false,"os":"Intel Mac OS X 10_15_7","ua":"Macintosh"}`,
		},
	},
}

var _ = AddExps(cValidJSON)
var cValidJSON = &FuncExample{
	FnName: FnValidJSONDesc.Name,
	Progs: []ProgCase{
		{
			Name: "valid_json",
			Script: `ok = valid_json("{\"a\": 1, \"b\": 2}")
if ok {
	printf("true")
}
`,
			Stdout: `true`,
		},
		{
			Name: "valid_json",
			Script: `ok = valid_json("1.1")
if ok {
	printf("true")
}
`,
			Stdout: `true`,
		},
		{
			Name: "valid_json",
			Script: `ok = valid_json("str_abc_def")
if !ok {
	printf("false")
}
`,
			Stdout: `false`,
		},
	},
}

var _ = AddExps(cValueType)
var cValueType = &FuncExample{
	FnName: FnValueTypeDesc.Name,
	Progs: []ProgCase{
		{
			Name: "value_type",
			Script: `v = value_type(1)
printf("%s", v)
`,
			Stdout: "int",
		},
		{
			Name: "value_type",
			Script: `printf("%s", value_type("abc"))
`,
			Stdout: "str",
		},
		{
			Name: "value_type",
			Script: `printf("%s", value_type(true))
`,
			Stdout: "bool",
		},
	},
}

var _ = AddExps(cXMLTest)
var cXMLTest = &FuncExample{
	FnName: FnXMLDesc.Name,
	Progs: []ProgCase{
		{
			Name: "xml_query",
			Script: `xml_data='''
<OrderEvent actionCode = "5">
 <OrderNumber>ORD12345</OrderNumber>
 <VendorNumber>V11111</VendorNumber>
 </OrderEvent>
'''
v, ok = xml_query(xml_data, "/OrderEvent/OrderNumber/text()")
if ok {
	printf("%s", v)
}
`,
			Stdout: `ORD12345`,
		},
		{
			Name: "xml_query",
			Script: `xml_data='''
<OrderEvent actionCode = "5">
 <OrderNumber>ORD12345</OrderNumber>
 <VendorNumber>V11111</VendorNumber>
 </OrderEvent>
'''
v, ok = xml_query(xml_data, "/OrderEvent/@actionCode")
if ok {
	printf("%s", v)
}
`,
			Stdout: `5`,
		},
	},
}
