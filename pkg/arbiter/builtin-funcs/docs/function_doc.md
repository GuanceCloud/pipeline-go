# Arbiter Built-In Function

## `append` {#fn-append}

Function prototype: `fn append(li: list, v: ...bool|int|float|str|list|map) -> list`

Function description: Appends a value to a list.

Function parameters:

- `li`: The list to append to.
- `v`: The value to append.


Function returns:

- `list`: The list with the appended value.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = [1, 2, 3]
    v = append(v, 4)
    printf("%v", v)
    
    ```

    Standard output:

    ```txt
    [1,2,3,4]
    ```

    
* Case 1:

    Script content:

    ```txt
    v = [1, 2, 3]
    v = append(v, "a", 1.1)
    printf("%v", v)
    
    ```

    Standard output:

    ```txt
    [1,2,3,"a",1.1]
    ```

    

## `b64dec` {#fn-b64dec}

Function prototype: `fn b64dec(data: str) -> (str, bool)`

Function description: Base64 decoding.

Function parameters:

- `data`: Data that needs to be base64 decoded.


Function returns:

- `str`: The decoded string.
- `bool`: Whether decoding is successful.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = "aGVsbG8sIHdvcmxk"
    v, ok = b64dec(v)
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    hello, world
    ```

    

## `b64enc` {#fn-b64enc}

Function prototype: `fn b64enc(data: str) -> (str, bool)`

Function description: Base64 encoding.

Function parameters:

- `data`: Data that needs to be base64 encoded.


Function returns:

- `str`: The encoded string.
- `bool`: Whether encoding is successful.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = "hello, world"
    v = b64enc(v)
    printf("%v", v)
    
    ```

    Standard output:

    ```txt
    aGVsbG8sIHdvcmxk
    ```

    

## `call_func` {#fn-call_func}

Function prototype: `fn call_func(name: str, kwargs: map = {}) -> map`

Function description: Calling remote functions on the Function platform

Function parameters:

- `name`: Remote function name.
- `kwargs`: Parameter map, corresponding to **kwargs in Python.


Function returns:

- `map`: 


Function examples:

* Case 0:

    Script content:

    ```txt
    result = call_func("echo", {
        "arg_1": "1",
    	"arg_2": [12,3],
    	"key3": "true"
    })
    printf("%v", result)
    ```

    Standard output:

    ```txt
    {"error":"","message":"","ok":true,"reason":"","result":{"arg_1":"1","arg_2":[12,3],"kwargs":{"key3":"true"}}}
    ```

    
* Case 1:

    Script content:

    ```txt
    result = call_func("echo", {"arg_2": [12,3]})
    printf("%v", result)
    ```

    Standard output:

    ```txt
    {"detail":{"exeception":"TypeError(\"echo() missing 1 required positional argument: 'arg_1'\")"},"error":"call function failed, api status code 599","message":"Func task failed","ok":false,"reason":"EFuncFailed"}
    ```

    

## `cast` {#fn-cast}

Function prototype: `fn cast(val: bool|int|float|str, typ: str) -> bool|int|float|str`

Function description: Convert the value to the target type.

Function parameters:

- `val`: The value of the type to be converted.
- `typ`: Target type. One of (`bool`, `int`, `float`, `str`).


Function returns:

- `bool|int|float|str`: The value after the conversion.


Function examples:

* Case 0:

    Script content:

    ```txt
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

* Case 0:

    Script content:

    ```txt
    ip = "192.0.2.233"
    if cidr(ip, "192.0.2.1/24") {
    	printf("%s", ip)
    }
    ```

    Standard output:

    ```txt
    192.0.2.233
    ```

    
* Case 1:

    Script content:

    ```txt
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

* Case 0:

    Script content:

    ```txt
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

Function prototype: `fn dql(query: str, qtype: str = "dql", limit: int = 10000, offset: int = 0, slimit: int = 0, time_range: list = []) -> map`

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


Function examples:

