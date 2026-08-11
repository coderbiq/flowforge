---
id: FEAT-CR26081101-dklt4yxctvsf
title: 重构 Subagent 角色拓扑、Prompt 与宿主适配
type: feature
status: done
importance: should
links:
    - target: DEC-CR26081101-dkltaxm2rgcf
      relation: references
    - target: FEAT-CR26081101-dklt4yui88kf
      relation: requires
    - target: PROP-CR26081101
      relation: belongs_to
created: 2026-08-11T03:51:19.512273666Z
updated: 2026-08-11T06:14:02.239310439Z
source: CR26081101
---

# 重构 Subagent 角色拓扑、Prompt 与宿主适配

## Summary

建立由轻量 Coordinator、旗舰 Design Analyst、通用 Investigator、Executor 和按风险启用 Reviewer 组成的稳定角色拓扑，并为 Codex/OpenCode 生成与职责一致的 Prompt、权限和模型 profile。

## Motivation

当前 Coordinator 把调查和设计全部路由给昂贵 Analyst，Analyst 又承担资料读取和综合；Prompt 还存在向用户追问与 Worker 禁止对话、Artifact 写入与宿主权限不一致等冲突。

## Design

### Key Decisions

- **Coordinator 是执行型调度器**：使用低成本模型，只执行 CLI 判定、用户交互和异常升级，因为它不应重新规划分析。
- **Design Analyst 是规划与综合者**：使用高能力模型，负责 FEATURE、Working Design、plan revision、证据采纳和 stage readiness。
- **Investigator 是通用调查者**：按 brief 使用代码、Library、网络或日志，不按数据源拆角色。
- **保持委托深度 1**：只有 Coordinator 可委派，因为两个宿主对嵌套委派支持不同且用户需要可见性。
- **Reviewer 独立且风险触发**：避免原 Analyst 自审产生系统性盲点。
- **模型 profile 与真实模型解耦**：策略定义能力/成本，adapter 或用户配置映射具体模型。

### Architecture

权限矩阵：Coordinator 可对话/委派但不写 Artifact；Design Analyst 可写设计 Artifact 但不改产品、不委派；Investigator 只读产品并只可编辑 work item 指定的 FIND 和 Journal return event；Executor 写产品并验证；Reviewer 只读语义复核。

Prompt 分三层：共享短章程声明权威数据、单层委托、Journal 恢复和停止规则；稳定角色契约声明唯一职责、允许/禁止动作和结果 schema；动态 brief 只包含 Proposal/FEATURE、cycle/revision/work ID、问题、范围、sources、证据、budget、done_when 和唯一可写 Artifact。

Coordinator 循环固定为：读取 `analysis status/ready/reentry` → 向用户展示后台动作 → 调用指定角色 → seal/登记结果 → 再次查询。除一次可恢复宿主错误重试外，不新增调查项、不解释证据。

外部网络默认关闭，只有 work item 显式声明 `sources: external` 后由宿主授权。Codex/OpenCode 都不依赖嵌套 subagent；无法硬限制的能力必须在 sync/status 标记 `enforcement: soft|unsupported`。

### Alternatives Considered

- 强 Coordinator：与低成本模型选择冲突。
- Design Analyst 委派：兼容性和可见性差。
- 每 Skill/工具一个角色：角色数量会失控。
- Investigator 修改正式设计：破坏证据收集与设计判断分离。

## Constraints

- 恰好一个可与用户交互的 Coordinator。
- Worker 不得委派或直接提问；用户决策由 Coordinator 转发。
- Artifact Write 与 Product Write 继续隔离。
- 外部网络必须由 work item 显式授权；不可用时返回 BLOCKED。
- sync 保留用户显式 model 和权限配置。

## Implementation Plan

### Step 1: 扩展中立角色策略

<!-- step-status: done -->

- **Goal**: 新增 Investigator 并表达规划权和执行型调度约束。
- **Files**: `internal/orchestration/config.go`, `internal/orchestration/config_test.go`
- **Approach**: 增加角色/profile/能力矩阵；维持仅 Coordinator 可委派；校验写域隔离。
- **Edge Cases**: 多 Coordinator、未知 Skill、Investigator 获得产品写权限、Reviewer disabled。
- **Dependencies**: FEAT-CR26081101-dklt4yui88kf。
- **Parallel**: yes。
- **Verification**: policy invariant 和默认拓扑测试通过。

### Step 2: 重写角色协议与结果契约

<!-- step-status: done -->

- **Goal**: 消除用户追问、CLI 权限和职责边界冲突。
- **Files**: `internal/orchestration/roles/*.md`, `internal/orchestration/render.go`
- **Approach**: 实现共享章程+角色契约；增加 INCONCLUSIVE、EVIDENCE_CONFLICT、USER_DECISION_REQUIRED。
- **Edge Cases**: 范围扩张、无网络、非法返回格式、用户决定暂停。
- **Dependencies**: Step 1 和 Journal JSON 契约。
- **Parallel**: no。
- **Verification**: prompt golden 包含职责、禁止项、re-entry 和输出 schema。

