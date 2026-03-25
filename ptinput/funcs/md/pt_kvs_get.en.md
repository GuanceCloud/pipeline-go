### `pt_kvs_get()` {#fn_pt_kvs_get}

Function prototype: `fn pt_kvs_get(name: str) -> any`

Function description: Return the value of the specified key in Point

Notes:

- If the field stored in Point is an array or map, the return value is preserved as `list` / `map`, so it can be used by script operations such as `append()`, `len()`, and `in`
- Tags are always returned as strings
- If the value was already stored as a JSON string earlier in the pipeline, the returned value is still a string

Function parameters:

- `name`: Key name

Example:

```python
host = pt_kvs_get("host")

nums = pt_kvs_get("nums")
nums = append(nums, 4)
pt_kvs_set("nums", nums)
```
