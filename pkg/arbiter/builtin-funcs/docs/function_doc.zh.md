# Arbiter 内置函数

## `append` {#fn-append}

函数原型： `fn append(li: list, v: ...bool|int|float|str|list|map) -> list`

函数描述： Appends a value to a list.

函数参数：

- `li`: The list to append to.
- `v`: The value to append.


函数返回值：

- `list`: The list with the appended value.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = [1, 2, 3]
    v = append(v, 4)
    printf("%v", v)
    
    ```

    标准输出：

    ```txt
    [1,2,3,4]
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    v = [1, 2, 3]
    v = append(v, "a", 1.1)
    printf("%v", v)
    
    ```

    标准输出：

    ```txt
    [1,2,3,"a",1.1]
    ```

    

## `b64dec` {#fn-b64dec}

函数原型： `fn b64dec(data: str) -> (str, bool)`

函数描述： Base64 decoding.

函数参数：

- `data`: Data that needs to be base64 decoded.


函数返回值：

- `str`: The decoded string.
- `bool`: Whether decoding is successful.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = "aGVsbG8sIHdvcmxk"
    v, ok = b64dec(v)
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    hello, world
    ```

    

## `b64enc` {#fn-b64enc}

函数原型： `fn b64enc(data: str) -> (str, bool)`

函数描述： Base64 encoding.

函数参数：

- `data`: Data that needs to be base64 encoded.


函数返回值：

- `str`: The encoded string.
- `bool`: Whether encoding is successful.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = "hello, world"
    v = b64enc(v)
    printf("%v", v)
    
    ```

    标准输出：

    ```txt
    aGVsbG8sIHdvcmxk
    ```

    

## `call_func` {#fn-call_func}

函数原型： `fn call_func(name: str, kwargs: map = {}) -> map`

函数描述： Calling remote functions on the Function platform

函数参数：

- `name`: Remote function name.
- `kwargs`: Parameter map, corresponding to **kwargs in Python.


函数返回值：

- `map`: 


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    result = call_func("echo", {
        "arg_1": "1",
    	"arg_2": [12,3],
    	"key3": "true"
    })
    printf("%v", result)
    ```

    标准输出：

    ```txt
    {"error":"","message":"","ok":true,"reason":"","result":{"arg_1":"1","arg_2":[12,3],"kwargs":{"key3":"true"}}}
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    result = call_func("echo", {"arg_2": [12,3]})
    printf("%v", result)
    ```

    标准输出：

    ```txt
    {"detail":{"exeception":"TypeError(\"echo() missing 1 required positional argument: 'arg_1'\")"},"error":"call function failed, api status code 599","message":"Func task failed","ok":false,"reason":"EFuncFailed"}
    ```

    

## `cast` {#fn-cast}

函数原型： `fn cast(val: bool|int|float|str, typ: str) -> bool|int|float|str`

函数描述： Convert the value to the target type.

函数参数：

- `val`: The value of the type to be converted.
- `typ`: Target type. One of (`bool`, `int`, `float`, `str`).


函数返回值：

- `bool|int|float|str`: The value after the conversion.


函数示例：

* 示例 0:

    脚本内容：

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

    标准输出：

    ```txt
    2.1; 2; -1; 2; false; 1.1; 1; true
    
    ```

    

## `cidr` {#fn-cidr}

函数原型： `fn cidr(ip: str, mask: str) -> bool`

函数描述： Check the IP whether in CIDR block

函数参数：

- `ip`: The ip address
- `mask`: The CIDR mask


函数返回值：