### Step 3: 更新 Codex/OpenCode Adapter

<!-- step-status: done -->

- **Goal**: 生成权限与角色实际工作一致，并报告 enforcement 等级。
- **Files**: `internal/orchestration/render.go`, `internal/orchestration/render_test.go`, `internal/command/sync.go`, `internal/command/sync_test.go`
- **Approach**: Analyst 放行 Artifact/FlowForge CLI；Investigator 只读产品并定向写 FIND；保留用户 model。
- **Edge Cases**: 宿主能力缺失、重复 sync、部分角色安装、自定义权限。
- **Dependencies**: Step 2。
- **Parallel**: no。
- **Verification**: 双宿主配置快照、幂等同步和用户配置保留测试通过。

### Step 4: 更新部署编排说明

<!-- step-status: done -->

- **Goal**: 目标项目清楚声明强 Analyst、弱 Coordinator 和单层委托。
- **Files**: `assets/AGENTS.md`, `internal/command/sync.go`, 相关测试
- **Approach**: 更新 FLOWFORGE managed block，要求委派前读取结构化分析状态并公开后台动作。
- **Edge Cases**: 旧 block、非 Codex/OpenCode、只安装部分角色。
- **Dependencies**: Step 3。
- **Parallel**: yes。
- **Verification**: init/sync 部署内容和包裹边界测试通过。

## Verification

- 简单请求最多调用一次 Analyst。
- 复杂请求由 Analyst 规划，Coordinator 可机械调度并按 re-entry 再调用 Analyst。
- 两个宿主均不依赖嵌套委派。
- Prompt 与文件/CLI 权限一致。
- 用户可观察后台调查的派发、返回和重新综合。

## History

- 2026-08-11：依据轻量 Coordinator 约束重定义角色权力和宿主适配。
- 2026-08-11：重新核对正文直接编辑与 CLI 结构化边界决策，角色权限与 DEC 一致。
- 2026-08-11T13:59:56+08:00 | progress | Step 1 完成：新增 Investigator、规划/调度/证据写能力与只读工具 profile，并补充单层委托和写域隔离校验；配置测试通过。
- 2026-08-11T14:01:54+08:00 | progress | Step 2 完成：共享章程和角色契约明确弱 Coordinator、强 Analyst、通用 Investigator、单层委托、re-entry 与外部网络边界，并加入 INCONCLUSIVE/EVIDENCE_CONFLICT/USER_DECISION_REQUIRED 结构化状态；orchestration 测试通过。
- 2026-08-11T14:08:33+08:00 | progress | Step 3 完成：Codex/OpenCode adapter 对齐 Analyst、Investigator、Executor 权限；sync 保留显式 model/权限配置并报告 hard/soft/unsupported enforcement；双宿主 render 与 sync 定向测试通过。
- 2026-08-11T14:11:28+08:00 | progress | Step 4 完成：部署 AGENTS 与动态 managed block 声明强 Analyst、弱 Coordinator、结构化分析状态、用户可见后台动作、通用 Investigator、单层委托和外部来源授权；init/sync 定向测试通过。
- 2026-08-11T14:12:29+08:00 | progress | Verification：go test ./internal/orchestration 与 go test ./internal/command -run TestSync 通过；go test ./internal/... 已运行但受并行分支中的 proposal/structure 断言失败和 sandbox 禁止 httptest 监听端口影响未全绿；flowforge validate all 已运行，25 个历史卡片链接错误与本 FEATURE 无关；risk-review=review_not_required。
- 2026-08-11T14:14:02+08:00 | blocked | PLAN_STALE：四个 Step 均已 done，但 card evolve 的实现要求 planned→in_progress→done，同时 CLI stage parser 拒绝 in_progress，导致无法通过结构化命令推进 FEATURE done；未直接编辑结构化状态。

## Links

### Outgoing

- [PROP-CR26081101](../../../03-proposal/CR26081101_复杂需求多代理分析编排.md) [proposal] - 复杂需求多代理分析编排
- [DEC-CR26081101-dkltaxm2rgcf](DEC-CR26081101-dkltaxm2rgcf_正文直接编辑与-cli-结构化操作边界.md) [decision] - 正文直接编辑与 CLI 结构化操作边界
- [FEAT-CR26081101-dklt4yui88kf](FEAT-CR26081101-dklt4yui88kf_定义复杂需求的迭代分析协议与-artifact-边界.md) [feature] - 定义复杂需求的迭代分析协议与 Artifact 边界

### Incoming

- [FEAT-CR26081101-dklt4yyrcz3r](FEAT-CR26081101-dklt4yyrcz3r_重构-flowforge-design-方法论与端到端评估.md) [feature] - 重构 flowforge-design 方法论与端到端评估

## Open Questions

None。

## Dependencies

- FEAT-CR26081101-dklt4yui88kf：分析循环。
- FEAT-CR26081101-dklt4yvz37gy：确定性调度接口。
