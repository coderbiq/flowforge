---
flowforge:
  schema: 1
  role: design
  id: fast-executor-reliability-design
  revision: 1
  consumes:
    requirements:
      fast-executor-reliability-requirements: 1
  areas:
    evidence-gate:
      revision: 1
      anchor: d-evidence-gate
    executor-behavior:
      revision: 1
      anchor: d-executor-behavior
    review-audit:
      revision: 1
      anchor: d-review-audit
    execution-unit:
      revision: 1
      anchor: d-execution-unit
---

<a id="fast-executor-reliability-design"></a>
# 快/弱执行者可靠交付方案

需求 authority：[快/弱执行者可靠交付需求](requirements.md#fast-executor-reliability-requirements)，修订版 1。

实证依据：[弱/快模型执行准确率调研](../../research/2026-09-12-weak-model-execution-accuracy.md)（下称"调研 A"）、[子代理 vs 常驻会话性能调研](../../research/2026-09-12-subagent-vs-persistent-session-performance.md)（下称"调研 B"）。衔接既有 authority：[lightweight-execution-contract](../lightweight-execution-contract/spec.md)（三层信息模型与 Fix: 协议保留；其"CLI 不解析段落/不强制"两条 out-of-scope 经需求批准修订）、[subagent-lifecycle](../subagent-lifecycle/design.md)（model 映射与权限表修订，见 d-execution-unit）。

## 一、evidence 四元组与 CLI 门禁 <a id="d-evidence-gate"></a>

### 格式（正文行内标记，不改 frontmatter schema）

每条勾选的 Change 必须携带缩进四元组，键固定、值为人可读：

```markdown
- [x] 1. `findPage` 增加 `sessionId: String?` 参数，非空时追加 `session_id = ?` 条件
    - cmd: `./gradlew :extension-impl:tangram-exposed:test --tests "*ExposedTimelineEventRepositoryTest*"`
    - exit: 0
    - output: "5 tests completed, 0 failed"
    - artifact: backend/extension-impl/tangram-exposed/src/.../ExposedTimelineEventRepository.kt
```

- 缩进与 Change 文本起始列对齐（CommonMark 子列表规范）；`cmd`、`exit`、`artifact` 必需，`output` 必需且限 1–3 行关键结果摘录（如 `N tests completed, 0 failed`），禁止粘贴全文（延续"不 dump 终端"精神，但摘录必须逐字）。
- `exit` 值必须为 `0`。非零退出码意味着该 Change 未完成：执行者不得勾选，票保持 open 并按失败即停规则处置（见 d-executor-behavior）。环境性失败同样不得勾选，走 BLOCKED 终态记录处置。
- 未勾选的 Change 带四元组不构成违例（修复轮中间态合法）。
- 决策依据：调研 A 观察四——BMAD "never infer success from chat output alone"、spec-kit converge 的命令锚定、OpenAI structured outputs 证明格式诚实可被约束而内容诚实只能靠退出码。

### 校验规则（`flowforge check` 扩展）

新增诊断码（接入 `internal/tracker` 现有诊断管线，`discoverArtifact` 拥有；`--strict` 下 warning 判失败，frontier 分类沿用现有 policy）：

| 诊断码 | 触发 | 严重度 |
|---|---|---|
| `evidence-missing` | `- [x]` 存在且其缩进块内无任何四元组键 | warning |
| `evidence-incomplete` | 四元组缺 `cmd`/`exit`/`output`/`artifact` 任一键 | warning |
| `evidence-exit-nonzero` | `exit` 值非 `0` | warning |
| `evidence-artifact-missing` | `artifact` 路径相对仓库根不存在 | warning |

- 解析范围：`role: ticket` 文件（含 legacy 兼容票）；纯文档票（无 Write set）不校验勾选证据（对齐 lightweight spec 的跳过先例）。
- 不做语义校验：`cmd` 是否真实运行过、`output` 是否属实属于内容诚实，归 review Round 0 与双轴（CLI 只解决格式诚实——调研 A 观察四的分工结论）。
- 豁免：`.flowforge/config.yaml` 新增 `evidence.exempt_proposals: [<dir>]`，命中则该 proposal 下全部 evidence 诊断不产生（存量渐进采用；不做 git 时间戳推断）。
- CLI 签名不变（复用现有 `--strict` 与诊断管线）；无网络、无 LLM 依赖，纯文件解析。

### 替代方案（被拒）

- **frontmatter evidence 字段**：结构最严谨，但触碰 Issue Schema 头规范（Ask-first 边界），且证据与 Change 在阅读上分离；正文标记与之相比零迁移成本。
- **独立 evidence 文件**（schema v1 `role: evidence`）：轻量场景过重，现有 spec 仅多环境验证才提升文件。
- **check warn-only / 纯 skill 纪律**：tangram-v2 实证弱执行者会违反无强制力纪律（lightweight 模式"不写 evidence"指令被违反），调研 A 观察六（hooks "deterministic and guarantee"）支持 strict。

## 二、implement 弱执行者行为契约 <a id="d-executor-behavior"></a>

对 `flowforge-implement` SKILL 的 lightweight mode 强化（full mode 面向强模型，保持现状）：

1. **editor 角色收窄**（aider editor 简化提示先例，调研 A 观察二）：正文显式声明"所有决策已在工单内完成，你只负责机械执行与如实报告"；歧义或与工单不符的事实的唯一合法出口是 `STATUS: BLOCKED` 返回，不是自行发挥。
2. **两段式执行**（aider 自我配对 +4.6~+10.3 点的免费收益）：动手前先在 Implementation note 头部产出"逐条 Change 复述 + 每条对应的验收命令"映射表；复述与工单有出入即 BLOCKED。
3. **失败即停**（SWE-agent：一次失败后恢复率 90.5%→57.2%，"succeed quickly fail slowly"）：单条 Change 验收命令失败重试上限 2 次；票级修复轮上限 5 轮（BMAD 参数移植）；超限写 BLOCKED 终态并保留现场，不做无限自愈。
4. **错误回喂**（Anthropic 工具设计原则）：验收命令失败输出逐字进入 Implementation note，禁止弱模型转述错误。
5. **勾选解耦**：只有四元组 `exit: 0` 写入后才允许将 `- [ ]` 改为 `- [x]`；两者必须同一次编辑完成（防"先勾后补"）。
6. **预检显式遍历**：开工前把每条 Constraint 与 Execution scenario 列为 checklist 逐条标注实现落点（弱模型不会隐式合规——tangram-v2 "global" 命名空间案）。

模式选择从执行者自评改为**派发方声明**：ticket 正文 `**Mode:** lightweight`（Plan/refine-ticket 产出时写入）或派发 prompt 显式指定；无声明时维持现有判据。

## 三、review Round 0 审计层 <a id="d-review-audit"></a>

`flowforge-review` SKILL 流程前置一个 phase（不新建 skill，复用现有双轴编排与 Fix: Change 回填机制）：

- **输入**：意图源 = ticket + linked requirement/design authorities（唯一意图源，不信 Implementation note 自述）；对象 = 固定 diff（fixed point）。
- **执行者**：tool-capable-read-only 档（中档模型）单代理，串行于双轴之前。
- **职责**：以工单逐条 Change 为核对单元做差距审计，gap 分类 `missing` / `partial` / `contradicts` / `unrequested`（spec-kit converge 分类）；每条 gap 必须带 file:line 证据；机械 gap 直接转 `Fix:` Change 回填（现有协议）。
- **升级条件**：Round 0 零 gap 才进入双轴复审；双轴只审 survivors，指令限定"不重复证伪 Round 0 已核对项"。
- **噪音控制**：指令限定只报影响验收条款的 gap（Claude Code："被要求找 gap 的 reviewer 总会报出一些"）。
- 串行而非并行的理由：并行会让双轴继续为机械 gap 付旗舰 token，恰是本设计要消除的成本；分钟级时延可接受（实测双轴子代理单轮 1.4–3.7 分钟）。

## 四、执行单元策略与部署物 <a id="d-execution-unit"></a>

写入 `assets/AGENTS.md` 部署区块新小节（复用 `applyAgentsBlock` 管线）：

1. **一票一上下文**：批量工单每 ticket 委派一个全新执行单元（子代理或新会话）；同批次同模型、同工具集（缓存友好，调研 B 观察三）；跨 ticket 状态只经工件（ticket 文件、STATUS 契约、`flowforge frontier`）传递，不经对话历史。
2. **入口必读工件清单**：执行单元启动时读 AGENTS.md + ticket + linked authorities + 前序 ticket 的 evidence，替代自由探索（冷启动税摊销走工件——调研 B 结论 6）。
3. **测试分离**：验收测试命令与 Expected tests 由 Plan/refine-ticket 预置；弱模式执行者禁改测试文件（宿主 permission 提示：opencode `permission.edit` fileRegex、Claude Code `disallowedTools` 示例写入部署文档）。

**subagent-lifecycle 关联修订**（该 design 同步升 revision 2）：model_profile 映射表中 `tool-capable` 在 OpenCode 下从"省略字段（继承主会话）"改为"`agents.models.tool-capable` 配置存在时写显式 `model` 字段，未配置回退继承"；implementer 角色权限表追加"plan 预置测试文件默认只读（可配置关闭）"。理由：继承主会话会使执行子代理落到旗舰模型（tangram-v2 部署现状实证），与"执行者钉快模型"的分工矛盾。

## Standards clauses

- must evidence 四元组保持人可读 Markdown 正文标记，执行者与 CLI 间不得引入传长文本的接口（源：`AGENTS.md` 核心设计原则 1，[Constraints]）。
- must 新诊断码为纯本地文件解析，无网络、无 LLM 调用，接入现有 check/frontier 诊断管线与 waiver/policy 语义（源：`AGENTS.md` 核心设计原则 2，[Constraints]）。
- must `assets/` 为权威源、`.agents/` 为部署快照，变更经 `make dev` 同步 `internal/command/assets/` 双拷贝后构建（源：`docs/proposals/lightweight-execution-contract/spec.md` Further Notes，[Conventions]）。
- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。

## 迁移与兼容

- 存量 proposal：默认 warning 不阻断；确需静默的加入 `evidence.exempt_proposals`。`flowforge upgrade` 输出一行提示说明新配置键。
- skill 文本变更随 `init --force` / `upgrade` 部署到既有项目；已部署项目不受 CLI 强制（evidence 校验由其本地 CLI 版本决定）。

## 验证策略

- Go 单测：四元组解析（合法/缺键/非零退出/artifact 缺失/嵌套缩进/CommonMark 对齐）、豁免列表、strict 语义、frontier 分级透传；用本仓库既有 proposal 票做回归（零新诊断预期：存量默认 warning 数量已知）。
- Skill 文本：结构断言（六项行为契约关键词存在性）+ 人工走查一张真实 repair ticket 的 dry-run。
- Round 0：在下一个真实 proposal 的首次交付上试跑校准（gap 分类准确率、双轴轮次下降），校准结论回写本设计修订。

## open items

- 无 gap 级未决事实。宿主级强制力（opencode/Claude Code hooks 生成）与 Round 0 prompt 实证校准列为后续增量，不阻塞本轮 Plan。
