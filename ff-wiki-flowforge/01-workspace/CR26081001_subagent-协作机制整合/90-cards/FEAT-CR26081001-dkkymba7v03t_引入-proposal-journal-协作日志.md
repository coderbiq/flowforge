---
id: FEAT-CR26081001-dkkymba7v03t
title: 引入 Proposal Journal 协作日志
type: feature
status: done
importance: should
links:
    - target: PROP-CR26081001
      relation: belongs_to
created: 2026-08-10T03:56:24.195033505Z
updated: 2026-08-10T07:44:04.812439976Z
source: CR26081001
---

# 引入 Proposal Journal 协作日志

## Summary

为每个 Proposal 增加一份 `JOURNAL.md`，按时序追加用户决定、subagent 工作摘要、正式文档引用、阻塞和下一步。它用于跨宿主 session 恢复协作进度，但不保存完整模型输出，也不成为第二套业务状态系统。

## Motivation

宿主 Agent session 不稳定且平台相关，FEATURE 又只应保存正式工程事实。缺少轻量协作日志时，新会话必须重新扫描全部 Proposal 或依赖旧 session；若把完整 handoff 写入 FEATURE，则会污染设计文档并重复存储内容。

## Design

### Key Decisions

- 名称定为 Proposal Journal，文件名为 `JOURNAL.md`；因为它是提案生命周期内的协作日志，而不是一次 runtime run。
- 每个 Proposal 默认只有一份、单 Markdown 文件、仅追加的 Journal；因为目标是恢复 Proposal 进度，不是复刻每个宿主 session。
- 正式产出直接写入 FEATURE、DEC、FIND、代码或 Verification，Journal 只保存摘要和引用；因为要避免同一事实多份副本。
- 恢复时读取尾部最近记录，并按引用加载正式文档；首版不引入 Invocation、Handoff、Result、checkpoint 或独立状态机。

### Architecture

`JOURNAL.md` 位于 Proposal 工作目录根部，与 `90-cards/` 并列。CLI 提供追加和读取尾部记录的结构化操作，Agent 不直接维护时间戳和格式。每条记录至少包含时间、actor、摘要和 next，references、status、verification、blocked reason 按需出现。新会话先读取 Journal 尾部，再根据引用和下一步恢复上下文；宿主 session ID 与实际模型 ID 仅作可选诊断信息，不参与恢复正确性。

### Alternatives Considered

- 名称使用 RUN：拒绝，因为会诱导增加执行实例、授权和 session 拓扑职责。
- 每次宿主 session 一个日志：拒绝，因为会产生恢复选择和文档碎片。
- 保存完整 subagent 输出：拒绝，因为正式结果已写入 artifact，完整输出高噪声且高 token。
- 将日志放入 SQLite：拒绝，因为协作记录需要可读、可审计并随 Proposal 提交。

## Constraints

- Journal 只追加，不能替代 `card evolve`、`card steps`、Verification 或正式文档更新。
- Journal 不保存 chain-of-thought、完整 Prompt 或默认完整模型输出。
- Journal 格式必须允许普通 Markdown 阅读，恢复不得依赖宿主 session、provider 或仍存活的 child session。
- Journal 不进入 library，不参与普通卡片搜索和卡片生命周期。
- 追加失败必须显式返回错误；正式 artifact 与 Journal 冲突时，以 artifact 和 CLI 状态为准。

## Implementation Plan

### Step 1: 定义 Journal 存储与记录格式

<!-- step-status: done -->

- **Goal**: Proposal 创建后具有可解析、可追加的 `JOURNAL.md`。
- **Files**: `internal/core/` Proposal 存储代码、Proposal 模板和单元测试
- **Approach**: 定义最小 entry 结构并渲染为固定 Markdown；新 Proposal 创建文件，旧 Proposal 首次追加时惰性创建。
- **Edge Cases**: 旧 Proposal、空 Journal、人工尾部文本、非 ASCII 标题与引用。
- **Dependencies**: Proposal v3 工作目录稳定。
- **Parallel**: no
- **Verification**: 测试创建、惰性创建、追加顺序、时间格式和 Markdown 转义。

