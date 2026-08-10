---
id: FEAT-CR26081001-dklbdy0lpa0t
title: 修复 OpenCode Subagent Prompt 与角色绑定
type: feature
status: done
importance: should
links:
    - target: PROP-CR26081001
      relation: belongs_to
created: 2026-08-10T13:56:42.830931634Z
updated: 2026-08-10T23:56:23.645296727+08:00
source: CR26081001
---

# 修复 OpenCode Subagent Prompt 与角色绑定

## Summary

修复 v3.1.0 subagent 仅生成骨架 Prompt、安装入口割裂和 AGENTS.md 未包含路由规则的问题。FlowForge 将自动探测 OpenCode/Codex，通过统一 `flowforge sync` 协调 Skills、完整角色 Prompt、宿主 Agent 文件、动态路由区块和 manifest。

## Motivation

当前 renderer 未嵌入角色协议或 Skill 工作流，生成的 subagent 无法可靠执行；adapter 旁路写 manifest，普通 assets update 会丢失管理记录；init/upgrade 不自动安装；AGENTS.md 只展示长安装命令而不引导主 Agent 调度。继续保留会让已发布能力名义存在但实际不可用。

## Design

### Key Decisions

- 正常用户入口收敛为 `flowforge init`、`flowforge sync`、`flowforge upgrade`；删除 `assets adapter` 命令，因为 adapter 是项目设施而非用户需理解的资产类别。
- `sync` 统一生成静态 assets、宿主 Agent 与 AGENTS.md orchestration 区块并保存完整 manifest，因为旁路追加会丢失管理关系。
- OpenCode 存在 `.opencode/`、`opencode.json` 或 `opencode.jsonc` 时自动启用；Codex 存在 `.codex/` 或 `.codex/config.toml` 时自动启用。
- manifest 中已管理的 host 即使当前探测不到仍继续更新，避免升级时误删；显式 `--without-host` 才卸载。
- Agent Prompt 必须组合共享规则、完整角色协议、Skill 执行契约、门禁、停止状态和输出格式，而不是仅引用 Skill 名称。
- AGENTS.md 使用独立 orchestration 托管区块，根据实际已安装 host/Agent 动态生成。

### Architecture

`BuildDesiredProjectState` 生成所有托管文件及内容，包含基础 assets、OpenCode/Codex Agent 文件和两个 AGENTS.md 区块；统一 reconcile 根据旧 manifest 与磁盘 hash 执行 added/updated/conflict/removed。Host detector 合并 manifest 已启用 host、项目探测和命令覆盖。OpenCode 生成 primary Coordinator 与 worker；Codex 生成 worker TOML 和协调章程。v3.1.0 遗留的 FlowForge Agent 文件仅在匹配已知历史 skeleton hash 时自动接管。

### Alternatives Considered

- 继续增强 `assets adapter`：拒绝，命令暴露内部结构且无法解决 init/upgrade/manifest 分裂。
- 只扩充 Prompt：拒绝，安装升级和 AGENTS.md 仍不可用。
- 自动修改 default_agent/provider/model：拒绝，属于用户宿主配置，仍需用户显式处理。

## Constraints

- 不覆盖未知或用户修改的 Agent 文件，不深度合并 OpenCode/Codex 配置。
- init 无宿主时只安装基础设施并提示后续运行 `flowforge sync`。
- Codex adapter 必须与 OpenCode 共享中立角色协议，宿主差异只存在于 renderer。
- Prompt 与生成物必须有 golden/静态解析测试；模型 smoke 可显式门禁但不阻塞普通测试。
- `assets adapter` 命令直接删除，不保留兼容 alias。

## Implementation Plan

### Step 1: 完整角色 Prompt 与双宿主 renderer

<!-- step-status: done -->

- **Goal**: OpenCode/Codex 生成物包含可执行角色协议、Skill 契约、门禁、停止和输出规则。
- **Files**: `internal/orchestration/` renderer、角色源、golden tests
- **Approach**: 以中立 policy 和角色 Markdown 为源，组合共享 workflow charter；OpenCode 生成 Coordinator/worker Markdown，Codex 生成 worker TOML/Prompt。
- **Edge Cases**: Reviewer 禁用、YAML/TOML 转义、Skill 或角色源缺失。
- **Dependencies**: 现有中立 policy。
- **Parallel**: no
- **Verification**: golden 和静态解析验证完整关键段落、角色权限和无 provider/model ID。

