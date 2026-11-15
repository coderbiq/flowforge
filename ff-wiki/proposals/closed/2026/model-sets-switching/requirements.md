# model-sets-switching — Requirements

> Type: proposal(小型)
> Status: implemented (2026-09-29,全部测试绿+本仓双向往返验收通过)
> 决策(2026-09-29 与用户对齐):A 增量覆盖 / B 切换自动重部署+失败回滚 / C set 内复用现有三层结构

## 问题

模型策略当前只有一套静态配置(`agents.models` / `models_by_name` / `models_by_host`)。
实际使用中同一台机器需要在多套模型方案间切换:

- 默认方案:glm-5.3 做旗舰(4 个决策/审查角色)
- 限额耗尽/高峰惩罚时段方案:mimo-2.6-pro 接替旗舰,其余角色不变

现状只能手工编辑 config.yaml 再手工重跑 `agents deploy` — 繁琐且易错
(忘重部署则切换不生效:模型是部署时编译进 agent 文件 frontmatter 的)。

## 需求

### named-model-sets

支持在 config 中声明多套命名模型方案(`agents.model_sets`),每套方案是
**增量覆盖**(overlay):只声明与基础层(顶层三层)的差异,未声明的角色
继续用基础层。set 内结构复用现有三层(`models` / `models_by_name` /
`models_by_host`),校验规则与顶层完全一致。

### quick-switch-command

提供 `flowforge model-set` 命令组:

- `list` — 列出全部 set 并标记当前激活(含基础层显示为 default)
- `use <name>` — 切换:写激活指针 → 自动重新 `agents deploy` → 失败则回滚指针并报错
- `show [name]` — 展示某 set(缺省为当前激活)与基础层**合并后**的有效模型表

### switch-atomicity

切换要么完整生效(指针写入 + 部署产物全部按新方案编译),要么完全不变
(部署失败 → 指针回滚,部署产物保持旧方案)。不允许出现"指针已切、
产物未切"的中间态。

## 非目标

- 自动检测限额/惩罚时段并自动切换(用户手动触发)
- 单角色临时 override 命令(现有 config 三层已覆盖)
- 跨机器同步(激活指针是 per-machine 状态,与 config.yaml 一样不进 git)
