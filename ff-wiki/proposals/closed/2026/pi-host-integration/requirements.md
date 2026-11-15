---
flowforge:
  schema: 1
  role: requirement
  id: pi-host-integration-requirements
  revision: 1
---

<a id="pi-host-integration-requirements"></a>
# PI 宿主集成需求

## 问题

FlowForge 的子代理部署体系（`flowforge agents deploy`）目前覆盖 opencode、Claude Code、Codex 三个宿主。PI（pi.dev）作为第四个宿主已在实际使用（本仓库即在 PI 会话中开发），但其能力模型不同：PI 核心不含子代理，由 pi-subagents 扩展提供（项目级定义在 `.pi/agents/*.md`）；PI 的项目级 skill 发现原生包含 `.agents/skills/`；编辑保护（test 文件 glob deny）在 pi-subagents frontmatter 中无对应字段（调研见 `docs/research/2026-09-19-pi-host-integration.md`）。当前用户在 PI 下只能靠正文提示获得弱化的 FlowForge 约束，且 `flowforge frontier` 的批次只能靠主会话逐个口头委派。

## 目标

1. **PI 成为第四个一等部署宿主**：`agents.hosts` 配置支持 `pi`；`deploy`/`remove`/`clean` 覆盖 `.pi/agents/` 目录，产物为 pi-subagents 可发现的原生 agent 文件；模型档位、只读权限、skill 绑定在 PI 下有等价或更强的表达。
2. **frontier 批次的运行时编排**：`flowforge frontier` 能输出一份 pi-subagents workflowScript，把"一 ticket 一个新鲜执行上下文"的执行单元政策从提示约定升级为运行时保证：每个 ready ticket 由一个 fresh-context 的 `flowforge-implementer` 子代理执行，子代理结束后由 host gate 运行结构校验。
3. **主机级 test 文件保护与原生工具**：项目级 pi extension 拦截对受保护 test 文件的写/编辑（glob 来源复用 `agents.test_file_globs` 单一真相），并把 `frontier`/`check` 注册为 LLM 原生工具。

## 范围与约束

- 不做 npm/pi-package 分发（第四步明确推迟，extension 与 agent 文件由仓库内 `agents deploy` 交付）。
- 编排首版**顺序执行**，不做 worktree 并行与合并回流（frontier ticket DAG 无阻塞 ≠ 无文件冲突，并行执行策略留待后续提案）。
- gate 仅运行 `flowforge check`（proposal 结构校验）；不把 ticket 正文的 "Expected tests" 结构化为 gate 命令（`Issue` 无机器可读验证命令字段，结构化属独立提案）。
- PI 的 AGENTS.md 与 `.agents/skills/` 兼容性为零改动项，不产生代码变更，仅作为设计依据。
- 不修改既有宿主的编译产物与接口签名（新增 flag/枚举值除外）。

## 可观察验收

1. 启用 `pi` 宿主并运行 `flowforge agents deploy` 后，`.pi/agents/` 下出现与定义一致的原生 agent 文件，且包含模型档位（`thinking`）、只读角色的工具白名单、skill 绑定（`skills`）字段；将 `pi` 移出启用集合并再次部署后，`.pi/agents/` 中受管文件被清理。
2. `flowforge frontier --pi-workflow` 输出的脚本文本符合 pi-subagents workflowScript 语法约束（顶层语句、显式 return、无嵌套函数声明），每个 ready ticket 恰好对应一个 `runs.run` 调用并携带 gate 命令；frontier 为空时输出明确的空批次提示。
3. PI 会话中加载项目 extension 后：对匹配 `agents.test_file_globs` 的文件执行 write/edit 被 block 并返回原因；LLM 可调用 `flowforge_frontier`/`flowforge_check` 工具并得到 CLI 输出。
