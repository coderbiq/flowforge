---
id: FEAT-CR26081001-dkkymf4v3m1u
title: 生成并部署宿主 Subagent Adapter
type: feature
status: done
importance: should
links:
    - target: FEAT-CR26081001-dkkymgqx8c63
      relation: requires
    - target: PROP-CR26081001
      relation: belongs_to
created: 2026-08-10T03:56:32.578399748Z
updated: 2026-08-10T13:05:33.520882047Z
source: CR26081001
---

# 生成并部署宿主 Subagent Adapter

## Summary

从 FlowForge 中立 subagent 策略生成 OpenCode/Codex 原生配置，并复用 managed assets manifest 完成项目级安装、升级、冲突检测和卸载。首版以 OpenCode 显式 opt-in 为目标，不覆盖用户已有默认 Agent 或模型配置。

## Motivation

AOK 已证明中立源生成 adapter 可行，但没有 installer；FlowForge 已有资产 hash、冲突与升级生命周期。另建安装器会重复能力，手工复制又无法安全升级或保留用户配置。

## Design

### Key Decisions

- adapter 从中立角色策略生成，生成物不作为手工维护事实源；因为要避免宿主 Prompt 和权限漂移。
- 复用 `ProjectManifest` 管理 Agent 文件和托管配置块；因为已有 SHA256、冲突与卸载语义。
- 首版仅支持 project scope OpenCode 且显式安装；因为该运行路径验证更充分。
- 不默认修改 `default_agent`，不写 provider/model ID；因为这些属于用户宿主环境。
- 无法安全合并配置时停止并展示差异，不静默覆盖。

### Architecture

中立角色源经 renderer 生成 `.opencode/agents/*.md` 和必要配置 fragment。CLI 按显式 target 部署，并把托管文件或块记录进 `.flowforge/manifest.yaml`。升级根据旧 hash 分类 added、updated、conflict、removed；卸载只移除仍与托管 hash 一致的内容。Codex 使用同一中立定义，但首版只保留兼容设计，不阻塞 OpenCode 交付。AOK smoke 的版本探测、隔离目录、session topology 和 filesystem diff 方法迁入兼容测试。

### Alternatives Considered

- 固定生成物直接放入 `assets/`：拒绝，因为角色启用和宿主版本需按环境渲染。
- 建设 AOK 独立 installer：拒绝，因为重复 FlowForge 生命周期。
- 首版同时支持双宿主和 global scope：拒绝，因为测试矩阵过大。
- 自动设置默认 manager：拒绝，因为会改变用户现有行为。

## Constraints

- Adapter 显式 opt-in，未启用不得部署 `.opencode/` 或 `.codex/` 内容。
- 不覆盖同名 Agent、默认 Agent、provider 或 model 配置。
- 所有部署生成物进入 managed asset/manifest 生命周期。
- 中立输入不含绝对路径、密钥、provider 凭证或实际模型 ID。
- 宿主不兼容时显式停止或降级，卸载不得删除用户修改内容。

## Implementation Plan

### Step 1: 实现中立策略 renderer

<!-- step-status: done -->

- **Goal**: 稳定生成 OpenCode Agent 文件和配置 fragment。
- **Files**: `internal/` renderer、模板、golden tests
- **Approach**: 使用确定性渲染函数，嵌入 source hash 和版本并拒绝未知能力。
- **Edge Cases**: 禁用角色、profile 未映射、YAML 转义。
- **Dependencies**: Subagent 中立 schema。
- **Parallel**: no
- **Verification**: golden、重复渲染一致性和 drift 检查。

### Step 2: 接入 managed asset 生命周期

<!-- step-status: done -->

- **Goal**: CLI 可预览、安装、更新和卸载 OpenCode adapter。
- **Files**: `internal/command/`、`internal/core/project_manifest.go`、部署测试
- **Approach**: manifest 支持可选生成资产和托管配置块，沿用 added/updated/conflict/removed 语义。
- **Edge Cases**: 同名文件、用户修改、JSON 格式差异、旧 adapter 残留。
- **Dependencies**: Step 1 renderer。
- **Parallel**: no
- **Verification**: 测试 install、幂等 update、conflict、preserve 和 uninstall。

### Step 3: 迁移宿主兼容 smoke tests

<!-- step-status: done -->

- **Goal**: 验证受支持 OpenCode 版本中的角色、权限、委托和 session 行为。
- **Files**: `tests/` 或 `evals/`、兼容说明、可选 CI
- **Approach**: 移植静态解析、隔离运行、事件拓扑和 filesystem diff；模型测试显式门禁。
- **Edge Cases**: event schema 变化、缺少凭证、不可观测导致 inconclusive。
- **Dependencies**: Step 2 adapter。
- **Parallel**: no
- **Verification**: 静态测试默认运行，模型测试区分 pass/fail/inconclusive 并清理。

## Verification

- 同一中立源生成的职责、Skill 和权限无漂移。
- 默认 init/upgrade 不安装未选择 adapter。
- 安装不覆盖 default agent、模型或同名 Agent。
- 升级区分托管文件和冲突，卸载保留用户内容。
- OpenCode smoke 覆盖 delegation depth、worker task deny、写权限和 active role。

## History

- 2026-08-10 [decision] AOK installer 不迁移，生命周期复用 FlowForge CLI。
- 2026-08-10 [decision] 首版交付项目级、显式 opt-in OpenCode adapter。
- 2026-08-10T21:05:33+08:00 | progress | 完成 OpenCode adapter：从中立 policy 渲染 FlowForge 命名 Agent 文件，显式 install/update/uninstall 使用 manifest hash 保护用户修改，不创建或合并 OpenCode 配置；增加静态 smoke。

## Links

### Outgoing

- [PROP-CR26081001](../../../03-proposal/CR26081001_subagent-协作机制整合.md) [proposal] - Subagent 协作机制整合
- [FEAT-CR26081001-dkkymgqx8c63](FEAT-CR26081001-dkkymgqx8c63_重构-subagent-路由与职责模型.md) [feature] - 重构 Subagent 路由与职责模型

### Incoming

- [FEAT-CR26081001-dkkylwkup2d3](FEAT-CR26081001-dkkylwkup2d3_整合-proposal-journal-与-subagent-协作层.md) [feature] - 整合 Proposal Journal 与 Subagent 协作层

## Open Questions

None

## Dependencies

- Subagent 路由与职责中立模型
- FlowForge managed asset manifest
