---
id: DES-安装版本检查与自动升级-djkddmpo0cqi
title: 卸载命令分层清理设计
type: design
status: draft
importance: should
links:
    - target: PROP-CR26062102_flowforge-安装版本检查与自动升级
      relation: belongs_to
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: DES-CR26062102-dji4ezbsk312
      relation: references
    - target: STR-djkddmozn3z9
      relation: indexes
created: 2026-06-28T03:41:42.047377969Z
updated: 2026-06-28T03:41:42.048162983Z
source: CR26062102_flowforge-安装版本检查与自动升级
---

## Goal

`flowforge uninstall` 实现分层清理：CLI 二进制 → `~/.flowforge/` 配置目录 → 目标项目托管文件（可选）。

## Decision

默认删除二进制和配置目录。`--keep-config` 仅删二进制。`--project <path>` 额外删除项目托管文件（`.agents/skills/`、`.flowforge/`、AGENTS.md 中的 FLOWFORGE 区块）。卸载前显示删除列表并等待确认（`--yes` 跳过）。幂等，重复卸载不报错。

## Constraints

- 不删除非 FlowForge 托管的用户文件
- 不自动修改 PATH 或环境变量
- 不依赖系统包管理器

## Links

### Outgoing

- `PROP-CR26062102_flowforge-安装版本检查与自动升级` [belongs_to]