- `bool`: Whether the IP is in CIDR block


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    ip = "192.0.2.233"
    if cidr(ip, "192.0.2.1/24") {
    	printf("%s", ip)
    }
    ```

    标准输出：

    ```txt
    192.0.2.233
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    ip = "192.0.2.233"
    if cidr(mask="192.0.1.1/24", ip=ip) {
    	printf("%s", ip)
    }
    ```

    标准输出：

    ```txt
    
    ```

    

## `delete` {#fn-delete}

函数原型： `fn delete(m: map, key: str)`

函数描述： Delete key from the map.

函数参数：

- `m`: The map for deleting key
- `key`: Key need delete from map.


函数示例：

* 示例 0:

    脚本内容：

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

    标准输出：

    ```txt
    result group 1: {"b":2}; [{"c":1.1}]
    result group 2: {}; {}
    
    ```

    

## `dql` {#fn-dql}

函数原型： `fn dql(query: str, qtype: str = "dql", limit: int = 10000, offset: int = 0, slimit: int = 0, time_range: list = []) -> map`

函数描述： Query data using dql or promql.

函数参数：

- `query`: DQL or PromQL query statements.
- `qtype`: Query language, One of `dql` or `promql`, default is `dql`.
- `limit`: Query limit.
- `offset`: Query offset.
- `slimit`: Query slimit.
- `time_range`: Query timestamp range, the default value can be modified externally by the script caller.


函数返回值：

- `map`: Query response.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = dql("M::cpu limit 2 slimit 2")
    v, ok = dump_json(v, "    ")
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

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

函数原型： `fn dql_series_get(series: map, name: str) -> list`

函数描述： get series data

函数参数：

- `series`: dql query result
- `name`: column or tag name


函数返回值：

- `list`: specified column or tag value for the series


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = dql("M::cpu limit 2 slimit 2") 
    
    hostLi = dql_series_get(v, "host")
    time_li = dql_series_get(v, "time")
    
    printf("%v", {"host": hostLi, "time": time_li})
    
    ```

    标准输出：

    ```txt
    {"host":[["172.16.241.111","172.16.241.111"],["172.16.242.112","172.16.242.112"]],"time":[[1744866108991,1744866103991],[1744866107975,1744866102975]]}
    ```

    

## `dql_timerange_get` {#fn-dql_timerange_get}

函数原型： `fn dql_timerange_get() -> list`

函数描述： Get the time range of the DQL query, which is passed in by the script caller or defaults to the last 15 minutes.

函数返回值：

- `list`: The time range. For example, `[1744214400000, 1744218000000]`, the timestamp precision is milliseconds


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    val = dql_timerange_get()
    printf("%v", val)
    ```

    标准输出：

    ```txt
    [1672531500000,1672532100000]
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    val = dql_timerange_get()
    printf("%v", val)
    ```

    标准输出：

    ```txt
    [1672531200000,1672532100000]
    ```

    

## `dump_json` {#fn-dump_json}

函数原型： `fn dump_json(v: str, indent: str = "") -> (str, bool)`

函数描述： Returns the JSON encoding of v.

函数参数：

- `v`: Object to encode.
- `indent`: Indentation prefix.


函数返回值：

- `str`: JSON encoding of v.
- `bool`: Whether decoding is successful.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = {"a": 1, "b": 2.1}
    v, ok = dump_json(v)
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    {"a":1,"b":2.1}
    
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    v = {"a": 1, "b": 2.1}
    v, ok = dump_json(v, "  ")
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    {
      "a": 1,
      "b": 2.1
    }
    
    ```

    

## `exit` {#fn-exit}

函数原型： `fn exit()`

函数描述： Exit the program

函数示例：

* 示例 0:

    脚本内容：

    ```txt
    printf("1\n")
    printf("2\n")
    exit()
    printf("3\n")
    	
    ```

    标准输出：

    ```txt
    1
    2
    
    ```

    

## `format_int` {#fn-format_int}

函数原型： `fn format_int(val: int, base: int) -> str`

函数描述： Formats an integer into a string.

函数参数：

- `val`: The integer to format.
- `base`: The base to use for formatting. Must be between 2 and 36.


函数返回值：

