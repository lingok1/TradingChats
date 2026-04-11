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
- ✅导航菜单更改 期货、期权、新闻、计划、持仓
- 定时刷新页面，根据系统配置来刷新
- 后端完成ai response后sse通知前端刷新界面
- 微信扫码登录

## 后端

- ✅ai\_responses表里面的prompt字段放到schedule\_logs表里面，schedule\_logs表添加字段区分是手动触发还是自动触发的任务
- ✅schedule\_configs表删除param1和param2，executeTask从system\_configs表里面获取param1和param2的数据
- ✅开启定时任务后不会定时运行
- ✅添加登录鉴权和多租户功能，修改和删除的接口需要进行鉴权操作，更新swagger接口文档，添加一个管理员和两个租户，给默认的账号密码
- ✅响应的结果截取
- ✅新增修改提示词需要给默认租户
- ✅定时任务修改接口
- ai\_responses表添加model\_api\_configs表中的id字段
- model\_api\_configs 模型是否启用功能
- 租户相互之间数据隔离
- ai\_responses表按天分表
- 服务器部署.md

