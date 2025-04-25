# SIEM Built-In Function

## `b64dec` {#fn-b64dec}

Function prototype: `fn b64dec(data: str) -> (str, bool)`

Function description: Base64 decoding.

Function parameters:

- `data`: Data that needs to be base64 decoded.

Function returns:

- `str`: The decoded string.
- `bool`: Whether decoding is successful.

## `b64enc` {#fn-b64enc}

Function prototype: `fn b64enc(data: str) -> (str, bool)`

Function description: Base64 encoding.

Function parameters:

- `data`: Data that needs to be base64 encoded.

Function returns:

- `str`: The encoded string.
- `bool`: Whether encoding is successful.

## `cast` {#fn-cast}

Function prototype: `fn cast(val: bool|int|float|str, typ: str) -> bool|int|float|str`

Function description: Convert the value to the target type.

Function parameters:

- `val`: The value of the type to be converted.
- `typ`: Target type. One of (`bool`, `int`, `float`, `str`).

Function returns:

- `bool|int|float|str`: The value after the conversion.

Function examples:

* CASE 0:

Script content:

```py
v1 = "1.1"
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
```

Standard output:

```txt
2.1; 2; -1; 2; false; 1.1; 1; true
```
## `cidr` {#fn-cidr}

Function prototype: `fn cidr(ip: str, mask: str) -> bool`

Function description: Check the IP whether in CIDR block

Function parameters:

- `ip`: The ip address
- `mask`: The CIDR mask

Function returns:

- `bool`: Whether the IP is in CIDR block

Function examples:

* CASE 0:

Script content:

```py
ip = "192.0.2.233"
if cidr(ip, "192.0.2.1/24") {
	printf("%s", ip)
}
```
Standard output:

```txt
192.0.2.233
```
* CASE 1:

Script content:

```py
ip = "192.0.2.233"
if cidr(mask="192.0.1.1/24", ip=ip) {
	printf("%s", ip)
}
```
Standard output:

```txt

```
## `delete` {#fn-delete}

Function prototype: `fn delete(m: map, key: str)`

Function description: Delete key from the map.

Function parameters:

- `m`: The map for deleting key
- `key`: Key need delete from map.

Function examples:

* CASE 0:

Script content:

```py
v = {
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
```

Standard output:

```txt
result group 1: {"b":2}; [{"c":1.1}]
result group 2: {}; {}
```
## `dql` {#fn-dql}

Function prototype: `fn dql(query: str, qtype: str = "dql", limit: int = 2000, offset: int = 0, slimit: int = 2000, time_range: list = []) -> (map, bool)`

Function description: Query data from the GuanceCloud using dql or promql.

Function parameters:

- `query`: DQL or PromQL query statements.
- `qtype`: Query language, One of `dql` or `promql`, default is `dql`.
- `limit`: Query limit.
- `offset`: Query offset.
- `slimit`: Query slimit.
- `time_range`: Query timestamp range, the default value can be modified externally by the script caller.

Function returns:

- `map`: Query response.
- `bool`: Query execution status

Function examples:

* CASE 0:

Script content:

```py
v, ok = dql("M::cpu limit 3 slimit 3")
if ok {
	v, ok = dump_json(v, "    ")
	if ok {
		printf("%v", v)
	}
}
```

Standard output:

```txt
{
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
}
```
## `dump_json` {#fn-dump_json}

Function prototype: `fn dump_json(v: str, indent: str = "") -> (str, bool)`

Function description: Returns the JSON encoding of v.

Function parameters:

- `v`: Object to encode.
- `indent`: Indentation prefix.

Function returns:

- `str`: JSON encoding of v.
- `bool`: Whether decoding is successful.

## `exit` {#fn-exit}

Function prototype: `fn exit()`

Function description: Exit the program
Function examples:

* CASE 0:

Script content:

```py
printf("1\n")
printf("2\n")
exit()
printf("3\n")
	
```
Standard output:

```txt
1
2
```
## `geoip` {#fn-geoip}

Function prototype: `fn geoip(ip: str) -> map`

Function description: GeoIP

Function parameters:

- `ip`: IP address.

Function returns:

- `map`: IP geographical information.

## `gjson` {#fn-gjson}

Function prototype: `fn gjson(input: str, json_path: str) -> (bool|int|float|str|list|map, bool)`

Function description: GJSON provides a fast and easy way to get values from a JSON document.

Function parameters:

- `input`: JSON format string to parse.
- `json_path`: JSON path.

Function returns:

- `bool|int|float|str|list|map`: Parsed result.
- `bool`: Parsed status.

## `grok` {#fn-grok}

Function prototype: `fn grok(input: str, pattern: str, extra_patterns: map = {}, trim_space: bool = true) -> (map, bool)`