- `str`: The formatted string.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = format_int(16, 16)
    printf("%s", v)
    
    ```

    标准输出：

    ```txt
    10
    ```

    

## `geoip` {#fn-geoip}

函数原型： `fn geoip(ip: str) -> map`

函数描述： GeoIP

函数参数：

- `ip`: IP address.


函数返回值：

- `map`: IP geographical information.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = geoip("127.0.0.1")
    printf("%v", v)
    
    ```

    标准输出：

    ```txt
    {"city":"","country":"","isp":"unknown","province":""}
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    ip_addr = "114.114.114.114"
    v, ok = dump_json(geoip(ip_addr), "    ");
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
     {
        "city": "Ji'an",
        "country": "CN",
        "isp": "chinanet",
        "province": "Jiangxi"
    }
    ```

    

## `gjson` {#fn-gjson}

函数原型： `fn gjson(input: str, json_path: str) -> (bool|int|float|str|list|map, bool)`

函数描述： GJSON provides a fast and easy way to get values from a JSON document.

函数参数：

- `input`: JSON format string to parse.
- `json_path`: JSON path.


函数返回值：

- `bool|int|float|str|list|map`: Parsed result.
- `bool`: Parsed status.


函数示例：

* 示例 0:

    脚本内容：

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

    标准输出：

    ```txt
    37
    ```

    
* 示例 1:

    脚本内容：

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

    标准输出：

    ```txt
    {"first": "Tom", "last": "Anderson"}
    ```

    
* 示例 2:

    脚本内容：

    ```txt
    v='''[
        {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
        {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
        {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
    ]'''
    net, ok = gjson(v, "0.nets.2")
    printf("%v", net)
    
    ```

    标准输出：

    ```txt
    tw
    ```

    

## `grok` {#fn-grok}

函数原型： `fn grok(input: str, pattern: str, extra_patterns: map = {}, trim_space: bool = true) -> (map, bool)`

函数描述： Extracts data from a string using a Grok pattern. Grok is based on regular expression syntax, and using regular (named) capture groups in a pattern is equivalent to using a pattern in a pattern. A valid regular expression is also a valid Grok pattern.

函数参数：

- `input`: The input string used to extract data.
- `pattern`: The pattern used to extract data.
- `extra_patterns`: Additional patterns for parsing patterns.
- `trim_space`: Whether to trim leading and trailing spaces from the parsed value.


函数返回值：

- `map`: The parsed result.
- `bool`: Whether the parsing was successful.


函数示例：

* 示例 0:

    脚本内容：

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

    标准输出：

    ```txt
    {
      "log_code_pos": "io/io.go:458",
      "log_cost": 6.87021,
      "log_level": "DEBUG",
      "log_time": "2021-01-11T17:43:51.887+0800"
    }
    ```

    

## `hash` {#fn-hash}

函数原型： `fn hash(text: str, method: str) -> str`

函数描述： 

函数参数：

- `text`: The string used to calculate the hash.
- `method`: Hash Algorithms, allowing values including `md5`, `sha1`, `sha256`, `sha512`.


函数返回值：

- `str`: The hash value.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    printf("%v", hash("abc", "md5"))
    
    ```

    标准输出：

    ```txt
    900150983cd24fb0d6963f7d28e17f72
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    printf("%v", hash("abc", "sha1"))
    
    ```

    标准输出：

    ```txt
    a9993e364706816aba3e25717850c26c9cd0d89d
    ```

    
* 示例 2:

    脚本内容：

    ```txt
    printf("%v", hash("abc", "sha256"))
    
    ```

    标准输出：

    ```txt
    ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad
    ```

    
* 示例 3:

    脚本内容：

    ```txt
    printf("%v", hash("abc", "sha512"))
    
    ```

    标准输出：

    ```txt
    ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f
    ```

    
* 示例 4:

    脚本内容：

    ```txt
    printf("%v", hash("abc", "xx"))
    
    ```

    标准输出：

    ```txt
    
    ```

    

## `http_request` {#fn-http_request}

函数原型： `fn http_request(method: str, url: str, body: bool|int|float|str|list|map = nil, headers: map = {}) -> map`

函数描述： Used to send http request.

函数参数：

- `method`: HTTP request method
- `url`: HTTP request url
- `body`: HTTP request body
- `headers`: HTTP request headers


函数返回值：

- `map`: HTTP response


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    resp = http_request("GET", "http://test-domain/test")
    delete(resp["headers"], "Date")
    resp_str, ok = dump_json(resp, "    ")
    printf("%s", resp_str)
    ```

    标准输出：

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

    
* 示例 1:

    脚本内容：

    ```txt
    resp = http_request("GET", "http://localhost:80/test")
    
    # Usually, access to private IPs will be blocked,
    # you need to contact the administrator.
    
    resp_str, ok = dump_json(resp, "    ")
    printf("%s", resp_str)
    ```

    标准输出：

    ```txt
    {
        "error": "Get \"http://localhost:80/test\": resolved IP 127.0.0.1 is blocked"
    }
    		
    ```

    

