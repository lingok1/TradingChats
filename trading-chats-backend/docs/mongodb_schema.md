# MongoDB数据库表结构文档

## 1. 集合：prompt_templates（提示词模版）

| 字段名 | 数据类型 | 描述 | 索引 |
|-------|---------|------|------|
| `_id` | ObjectID | 文档唯一标识符 | 主键索引 |
| `name` | String | 提示词模版名称 | 普通索引 |
| `content` | String | 提示词内容 | - |
| `tags` | Array<String> | 标签，如股票、期货、期权等 | 多键索引 |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |

**用途**：存储AI提示词模版，支持按标签分类和检索。

## 2. 集合：model_api_configs（模型与API配置）

| 字段名 | 数据类型 | 描述 | 索引 |
|-------|---------|------|------|
| `_id` | ObjectID | 文档唯一标识符 | 主键索引 |
| `name` | String | 配置名称 | 普通索引 |
| `api_url` | String | API接口地址 | - |
| `api_key` | String | API密钥 | - |
| `models` | Array<String> | 模型名称列表 | - |
| `provider` | String | 提供商，如anthropic、openai | 普通索引 |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |

**用途**：存储AI模型的API配置信息，支持不同提供商的配置管理。

## 3. 集合：ai_responses（AI模型响应信息）

| 字段名 | 数据类型 | 描述 | 索引 |
|-------|---------|------|------|
| `_id` | ObjectID | 文档唯一标识符 | 主键索引 |
| `batch_id` | String | 批次ID | 普通索引 |
| `prompt` | String | 发送的提示词 | - |
| `response` | String | AI模型的响应 | - |
| `model_name` | String | 使用的模型名称 | 普通索引 |
| `provider` | String | 提供商 | 普通索引 |
| `status` | String | 任务状态：pending, completed, failed | 普通索引 |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |
| `completed_at` | Date | 完成时间 | - |

**用途**：存储AI模型的响应数据，支持按批次、状态和提供商检索。

## 4. Redis存储：task_status（任务状态）

| 字段名 | 数据类型 | 描述 |
|-------|---------|------|
| `task_id` | String | 任务ID |
| `status` | String | 任务状态 |
| `created_at` | Date | 创建时间 |
| `updated_at` | Date | 更新时间 |

**用途**：在Redis中存储任务状态，用于实时跟踪任务执行情况。

## 索引建议

1. **prompt_templates集合**：
   - `name`字段创建普通索引，加速按名称查询
   - `tags`字段创建多键索引，加速按标签查询

2. **model_api_configs集合**：
   - `name`字段创建普通索引，加速按名称查询
   - `provider`字段创建普通索引，加速按提供商查询

3. **ai_responses集合**：
   - `batch_id`字段创建普通索引，加速按批次查询
   - `status`字段创建普通索引，加速按状态查询
   - `model_name`和`provider`字段创建组合索引，加速按模型和提供商查询

## 数据关系

- `prompt_templates`与`ai_responses`：一对多关系，一个提示词模版可以生成多个AI响应
- `model_api_configs`与`ai_responses`：一对多关系，一个模型配置可以产生多个AI响应

## 存储策略

- **MongoDB**：存储主要业务数据，包括提示词模版、模型配置和AI响应
- **Redis**：存储临时任务状态，提供快速的状态查询和更新

此数据库设计支持系统的核心功能，包括提示词管理、模型配置管理和AI响应存储，同时通过合理的索引设计提高查询性能。