* Case 0:

    Script content:

    ```txt
    v = dql("M::cpu limit 2 slimit 2")
    v, ok = dump_json(v, "    ")
    if ok {
    	printf("%v", v)
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
                        "test_site": "testing",
                        "host": "172.16.241.111",
                        "host_ip": "172.16.241.111",
                        "name": "cpu",
                        "project": "demo-testing"
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
                        "test_site": "testing",
                        "host": "172.16.241.111",
                        "host_ip": "172.16.241.111",
                        "name": "cpu",
                        "project": "demo-testing"
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
                        "test_site": "testing",
                        "host": "172.16.242.112",
                        "host_ip": "172.16.242.112",
                        "name": "cpu",
                        "project": "demo-testing"
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
                        "test_site": "testing",
                        "host": "172.16.242.112",
                        "host_ip": "172.16.242.112",
                        "name": "cpu",
                        "project": "demo-testing"
                    }
                }
            ]
        ],
        "status_code": 200
    }
    ```

    

## `dql_series_get` {#fn-dql_series_get}

Function prototype: `fn dql_series_get(series: map, name: str) -> list`

Function description: get series data

Function parameters:

- `series`: dql query result
- `name`: column or tag name


Function returns:

- `list`: specified column or tag value for the series


Function examples:

* Case 0:

    Script content:

    ```txt
    v = dql("M::cpu limit 2 slimit 2") 
    
    hostLi = dql_series_get(v, "host")
    time_li = dql_series_get(v, "time")
    
    printf("%v", {"host": hostLi, "time": time_li})
    
    ```

    Standard output:

    ```txt
    {"host":[["172.16.241.111","172.16.241.111"],["172.16.242.112","172.16.242.112"]],"time":[[1744866108991,1744866103991],[1744866107975,1744866102975]]}
    ```

    

## `dql_timerange_get` {#fn-dql_timerange_get}

Function prototype: `fn dql_timerange_get() -> list`

Function description: Get the time range of the DQL query, which is passed in by the script caller or defaults to the last 15 minutes.

Function returns:

- `list`: The time range. For example, `[1744214400000, 1744218000000]`, the timestamp precision is milliseconds


Function examples:

* Case 0:

    Script content:

    ```txt
    val = dql_timerange_get()
    printf("%v", val)
    ```

    Standard output:

    ```txt
    [1672531500000,1672532100000]
    ```

    
* Case 1:

    Script content:

    ```txt
    val = dql_timerange_get()
    printf("%v", val)
    ```

    Standard output:

    ```txt
    [1672531200000,1672532100000]
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


Function examples:

* Case 0:

    Script content:

    ```txt
    v = {"a": 1, "b": 2.1}
    v, ok = dump_json(v)
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    {"a":1,"b":2.1}
    
    ```

    
* Case 1:

    Script content:

    ```txt
    v = {"a": 1, "b": 2.1}
    v, ok = dump_json(v, "  ")
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    {
      "a": 1,
      "b": 2.1
    }
    
    ```

    

## `exit` {#fn-exit}

Function prototype: `fn exit()`

Function description: Exit the program

Function examples:

* Case 0:

    Script content:

    ```txt
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

    

## `format_int` {#fn-format_int}

Function prototype: `fn format_int(val: int, base: int) -> str`

Function description: Formats an integer into a string.

Function parameters:

- `val`: The integer to format.
- `base`: The base to use for formatting. Must be between 2 and 36.


Function returns:

- `str`: The formatted string.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = format_int(16, 16)
    printf("%s", v)
    
    ```

    Standard output:

    ```txt
    10
    ```

    

## `geoip` {#fn-geoip}

Function prototype: `fn geoip(ip: str) -> map`

Function description: GeoIP

Function parameters:

- `ip`: IP address.


Function returns:

- `map`: IP geographical information.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = geoip("127.0.0.1")
    printf("%v", v)
    
    ```

    Standard output:

    ```txt
    {"city":"","country":"","isp":"unknown","province":""}
    ```

    
* Case 1:

    Script content:

    ```txt
    ip_addr = "114.114.114.114"
    v, ok = dump_json(geoip(ip_addr), "    ");
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
     {
        "city": "Ji'an",
        "country": "CN",
        "isp": "chinanet",
        "province": "Jiangxi"
    }
    ```

    

## `gjson` {#fn-gjson}

Function prototype: `fn gjson(input: str, json_path: str) -> (bool|int|float|str|list|map, bool)`

Function description: GJSON provides a fast and easy way to get values from a JSON document.

Function parameters:

- `input`: JSON format string to parse.
- `json_path`: JSON path.