## `len` {#fn-len}

函数原型： `fn len(val: map|list|str) -> int`

函数描述： Get the length of the value. If the value is a string, returns the length of the string. If the value is a list or map, returns the length of the list or map.

函数参数：

- `val`: The value to get the length of.


函数返回值：

- `int`: The length of the value.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    printf("%v", len("abc"))
    
    ```

    标准输出：

    ```txt
    3
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    printf("%v", len([1, 2, 3]))
    
    ```

    标准输出：

    ```txt
    3
    ```

    
* 示例 2:

    脚本内容：

    ```txt
    printf("%v", len({"a": 1, "b": 2, "c": 3}))
    
    ```

    标准输出：

    ```txt
    3
    ```

    

## `load_json` {#fn-load_json}

函数原型： `fn load_json(val: str) -> (bool|int|float|str|list|map, bool)`

函数描述： Unmarshal json string

函数参数：

- `val`: JSON string.


函数返回值：

- `bool|int|float|str|list|map`: Unmarshal result.
- `bool`: Unmarshal status.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    jstr = '{"a": 1, "b": 2, "c": 3}'
    v, ok = load_json(jstr)
    if ok {
    	printf("%v", v["b"])
    }
    
    ```

    标准输出：

    ```txt
    2
    ```

    

## `lowercase` {#fn-lowercase}

函数原型： `fn lowercase(val: str) -> str`

函数描述： Converts a string to lowercase.

函数参数：

- `val`: The string to convert.


函数返回值：

- `str`: Returns the lowercase value.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    printf("%s", lowercase("ABC"))
    
    ```

    标准输出：

    ```txt
    abc
    ```

    

## `match` {#fn-match}

函数原型： `fn match(val: str, pattern: str, n: int = 1) -> (list, bool)`

函数描述： Regular expression matching.

函数参数：

- `val`: The string to match.
- `pattern`: Regular expression pattern.
- `n`: The number of matches to return. Defaults to 1, -1 for all matches.


函数返回值：

- `list`: Returns the matched value.
- `bool`: Returns true if the regular expression matches.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    text="abc def 123 abc def 123"
    v, ok = match(text, "(abc) (?:def) (?P<named_group>123)")
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    [["abc def 123","abc","123"]]
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    text="abc def 123 abc def 123"
    v, ok = match(text, "(abc) (?:def) (?P<named_group>123)", -1)
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    [["abc def 123","abc","123"],["abc def 123","abc","123"]]
    ```

    

## `parse_date` {#fn-parse_date}

函数原型： `fn parse_date(date: str, timezone: str = "") -> (int, bool)`

函数描述： Parses a date string to a nanoseconds timestamp, support multiple date formats. If the date string not include timezone and no timezone is provided, the local timezone is used.

函数参数：

- `date`: The key to use for parsing.
- `timezone`: The timezone to use for parsing. If 


函数返回值：