### Step 2: 统一 host detection 与 sync 生命周期

<!-- step-status: done -->

- **Goal**: init/sync/upgrade 自动协调基础和宿主设施，manifest 不再丢失 adapter。
- **Files**: `internal/command/sync.go`、managed state/reconcile、init/upgrade、迁移测试
- **Approach**: detector 合并已管理 host 与项目证据；统一 desired state/reconcile；识别并接管已知 v3.1.0 skeleton。
- **Edge Cases**: 双宿主、探测文件消失、用户修改托管文件、未知旧 skeleton、无 host。
- **Dependencies**: Step 1 renderer。
- **Parallel**: no
- **Verification**: init 自动安装、sync 幂等、upgrade 保持 host、冲突保留、显式卸载和 v3.1.0 迁移测试。

### Step 3: 动态 AGENTS.md 路由与 CLI 收敛

<!-- step-status: done -->

- **Goal**: 主 Agent 获得准确的 subagent 调度规则，用户只使用短入口。
- **Files**: AGENTS block 核心逻辑、`assets/AGENTS.md`、root/init/upgrade 命令和测试
- **Approach**: 基础区块与 orchestration 区块分离；sync 按安装 host/Agent 生成路由；删除 assets adapter 命令并更新提示。
- **Edge Cases**: 无 host 时移除 orchestration 区块但保留基础规则；用户自有 AGENTS 内容保持不变。
- **Dependencies**: Step 2 sync。
- **Parallel**: no
- **Verification**: AGENTS 保留用户内容、区块增删更新、CLI help 不再出现废弃命令。

## Verification

- OpenCode/Codex Prompt 包含完整角色和工作流契约，Coordinator 能路由 worker。
- `flowforge init` 在宿主项目自动安装，`flowforge upgrade` 自动同步，`flowforge sync` 可修复状态。
- manifest 始终保留基础和 host 资产，用户修改冲突不被覆盖。
- AGENTS.md 动态列出当前可用 Agent 和 preflight/risk-review/Journal 调度规则。
- 旧 `assets adapter` 命令不存在。

## History

- 2026-08-10 [bug] v3.1.0 renderer 只输出一句通用 Prompt，未嵌入角色或 Skill 协议。
- 2026-08-10 [bug] adapter 安装与 assets update 分裂，manifest 管理关系会丢失。
- 2026-08-10 [decision] 删除旧命令，统一由自动探测和 sync 管理。
- 2026-08-10T23:03:22+08:00 | progress | 完成 OpenCode/Codex 完整 Prompt、Coordinator、角色权限映射、结果契约与静态解析测试。
- 2026-08-10T23:03:22+08:00 | progress | 完成统一 sync 生命周期、host 探测、持久禁用、dry-run/adopt、manifest 冲突保留与 v3.1.0 精确接管。
- 2026-08-10T23:04:44+08:00 | progress | 完成动态 AGENTS.md orchestration 区块、init/upgrade 自动同步，并彻底删除 assets adapter、assets update、skill update 命令入口。
- 2026-08-10T23:16:36+08:00 | progress | 同步升级已保留受管 OpenCode/Codex subagent 的显式 model 配置，其余生成内容随新版本替换；新增两种格式的回归测试。
- 2026-08-10T23:35:35+08:00 | bug | v3.1.x 升级器在替换为 v3.2.0 后仍调用 assets update，因新二进制删除该入口导致后置 sync 失败；新增隐藏跨版本兼容桥接并转发到 sync。
- 2026-08-10T23:47:12+08:00 | bug | v3.2.1 对旧 OpenCode frontmatter 的严格 YAML 解析仍会阻断 sync，且 v3.2.0 manifest 缺失的已生成 Agent 被误报未管理；改为容错 model 提取并自动接管带 FlowForge 生成签名的 Agent。
- 2026-08-10T23:56:23+08:00 | bug | 真实项目发现基础 FlowForge 区块与独立 orchestration 区块被 Beads 内容分隔，导致 AGENTS.md 规则视觉和结构混排；已将 orchestration 合并进 FLOWFORGE 主区块并保留用户/Beads 内容及冲突修改。

## Open Questions

None

## Dependencies

- `FEAT-CR26081001-dkkymf4v3m1u`