Function returns:

- `bool|int|float|str|list|map`: Parsed result.
- `bool`: Parsed status.


Function examples:

* Case 0:

    Script content:

    ```txt
    v='''{
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
    
    ```

    Standard output:

    ```txt
    37
    ```

    
* Case 1:

    Script content:

    ```txt
    v='''{
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
    
    ```

    Standard output:

    ```txt
    {"first": "Tom", "last": "Anderson"}
    ```

    
* Case 2:

    Script content:

    ```txt
    v='''[
        {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
        {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
        {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
    ]'''
    net, ok = gjson(v, "0.nets.2")
    printf("%v", net)
    
    ```

    Standard output:

    ```txt
    tw
    ```

    

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


Function examples:

* Case 0:

    Script content:

    ```txt
    app_log="2021-01-11T17:43:51.887+0800  DEBUG io  io/io.go:458  post cost 6.87021ms"
    
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
    
    ```

    Standard output:

    ```txt
    {
      "log_code_pos": "io/io.go:458",
      "log_cost": 6.87021,
      "log_level": "DEBUG",
      "log_time": "2021-01-11T17:43:51.887+0800"
    }
    ```

    

## `hash` {#fn-hash}

Function prototype: `fn hash(text: str, method: str) -> str`

Function description: 

Function parameters:

- `text`: The string used to calculate the hash.
- `method`: Hash Algorithms, allowing values including `md5`, `sha1`, `sha256`, `sha512`.


Function returns:

- `str`: The hash value.


Function examples:

* Case 0:

    Script content:

    ```txt
    printf("%v", hash("abc", "md5"))
    
    ```

    Standard output:

    ```txt
    900150983cd24fb0d6963f7d28e17f72
    ```

    
* Case 1:

    Script content:

    ```txt
    printf("%v", hash("abc", "sha1"))
    
    ```

    Standard output:

    ```txt
    a9993e364706816aba3e25717850c26c9cd0d89d
    ```

    
* Case 2:

    Script content:

    ```txt
    printf("%v", hash("abc", "sha256"))
    
    ```

    Standard output:

    ```txt
    ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad
    ```

    
* Case 3:

    Script content:

    ```txt
    printf("%v", hash("abc", "sha512"))
    
    ```

    Standard output:

    ```txt
    ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f
    ```

    
* Case 4:

    Script content:

    ```txt
    printf("%v", hash("abc", "xx"))
    
    ```

    Standard output:

    ```txt
    
    ```

    

## `http_request` {#fn-http_request}

Function prototype: `fn http_request(method: str, url: str, body: bool|int|float|str|list|map = nil, headers: map = {}) -> map`

Function description: Used to send http request.

Function parameters:

- `method`: HTTP request method
- `url`: HTTP request url
- `body`: HTTP request body
- `headers`: HTTP request headers


Function returns:

- `map`: HTTP response


Function examples:

* Case 0:

    Script content:

    ```txt
    resp = http_request("GET", "http://test-domain/test")
    delete(resp["headers"], "Date")
    resp_str, ok = dump_json(resp, "    ")
    printf("%s", resp_str)
    ```

    Standard output:

    ```txt
    {
        "body": "{\"code\":200, \"message\":\"success\"}",
        "headers": {
            "Content-Length": "33",
            "Content-Type": "application/json"
        },
        "status_code": 200
    }
    ```

    
* Case 1:

    Script content:

    ```txt
    resp = http_request("GET", "http://localhost:80/test")
    
    # Usually, access to private IPs will be blocked,
    # you need to contact the administrator.
    
    resp_str, ok = dump_json(resp, "    ")
    printf("%s", resp_str)
    ```

    Standard output:

    ```txt
    {
        "error": "Get \"http://localhost:80/test\": resolved IP 127.0.0.1 is blocked"
    }
    		
    ```

    

## `len` {#fn-len}

Function prototype: `fn len(val: map|list|str) -> int`

Function description: Get the length of the value. If the value is a string, returns the length of the string. If the value is a list or map, returns the length of the list or map.

Function parameters:

- `val`: The value to get the length of.


Function returns:

- `int`: The length of the value.


Function examples:

* Case 0:

    Script content:

    ```txt
    printf("%v", len("abc"))
    
    ```

    Standard output:

    ```txt
    3
    ```

    