- `int`: The parsed timestamp in nanoseconds.
- `bool`: Whether the parsing was successful.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v, ok = parse_date("2021-12-2T11:55:43.123+0800")
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    1638417343123000000
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    v, ok = parse_date("2021-12-2T11:55:43.123", "+8")
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    1638417343123000000
    ```

    
* 示例 2:

    脚本内容：

    ```txt
    v, ok = parse_date("2021-12-2T11:55:43.123", "Asia/Shanghai")
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    1638417343123000000
    ```

    

## `parse_duration` {#fn-parse_duration}

函数原型： `fn parse_duration(s: str) -> (int, bool)`

函数描述： Parses a golang duration string into a duration. A duration string is a sequence of possibly signed decimal numbers with optional fraction and unit suffixes for each number, such as `300ms`, `-1.5h` or `2h45m`. Valid units are `ns`, `us` (or `μs`), `ms`, `s`, `m`, `h`. 

函数参数：

- `s`: The string to parse.


函数返回值：

- `int`: The duration in nanoseconds.
- `bool`: Whether the duration is valid.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v, ok = parse_duration("1s")
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    1000000000
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    v, ok = parse_duration("100ns")
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    100
    ```

    

## `parse_int` {#fn-parse_int}

函数原型： `fn parse_int(val: str, base: int) -> (int, bool)`

函数描述： Parses a string into an integer.

函数参数：

- `val`: The string to parse.
- `base`: The base to use for parsing. Must be between 2 and 36.


函数返回值：

- `int`: The parsed integer.
- `bool`: Whether the parsing was successful.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v, ok = parse_int("123", 10)
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    123
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    v, ok = parse_int("123", 16)	
    if ok {
    	printf("%v", v)
    }
    
    ```

    标准输出：

    ```txt
    291
    ```

    

## `printf` {#fn-printf}

函数原型： `fn printf(format: str, args: ...str|bool|int|float|list|map)`

函数描述： Output formatted strings to the standard output device.

函数参数：

- `format`: String format.
- `args`: Argument list, corresponding to the format specifiers in the format string.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    printf("hello, %s", "world")
    
    ```

    标准输出：

    ```txt
    hello, world
    ```

    

## `replace` {#fn-replace}

函数原型： `fn replace(input: str, pattern: str, replacement: str) -> (str, bool)`

函数描述： Replaces text in a string.

函数参数：

- `input`: The string to replace text in.
- `pattern`: Regular expression pattern.
- `replacement`: Replacement text to use.


函数返回值：

- `str`: The string with text replaced.
- `bool`: True if the pattern was found and replaced, false otherwise.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v, ok = replace("abcdef", "bc", "123")
    printf("%s", v)
    
    ```

    标准输出：

    ```txt
    a123def
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    v, ok = replace("bonjour; 你好", "[\u4e00-\u9fa5]+", "hello")
    printf("%s", v)
    
    ```

    标准输出：

    ```txt
    bonjour; hello
    ```

    

## `setopt` {#fn-setopt}

函数原型： `fn setopt(trigger_keepalive: int)`

函数描述： Set option.

函数参数：

- `trigger_keepalive`: Set the trigger keepalive time (seconds).


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    setopt(60 * 15)
    ```

    标准输出：

    ```txt
    
    ```

    

## `sql_cover` {#fn-sql_cover}

函数原型： `fn sql_cover(val: str) -> (str, bool)`

函数描述： Obfuscate SQL query string.

函数参数：

- `val`: The sql to obfuscate.


函数返回值：

- `str`: The obfuscated sql.
- `bool`: The obfuscate status.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v, ok = sql_cover("select abc from def where x > 3 and y < 5")
    if ok {
    	printf("%s",v)
    }
    
    ```

    标准输出：

    ```txt
    select abc from def where x > ? and y < ?
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    v, ok = sql_cover("SELECT $func$INSERT INTO table VALUES ('a', 1, 2)$func$ FROM users")
    if ok {
    	printf("%s",v)
    }
    
    ```

    标准输出：

    ```txt
    SELECT ? FROM users
    ```

    
* 示例 2:

    脚本内容：

    ```txt
    v, ok = sql_cover("SELECT ('/uffd')")
    if ok {
    	printf("%s",v)
    }
    
    ```

    标准输出：

    ```txt
    SELECT ( ? )
    ```

    

## `str_join` {#fn-str_join}

函数原型： `fn str_join(li: list, sep: str) -> str`

函数描述： String join.

函数参数：

- `li`: List to be joined with separator. The elements type need to be string, if not, they will be ignored.
- `sep`: Separator to be used between elements.


函数返回值：

- `str`: Joined string.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = str_join(["a", "b", "c"], "##")
    printf("%s", v)
    
    ```

    标准输出：

    ```txt
    a##b##c
    ```

    

