 # 凌期AI辅助期货挑选品种与持仓止盈止损分析（前端需求）
 
 ## 前端技术架构
 
- Vue 3 + Vite
- TypeScript：前端代码使用 TS 语法编写
 - Element Plus
 - Axios：统一请求封装与错误处理
 - Markdown 渲染：markdown-it（渲染 AI 返回的 Markdown 表格）

## Swagger 驱动开发（必须）

- 前端接口对接必须以 Swagger 为准（先看 Swagger 再写代码），避免根据实现细节猜字段
- Swagger 文件来源（以此为准）：
  - trading-chats-backend/docs/swagger.yaml
  - trading-chats-backend/docs/swagger.json
- basePath：/api（前端请求 URL 需要带 /api 前缀）
- 统一响应包裹：models.Response（code/msg/data）
- 若 Swagger 与实际实现不一致：先修正 Swagger，再按 Swagger 开发前端
 
 ## 核心页面：AI 分析报告面板
 
 ### 1) 数据来源
 
 - 首屏默认：GET /api/ai-responses/latest（获取最新一次成功批次的全量模型结果）
 - 历史查看（可选）：GET /api/ai-responses/batch?batch_id=...
 
### 1.1) 数据契约与渲染优先级（与后端对齐）

- 前端期望后端返回“批次 + models[]”的统一结构，其中每个模型包含：
  - response_markdown：原始 Markdown 表格（PC 端主渲染、移动端详情兜底）
  - signals：结构化的信号列表（移动端信号层优先使用）
- 兼容现状：
  - 若 /latest、/batch 暂时返回 AIResponse 数组，前端以数组第一条的 batch_id 作为当前批次，并按 model_name 分组渲染；signals 不存在则按空数组处理
- 渲染优先级：
  - PC：优先渲染 response_markdown（多模型卡片纵向排列对比）
  - 移动端：优先渲染 signals（信号层列表），点击后展示详情；若 signals 为空，则降级为直接展示 response_markdown

### 1.2) 批次信息展示（建议）

- 页面顶部（或批次信息栏）展示：
  - 批次号（batch_id）
  - 批次时间（created_at）
  - 批次状态（completed/partial_failed/failed/running）
  - 成功模型数 / 失败模型数
- 提供一键复制：
  - 复制 batch_id
  - 复制当前模型的 Markdown（PC 端）

### 1.3) 手动触发生成（建议）

- PC 端提供“立即生成”入口（调用 /api/ai-responses/generate）
- 移动端将“立即生成”入口放入右侧抽屉（避免干扰阅读）
- 交互约定：
  - 生成中：展示 loading，并允许用户刷新 latest
  - 生成成功：自动刷新并显示新批次
  - 生成失败：提示错误原因与重试入口

### 1.4) 模型对比辅助（建议）

- 模型卡片支持：
  - 折叠/展开（失败模型默认折叠但可查看 error）
  - 仅看成功模型/显示全部模型切换
  - 模型排序（按完成时间或按模型名称）

 ### 2) PC 端展示方式
 
 - 页面按“模型”分块展示：
   - 每个模型一个卡片（Card）
   - 卡片表头显示模型名称与状态（completed/failed）
   - 卡片内容区直接渲染该模型返回的 Markdown 表格
 - 多模型对比：多个模型卡片纵向排列，便于上下滚动对比
 
 ## 移动端展示方式（以可用性优先）
 
 ### 1) 目标
 
 - 避免 16 列表格在手机上横向滚动带来的低可用性
 - 手机端先看“信号”，点击再看“细节”
 
 ### 2) 信息分层策略
 
 - 信号层（列表/卡片直接展示）：
   - 品种（代码）
   - 多空方向
   - 入场区间
 - 详情层（点击后展示）：
   - 止盈、止损、建议持仓时间
   - 博弈多头逻辑、博弈空头逻辑
   - 技术要点、基本面要点、市场情绪、资金流入/流出、持仓量变化等
 
 ### 3) 交互形式
 
 - 移动端点击单条品种后，以“抽屉/全屏弹层”展示详情内容
 - 详情区域采用纵向布局（描述列表/分组模块），避免横滑阅读

