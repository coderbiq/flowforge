---
id: FIND-安装版本检查与自动升级-djkddmqxcs12
title: 项目制品 manifest 追踪范围
type: finding
status: draft
importance: should
links:
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: LOG-CR26062102-dji535saoxml
      relation: references
    - target: PROP-CR26062102_flowforge-安装版本检查与自动升级
      relation: belongs_to
    - target: STR-djkddmq0cabv
      relation: indexes
created: 2026-06-28T03:41:42.123662323Z
updated: 2026-06-28T03:41:42.124325057Z
source: CR26062102_flowforge-安装版本检查与自动升级
---

## Finding

`.flowforge/manifest.yaml` 的追踪范围确定为 4 个部署路径源：
- `assets/skills/` → `.agents/skills/`
- `assets/templates/` → `.flowforge/templates/`
- `assets/wiki/` → 项目 wiki 根目录
- `assets/AGENTS.md` → 根目录 AGENTS.md

排除 `.gitkeep` 等占位文件。目标项目托管文件总数 < 20，全量比较策略可行。

## Links

### Outgoing

- `PROP-CR26062102_flowforge-安装版本检查与自动升级` [belongs_to]

