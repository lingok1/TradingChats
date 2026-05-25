## 前端

- ✅响应的结果前端展示需要带正则匹配（例如| 序号 |和\n| 1 |）才展示
- ✅模版生成，复制模版生成的内容到粘贴板（兼容pc和移动端）
- ✅模型配置的测试不知道测试**单一**还是全部接口
- ✅设置接入鉴权功能，修改和删除的接口需要进行鉴权操作
- ✅定时任务开关路由不正确
- ✅failed to get prompt template: mongo: no documents in result
- ✅新增定时任务，param1（URL）删除，模版下拉选择
- ✅移动端菜单按钮未显示，向上滑会跳转到其他菜单页面（新闻页面），向右滑动会加载当日分析页面
- ✅移动端左右滑动问题切换标签页
- ✅上划到顶部（兼容pc和移动端）
- ✅前端移动端展示，止盈止损占用一行
- ✅导航菜单更改 期货、期权、新闻、计划、持仓、关于
- ✅批次：删除，时间改为数据更新
- ✅去掉持仓时间，天
- ✅页面设计期权、新闻、计划、持仓、关于
- ✅移动端将菜单切换放到后面，用户登录放到前面
- ✅刷新按钮放到main里面
- ✅右下角默认显示明暗切换，明暗切换和返回顶部放一起，顶部上明暗下
- ✅移动端，颜色为灰色，需要为白色更清晰
- ✅最近刷新和成功刷新放在一行
- ✅关于，文字打字机的效果
- ✅定时刷新页面，根据系统配置来刷新
- ✅后端完成ai response后sse通知前端刷新界面
- ✅顶部标签页切换后获取最新数据
- ✅参数，模型等添加搜索
- ✅模型的启用禁用
- ✅期权调用全部AI，新闻调用一个就够了
- ✅期权和新闻提示词
- ✅内容表格右滑，显示品种(弃用)
- ✅模型列表变长一点，样式可以重新设计
- ✅定时任务日志不显示还是单独一张表存储prompt
- ✅交易计划前端展示优化
- ✅主页介绍仿站 [PnLClaw - Crypto Quantitative Trading Platform](https://www.pnlclaw.com/)
- ✅主页介绍仿站[ ](https://www.pnlclaw.com/)[LUCKYQUANT PRO](https://luckyquant.top/)
- ✅移动端，设置重新设计，提醒用户左右滑动快速切换（弃用）
- ✅前端icon修改
- ✅移动端顶部（凌期ai旁边）添加当前界面
- ✅注释隐藏持仓和新闻tab页，不删除代码
- ✅token长时间未自动刷新(过期时间改为一周)
- ✅期货涨跌平判断市场情绪
- ✅计划前端页面调整
- ✅计划和期权更换icon
- ✅股票模块提示词，中金所4个，国债4个
- ✅动态参数输入框，输入一个字符后会失去焦点，我需要能输入多个字符
- ✅期货优选品种推荐（先获取已经分析好的数据让ai根据提示词进行推荐）
- ✅前端新增首页tab页功能（放在期货的左边），里面可以放期货、期权、股票优选推荐,优选推荐内容完成后前端自动刷新
- <br />
- <br />
- 期权tab页添加策略子tab页
- 商品期权策略（数据获取和提示词）
- 实时刷新已连接不显示（合并到成功n/n里面，黄色和红色来区分功能和失败）
- 统计模型正确率功能
- 前端底部备案信息改为一行
- 股票优选入库的时候没有tag：stock导致无法再前端显示
- 股票改为股指
- **首页-期权优选-未自动刷新**
- 定时任务里面的Tab 页签可修改
- futures-sentiment和futures-top-movers接口从动态参数获取数据
- 涨跌幅榜优化，日盘夜盘交易时间内才调用
- 新闻暂不实现（事件驱动会有滞后问题）
- 微信扫码登录
- <br />
- <br />
- <br />

## 后端

- ✅ai\_responses表里面的prompt字段放到schedule\_logs表里面，schedule\_logs表添加字段区分是手动触发还是自动触发的任务
- ✅schedule\_configs表删除param1和param2，executeTask从system\_configs表里面获取param1和param2的数据
- ✅开启定时任务后不会定时运行
- ✅添加登录鉴权和多租户功能，修改和删除的接口需要进行鉴权操作，更新swagger接口文档，添加一个管理员和两个租户，给默认的账号密码
- ✅响应的结果截取
- ✅新增修改提示词需要给默认租户
- ✅定时任务修改接口
- ✅服务器部署.md
- ✅model\_api\_configs表添加tab页标签（期货、期权、新闻、持仓）和tab页标签模型是否启用
- ✅ai\_responses表添加model\_api\_name和model\_api\_id两个字段（对应model\_api\_configs表的id和name）
- ✅根据ai\_responses表新增期权、新闻、持仓3张新表
- ✅后端完成ai response后sse通知前端刷新界面
- ✅交易计划表设计和功能实现
- ✅定时任务跑多个无法运行（时区原因）
- ✅model\_api\_configs表tab\_tag + tab\_enable 兼容字段删除
- <br />
- <br />
- 租户相互之间数据隔离
- schedule\_logs日志表提示词导致数据量大
- ai\_responses表按天分表
- <br />

<br />

## 期货期权持仓前后端功能设计与实现

### 整体流程

```
截图上传 → AI多模态OCR识别 → 可编辑预览(增删改) → 确认保存 → 持仓列表 → AI风险分析
```

### 一、数据模型（Position 期货期权合一表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 主键 |
| tenant_id | string | 租户隔离 |
| position_type | enum | `futures` / `options` |
| symbol | string | 合约代码，如 IF2406、m2407-C-3200 |
| direction | enum | `long` / `short` |
| volume | int | 手数 |
| open_price | float | 开仓均价（计算基准） |
| margin | float | 保证金（期货） |
| option_type | enum | `call` / `put`（期权专用） |
| strike_price | float | 行权价（期权专用） |
| premium | float | 开仓权利金（期权，计算基准） |
| underlying_price | float | 标的价格（期权专用） |
| expire_date | date | 到期日（期权专用） |
| contract_multiplier | int | **新增**，合约乘数（如 IF=300） |
| status | enum | `holding` / `closed` |
| source | enum | `manual` / `ocr` |
| screenshot_url | string | 原始截图地址（可选） |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

> **注意**：`current_price` 和 `profit_loss` 不入库，前端实时获取行情数据并计算

### 二、前端动态计算

| 展示字段 | 计算方式 |
|---------|---------|
| current_price | 实时行情接口获取 |
| profit_loss | `(current - open) × volume × multiplier × direction` |
| profit_loss_pct | `profit_loss / (open_price × volume × multiplier)` |

**行情数据方案**：后端代理行情请求（避免前端跨域，可加缓存），后续可升级 WebSocket/SSE 推送

### 三、API 设计

融入现有 `tab_tag` 路由体系，`tab_tag=position`

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/positions` | GET | 查询持仓列表（支持 status 筛选） |
| `/api/positions` | POST | 手动新增单条持仓 |
| `/api/positions/:id` | PUT | 编辑持仓 |
| `/api/positions/:id` | DELETE | 删除持仓 |
| `/api/positions/:id/close` | POST | 平仓（status → closed） |
| `/api/positions/batch` | POST | 批量保存（OCR 确认后，覆盖模式） |
| `/api/positions/ocr` | POST | 上传截图，返回 OCR 识别结果（不入库） |
| `/api/positions/analyze` | POST | 触发 AI 风险分析 |

**OCR 接口流程**：
- `POST /api/positions/ocr` → 接收图片 → 调用 AI 多模态识别 → 返回结构化 JSON（不入库）
- `POST /api/positions/batch` → 接收 Position[] + mode(overwrite/append) → 覆盖模式清空 holding 后写入；追加模式直接插入

**AI 分析接口**：
- `POST /api/positions/analyze` → 输入当前 tenant 的 holding 持仓 → 组装持仓 + 行情 + 已有市场分析 → Markdown 输出 → 存入 `ai_responses_position`

### 四、前端交互设计

**页面布局**：

```
┌─────────────────────────────────────────┐
│  [上传截图]  [手动添加]  [AI分析]        │
├─────────────────────────────────────────┤
│  AI 分析结果面板（折叠/展开）            │
│  - 风险评估 / 卖出建议                   │
├─────────────────────────────────────────┤
│  Tab: 期货持仓 │ 期权持仓 │ 已平仓       │
├─────────────────────────────────────────┤
│  持仓列表表格                            │
│  - 合约 | 方向 | 手数 | 开仓价 | 浮盈    │
│  - [编辑] [平仓] [删除]                  │
└─────────────────────────────────────────┘
```

**截图上传确认流程**：
1. 上传截图
2. 显示 loading，调用 /ocr 接口
3. 返回可编辑表单（预填 OCR 结果），每行一个持仓，可增删改
4. 用户确认 → 调用 /batch 接口（覆盖模式）
5. 刷新持仓列表

### 五、AI 风险分析设计

**输入组装**：
1. 当前持仓数据（holding 状态）
2. 行情数据（current_price，后端代理获取）
3. 已有市场分析（ai_responses 中最新内容）

**Prompt 方向**：
- 单持仓风险：盈亏比、止损位置、持仓占比
- 组合风险：多空对冲情况、总保证金占用、集中度
- 卖出参考：结合市场情绪给出减仓/平仓/持有的建议

**输出**：Markdown 格式，存入 `ai_responses_position`，前端折叠面板展示（默认折叠，有结果时展开）

### 六、实现优先级

| 阶段 | 内容 |
|------|------|
| P0 | 数据模型 + CRUD API + 前端持仓列表（手动增删改查） + 行情代理接口 |
| P1 | OCR 上传识别 + 确认流程 |
| P2 | AI 风险分析集成 |
| P3 | 平仓历史、与交易计划关联、实时行情 WebSocket/SSE 推送 |

