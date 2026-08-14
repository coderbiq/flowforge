# v3 实施计划

> 文档状态：`planned`。本文是执行计划，不是实现事实、兼容承诺或迁移方案。

## 目标

- 固化 PROP、FEATURE、CONV、DEC、MOD、FIND 六类当前模型。
- 由 FEATURE 阶段和步骤承载设计到实施的连续上下文。
- 由 `proposal inspect` 聚合 PROP control-plane metadata。
- 让普通卡片查询和派生索引只处理 current-v3 数据。

## 计划边界

实现前必须有已批准 FEATURE Step、明确写集和验证要求。Step 完成后更新 Step、History、Verification，再运行卡片和 Proposal 检查。

旧 task、structure、log create、requirement CLI 直接删除。STR 仅保留为 Proposal control-plane metadata。计划不包含旧 wiki、旧 ID/links 的迁移，不包含 UI，也不把历史数据改写为 v3。

## 验收

每个实施 Step 至少提供相关测试或扫描证据，并确认未修改批准范围外文件。计划中的条目只有在对应代码、测试和 Verification 证据存在时才能标记完成。