* Case 1:

    Script content:

    ```txt
    printf("%v", len([1, 2, 3]))
    
    ```

    Standard output:

    ```txt
    3
    ```

    
* Case 2:

    Script content:

    ```txt
    printf("%v", len({"a": 1, "b": 2, "c": 3}))
    
    ```

    Standard output:

    ```txt
    3
    ```

    

## `load_json` {#fn-load_json}

Function prototype: `fn load_json(val: str) -> (bool|int|float|str|list|map, bool)`

Function description: Unmarshal json string

Function parameters:

- `val`: JSON string.


Function returns:

- `bool|int|float|str|list|map`: Unmarshal result.
- `bool`: Unmarshal status.


Function examples:

* Case 0:

    Script content:

    ```txt
    jstr = '{"a": 1, "b": 2, "c": 3}'
    v, ok = load_json(jstr)
    if ok {
    	printf("%v", v["b"])
    }
    
    ```

    Standard output:

    ```txt
    2
    ```

    

## `lowercase` {#fn-lowercase}

Function prototype: `fn lowercase(val: str) -> str`

Function description: Converts a string to lowercase.

Function parameters:

- `val`: The string to convert.


Function returns:

- `str`: Returns the lowercase value.


Function examples:

* Case 0:

    Script content:

    ```txt
    printf("%s", lowercase("ABC"))
    
    ```

    Standard output:

    ```txt
    abc
    ```

    

## `match` {#fn-match}

Function prototype: `fn match(val: str, pattern: str, n: int = 1) -> (list, bool)`

Function description: Regular expression matching.

Function parameters:

- `val`: The string to match.
- `pattern`: Regular expression pattern.
- `n`: The number of matches to return. Defaults to 1, -1 for all matches.


Function returns:

- `list`: Returns the matched value.
- `bool`: Returns true if the regular expression matches.


Function examples:

* Case 0:

    Script content:

    ```txt
    text="abc def 123 abc def 123"
    v, ok = match(text, "(abc) (?:def) (?P<named_group>123)")
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    [["abc def 123","abc","123"]]
    ```

    
* Case 1:

    Script content:

    ```txt
    text="abc def 123 abc def 123"
    v, ok = match(text, "(abc) (?:def) (?P<named_group>123)", -1)
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    [["abc def 123","abc","123"],["abc def 123","abc","123"]]
    ```

    

## `parse_date` {#fn-parse_date}

Function prototype: `fn parse_date(date: str, timezone: str = "") -> (int, bool)`

Function description: Parses a date string to a nanoseconds timestamp, support multiple date formats. If the date string not include timezone and no timezone is provided, the local timezone is used.

Function parameters:

- `date`: The key to use for parsing.
- `timezone`: The timezone to use for parsing. If 


Function returns:

- `int`: The parsed timestamp in nanoseconds.
- `bool`: Whether the parsing was successful.


Function examples:

* Case 0:

    Script content:

    ```txt
    v, ok = parse_date("2021-12-2T11:55:43.123+0800")
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    1638417343123000000
    ```

    
* Case 1:

    Script content:

    ```txt
    v, ok = parse_date("2021-12-2T11:55:43.123", "+8")
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    1638417343123000000
    ```

    
* Case 2:

    Script content:

    ```txt
    v, ok = parse_date("2021-12-2T11:55:43.123", "Asia/Shanghai")
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    1638417343123000000
    ```

    

## `parse_duration` {#fn-parse_duration}

Function prototype: `fn parse_duration(s: str) -> (int, bool)`

Function description: Parses a golang duration string into a duration. A duration string is a sequence of possibly signed decimal numbers with optional fraction and unit suffixes for each number, such as `300ms`, `-1.5h` or `2h45m`. Valid units are `ns`, `us` (or `μs`), `ms`, `s`, `m`, `h`. 

Function parameters:

- `s`: The string to parse.


Function returns:

- `int`: The duration in nanoseconds.
- `bool`: Whether the duration is valid.


Function examples:

* Case 0:

    Script content:

    ```txt
    v, ok = parse_duration("1s")
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    1000000000
    ```

    
* Case 1:

    Script content:

    ```txt
    v, ok = parse_duration("100ns")
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    100
    ```

    

