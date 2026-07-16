---
id: DES-安装版本检查与自动升级-djkddmq8tkh7
title: 项目制品 manifest.yaml 结构与升级策略
type: design
status: draft
importance: should
links:
    - target: DES-CR26062102-dji543o8ff5s
      relation: references
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: PROP-CR26062102_flowforge-安装版本检查与自动升级
      relation: belongs_to
    - target: STR-djkddmq0cabv
      relation: indexes
created: 2026-06-28T03:41:42.082486406Z
updated: 2026-06-28T03:41:42.083118096Z
source: CR26062102_flowforge-安装版本检查与自动升级
---

## Goal

`.flowforge/manifest.yaml` 记录所有部署到目标项目的文件，`flowforge upgrade` 根据 manifest 比较执行四类文件处理。

## Decision

manifest.yaml 记录 `source_path` / `target_path` / `sha256` / `file_type`。升级时全量比较 CLi 内嵌 assets 与目标项目文件，分为四类：
- **conflict**：双方都改变 → 标记不覆盖
- **added**：仅源新增 → 自动添加
- **updated**：源变目标不变 → 自动更新
- **agents_block**：AGENTS.md 专用 → 仅替换 FLOWFORGE 标记间内容

升级前备份到 `.flowforge/backup/<version>/`，升级后运行 `validate all`。

## Constraints

- manifest.yaml 排除 `.gitkeep` 文件
- 全量比较简单可靠（目标文件数 < 20）
- AGENTS.md 的 sha256 仅对区块内容计算（不含标记行）

## Links

### Outgoing

- `PROP-CR26062102_flowforge-安装版本检查与自动升级` [belongs_to]

