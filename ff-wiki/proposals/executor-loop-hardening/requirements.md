---
flowforge:
  schema: 1
  role: requirement
  id: executor-loop-hardening-requirements
  revision: 2
---

<a id="executor-loop-hardening-requirements"></a>
# 执行者循环硬化需求

## 问题

fast-executor-reliability 交付后，tangram-v2 实测（2026-09-13，opencode 会话 DB 取证）暴露新的结构性失败：一张票由 gemini-3.8-flash 执行 86.4 分钟、893 次工具调用，其中同一条 `python3 -c …javap…` 命令**原样重试 290 次**、同一条 gradle 集成测试命令执行 49 次；会话累计输入 1590 万 tokens（cache 读 6370 万），无任何长停顿——慢的机制不是基础设施，是"微步爆炸 × 每步往返延迟 × 上下文持续膨胀"。对照同模型相邻票仅 51-149 次调用，属弱执行者重试循环离群。

三个断层被逐一证实：

1. **无执行预算**：宿主 opencode 原生支持 agent 级 `steps` 上限（到限注入总结 system prompt 优雅收束），flowforge 编译器不产出该字段，弱执行者可以无限迭代。
2. **契约未固化**：fail-fast（同一命令失败 2 次即上报 BLOCKED）等六项契约只存在于 SKILL.md，按需加载；涉事会话确实加载过 1 次 skill，但 290 次重试证明弱模型加载后不遵守。业界一手结论（docs/research/2026-09-13-agent-loop-hardening-prior-art.md）：持久规则必须在"每次请求都在的层"（Claude SDK 原话 "re-injected on every request"；Codex/Gemini/opencode 的 AGENTS.md 均每会话注入，源码实证），而 flowforge-implementer 的固化 prompt 仅 43 行、不含 fail-fast。
3. **无事后循环信号**：evidence 四元组已机器解析（cmd/exit），但"同一命令非零退出反复出现"这一最强烈的空转信号不产生任何诊断，review 需人眼发现。

业界对照（同研究文档）：Gemini CLI 对工具名+参数 SHA-256 指纹连续 5 次判循环（首检注入反馈、复检熔断）；opencode 内建 `DOOM_LOOP_THRESHOLD=3`（同工具同参数连续 3 次）+ `steps` 预算；OpenAI Agents SDK 默认 `max_turns=10`。查证确认无任何宿主内置"验收命令业务失败连续 N 次硬停"——该层是 flowforge 可自建的增量。

## 目标

1. 执行预算可配置且编译产物可表达：`agents.max_steps` 配置后 implementer 的 OpenCode 编译产物携带 `steps` 字段；未配置使用保守默认；可显式关闭。
2. 行为契约固化到每次注入的层：implementer 的固化 prompt（agent 定义 body）携带不可协商的防循环摘要（fail-fast 上限、禁止原样重试、预算到限收束），不依赖 SKILL.md 按需加载。
3. 空转信号机器可识别：`flowforge check` 对已勾选 Change 的四元组中"同一命令非零退出出现 ≥3 次"产生 `evidence-repeat-failure` 诊断（默认 warning，strict 判失败）。

## 范围与约束

- 宿主 enforcement 仅 OpenCode（`steps` 为 opencode agent 配置字段；Claude Code/Codex 无对应 frontmatter 能力，不做模拟）。
- 不在 flowforge CLI 内实现运行时循环检测（CLI 不在场于宿主执行循环；预算与断路归宿主，事后信号归 check）。
- 不修改 Issue Schema 头规范；新诊断沿用现有 evidence 豁免机制（`evidence.exempt_proposals`）。
- prompt 摘要是 SKILL.md 契约的固化引用，不是第二事实源：SKILL.md 仍拥有完整契约文本，摘要与其措辞保持一致（结构断言测试锁定关键词）。

## 可观察验收

1. 配置 `agents.max_steps: 200` 并部署后，`.opencode/agent/flowforge-implementer.md` frontmatter 含 `steps: 200`；`agents status` 报 current。
2. 不配置时，implementer 产物含 `steps: <默认值>`；配置 `-1` 时产物不含 `steps` 字段。
3. 非 implementer 的 subagent（analyst/planner/architect/reviewer/investigator）产物不受 `max_steps` 影响。
4. `assets/subagents/flowforge-implementer.md` body 含 fail-fast 上限与禁止原样重试的关键句；部署后的 `.opencode/agent/flowforge-implementer.md` 同样含（每会话注入层固化）。
5. 构造一张票：同一命令非零退出四元组出现 3 次 → `flowforge check` 报 `evidence-repeat-failure` warning，`--strict` 退出非零；出现 2 次不报。
6. 存量票（无四元组或已豁免）不受新诊断影响。
7. 升级已部署项目（无新配置）后：implementer 产物与旧版差异仅为新增 `steps` 默认值行与 prompt 摘要小节，其余逐字节一致。

## 术语

- **执行预算（execution budget）**：宿主级单会话代理迭代上限（opencode `steps`），到限注入总结提示、强制文本收束，非硬杀。
- **契约固化（contract pinning）**：把行为契约放进每次请求必然注入的层（agent 定义 body），与按需加载（skills）相对。
- **重复失败信号（repeat-failure signal）**：同一规范化命令在已勾选 Change 的四元组证据中以非零退出重复出现 ≥3 次。