## `parse_int` {#fn-parse_int}

Function prototype: `fn parse_int(val: str, base: int) -> (int, bool)`

Function description: Parses a string into an integer.

Function parameters:

- `val`: The string to parse.
- `base`: The base to use for parsing. Must be between 2 and 36.


Function returns:

- `int`: The parsed integer.
- `bool`: Whether the parsing was successful.


Function examples:

* Case 0:

    Script content:

    ```txt
    v, ok = parse_int("123", 10)
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    123
    ```

    
* Case 1:

    Script content:

    ```txt
    v, ok = parse_int("123", 16)	
    if ok {
    	printf("%v", v)
    }
    
    ```

    Standard output:

    ```txt
    291
    ```

    

## `printf` {#fn-printf}

Function prototype: `fn printf(format: str, args: ...str|bool|int|float|list|map)`

Function description: Output formatted strings to the standard output device.

Function parameters:

- `format`: String format.
- `args`: Argument list, corresponding to the format specifiers in the format string.


Function examples:

* Case 0:

    Script content:

    ```txt
    printf("hello, %s", "world")
    
    ```

    Standard output:

    ```txt
    hello, world
    ```

    

## `replace` {#fn-replace}

Function prototype: `fn replace(input: str, pattern: str, replacement: str) -> (str, bool)`

Function description: Replaces text in a string.

Function parameters:

- `input`: The string to replace text in.
- `pattern`: Regular expression pattern.
- `replacement`: Replacement text to use.


Function returns:

- `str`: The string with text replaced.
- `bool`: True if the pattern was found and replaced, false otherwise.


Function examples:

* Case 0:

    Script content:

    ```txt
    v, ok = replace("abcdef", "bc", "123")
    printf("%s", v)
    
    ```

    Standard output:

    ```txt
    a123def
    ```

    
* Case 1:

    Script content:

    ```txt
    v, ok = replace("bonjour; 你好", "[\u4e00-\u9fa5]+", "hello")
    printf("%s", v)
    
    ```

    Standard output:

    ```txt
    bonjour; hello
    ```

    

## `sql_cover` {#fn-sql_cover}

Function prototype: `fn sql_cover(val: str) -> (str, bool)`

Function description: Obfuscate SQL query string.

Function parameters:

- `val`: The sql to obfuscate.


Function returns:

- `str`: The obfuscated sql.
- `bool`: The obfuscate status.


Function examples:

* Case 0:

    Script content:

    ```txt
    v, ok = sql_cover("select abc from def where x > 3 and y < 5")
    if ok {
    	printf("%s",v)
    }
    
    ```

    Standard output:

    ```txt
    select abc from def where x > ? and y < ?
    ```

    
* Case 1:

    Script content:

    ```txt
    v, ok = sql_cover("SELECT $func$INSERT INTO table VALUES ('a', 1, 2)$func$ FROM users")
    if ok {
    	printf("%s",v)
    }
    
    ```

    Standard output:

    ```txt
    SELECT ? FROM users
    ```

    
* Case 2:

    Script content:

    ```txt
    v, ok = sql_cover("SELECT ('/uffd')")
    if ok {
    	printf("%s",v)
    }
    
    ```

    Standard output:

    ```txt
    SELECT ( ? )
    ```

    

## `str_join` {#fn-str_join}

Function prototype: `fn str_join(li: list, sep: str) -> str`

Function description: String join.

Function parameters:

- `li`: List to be joined with separator. The elements type need to be string, if not, they will be ignored.
- `sep`: Separator to be used between elements.


Function returns:

- `str`: Joined string.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = str_join(["a", "b", "c"], "##")
    printf("%s", v)
    
    ```

    Standard output:

    ```txt
    a##b##c
    ```

    

## `strfmt` {#fn-strfmt}

Function prototype: `fn strfmt(format: str, args: ...bool|int|float|str|list|map) -> str`

Function description: 

Function parameters:

- `format`: String format.
- `args`: Parameters to replace placeholders.


Function returns:

- `str`: String.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = strfmt("abc %s def %d", "123", 456)
    printf("%s", v)
    
    ```

    Standard output:

    ```txt
    abc 123 def 456
    ```

    

## `time_now` {#fn-time_now}

Function prototype: `fn time_now(precision: str = "ns") -> int`

