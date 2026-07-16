---
id: DES-CR26062102-dji543o8ff5s
title: 项目制品升级 manifest 结构与升级策略设计
type: design
status: draft
importance: should
links:
- target: PROP-CR26062102
  relation: belongs_to
- target: REQ-CR26062102-djeu2wuos60w
  relation: implements
created: 2026-06-25 12:47:46.532584+00:00
updated: 2026-06-25 12:47:46.532843+00:00
source: CR26062102
slug: 项目制品升级-manifest-结构与升级策略设计
---

# 项目制品升级 manifest 结构与升级策略设计

## Goal

设计 .flowforge/manifest.yaml 的数据结构和项目制品升级的文件处理策略，包括 AGENTS.md 的区块替换。

## Decision

manifest.yaml 记录所有部署到目标项目的文件（4 个路径源），排除 .gitkeep 占位文件。每条记录包含 source_path、target_path、sha256、file_type。升级采用全量比较策略：按 source_path 对比 manifest，分为四类处理：

- **conflict**：双方都改变 → 标记，不覆盖
- **added**：仅源新增 → 自动添加
- **updated**：源变目标不变 → 自动更新
- **agents_block**：AGENTS.md 专用 → 仅替换 FLOWFORGE 标记间内容，sha256 对区块内容计算

manifest.yaml 结构：
```yaml
version: 1
cli_version: 0.1.0
files:
  - source: assets/skills/flowforge-design/SKILL.md
    target: .agents/skills/flowforge-design/SKILL.md
    sha256: abc123...
    type: skill
  - source: assets/AGENTS.md
    target: AGENTS.md
    sha256: def456...
    type: agents_block
    markers:
      start: "<!-- FLOWFORGE:START -->"
      end: "<!-- FLOWFORGE:END -->"
```

## Rationale

- 全量比较简单可靠：目标项目托管文件数量少（< 20），性能不是瓶颈
- 四分类策略覆盖所有变更场景
- AGENTS.md 独立 file_type 区分其特殊处理逻辑
- 按 source_path 而非 target_path 索引

## Constraints

- manifest.yaml 由 `flowforge init` 首次创建，`flowforge upgrade` 更新
- backup 目录存储升级前的完整制品副本：`.flowforge/backup/<cli_version>/`
- conflict 文件标记但不覆盖，提示用户手动处理
- .gitkeep 文件不编入 manifest
- AGENTS.md 的 sha256 仅对标记间区块内容计算，不含标记行

## Impact

- 新增 internal/core/manifest.go 实现 manifest 读写和比较
- 新增 internal/core/agents_block.go 实现区块替换逻辑
- flowforge init 时写入初始 manifest.yaml（含 AGENTS.md 区块处理）
- flowforge upgrade 时读取 manifest.yaml 执行四类文件处理

## Verification

- init 后 manifest.yaml 包含所有部署文件记录，AGENTS.md 为 agents_block 类型
- upgrade 无变更时：manifest 无 diff，报告 "已是最新"
- upgrade 有新增文件时：自动添加，报告新增列表
- upgrade 有变更文件时：自动覆盖，报告更新列表
- upgrade AGENTS.md 区块变更时：仅替换标记间内容，用户内容无损
- upgrade 有冲突时：标记 conflict，不覆盖，报告冲突列表

## Follow-up Tasks

- 实现 manifest.go 读写和比较逻辑
- 实现 agents_block.go 区块替换逻辑
- 实现四类文件处理策略
- 集成到 init 和 upgrade 命令

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [REQ-CR26062102-djeu2wuos60w](REQ-CR26062102-djeu2wuos60w_目标项目制品升级.md) [requirement] - 目标项目制品升级

### Incoming

#### implements
- [TASK-CR26062102-i-dji5li6ksix9](TASK-CR26062102-i-dji5li6ksix9_实现-project-manifestyaml-读写与文件比较逻辑.md) [task] - 实现 project manifest.yaml 读写与文件比较逻辑
- [TASK-CR26062102-i-dji5ln67galh](TASK-CR26062102-i-dji5ln67galh_实现-agentsmd-区块替换与四类文件处理策略.md) [task] - 实现 AGENTS.md 区块替换与四类文件处理策略
- [TASK-CR26062102-i-dji5lsjfsi1c](TASK-CR26062102-i-dji5lsjfsi1c_集成制品升级到-upgrade-和-init-命令.md) [task] - 集成制品升级到 upgrade 和 init 命令 — 备份、验证、报告
- [DES-CR26062102-dji5hnjgds9i](DES-CR26062102-dji5hnjgds9i_agentsmd-区块包裹部署规范.md) [design] - AGENTS.md 区块包裹部署规范

