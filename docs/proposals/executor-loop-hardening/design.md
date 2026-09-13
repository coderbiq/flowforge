---
flowforge:
  schema: 1
  role: design
  id: executor-loop-hardening-design
  revision: 2
  consumes:
    requirements:
      executor-loop-hardening-requirements: 2
---

<a id="executor-loop-hardening-design"></a>
# 执行者循环硬化方案

依据：[执行者循环硬化需求](requirements.md#executor-loop-hardening-requirements)；一手研究 `docs/research/2026-09-13-agent-loop-hardening-prior-art.md`（Gemini CLI 循环检测、opencode `steps`/`doom_loop`、契约固化层级官方表述）；事故取证（tangram-v2 票 06 会话：893 调用/290 次同命令重试/86.4 分钟）。

## 方位：三个独立层，三张零边票

研究结论决定分层：预算断路在**宿主**（flowforge 编译产物声明，宿主执行）、契约固化在**编译层**（agent 定义 body，每会话注入）、重复失败信号在**CLI 事后**（evidence 门禁扩展）。三层无接口耦合，各自独立可验证——零 DAG 边并行。

## <a id="d-execution-budget"></a>d-execution-budget：执行预算通道

- `AgentsConfig` 新增 `MaxSteps int`（`yaml:"max_steps,omitempty"`）。解析约定：`0`（未配置）→ 采用默认 200；`-1` → 不设限（不产出字段）；`>0` → 原值。负值仅 `-1` 合法，其余负数报配置错误（不静默）。
- 默认值 200 的依据：事故离群会话 893 步；正常 flash 执行 51-149 步；Qwen-Code 默认软上限 100、绝对 10×；200 覆盖正常分布上界同时截断离群尾部。
- `CompileOptions` 新增 `MaxSteps *int`（指针区分"未设"与"0"）；`compile_opencode.go` frontmatter 增 `Steps *int yaml:"steps,omitempty"`。仅 OpenCode：Claude Code/Codex 编译器无对应字段能力（需求范围条款），不动。
- 作用域钉定 implementer（`def.Name == "flowforge-implementer"`），与 test guard 同一收窄点（agents.go resolveCompileOptions）；非 implementer 不受 `max_steps` 影响（验收 3）。
- deploy 与 status 共用既有 resolveCompileOptions helper，无新调用点。

## <a id="d-prompt-permanence"></a>d-prompt-permanence：契约固化

- `assets/subagents/flowforge-implementer.md` body 在 Identity 之后新增小节（示意标题 `## Non-negotiables`，英文，与定义 body 语言一致），内容为 SKILL.md 六项契约的固化摘要，锚点为 SKILL.md 英文原句逐字关键词（`at most 2 times`、`STATUS: BLOCKED`、`failed repair rounds`）：
  - fail-fast：任何验证命令至多运行 `at most 2 times`，第二次失败即以 `STATUS: BLOCKED` 上报——绝不做第三次相同尝试（此句是 fail-fast 条款的直接推论，非新规则）；
  - 修复上限：`5 failed repair rounds` 后停止并上报；
  - 预算收束：迭代预算将尽时总结已完成与剩余任务、以 STATUS 终态收束（与宿主 `steps` 到限行为语义一致）。
- 单一事实源约束：摘要不引入 SKILL.md 之外的新规则；`assets_deploy_test.go` 结构断言锁定三个英文锚点在 SKILL.md、定义源文件、部署产物三处同时存在，SKILL.md 措辞变更时测试同步——措辞漂移在 CI 可见。
- 不动 frontmatter 元数据（name/model_profile/permission 等已稳定）；deploy 流程零改动（body 变更天然随编译产物下发）。

## <a id="d-repeat-failure-gate"></a>d-repeat-failure-gate：重复失败诊断

- `catalog_evidence.go` 在既有四元组解析上聚合：对单张票全部已勾选 Change 的四元组，取规范化 cmd（去首尾空白、压连续空白）中 `exit != 0` 者，按 cmd 计数；任一 cmd 计数 ≥3 → 产生 `evidence-repeat-failure`（warning；strict 判失败）。
- 与既有 `evidence-exit-nonzero` 关系：后者管单条非零退出（勾选了失败证据），前者管**重复模式**（同命令反复失败仍继续/仍勾选）——两诊断可并存，语义不同。
- 复用 `EvidenceConfig.ExemptProposals` 豁免通道与 strict 语义，不新增配置面。
- 常量 `DiagEvidenceRepeatFailure` 入 `catalog.go` 诊断码族；`check.go` help 文案补一行。

## Standards clauses

- must evidence 四元组与新诊断保持人可读 Markdown 正文标记，执行者与 CLI 间不得引入传长文本的接口（源：`AGENTS.md` 核心设计原则 1，[Constraints]）。
- must 新诊断码为纯本地文件解析，无网络、无 LLM 调用，接入现有 check/frontier 诊断管线与 waiver/policy 语义（源：`AGENTS.md` 核心设计原则 2，[Constraints]）。
- must `assets/` 为权威源，变更经 `make dev` 同步 `internal/command/assets/` 双拷贝后构建（源：`docs/proposals/lightweight-execution-contract/spec.md` Further Notes，[Conventions]）。
- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- must `max_steps` 仅对 flowforge-implementer 生效，其余 subagent 编译产物零变化（源：本需求验收 3，[Constraints]）。

## 兼容与迁移

- 已部署项目升级后（无新配置）：implementer 产物差异 = 新增 `steps: 200` 行 + prompt 摘要小节（验收 7）；其余 subagent 产物不变。
- `max_steps: -1` 是显式退出开关（回到现状无限迭代），upgrade 提示不追加新行（evidence hint 已占位，避免提示噪音）。
- 新诊断默认 warning：存量含失败四元组的票不受阻断；strict 用户按既有豁免机制处理。

## 验证策略

- 预算通道：单元测试覆盖 0/-1/N 三态的编译产物（frontmatter 含/不含 steps）、resolveCompileOptions 的 implementer 收窄、deploy-status 一致性；向后兼容用既有三宿主回归断言。
- 契约固化：源文件与部署产物双端结构断言（关键词组），加 SKILL↔摘要一致性断言。
- 诊断门禁：构造 3 次同命令非零退出 → 报；2 次 → 不报；豁免名单 → 不报；与 exit-nonzero 并存用例。
