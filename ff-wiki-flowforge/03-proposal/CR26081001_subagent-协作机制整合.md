---
id: PROP-CR26081001
title: Subagent 协作机制整合
type: proposal
status: active
importance: should
links:
    - target: STR-CR26081001-REQ
      relation: indexes
created: 2026-08-10T11:55:44.981896943+08:00
updated: 2026-08-10T11:55:44.98194203+08:00
source: CR26081001
proposal_id: CR26081001
dir_name: CR26081001_subagent-协作机制整合
slug: subagent-协作机制整合
---

# Subagent 协作机制整合

## Purpose

将 agent-orchestration-kit 的有效设计重写并整合为 FlowForge 的 subagent 协作层，使需求分析所使用的角色、模型 profile、Skill、文档门禁与宿主部署形成一个统一工作流。

## Entries

- [STR-CR26081001-REQ](STR-CR26081001-REQ.md) (structure, active) - Requirement index

## Summary

本提案以 FlowForge 为唯一 artifact 和生命周期事实源。每个 Proposal 使用一份 `JOURNAL.md` 记录协作摘要与文档引用；Coordinator、Design Analyst、Executor 和按风险调用的 Reviewer 通过 Journal 与正式 artifact 交接；OpenCode/Codex adapter 从中立策略生成，并复用 FlowForge CLI 安装升级。

## Feature Map

- `FEAT-CR26081001-dkkylwkup2d3`：定义总体架构、AOK 迁移边界和整合顺序。
- `FEAT-CR26081001-dkkymba7v03t`：定义 Proposal Journal 的最小日志语义和恢复流程。
- `FEAT-CR26081001-dkkymgqx8c63`：定义角色拓扑、路由、Skill、模型 profile 与门禁映射。
- `FEAT-CR26081001-dkkymf4v3m1u`：定义宿主 adapter 生成以及安装升级生命周期。

## Architecture Overview

正式设计与状态保存在 FEATURE、DEC、FIND、Step 和 Verification；Journal 仅追加协作摘要、引用和下一步；中立 orchestration 策略负责角色与路由；宿主 adapter 负责模型、session、工具权限与 sandbox。四层之间保持单向依赖，FlowForge 不实现模型 runtime。

## Key Constraints

- 不原样迁移 AOK 固定四角色、完整 handoff 或独立 installer。
- 不把 Journal 发展为第二套工作流状态机或完整 session archive。
- 不把 Prompt 规则描述成 runtime 硬授权。
- 未显式安装 adapter 时，FlowForge 现有单 Agent 工作流必须保持可用。