### 4) 异常态与降级策略（必须）

- 无最新成功批次：
  - 展示空态说明与“重新生成/刷新”入口
- 部分模型失败：
  - 对失败模型卡片显示错误摘要（error），并保留成功模型可读内容
- Markdown 渲染失败或内容为空：
  - 以纯文本方式展示 response_markdown，并提示复制/反馈
- 网络超时：
  - 展示可重试按钮，并允许用户手动刷新 latest

### 5) 刷新与历史（可选）

- 提供“刷新最新”按钮（重新拉取 /latest）
- 提供“按 batch_id 查看历史”的入口（输入 batch_id）
- 可选增强：提供最近 N 条批次列表（按时间倒序），点击进入后展示该批次下多模型对比

### 6) 移动端信号层增强（可选）

- 筛选与搜索：
  - 只看做多 / 只看做空
  - 关键字搜索（品种代码/名称）
- 排序：
  - 按序号
  - 按方向（做多/做空分组）

### 7) 详情层阅读体验（建议）

- 详情弹层/抽屉内的内容按模块分段展示：
  - 交易参数（入场/止盈/止损/持仓建议）
  - 多头逻辑 / 空头逻辑
  - 技术要点 / 基本面要点
  - 市场情绪 / 资金面
- 详情层提供快捷操作：
  - 复制当前品种的“信号摘要”
  - 复制当前品种的“完整逻辑”
 
 ## 设置与管理入口（移动端右侧抽屉）
 
 - 移动端设置入口采用右侧抽屉（Right Drawer）：
   - 明暗主题切换
   - 模板管理入口
   - 模型配置入口
   - 定时任务管理入口

### 抽屉交互约定

- 移动端：
  - 右侧抽屉用于“设置/管理”（全局功能）
  - 底部抽屉或全屏弹层用于“单条品种详情”（内容阅读）

### 抽屉内容建议（可选）

- 阅读与展示偏好：
  - 字号大小（标准/大号）
  - 布局密度（标准/紧凑）
  - 默认展示（仅成功模型/全部模型）
- 刷新策略：
  - 自动刷新开关
  - 自动刷新间隔（仅做需求约定）
- 关于与免责声明入口：
  - 集中放入抽屉底部，避免占用主页面空间
 
 ## 主题与视觉规范
 
 - 明暗主题切换：
   - 明亮模式：适合白天
   - 暗黑模式：适合夜盘
- 主题切换要求：
  - 主题选择需要持久化（本地存储），刷新后保持用户选择
  - 默认策略：优先读取用户选择；无选择时跟随系统 prefers-color-scheme
 - 语义化颜色：
   - 做多：红色系
   - 做空：绿色系
 - Markdown 表格样式：
   - PC 端优先保证排版整齐与可读性
   - 移动端以“信号层 + 详情层”替代大表格直出

## 性能与渲染策略（建议）

- Markdown 渲染可能内容较大：
  - PC 端建议按模型卡片做懒渲染（滚动进入可视区再渲染）
  - 移动端详情弹层建议延迟渲染 Markdown（先展示信号与结构化字段）

## API 对接清单（从 Swagger 提取）

以下接口以 Swagger（basePath=/api）为准：

- AI 响应
  - GET /ai-responses/latest
  - GET /ai-responses/batch?batch_id=...
  - POST /ai-responses/generate（body：template_id,param1,param2）
  - GET /ai-responses
  - GET /ai-responses/{id}
- 提示词模板
  - POST/GET /prompt-templates
  - GET/PUT/DELETE /prompt-templates/{id}
  - GET /prompt-templates/tag?tag=...
  - POST /prompt-templates/generate
- 模型配置
  - POST/GET /model-api-configs
  - GET/PUT/DELETE /model-api-configs/{id}
  - GET /model-api-configs/provider?provider=...
  - POST /model-api-configs/{id}/test
- 定时任务
  - POST/GET /schedules
  - PUT /schedules/{id}/status（body：{"status":"active|paused"}）
  - GET /schedules/{id}/logs
  - DELETE /schedules/{id}
