# MyAPI Demo

一个基于 Go + Gin + GORM 的简单模型管理 API Demo，支持模型的增删改查（CRUD），并支持部分字段更新。

## 功能简介

- 支持 AI 模型的注册、查询、更新、删除
- 更新接口支持只更新部分字段，未提供的字段不会被清空
- RESTful API 设计，返回统一 JSON 响应格式
- 自动创建数据库表结构
- 适合学习、二次开发和接口测试

## 快速开始

### 1. 启动服务

```bash
go build -o bin/myapi .
./bin/myapi -c etc/config.yaml
```

### 2. API 说明

#### 创建模型
```bash
curl -X POST http://localhost:3000/api/v1/models/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "GPT-4",
    "endpoint": "https://api.openai.com/v1/chat/completions",
    "api_key": "sk-xxxxxxxxxxxxxxxxxxxxxxxx",
    "timeout": 30,
    "type": "chat",
    "dimension": 1536
  }'
```

#### 获取模型列表
```bash
curl -X GET http://localhost:3000/api/v1/models/get
```

#### 获取单个模型
```bash
curl -X GET http://localhost:3000/api/v1/models/<model_id>
```

#### 更新模型（支持部分字段更新）
- 只更新 name 字段：
```bash
curl -X PUT http://localhost:3000/api/v1/models/<model_id> \
  -H "Content-Type: application/json" \
  -d '{"name": "新名称"}'
```
- 只更新 timeout 字段：
```bash
curl -X PUT http://localhost:3000/api/v1/models/<model_id> \
  -H "Content-Type: application/json" \
  -d '{"timeout": 60}'
```
- 同时更新多个字段：
```bash
curl -X PUT http://localhost:3000/api/v1/models/<model_id> \
  -H "Content-Type: application/json" \
  -d '{"name": "新名称", "timeout": 60, "api_key": "sk-new-key"}'
```

#### 删除模型
```bash
curl -X DELETE http://localhost:3000/api/v1/models/<model_id>
```

## 响应格式示例

```json
{
  "status": 200,
  "data": {
    "model_id": "...",
    "name": "GPT-4",
    "endpoint": "...",
    "api_key": "...",
    "timeout": 30,
    "type": "chat",
    "dimension": 1536,
    "created_at": "2025-06-28T18:30:00+08:00",
    "updated_at": "2025-06-28T18:30:00+08:00"
  },
  "msg": ""
}
```

## 错误处理示例

- 名称重复：
```json
{
  "status": 409,
  "data": null,
  "msg": "模型名称已存在"
}
```
- 参数错误：
```json
{
  "status": 400,
  "data": null,
  "msg": "参数错误: ..."
}
```
- 未找到模型：
```json
{
  "status": 404,
  "data": null,
  "msg": "模型不存在"
}
```

## 数据库表结构（自动迁移）
| 字段名      | 类型         | 说明         |
| ----------- | ------------ | ------------ |
| model_id    | varchar(64)  | 主键，UUID   |
| name        | varchar(255) | 唯一，必填   |
| endpoint    | varchar(255) | 必填         |
| api_key     | varchar(255) | 必填         |
| timeout     | int          | 必填         |
| type        | varchar(255) | 必填         |
| dimension   | int          | 必填         |
| created_at  | timestamp    | 创建时间     |
| updated_at  | timestamp    | 更新时间     |

---

如有问题欢迎提 issue 或二次开发！ 