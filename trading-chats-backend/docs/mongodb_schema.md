# MongoDB 数据库表结构文档

## 1. 集合：`prompt_templates`（提示词模板）

| 字段名 | 数据类型 | 说明 | 索引 |
| --- | --- | --- | --- |
| `_id` | ObjectID | 文档唯一标识 | 主键索引 |
| `name` | String | 模板名称 | 普通索引 |
| `content` | String | 模板内容 | - |
| `tags` | Array<String> | 模板标签，例如期货、期权 | 多键索引 |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |

用途：存储 AI 提示词模板，支持按名称和标签管理。

## 2. 集合：`model_api_configs`（模型与 API 配置）

| 字段名 | 数据类型 | 说明 | 索引 |
| --- | --- | --- | --- |
| `_id` | ObjectID | 文档唯一标识 | 主键索引 |
| `name` | String | 配置名称 | 普通索引 |
| `api_url` | String | API 地址 | - |
| `api_key` | String | API Key | - |
| `models` | Array<String> | 模型名称列表 | - |
| `provider` | String | 提供商，例如 `openai`、`anthropic` | 普通索引 |
| `tab_settings` | Array<Object> | Tab 页配置数组，每个对象包含 `tab_tag` 和 `enabled` | - |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |

用途：存储 AI 模型连接配置，并通过 `tab_settings` 数组标记该配置归属的业务 Tab 以及是否启用。

## 3. 集合：`ai_responses`（期货 AI 响应）

| 字段名 | 数据类型 | 说明 | 索引 |
| --- | --- | --- | --- |
| `_id` | ObjectID | 文档唯一标识 | 主键索引 |
| `tenant_id` | String | 租户 ID | 普通索引 |
| `batch_id` | String | 批次 ID | 普通索引 |
| `prompt` | String | 发送给模型的提示词 | - |
| `response` | String | 模型响应内容 | - |
| `model_api_id` | ObjectID | 对应 `model_api_configs` 的 `_id` | 普通索引 |
| `model_api_name` | String | 对应 `model_api_configs` 的 `name` | 普通索引 |
| `model_name` | String | 模型名称 | 普通索引 |
| `provider` | String | 提供商 | 普通索引 |
| `status` | String | 状态：`pending`、`completed`、`failed` | 普通索引 |
| `error` | String | 错误信息 | - |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |
| `completed_at` | Date | 完成时间 | - |

用途：存储期货 Tab 的模型生成结果。

## 4. 集合：`ai_responses_options`（期权 AI 响应）

字段结构与 `ai_responses` 完全一致：

| 字段名 | 数据类型 | 说明 | 索引 |
| --- | --- | --- | --- |
| `_id` | ObjectID | 文档唯一标识 | 主键索引 |
| `tenant_id` | String | 租户 ID | 普通索引 |
| `batch_id` | String | 批次 ID | 普通索引 |
| `prompt` | String | 发送给模型的提示词 | - |
| `response` | String | 模型响应内容 | - |
| `model_api_id` | ObjectID | 对应 `model_api_configs` 的 `_id` | 普通索引 |
| `model_api_name` | String | 对应 `model_api_configs` 的 `name` | 普通索引 |
| `model_name` | String | 模型名称 | 普通索引 |
| `provider` | String | 提供商 | 普通索引 |
| `status` | String | 状态：`pending`、`completed`、`failed` | 普通索引 |
| `error` | String | 错误信息 | - |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |
| `completed_at` | Date | 完成时间 | - |

用途：存储期权 Tab 的模型生成结果。

## 5. 集合：`ai_responses_news`（新闻 AI 响应）

字段结构与 `ai_responses` 完全一致：

