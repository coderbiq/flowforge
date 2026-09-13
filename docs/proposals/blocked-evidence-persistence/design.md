---
flowforge:
  schema: 1
  role: design
  id: blocked-evidence-persistence-design
  revision: 1
  consumes:
    requirements:
      blocked-evidence-persistence-requirements: 1
---

<a id="blocked-evidence-persistence-design"></a>
# BLOCKED 证据工件化方案

依据：[BLOCKED 证据工件化需求](requirements.md#blocked-evidence-persistence-requirements)；executor-loop-hardening 交付后的缺口分析（BLOCKED 会话返回不落工件，重派发上下文失忆）；tangram-v2 票 06 事故（290 次重试 = 无失败记忆的必然结果）。

## 方位：一个状态语义，一个写入契约，一个消费循环

核心是一个**瞬态票内小节** `## Blocked evidence`：执行者写入（写入契约）、tracker 识别（状态语义）、refine 消费（闭环）。三张零边票各自独立可验证。不新增 CLI 面——诊断复用 evidence 族管线，frontier 归类复用 classifyReady 既有 warning 通道。

## <a id="d-blocked-artifact"></a>d-blocked-artifact：状态语义与诊断

- 解析层：`catalog.go`（或 catalog_semantics.go，与 hasCompletionEvidence 同层）新增标题匹配 `(?mi)^##\s+Blocked evidence\s*$`——open 票（Status 解析为 open）命中即产出 `blocked-evidence-present`（severity warning，源指向票文件）。
- 消息措辞面向派发方："open ticket carries blocked evidence; run flowforge-refine-ticket to consume it before redispatch"。
- strict 语义：warning → `--strict` 判失败（与 evidence 族一致；阻断状态需要人介入，strict 失败是正确信号）。
- 清除条件天然成立：小节移除或 Status 转 closed 后解析不再命中，无状态迁移代码。
- frontier 零改动：classifyReady 按诊断 severity 归类，warning → ready_with_warnings 既有桶（验收 4 复用现机制，需测试锁定防回归）。
- 小节位置约定：执行者追加在票末尾（Implementation note 之后）；解析不依赖位置（标题匹配），位置约定只为 diff 可读性。

## <a id="d-blocked-write"></a>d-blocked-write：执行者写入契约

- `assets/skills/flowforge-implement/SKILL.md`：BLOCKED 出口条款（既有六处 BLOCKED 语义不改）统一增补一句：返回 `STATUS: BLOCKED` 前，先在票末尾追加 `## Blocked evidence` 小节，内容三要素——逐字错误（verbatim command output）、已尝试命令与退出码列表、下一步假设（若有）。双通道声明：会话返回照旧，工件为权威。
- `assets/subagents/flowforge-implementer.md` Non-negotiables 增第四句 block-and-record（英文，锚点 `## Blocked evidence`）："When you must return STATUS: BLOCKED, first append a `## Blocked evidence` section to the ticket (verbatim error, commands tried with exit codes, next hypothesis)."
- 既有三锚点（`at most 2 times`/`STATUS: BLOCKED`/`failed repair rounds`）不动；`TestImplementerPromptPinsLoopContracts` 用 contains 断言，增句不破坏。
- 结构断言：SKILL.md 含 `## Blocked evidence` 指令句、定义源与部署产物含第四句锚点。

## <a id="d-blocked-consumption"></a>d-blocked-consumption：refine 消费循环

- `assets/skills/flowforge-refine-ticket/SKILL.md` 流程前置一步：票含 `## Blocked evidence` 时——
  1. 逐条转写其失败事实为 Verified contracts 条目（"命令 X 以方式 Y 失败；正确调用为 Z"——Z 需验证，验证不了则保留为失败事实不含猜测）；
  2. 移除已消费的 `## Blocked evidence` 小节（瞬态即焚，知识归宿是 contracts）；
  3. 复跑 `flowforge check --dir`，确认 `blocked-evidence-present` 消失。
- 顺序：消费在填充五节契约之前（阻断票先消化失败史再补契约）。
- 该步纯方法论流程，不加 CLI 强制；机器可见性由诊断承担。

## Standards clauses

- must 新诊断码为纯本地文件解析，无网络、无 LLM 调用，接入现有 check/frontier 诊断管线与 waiver/policy 语义（源：`AGENTS.md` 核心设计原则 2，[Constraints]）。
- must `assets/` 为权威源，变更经 `make dev` 同步 `internal/command/assets/` 双拷贝后构建（源：`docs/proposals/lightweight-execution-contract/spec.md` Further Notes，[Conventions]）。
- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- must Non-negotiables 既有三锚点逐字不变，第四句为纯增补（源：executor-loop-hardening design rev2 d-prompt-permanence 的单一事实源约束，[Constraints]）。

## 兼容与迁移

- 存量票无该小节 → 解析不命中 → 零影响（验收 8）。
- 本仓自用：future 阻断票会携带小节，refine 流程照常消费；不迁移历史。
- 与 evidence 四元组正交：Blocked evidence 记录失败，四元组记录成功；同一票内共存无冲突（小节在标题层，四元组在列表项层）。

## 验证策略

- 诊断：表驱动测试覆盖 open+小节→报、closed+小节→不报、移除→不报、strict→非零；frontier 归类测试（ready_with_warnings 桶）。
- 写入契约：三端结构断言（SKILL.md 指令句 / 定义源第四句 / 部署产物第四句）+ 既有三锚点回归。
- 消费循环：结构断言（refine-ticket skill 含三步指令与关键词 `Blocked evidence`/`Verified contracts`）。
