---
flowforge:
  schema: 1
  role: requirement
  id: fast-executor-reliability-requirements
  revision: 1
---

<a id="fast-executor-reliability-requirements"></a>
# 快/弱执行者可靠交付需求

## 问题

FlowForge 的实际使用形态是"旗舰模型做 proposal 分析与设计、廉价快模型（flash 级）做 ticket 执行"。tangram-v2 的实践复盘（2026-09-12，基于本机 opencode 会话记录与 ff-wiki review 轮次统计：39 张票 R1 之后仍有 17 张到 R2、8 张到 R3、最多 7 轮）暴露五类反复出现的弱执行者失败：

1. 勾选 `[x]` 但对应 Change 未实现（假完成）；
2. 测试存在但断言空洞或断言了错误对象（vacuous test）；
3. Completion evidence 的文件名、测试数量与实际不符（幻觉式总结）；
4. 无视票内 Failure scenario / Conventions；
5. 用自创机制替代指定的先例机制。

这些失败全部落在旗舰 review 的证伪工作量里，导致一个 proposal 需要 2–7 轮双轴复审，旗舰模型 token 被全量证伪消耗。同时，批量执行的上下文策略缺失：单常驻会话连做整批 ticket 会积累巨型上下文并出现遵从度衰减（实测 54M token 会话后期 ticket 的 review 发现数显著高于前期），而业界一手资料（`docs/research/2026-09-12-subagent-vs-persistent-session-performance.md`）确认常驻与全新执行单元各有结构性成本，收敛答案是对批量工单采用"每任务一个全新执行单元 + 上下文复用走工件"。

业界实证（`docs/research/2026-09-12-weak-model-execution-accuracy.md`）给出的对应解法方向：接口收窄与格式简单化（aider edit format 基准）、强规划+弱执行分工（aider architect 配对 +10.3 点）、确定性门禁替代 LLM 轮次（Claude Code hooks / spec-kit checklist gate）、命令锚定的机器可校验证据（BMAD / spec-kit converge）、失败即停（SWE-agent 恢复率数据）、测试作者与实现者分离（spec-kit / Claude Code TDD 范式）。

## 目标

1. "勾了 `[x]` 但没有命令锚定证据"可被 `flowforge check` 机械识别：默认 warning，`--strict` 下判失败——把假完成从 review 问题降为 lint 问题。
2. `flowforge-implement` Skill 对弱执行者提供明确行为契约：editor 角色收窄、两段式执行（先复述后动手）、失败即停与重试上限、勾选与命令退出码解耦。
3. `flowforge-review` 增加廉价审计前置层（converge 式差距审计）：以工单为唯一意图源核对代码，机械类 gap 在进入旗舰双轴复审前被拦截或修复。
4. 执行单元策略成文：批量工单每 ticket 一个全新执行上下文，跨 ticket 状态经工件（ticket 文件 + STATUS 契约）传递；执行者子代理的模型钉定与测试文件保护在部署物中可表达。

## 范围与约束

- 用户已批准三项边界决策（2026-09-12 会话）：evidence 门禁采用 **CLI strict 校验**（扩展 `flowforge check`，不新增子命令签名）；evidence 载体采用 **ticket 正文行内标记**（不修改 Issue Schema frontmatter）；本轮纳入 review 审计层与 subagent-lifecycle 映射修订。
- 本需求**修订** `docs/proposals/lightweight-execution-contract/spec.md` 的两条 out-of-scope 决定（"CLI 不解析新段落"、"不自动强制"）：evidence 四元组引入机器解析与 strict 强制。该 spec 的三层信息模型、Fix: Change 协议、lightweight/full 模式边界保留不动。
- 向后兼容：存量无 evidence 标记的 ticket 默认只产生 warning；豁免名单经配置管理，不做自动迁移。
- 宿主差异（opencode / Claude Code / Codex 的 hooks、permission 能力）不阻塞本需求核心：强制力落在 CLI 校验，宿主层能力只作为增量提示。

## 可观察验收

- 给定一张含 `- [x]` Change 但无 evidence 四元组的 ticket，`flowforge check --strict` 报违例并以非零退出；默认模式仅 warning。
- 给定四元组中 `exit` 非 0、缺必需键、或 `artifact` 路径在仓库不存在，check 报对应违例。
- 给定完整合规四元组的 ticket，check 通过且不产生新诊断。
- 存量 proposal 加入豁免名单后，其 evidence warning 消失。
- 更新后的 `flowforge-implement` SKILL 文本包含：editor 角色声明（歧义唯一出口是 BLOCKED）、两段式执行步骤、Change 级重试上限与票级修复轮上限、"退出码非 0 不得勾选"规则。
- 更新后的 `flowforge-review` SKILL 文本包含：Round 0 审计 phase 的输入（意图源 = 工单 + linked authorities；对象 = 固定 diff）、gap 分类（missing / partial / contradicts / unrequested）、升入双轴复审的条件。
- `flowforge agents deploy` 生成的 OpenCode implementer 定义在配置了执行者模型时携带显式 `model` 字段；未配置时行为与现状一致（继承主会话）。

## 非目标

- 不为 opencode/Claude Code 生成宿主 hooks（远期增量，列为后续 open item）。
- 不实现自动重试调度、代理间消息传递或执行编排器。
- 不改造用户的多窗口使用习惯；执行单元策略是规范与部署物支持，不是强制。
- 不引入 evidence 独立文件 role 的通用提升（schema v1 `role: evidence` 维持现状）。
- 不做跨会话记忆或学习机制。
