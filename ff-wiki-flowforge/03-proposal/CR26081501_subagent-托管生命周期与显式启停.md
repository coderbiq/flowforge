---
id: PROP-CR26081501
title: Subagent 托管生命周期与显式启停
type: proposal
status: active
importance: should
links:
    - target: STR-CR26081501-REQ
      relation: indexes
created: 2026-08-15T11:01:56.192741+08:00
updated: 2026-08-15T11:01:56.192742+08:00
source: CR26081501
proposal_id: CR26081501
dir_name: CR26081501_subagent-托管生命周期与显式启停
slug: subagent-托管生命周期与显式启停
---

# Subagent 托管生命周期与显式启停

## Purpose

为 FlowForge 设计 subagent 的显式托管生命周期：新项目默认 `non_subagent`，宿主检测只用于 `status`；只有 `subagent enable` 才启用，`disable` 负责按 manifest 登记项备份并卸载。本文 Proposal 是独立设计面，不修改或重开 `CR26081401`、`CR26081001`。

## Entries

- [STR-CR26081501-REQ](../01-workspace/CR26081501_subagent-托管生命周期与显式启停/STR-CR26081501-REQ.md) (structure, active) - Requirement index

## Summary

本 Proposal 将现有 manifest/hash/marker/renderer 能力收敛为 schema v2 与 host intent，并明确 OpenCode/Codex 文件集合差异。命令面新增 `subagent enable|disable|status`，保留 `sync` 作为按 manifest 同步；disable 无论 hash 是否变化都先备份到 `.flowforge/backups/subagent-disable/<timestamp>/` 后删除登记的 subagent 文件，未登记文件绝不删除。`AGENTS.md` 仅移除 orchestration managed block，保留基础 FLOWFORGE block、用户内容和其他工具 block。

## Non-goals

- 不修改或重开 `CR26081401`、`CR26081001`。
- 不让 host detection 自动启用，不在 non_subagent 模式生成 host files/block。
- 不接管 manifest 未登记文件，不覆盖用户配置，不在本设计阶段修改产品代码。

## Feature Map

- `FEAT-CR26081501-003`：定义 manifest v2、v1 读取迁移、host intent 与 Codex/OpenCode renderer/file set。
- `FEAT-CR26081501-004`：定义 `subagent enable|disable|status`，显式卸载授权、备份/dry-run、AGENTS block 保留与幂等。
- `FEAT-CR26081501-002`：收敛 `sync`、upgrade、uninstall 的 intent-driven 生命周期、冲突和 manifest 删除边界。
- `FEAT-CR26081501-005`：建立 schema/双宿主/命令/生命周期/文档的回归测试矩阵与 CLI help 验收。

## Design Gates

- REQ `REQ-CR26081501-001` 是四个 FEATURE 的共同实现追踪入口，已由 `STR-CR26081501-REQ` 索引。
- 执行顺序：003 → 004 → 002 → 005；每个 FEATURE 必须先通过 `designed`、再通过 `planned`，实施需另有明确用户意图和 preflight allow。
- 设计完成验收：`proposal inspect` 无 error/warn、四个 FEATURE 全部 planned、REQ/FEATURE/STR links 可追踪、validate 与 Journal 结果已记录。
