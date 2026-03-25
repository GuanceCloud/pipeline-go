### `pt_kvs_set()` {#fn_pt_kvs_set}

Function prototype: `fn pt_kvs_set(name: str, value: any, as_tag: bool = false) -> bool`

Function description: Add a key to a Point or modify the value of a key in a Point

Notes:

- When setting a field, `list` values are preserved as native Point array fields whenever possible
- When setting a field, `map` values are preserved as native Point dict fields whenever possible
- When setting a tag, the value is converted to a string

Function parameters:

- `name`: The name of the field or label to be added or modified
- `value`: The value of a field or label
- `as_tag`: Set as tag or not

Example:

```python
kvs = {
    "a": 1,
    "b": 2
}

for k in kvs {
    pt_kvs_set(k, kvs[k])
}

nums = pt_kvs_get("nums")
nums = append(nums, 4)
pt_kvs_set("nums", nums)
```