Function description: Extracts data from a string using a Grok pattern. Grok is based on regular expression syntax, and using regular (named) capture groups in a pattern is equivalent to using a pattern in a pattern. A valid regular expression is also a valid Grok pattern.

Function parameters:

- `input`: The input string used to extract data.
- `pattern`: The pattern used to extract data.
- `extra_patterns`: Additional patterns for parsing patterns.
- `trim_space`: Whether to trim leading and trailing spaces from the parsed value.

Function returns:

- `map`: The parsed result.
- `bool`: Whether the parsing was successful.

## `hash` {#fn-hash}

Function prototype: `fn hash(text: str, method: str) -> (str, bool)`

Function description: 

Function parameters:

- `text`: The string used to calculate the hash.
- `method`: Hash Algorithms, allowing values including `md5`, `sha1`, `sha256`, `sha512`.

Function returns:

- `str`: The hash value.
- `bool`: Hash calculation status.

## `len` {#fn-len}

Function prototype: `fn len(val: map|list|str) -> int`

Function description: Get the length of the value. If the value is a string, returns the length of the string. If the value is a list or map, returns the length of the list or map. If it is neither, returns -1.

Function parameters:

- `val`: The value to get the length of.

Function returns:

- `int`: The length of the value.

## `load_json` {#fn-load_json}

Function prototype: `fn load_json(val: str) -> (bool|int|float|str|list|map, bool)`

Function description: Unmarshal json string

Function parameters:

- `val`: JSON string.

Function returns:

- `bool|int|float|str|list|map`: Unmarshal result.
- `bool`: Unmarshal status.

## `lowercase` {#fn-lowercase}

Function prototype: `fn lowercase(val: str) -> str`

Function description: Converts a string to lowercase.

Function parameters:

- `val`: The string to convert.

Function returns:

- `str`: Returns the lowercase value.

## `match` {#fn-match}

Function prototype: `fn match(val: str, pattern: str) -> (list, bool)`

Function description: Regular expression matching.

Function parameters:

- `val`: The string to match.
- `pattern`: Regular expression pattern.

Function returns:

- `list`: Returns the matched value.
- `bool`: Returns true if the regular expression matches.

## `parse_date` {#fn-parse_date}

Function prototype: `fn parse_date(date: str, timezone: str = "") -> (int, bool)`

Function description: Parses a date string to a nanoseconds timestamp, support multiple date formats. If the date string not include timezone and no timezone is provided, the local timezone is used.

Function parameters:

- `date`: The key to use for parsing.
- `timezone`: The timezone to use for parsing. If 

Function returns:

- `int`: The parsed timestamp in nanoseconds.
- `bool`: Whether the parsing was successful.

## `parse_duration` {#fn-parse_duration}

Function prototype: `fn parse_duration(s: str) -> (int, int)`

Function description: Parses a golang duration string into a duration. A duration string is a sequence of possibly signed decimal numbers with optional fraction and unit suffixes for each number, such as `300ms`, `-1.5h` or `2h45m`. Valid units are `ns`, `us` (or `μs`), `ms`, `s`, `m`, `h`. 

Function parameters:

- `s`: The string to parse.

Function returns:

- `int`: The duration in nanoseconds.
- `int`: The duration in nanoseconds.

## `parse_int` {#fn-parse_int}

Function prototype: `fn parse_int(val: str, base: int) -> (int, bool)`

Function description: Parses a string into an integer.

Function parameters:

- `val`: The string to parse.
- `base`: The base to use for parsing.

Function returns:

- `int`: The parsed integer.
- `bool`: Whether the parsing was successful.

## `printf` {#fn-printf}

Function prototype: `fn printf(format: str, args: ...str|bool|int|float|list|map)`

Function description: Output formatted strings to the standard output device.

Function parameters:

- `format`: String format.
- `args`: Argument list, corresponding to the format specifiers in the format string.

## `replace` {#fn-replace}

Function prototype: `fn replace(input: str, pattern: str, replacement: str) -> str`

Function description: Replaces text in a string.

Function parameters:

- `input`: The string to replace text in.
- `pattern`: Regular expression pattern.
- `replacement`: Replacement text to use.

Function returns:

- `str`: The string with text replaced.

## `sql_cover` {#fn-sql_cover}

Function prototype: `fn sql_cover(val: str) -> (str, bool)`

Function description: Obfuscate SQL query string.

Function parameters:

- `val`: The sql to obfuscate.

Function returns:

- `str`: The obfuscated sql.
- `bool`: The obfuscate status.

## `strfmt` {#fn-strfmt}

Function prototype: `fn strfmt(format: str, args: ...bool|int|float|str|list|map)`