| 字段名 | 数据类型 | 说明 | 索引 |
| --- | --- | --- | --- |
| `_id` | ObjectID | 文档唯一标识 | 主键索引 |
| `tenant_id` | String | 租户 ID | 普通索引 |
| `batch_id` | String | 批次 ID | 普通索引 |
| `prompt` | String | 发送给模型的提示词 | - |
| `response` | String | 模型响应内容 | - |
| `model_api_id` | ObjectID | 对应 `model_api_configs` 的 `_id` | 普通索引 |
| `model_api_name` | String | 对应 `model_api_configs` 的 `name` | 普通索引 |
| `model_name` | String | 模型名称 | 普通索引 |
| `provider` | String | 提供商 | 普通索引 |
| `status` | String | 状态：`pending`、`completed`、`failed` | 普通索引 |
| `error` | String | 错误信息 | - |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |
| `completed_at` | Date | 完成时间 | - |

用途：存储新闻 Tab 的模型生成结果。

## 6. 集合：`ai_responses_position`（持仓 AI 响应）

字段结构与 `ai_responses` 完全一致：

| 字段名 | 数据类型 | 说明 | 索引 |
| --- | --- | --- | --- |
| `_id` | ObjectID | 文档唯一标识 | 主键索引 |
| `tenant_id` | String | 租户 ID | 普通索引 |
| `batch_id` | String | 批次 ID | 普通索引 |
| `prompt` | String | 发送给模型的提示词 | - |
| `response` | String | 模型响应内容 | - |
| `model_api_id` | ObjectID | 对应 `model_api_configs` 的 `_id` | 普通索引 |
| `model_api_name` | String | 对应 `model_api_configs` 的 `name` | 普通索引 |
| `model_name` | String | 模型名称 | 普通索引 |
| `provider` | String | 提供商 | 普通索引 |
| `status` | String | 状态：`pending`、`completed`、`failed` | 普通索引 |
| `error` | String | 错误信息 | - |
| `created_at` | Date | 创建时间 | - |
| `updated_at` | Date | 更新时间 | - |
| `completed_at` | Date | 完成时间 | - |

用途：存储持仓 Tab 的模型生成结果。

## 7. Redis：`task_status`（任务状态）

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| `task_id` | String | 任务 ID |
| `status` | String | 任务状态 |
| `created_at` | Date | 创建时间 |
| `updated_at` | Date | 更新时间 |

用途：存储临时任务执行状态，用于快速查询任务进度。

## 索引建议

1. `prompt_templates`
   - `name`
   - `tags`
2. `model_api_configs`
   - `name`
   - `provider`
3. `ai_responses`
   - `tenant_id`
   - `batch_id`
   - `status`
   - `model_api_id`
   - `model_api_name`
   - 组合索引：`model_name + provider`
4. `ai_responses_options`
   - 与 `ai_responses` 相同
5. `ai_responses_news`
   - 与 `ai_responses` 相同

6. `ai_responses_position`
   - 与 `ai_responses` 相同

## 数据关系

- `prompt_templates` 与各类 `*_responses`：一对多。
- `model_api_configs` 与各类 `*_responses`：一对多。

## 存储策略

- MongoDB：存储提示词模板、模型配置、各业务 Tab 的 AI 响应等业务数据。
- Redis：存储临时任务状态，提供更快的状态读取能力。

## Trade Plans

集合名：`trade_plans`

用途：
- 存储期货、期权交易计划
- 使用 `tab_tag` 区分 `futures` / `options`
- 按 `tenant_id` 做租户隔离

字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `_id` | ObjectId | 主键 |
| `tenant_id` | String | 租户 ID |
| `tab_tag` | String | `futures` 或 `options` |
| `name` | String | 计划名称 |
| `symbol` | String | 合约代码 |
| `strategy` | String | 策略说明，如 `breakout_follow`、`buy_call` |
| `direction` | String | 方向，如 `long`、`short`、`bullish`、`bearish` |
| `entry_price` | Number | 入场价 |
| `take_profit` | Number | 止盈价 |
| `stop_loss` | Number | 止损价 |
| `open_time` | String | 计划开仓时间 |
| `close_time` | String | 计划平仓时间 |
| `status` | String | `planned` / `active` / `closed` / `cancelled` |
| `remark` | String | 备注 |
| `created_at` | Date/String | 创建时间 |
| `updated_at` | Date/String | 更新时间 |

索引建议：
- `tenant_id + tab_tag + status + updated_at`
- `tenant_id + symbol + tab_tag`