Function description: Get current timestamp with the specified precision.

Function parameters:

- `precision`: The precision of the timestamp. Supported values: `ns`, `us`, `ms`, `s`.


Function returns:

- `int`: Returns the current timestamp.


Function examples:

* Case 0:

    Script content:

    ```txt
    printf("%v", time_now("s"))
    
    ```

    Standard output:

    ```txt
    1745823860
    ```

    

## `trigger` {#fn-trigger}

Function prototype: `fn trigger(result: int|float|bool|str, status: str = "", dimension_tags: map = {}, related_data: map = {})`

Function description: Trigger a security event.

Function parameters:

- `result`: Event check result.
- `status`: Event status. One of: (`critical`, `high`, `medium`, `low`, `info`).
- `dimension_tags`: Dimension tags.
- `related_data`: Related data.


Function examples:

* Case 0:

    Script content:

    ```txt
    trigger(1, "critical", {"tag_abc":"1"}, {"a":"1", "a1":2.1})
    
    trigger(result=2, dimension_tags={"a":"1", "b":"2"}, related_data={"b": {}})
    
    trigger(false, related_data={"a":1, "b":2}, status="critical")
    
    trigger("hello", dimension_tags={}, related_data={"a":1, "b":[1]}, status="critical")
    
    ```

    Standard output:

    ```txt
    
    ```

    
    Trigger output:
    ```json
    [
        {
            "result": 1,
            "status": "critical",
            "dimension_tags": {
                "tag_abc": "1"
            },
            "related_data": {
                "a": "1",
                "a1": 2.1
            }
        },
        {
            "result": 2,
            "status": "",
            "dimension_tags": {
                "a": "1",
                "b": "2"
            },
            "related_data": {
                "b": {}
            }
        },
        {
            "result": false,
            "status": "critical",
            "dimension_tags": {},
            "related_data": {
                "a": 1,
                "b": 2
            }
        },
        {
            "result": "hello",
            "status": "critical",
            "dimension_tags": {},
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


Function examples:

* Case 0:

    Script content:

    ```txt
    printf("%s", trim(" abcdef "))
    
    ```

    Standard output:

    ```txt
    abcdef
    ```

    
* Case 1:

    Script content:

    ```txt
    printf("%s", trim("#-abcdef-#", "-#", 2))
    
    ```

    Standard output:

    ```txt
    #-abcdef
    ```

    
* Case 2:

    Script content:

    ```txt
    printf("%s", trim("#-abcdef-#", "-#", 1))
    
    ```

    Standard output:

    ```txt
    abcdef-#
    ```

    
* Case 3:

    Script content:

    ```txt
    printf("%s", trim("#-abcdef-#", side=0, cutset="-#"))
    
    ```

    Standard output:

    ```txt
    abcdef
    ```

    

## `uppercase` {#fn-uppercase}

Function prototype: `fn uppercase(val: str) -> str`

Function description: Converts a string to uppercase.

Function parameters:

- `val`: The string to convert.


Function returns:

- `str`: Returns the uppercase value.


Function examples:

* Case 0:

    Script content:

    ```txt
    printf("%s", uppercase("abc"))
    
    ```

    Standard output:

    ```txt
    ABC
    ```

    

## `url_decode` {#fn-url_decode}

Function prototype: `fn url_decode(val: str) -> (str, bool)`

Function description: Decodes a URL-encoded string.

Function parameters:

- `val`: The URL-encoded string to decode.


Function returns:

- `str`: The decoded string.
- `bool`: The decoding status.


Function examples:

* Case 0:

    Script content:

    ```txt
    v, ok = url_decode("https:%2F%2Fkubernetes.io%2Fdocs%2Freference%2Faccess-authn-authz%2Fbootstrap-tokens%2F")
    if ok {
    	printf("%s", v)
    }
    
    ```

    Standard output:

    ```txt
    https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/
    ```

    

## `url_parse` {#fn-url_parse}

Function prototype: `fn url_parse(url: str) -> (map, bool)`

Function description: Parses a URL and returns it as a map.

Function parameters:

- `url`: The URL to parse.


Function returns:

- `map`: Returns the parsed URL as a map.
- `bool`: Returns true if the URL is valid.


Function examples:

* Case 0:

    Script content:

    ```txt
    v, ok = url_parse("http://www.example.com:8080/path/to/file?query=abc")
    if ok {
    	v, ok = dump_json(v, "  ")
    	if ok {
    		printf("%v", v)
    	}
    }
    
    ```

    Standard output:

    ```txt
    {
      "host": "www.example.com:8080",
      "params": {
        "query": [
          "abc"
        ]
      },
      "path": "/path/to/file",
      "port": "8080",
      "scheme": "http"
    }
    ```

    

## `user_agent` {#fn-user_agent}

Function prototype: `fn user_agent(header: str) -> map`

Function description: Parses a User-Agent header.

Function parameters:

- `header`: The User-Agent header to parse.


Function returns:

- `map`: Returns the parsed User-Agent header as a map.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = user_agent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
    printf("%s", v)
    
    ```

    Standard output:

    ```txt
    {"browser":"Chrome","browserVer":"96.0.4664.110","engine":"AppleWebKit","engineVer":"537.36","isBot":false,"isMobile":false,"os":"Intel Mac OS X 10_15_7","ua":"Macintosh"}
    ```

    