### Step 2: 增加 Journal CLI 操作

<!-- step-status: done -->

- **Goal**: Agent 可通过 CLI 追加记录并读取用于恢复的尾部条目。
- **Files**: `internal/command/` Journal 命令、命令测试、`assets/AGENTS.md`
- **Approach**: 提供最小 append/recent 操作，参数覆盖 actor、message、references、status、next 和 blocked reason。
- **Edge Cases**: Proposal 不存在、引用不存在、读取越界、并发 subagent 同时追加。
- **Dependencies**: Step 1 存储格式。
- **Parallel**: no
- **Verification**: 测试文本/JSON 输出、错误处理和并发追加不丢记录。

### Step 3: 接入协作与恢复工作流

<!-- step-status: done -->

- **Goal**: subagent 结束工作时追加 Journal，新会话从最近记录恢复。
- **Files**: `assets/skills/flowforge-*`、中立角色源、协议测试
- **Approach**: 正式结果先写 artifact 并通过 CLI，再追加摘要、引用与下一步。
- **Edge Cases**: artifact 更新后 Agent 中断，恢复时应以实际状态纠正过期 Journal。
- **Dependencies**: Step 2 和 subagent 路由模型。
- **Parallel**: no
- **Verification**: 场景覆盖设计交接、执行阻塞、用户决定、跨 session 恢复和 Proposal 完成。

## Verification

- 每个 Proposal 最多一份 `JOURNAL.md`，从最近记录和引用可定位当前 FEATURE、阻塞与下一步。
- Journal 不含正式设计或完整模型输出副本，删除宿主 session 后仍可恢复。
- Journal 与 artifact 冲突时以 artifact 为准，并能追加纠正记录。

## History

- 2026-08-10 [decision] 将早期名称 RUN 改为 Proposal Journal。
- 2026-08-10 [decision] 取消复杂状态机、Invocation、完整 handoff 和多文档拆分。
- 2026-08-10T13:49:24+08:00 | progress | 完成 Proposal Journal 存储层：新 Proposal 创建 JOURNAL.md，旧 Proposal 首次追加时惰性创建；补充创建、追加、读取尾部、格式规范化和损坏条目测试。go test ./internal/core/... 通过；go test ./internal/... 仍受既有 command 测试失败阻断。
- 2026-08-10T14:41:30+08:00 | progress | 完成 Journal CLI：新增 journal append/recent，支持当前或显式 Proposal、引用校验、文本/JSON 输出、limit、错误处理；增加并发追加保护和命令测试。聚焦 Journal 测试通过；go test ./internal/... 仍仅受既有 command 测试失败阻断。
- 2026-08-10T15:44:04+08:00 | progress | 完成 Journal 工作流接入：design、implement、feedback 和 proposal archive curate 流程均在开始时读取最近协作记录，并在正式 artifact 更新和验证后追加摘要；同步修复嵌入 assets 与根 assets 的 v2/v3 漂移。聚焦 Journal/初始化测试通过；go test ./internal/... 仍仅受既有 Proposal/Structure 测试失败阻断。

## Links

### Outgoing

- [PROP-CR26081001](../../../03-proposal/CR26081001_subagent-协作机制整合.md) [proposal] - Subagent 协作机制整合

### Incoming

#### requires
- [FEAT-CR26081001-dkkylwkup2d3](FEAT-CR26081001-dkkylwkup2d3_整合-proposal-journal-与-subagent-协作层.md) [feature] - 整合 Proposal Journal 与 Subagent 协作层
- [FEAT-CR26081001-dkkymgqx8c63](FEAT-CR26081001-dkkymgqx8c63_重构-subagent-路由与职责模型.md) [feature] - 重构 Subagent 路由与职责模型

## Open Questions

None

## Dependencies

- Proposal 工作目录和 v3 artifact 模型