Function description: 

Function parameters:

- `format`: String format.
- `args`: Parameters to replace placeholders.

## `time_now` {#fn-time_now}

Function prototype: `fn time_now(precision: str = "ns") -> int`

Function description: Get current timestamp with the specified precision.

Function parameters:

- `precision`: The precision of the timestamp. Supported values: `ns`, `us`, `ms`, `s`.

Function returns:

- `int`: Returns the current timestamp.

## `trigger` {#fn-trigger}

Function prototype: `fn trigger(result: int|float|bool|str, level: str = "", dim_tags: map = {}, related_data: map = {})`

Function description: Trigger a security event.

Function parameters:

- `result`: Event check result.
- `level`: Event level. One of: (`critical`, `high`, `medium`, `low`, `info`).
- `dim_tags`: Dimension tags.
- `related_data`: Related data.

Function examples:

* CASE 0:

Script content:

```py
trigger(1, "critical", {"tag_abc": "1"}, {"a": "1", "a1": 2.1})

trigger(2, dim_tags={"a": "1", "b": "2"}, related_data={"b": {}})

trigger(false, related_data={"a": 1, "b": 2}, level="critical")

trigger("hello",  dim_tags={}, related_data={"a": 1, "b": [1]}, level="critical")
```

Standard output:

```txt

```
Trigger output:
```json
[
    {
        "result": 1,
        "level": "critical",
        "dim_tags": {
            "tag_abc": "1"
        },
        "related_data": {
            "a": "1",
            "a1": 2.1
        }
    },
    {
        "result": 2,
        "level": "",
        "dim_tags": {
            "a": "1",
            "b": "2"
        },
        "related_data": {
            "b": {}
        }
    },
    {
        "result": false,
        "level": "critical",
        "dim_tags": {},
        "related_data": {
            "a": 1,
            "b": 2
        }
    },
    {
        "result": "hello",
        "level": "critical",
        "dim_tags": {},
        "related_data": {
            "a": 1,
            "b": [
                1
            ]
        }
    }
]

```
## `trim` {#fn-trim}

Function prototype: `fn trim(val: str, cutset: str = "", side: int = 0) -> str`

Function description: Removes leading and trailing whitespace from a string.

Function parameters:

- `val`: The string to trim.
- `cutset`: Characters to remove from the beginning and end of the string. If not specified, whitespace is removed.
- `side`: The side to trim from. If value is 0, trim from both sides. If value is 1, trim from the left side. If value is 2, trim from the right side.

Function returns:

- `str`: The trimmed string.

## `uppercase` {#fn-uppercase}

Function prototype: `fn uppercase(val: str) -> str`

Function description: Converts a string to uppercase.

Function parameters:

- `val`: The string to convert.

Function returns:

- `str`: Returns the uppercase value.

## `url_decode` {#fn-url_decode}

Function prototype: `fn url_decode(val: str) -> (str, bool)`

Function description: Decodes a URL-encoded string.

Function parameters:

- `val`: The URL-encoded string to decode.

Function returns:

- `str`: The decoded string.
- `bool`: The decoding status.

## `url_parse` {#fn-url_parse}

Function prototype: `fn url_parse(url: str) -> (map, bool)`

Function description: Parses a URL and returns it as a map.

Function parameters:

- `url`: The URL to parse.

Function returns:

- `map`: Returns the parsed URL as a map.
- `bool`: Returns true if the URL is valid.

## `user_agent` {#fn-user_agent}

Function prototype: `fn user_agent(header: str) -> map`

Function description: Parses a User-Agent header.

Function parameters:

- `header`: The User-Agent header to parse.

Function returns:

- `map`: Returns the parsed User-Agent header as a map.

## `valid_json` {#fn-valid_json}

Function prototype: `fn valid_json(val: str) -> bool`

Function description: Returns true if the value is a valid JSON.

Function parameters:

- `val`: The value to check.

Function returns:

- `bool`: Returns true if the value is a valid JSON.

## `value_type` {#fn-value_type}

Function prototype: `fn value_type(val: str) -> str`

Function description: Returns the type of the value.

Function parameters:

- `val`: The value to get the type of.

Function returns:

- `str`: Returns the type of the value. One of (`bool`, `int`, `float`, `str`, `list`, `map`, `nil`). If the value and the type is nil, returns `nil`.

## `xml` {#fn-xml}

Function prototype: `fn xml(input: str, xpath: str) -> (str, bool)`

Function description: Returns the value of an XML field.

Function parameters:

- `input`: The XML input to get the value of.
- `xpath`: The XPath expression to get the value of.

Function returns:

- `str`: Returns the value of the XML field.
- `bool`: Returns true if the field exists, false otherwise.
