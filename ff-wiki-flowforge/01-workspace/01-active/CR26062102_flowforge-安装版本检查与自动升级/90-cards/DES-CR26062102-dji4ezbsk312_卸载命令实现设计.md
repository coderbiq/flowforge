---
id: DES-CR26062102-dji4ezbsk312
title: 卸载命令实现设计
type: design
status: draft
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wuow388
      relation: implements
created: 2026-06-25T12:14:57.969102999Z
updated: 2026-06-25T12:16:13.449550847Z
source: CR26062102
---

# 卸载命令实现设计

## Goal

设计 `flowforge uninstall` 命令的实现方案，包含二进制清理、配置删除、项目制品可选清理和确认机制。

## Decision

`flowforge uninstall` 默认可选清理：1) 删除 CLI 二进制 2) 删除 `~/.flowforge/` 配置和状态目录。`--project <path>` 额外清理目标项目托管文件。不自动修改 PATH/环境变量。提供 `--keep-config` 选项仅删除二进制。

## Rationale

- 显式确认和分层清理避免误删用户数据
- `--keep-config` 支持仅卸载二进制保留配置（用于重新安装场景）
- PATH 清理由用户负责，避免跨 shell/跨平台兼容性风险
- 卸载是低频操作，安全优先于便利

## Constraints

- 卸载前必须显示将要删除的文件列表
- `--yes` 跳过确认提示（用于自动化）
- 不删除非 FlowForge 的文件
- 项目制品清理仅删除已知部署目标（`.agents/skills/`、`.flowforge/templates/`、`AGENTS.md`）

## Impact

- 新增 `internal/command/uninstall.go` 实现 CLI 命令
- 新增 `internal/uninstall/cleaner.go` 实现清理逻辑
- 退出码区分：0=成功，1=一般错误

## Verification

- 完整卸载：删除二进制 + 删除 ~/.flowforge/ → 输出摘要
- 仅卸载二进制：`--keep-config` → 仅删除二进制 → 配置保留
- 项目清理：`--project <path>` → 删除目标项目托管文件 → 保留其他文件
- 确认提示：无 `--yes` 时列出文件并要求确认
- 幂等性：重复卸载不报错

## Follow-up Tasks

- 实现 uninstall 命令结构
- 实现 cleaner.go 清理逻辑
- 实现项目制品清理逻辑

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [REQ-CR26062102-djeu2wuow388](REQ-CR26062102-djeu2wuow388_cli-卸载命令.md) [requirement] - CLI 卸载命令

### Incoming

- [TASK-CR26062102-i-dji5lzhbe4fr](TASK-CR26062102-i-dji5lzhbe4fr_实现-flowforge-uninstall-命令-cleaner.md) [task] - 实现 flowforge uninstall 命令 — cleaner + 项目制品清理 + AGENTS.md 区块移除