## `valid_json` {#fn-valid_json}

Function prototype: `fn valid_json(val: str) -> bool`

Function description: Returns true if the value is a valid JSON.

Function parameters:

- `val`: The value to check.


Function returns:

- `bool`: Returns true if the value is a valid JSON.


Function examples:

* Case 0:

    Script content:

    ```txt
    ok = valid_json("{\"a\": 1, \"b\": 2}")
    if ok {
    	printf("true")
    }
    
    ```

    Standard output:

    ```txt
    true
    ```

    
* Case 1:

    Script content:

    ```txt
    ok = valid_json("1.1")
    if ok {
    	printf("true")
    }
    
    ```

    Standard output:

    ```txt
    true
    ```

    
* Case 2:

    Script content:

    ```txt
    ok = valid_json("str_abc_def")
    if !ok {
    	printf("false")
    }
    
    ```

    Standard output:

    ```txt
    false
    ```

    

## `value_type` {#fn-value_type}

Function prototype: `fn value_type(val: str) -> str`

Function description: Returns the type of the value.

Function parameters:

- `val`: The value to get the type of.


Function returns:

- `str`: Returns the type of the value. One of (`bool`, `int`, `float`, `str`, `list`, `map`, `nil`). If the value and the type is nil, returns `nil`.


Function examples:

* Case 0:

    Script content:

    ```txt
    v = value_type(1)
    printf("%s", v)
    
    ```

    Standard output:

    ```txt
    int
    ```

    
* Case 1:

    Script content:

    ```txt
    printf("%s", value_type("abc"))
    
    ```

    Standard output:

    ```txt
    str
    ```

    
* Case 2:

    Script content:

    ```txt
    printf("%s", value_type(true))
    
    ```

    Standard output:

    ```txt
    bool
    ```

    

## `xml_query` {#fn-xml_query}

Function prototype: `fn xml_query(input: str, xpath: str) -> (str, bool)`

Function description: Returns the value of an XML field.

Function parameters:

- `input`: The XML input to get the value of.
- `xpath`: The XPath expression to get the value of.


Function returns:

- `str`: Returns the value of the XML field.
- `bool`: Returns true if the field exists, false otherwise.


Function examples:

* Case 0:

    Script content:

    ```txt
    xml_data='''
    <OrderEvent actionCode = "5">
     <OrderNumber>ORD12345</OrderNumber>
     <VendorNumber>V11111</VendorNumber>
     </OrderEvent>
    '''
    v, ok = xml_query(xml_data, "/OrderEvent/OrderNumber/text()")
    if ok {
    	printf("%s", v)
    }
    
    ```

    Standard output:

    ```txt
    ORD12345
    ```

    
* Case 1:

    Script content:

    ```txt
    xml_data='''
    <OrderEvent actionCode = "5">
     <OrderNumber>ORD12345</OrderNumber>
     <VendorNumber>V11111</VendorNumber>
     </OrderEvent>
    '''
    v, ok = xml_query(xml_data, "/OrderEvent/@actionCode")
    if ok {
    	printf("%s", v)
    }
    
    ```

    Standard output:

    ```txt
    5
    ```

    

