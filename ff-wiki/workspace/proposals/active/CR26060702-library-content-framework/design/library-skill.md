---
doc_type: design
title: flowforge-library-doctor/surgeon/keeper 三 SKILL 完整设计
status: draft
created: 2026-06-07T06:30:00Z
updated: 2026-06-07T06:30:00Z
domain:
  scope: system
  type: design
---

# flowforge-library-doctor/surgeon/keeper 三 SKILL 设计

## 拆分理由

| 维度 | doctor | surgeon | keeper |
|------|--------|---------|--------|
| 操作性质 | 只读诊断 | 写操作（需确认） | 自动化维护 |
| 安全边界 | 无风险 | .bak 回滚 | 需信任 |
| 触发信号 | "检查"、"有没有问题" | "修复"、"合并"、"重组" | "定期维护"、"自动监控" |

## flowforge-library-doctor（诊断引擎）

### Description

```yaml
description: |
  FlowForge Library 诊断引擎。在不依赖 proposal 的情况下，
  对 library 进行健康检查、过期诊断、合规报告。

  必须在以下场景激活：
  - "检查 library"、"library 健康"、"有没有过期内容"
  - "断链检测"、"重复内容"、"孤立文档"
  - "library 合规报告"、"多少 seed"、"maturity 分布"
  - Agent 自主识别到距上次检查超过 7 天

  不要在以下情况激活：
  - 代码/测试问题 → flowforge-feedback
  - proposal 探索中查阅 library → flowforge-design
  - 持续监控 library → flowforge-library-keeper
  - "修复"、"合并"、"重组" → flowforge-library-surgeon
```

### 工作流

```
识别场景 → 运行诊断 → 输出报告 → 建议操作

诊断命令:
  flowforge library check --staleness    # 过期
  flowforge library check --broken-refs  # 断链
  flowforge library check --duplicates   # 重复
  flowforge library check --orphans      # 孤立
  flowforge library check --validate-all # 合规
  flowforge library graph backlinks <p>  # 反向链接
  flowforge library graph blast-radius <p> # 影响范围

自主提示 (S14):
  距上次 check > 7 天 → 提示用户
```

## flowforge-library-surgeon（修复引擎）

### Description

```yaml
description: |
  FlowForge Library 修复引擎。在不依赖 proposal 的情况下，
  直接探索记录、优化内容、重組目录、修复问题。

  必须在以下场景激活：
  - 无活跃 proposal 时，"探索 X 并记录到 library"、"扫描代码更新 library"
  - "合并文档"、"拆分文档"、"提取章节"
  - "重命名并更新引用"、"重组 library 目录"
  - "修复断链"、"标准化 frontmatter"
  - "搜索 library"、"library 里有没有 X"
  - "刷新索引"、"更新 INDEX"
  - "标记为废弃"、"提升为 must"、"清理过期"

  不要在以下情况激活：
  - proposal 探索中写入 → flowforge-design
  - proposal 归档时合成 → flowforge-archive
  - "检查 library"、"有没有问题" → flowforge-library-doctor
```

### 工作流

```
识别操作 → 生成变更计划(dry-run) → 预览 → 确认 → 执行 → 校验

所有写操作遵循五步:
  Plan → Preview → Confirm → Execute → Verify

回滚: .bak 备份 + flowforge library surgeon rollback
```

## flowforge-library-keeper（维护引擎）

### Description

```yaml
description: |
  FlowForge Library 维护引擎。持续监控 library 状态，
  自动检测过期、建议更新、维护索引。

  必须在以下场景激活：
  - "定期维护 library"、"自动监控"、"持续检测"
  - "设置 library 自动检查"
  - "生成维护报告"

  不要在以下情况激活：
  - 一次性检查 → flowforge-library-doctor
  - proposal 相关 → flowforge-archive
```

### 工作流

```
设置监控 → 定期触发 → 生成报告 → 建议维护

核心命令:
  flowforge library keeper watch     # 启动监控（文件变更触发检查）
  flowforge library keeper report    # 生成维护摘要
  flowforge library keeper drift     # 过期检测 + 自动标记
```

## 场景覆盖

| 场景 | SKILL |
|------|-------|
| S11 独立健康检查 | doctor |
| S12 直接探索记录 | surgeon |
| S13 手动维护 | surgeon |
| S14 自主提示 | doctor |
| S15 内容提取 | surgeon |
| S16 内容合并 | surgeon |
| S17 目录重组 | surgeon |
| S18 孤立检测 | doctor |
| S19 反向链接 | doctor |
| S20 关系图谱 | doctor |
| S21 成熟度仪表盘 | doctor |