## `strfmt` {#fn-strfmt}

函数原型： `fn strfmt(format: str, args: ...bool|int|float|str|list|map) -> str`

函数描述： 

函数参数：

- `format`: String format.
- `args`: Parameters to replace placeholders.


函数返回值：

- `str`: String.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = strfmt("abc %s def %d", "123", 456)
    printf("%s", v)
    
    ```

    标准输出：

    ```txt
    abc 123 def 456
    ```

    

## `time_now` {#fn-time_now}

函数原型： `fn time_now(precision: str = "ns") -> int`

函数描述： Get current timestamp with the specified precision.

函数参数：

- `precision`: The precision of the timestamp. Supported values: `ns`, `us`, `ms`, `s`.


函数返回值：

- `int`: Returns the current timestamp.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    printf("%v", time_now("s"))
    
    ```

    标准输出：

    ```txt
    1745823860
    ```

    

## `trigger` {#fn-trigger}

函数原型： `fn trigger(result: int|float|bool|str, status: str = "", dimension_tags: map = {}, related_data: map = {})`

函数描述： Trigger a security event.

函数参数：

- `result`: Event check result.
- `status`: Event status. One of: (`critical`, `high`, `medium`, `low`, `info`).
- `dimension_tags`: Dimension tags.
- `related_data`: Related data.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    trigger(1, "critical", {"tag_abc":"1"}, {"a":"1", "a1":2.1})
    
    trigger(result=2, dimension_tags={"a":"1", "b":"2"}, related_data={"b": {}})
    
    trigger(false, related_data={"a":1, "b":2}, status="critical")
    
    trigger("hello", dimension_tags={}, related_data={"a":1, "b":[1]}, status="critical")
    
    ```

    标准输出：

    ```txt
    
    ```

    
    触发器输出：
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

函数原型： `fn trim(val: str, cutset: str = "", side: int = 0) -> str`

函数描述： Removes leading and trailing whitespace from a string.

函数参数：

- `val`: The string to trim.
- `cutset`: Characters to remove from the beginning and end of the string. If not specified, whitespace is removed.
- `side`: The side to trim from. If value is 0, trim from both sides. If value is 1, trim from the left side. If value is 2, trim from the right side.


函数返回值：

- `str`: The trimmed string.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    printf("%s", trim(" abcdef "))
    
    ```

    标准输出：

    ```txt
    abcdef
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    printf("%s", trim("#-abcdef-#", "-#", 2))
    
    ```

    标准输出：

    ```txt
    #-abcdef
    ```

    
* 示例 2:

    脚本内容：

    ```txt
    printf("%s", trim("#-abcdef-#", "-#", 1))
    
    ```

    标准输出：

    ```txt
    abcdef-#
    ```

    
* 示例 3:

    脚本内容：

    ```txt
    printf("%s", trim("#-abcdef-#", side=0, cutset="-#"))
    
    ```

    标准输出：

    ```txt
    abcdef
    ```

    

## `uppercase` {#fn-uppercase}

函数原型： `fn uppercase(val: str) -> str`

函数描述： Converts a string to uppercase.

函数参数：

- `val`: The string to convert.


函数返回值：

