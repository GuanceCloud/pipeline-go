### `pt_kvs_get()` {#fn_pt_kvs_get}

函数原型：`fn pt_kvs_get(name: str) -> any`

函数说明：返回 Point 中指定 key 的值

说明：

- 当 Point 中的字段是数组或字典时，返回值会保留为 `list` / `map`，可继续参与 `append()`、`len()`、`in` 等脚本运算
- tag 始终按字符串返回
- 如果该值已在更早的处理链路中被存成 JSON 字符串，则返回的仍然是字符串

函数参数：

- `name`: Key 名

示例：

```python
host = pt_kvs_get("host")

nums = pt_kvs_get("nums")
nums = append(nums, 4)
pt_kvs_set("nums", nums)
```
