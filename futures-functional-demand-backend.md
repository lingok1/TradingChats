 # 凌期AI辅助期货挑选品种与持仓止盈止损分析（后端需求）
 
 ## 后端技术架构
 
 - Golang
 - Gin：RESTful API + Swagger
 - MongoDB：保存 AI 请求与响应结果（按批次存储）
 - Redis：异步任务状态跟踪
 
 ## 环境与连接信息
 
 Redis 连接（示例）：
 
redis-cli -h 150.158.18.92 -p 6379 -a redis123
 
 MongoDB 连接（示例）：
 
mongodb://admin:mongo123@150.158.18.92:27017
 
 ## 业务需求
 
 ⭐️业务需求（优先实现）：
 
 1. Redis 和 MongoDB 的连接与增删改查测试。
 2. 默认提示词模板（MongoDB）设计：支持标签（股票、期货、期权等），包含增删改查功能。
 3. 获取提示词模板并动态生成提示词：
    - 前端传入 2 个 HTTPS 接口参数
    - 需要从两个 HTTPS 接口获取 JSON 数据并追加到提示词中
    - 获取当前北京时间，开仓时间为当前北京时间 10 分钟后
 4. 模型与 API 配置（MongoDB）设计：
    - 支持自定义 API URL、密钥、模型列表（一次可添加多个模型名称）
    - 支持连通性测试
    - 协议支持：
      - anthropic：/v1/messages（POST）
      - openai：/v1/chat/completions（POST）
 5. 定时任务（交易时段）：
    - 周一到周五日盘/夜盘交易时间段内
    - 从 MongoDB 获取动态生成提示词
    - 10 分钟循环异步调用已配置的模型与 API
    - 响应信息按批次保存到 MongoDB，便于查询
 
 ⭐️业务需求（暂不实现）：
 
 1. 帐户持仓止盈止损建议：截图转文字后由 agent 给出分析建议。
 2. 帐户持仓实时风险监控面板。
 3. 统计挑选品种盈利正确率（下午三点后）：根据收盘价 > 入场区间最大值。
 4. 期货波动率预测。
 
 ## API 设计（与前端联动关键接口）
 
 - 生成批次（触发异步多模型）：POST /api/ai-responses/generate
 - 查询批次（用于历史查看）：GET /api/ai-responses/batch?batch_id=...
 - 查询最新成功批次（前端首屏默认入口）：GET /api/ai-responses/latest
 - 模板管理：/api/prompt-templates（CRUD + generate）
 - 模型配置：/api/model-api-configs（CRUD + test）
 - 定时任务：/api/schedules（CRUD + status + logs）

## 数据与状态契约（前后端对齐）

### 当前接口返回（现状）

- 现阶段 /api/ai-responses/latest、/api/ai-responses/batch 的 data 可能直接返回 AIResponse 数组（每条包含 batch_id、model_name、status、response 等字段）
- 前端可以用 batch_id 对数组做聚合，并按 model_name 分块展示

### 状态机定义

- 模型级状态：queued | running | completed | failed
- 批次级状态：running | completed | partial_failed | failed
- latest 的筛选规则：返回最近一次 batch_status 为 completed 或 partial_failed 的批次；若仅有 failed/空数据则返回 404

### AI 批次返回结构（建议目标）

GET /api/ai-responses/latest 与 GET /api/ai-responses/batch 返回 data 结构建议统一：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "batch_id": "string",
    "batch_status": "completed",
    "template_id": "string",
    "params": {
      "param1": "string",
      "param2": "string"
    },
    "created_at": "YYYY-MM-DD HH:mm:ss",
    "models": [
      {
        "model_name": "string",
        "provider": "openai|anthropic",
        "status": "completed",
        "error": "",
        "response_markdown": "string",
        "signals": [
          {
            "index": 1,
            "symbol": "甲醇(MA605)",
            "direction": "做多|做空",
            "entry_range": "2910-2940"
          }
        ]
      }
    ]
  }
}
```

说明：

- response_markdown：用于 PC 端直接渲染，作为移动端详情兜底
- signals：用于移动端“信号层”列表优先展示；若解析失败可为空数组，但不得影响 response_markdown 返回

## 异步与任务可观测性（Redis）

- 生成批次时应具备幂等能力（同一组输入避免短时间重复生成造成并发风暴）
- 建议生成入口返回 batch_id（或 job_id），并提供任务状态查询能力（轮询/刷新）
- 需要定义超时、重试、并发度限制与失败原因落库规则（模型级 error、批次级 batch_status）

## 交易时段与时间规则（待细化）

- 明确日盘/夜盘的具体起止时间（北京时间）
- 跨日夜盘归属（当日/次日）与节假日策略（停盘、临时休市）
- “开仓时间 = 当前北京时间 + 10 分钟”的边界处理（接近收盘/夜盘结束时的策略）

## 安全与配置规范（必须）

- 文档中的 Redis/Mongo 连接信息为示例；生产环境需使用环境变量或密钥管理，不在仓库与文档中保存明文口令
- 模型 API Key 必须加密存储与脱敏展示，日志与 Swagger 不得输出密钥
- 需要规划接口鉴权、CORS、限流与审计（尤其是 generate、schedules 等写操作）