- `str`: Returns the uppercase value.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    printf("%s", uppercase("abc"))
    
    ```

    标准输出：

    ```txt
    ABC
    ```

    

## `url_decode` {#fn-url_decode}

函数原型： `fn url_decode(val: str) -> (str, bool)`

函数描述： Decodes a URL-encoded string.

函数参数：

- `val`: The URL-encoded string to decode.


函数返回值：

- `str`: The decoded string.
- `bool`: The decoding status.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v, ok = url_decode("https:%2F%2Fkubernetes.io%2Fdocs%2Freference%2Faccess-authn-authz%2Fbootstrap-tokens%2F")
    if ok {
    	printf("%s", v)
    }
    
    ```

    标准输出：

    ```txt
    https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/
    ```

    

## `url_parse` {#fn-url_parse}

函数原型： `fn url_parse(url: str) -> (map, bool)`

函数描述： Parses a URL and returns it as a map.

函数参数：

- `url`: The URL to parse.


函数返回值：

- `map`: Returns the parsed URL as a map.
- `bool`: Returns true if the URL is valid.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v, ok = url_parse("http://www.example.com:8080/path/to/file?query=abc")
    if ok {
    	v, ok = dump_json(v, "  ")
    	if ok {
    		printf("%v", v)
    	}
    }
    
    ```

    标准输出：

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

函数原型： `fn user_agent(header: str) -> map`

函数描述： Parses a User-Agent header.

函数参数：

- `header`: The User-Agent header to parse.


函数返回值：

- `map`: Returns the parsed User-Agent header as a map.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = user_agent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
    printf("%s", v)
    
    ```

    标准输出：

    ```txt
    {"browser":"Chrome","browserVer":"96.0.4664.110","engine":"AppleWebKit","engineVer":"537.36","isBot":false,"isMobile":false,"os":"Intel Mac OS X 10_15_7","ua":"Macintosh"}
    ```

    

## `valid_json` {#fn-valid_json}

函数原型： `fn valid_json(val: str) -> bool`

函数描述： Returns true if the value is a valid JSON.

函数参数：

- `val`: The value to check.


函数返回值：

- `bool`: Returns true if the value is a valid JSON.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    ok = valid_json("{\"a\": 1, \"b\": 2}")
    if ok {
    	printf("true")
    }
    
    ```

    标准输出：

    ```txt
    true
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    ok = valid_json("1.1")
    if ok {
    	printf("true")
    }
    
    ```

    标准输出：

    ```txt
    true
    ```

    
* 示例 2:

    脚本内容：

    ```txt
    ok = valid_json("str_abc_def")
    if !ok {
    	printf("false")
    }
    
    ```

    标准输出：

    ```txt
    false
    ```

    

## `value_type` {#fn-value_type}

函数原型： `fn value_type(val: str) -> str`

函数描述： Returns the type of the value.

函数参数：

- `val`: The value to get the type of.


函数返回值：

- `str`: Returns the type of the value. One of (`bool`, `int`, `float`, `str`, `list`, `map`, `nil`). If the value and the type is nil, returns `nil`.


函数示例：

* 示例 0:

    脚本内容：

    ```txt
    v = value_type(1)
    printf("%s", v)
    
    ```

    标准输出：

    ```txt
    int
    ```

    
* 示例 1:

    脚本内容：

    ```txt
    printf("%s", value_type("abc"))
    
    ```

    标准输出：

    ```txt
    str
    ```

    
* 示例 2:

    脚本内容：

    ```txt
    printf("%s", value_type(true))
    
    ```

    标准输出：

    ```txt
    bool
    ```

    

## `xml_query` {#fn-xml_query}

函数原型： `fn xml_query(input: str, xpath: str) -> (str, bool)`

函数描述： Returns the value of an XML field.

函数参数：

- `input`: The XML input to get the value of.
- `xpath`: The XPath expression to get the value of.


函数返回值：

- `str`: Returns the value of the XML field.
- `bool`: Returns true if the field exists, false otherwise.


函数示例：

* 示例 0:

    脚本内容：

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

    标准输出：

    ```txt
    ORD12345
    ```

    
* 示例 1:

    脚本内容：

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

    标准输出：

    ```txt
    5
    ```

    

