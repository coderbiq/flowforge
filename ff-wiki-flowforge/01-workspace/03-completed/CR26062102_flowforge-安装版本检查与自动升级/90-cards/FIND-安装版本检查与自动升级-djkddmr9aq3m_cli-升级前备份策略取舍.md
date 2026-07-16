---
id: FIND-安装版本检查与自动升级-djkddmr9aq3m
title: CLI 升级前备份策略取舍
type: finding
status: draft
importance: should
links:
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: LOG-CR26062102-dji538oomcky
      relation: references
    - target: PROP-CR26062102_flowforge-安装版本检查与自动升级
      relation: belongs_to
    - target: STR-djkddmq0cabv
      relation: indexes
created: 2026-06-28T03:41:42.143747315Z
updated: 2026-06-28T03:41:42.14438559Z
source: CR26062102_flowforge-安装版本检查与自动升级
---

## Finding

CLI 自更新前不执行全量配置备份。理由：
- binary `.old` 副本已提供二进制回滚能力
- sqlite schema 采用增量/加法式变更保证向前兼容
- ConfigService yaml 配置通过 `schema_version` 字段实现自动迁移

## Impact

简化升级流程，避免每次升级的 I/O 开销。如需降级，直接运行旧版本二进制，ConfigService 检测 schema 兼容性。

## Links

### Outgoing

- `PROP-CR26062102_flowforge-安装版本检查与自动升级` [belongs_to]

