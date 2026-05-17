# {模块名称} - API 文档

**基础路径**: `/api/v1/{module}`

---

## 接口列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /resources | 获取资源列表 |
| POST | /resources | 创建资源 |
| GET | /resources/:id | 获取单个资源 |
| PUT | /resources/:id | 更新资源 |
| DELETE | /resources/:id | 删除资源 |

---

## GET /resources

获取资源列表

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | number | 否 | 页码，默认 1 |
| limit | number | 否 | 每页数量，默认 20 |

### 响应

```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

---

## POST /resources

创建资源

### 请求体

```json
{
  "field1": "value1",
  "field2": "value2"
}
```

### 响应

```json
{
  "id": "xxx",
  "field1": "value1",
  "field2": "value2",
  "createdAt": "2026-05-17T00:00:00Z"
}
```

---

## 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |
