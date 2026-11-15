---
flowforge:
  schema: 1
  role: requirement
  id: standards-injection
  revision: 2
  areas:
    standards-flow:
      revision: 2
      anchor: standards-flow
---

<a id="standards-injection"></a>
# 项目规范注入任务卡片

## 问题陈述

FlowForge 当前 Align → Solution Design → Plan 流程只解决需求目标如何实现，不解决实现必须遵守哪些项目规范。在一个真实项目中，实现一个需求除需求目标外还必须遵守项目规范（分层依赖规则、接口契约、测试策略等），这些规范根据每张卡片要实现的分层和具体需求场景决定。

tangram-v2 的 `backend-architecture-refactoring-wayfinding` 是真实参照：这批卡片今天靠人工把 `docs/112-backend-dependency-rules.md` 的分层依赖规则手工抄进 Constraints（如 `isTransitive = false`、`AgentRuntime 公开签名不得引用 agentscope 包`）。这是可行的，但依赖执行 agent 恰好知道这些规范、或在 review 阶段才发现违规。需求是：让规范在分析和设计阶段就进入，使设计本身合规，最终卡片中的规范是 Align/Design 共同理解的一致产物。

## 可观察结果

<a id="standards-flow"></a>
1. 部署 FlowForge 到一个项目后，该项目有一份说明文档，描述本项目的规范文档在哪、以及如何为一张卡片判断它适用哪些规范。
2. Align 阶段读取提取说明，识别当前需求适用哪些项目规范，并传递给 Solution Design。
3. Solution Design 设计合规方案，并将规范转换成 `must` / `must not` 陈述（含 tier 归属和规范源链接）写入设计 authority。
4. Plan 从设计 authority 机械转写 `must`/`must not` 到卡片对应 tier；Plan 不读提取说明、不做规范语义转换。
5. 实现执行前可以检查卡片是否已含 `must`/`must not`；缺失时退回 Plan 而非开始执行。
6. Review 的 Standards 轴只检查实现是否符合卡片中已有的规范，不重跑提取逻辑、不检查遗漏。

## 范围

### In scope

- 项目级规范提取说明的内置默认版本与部署。
- Align 阶段读取提取说明、识别适用规范、传递给 Solution Design。
- Solution Design 设计合规方案并将规范转换成 `must`/`must not`（含 tier 归属）写入设计 authority。
- Plan 从设计 authority 机械转写 `must`/`must not` 到卡片；不读提取说明、不做语义转换。
- ticket 契约中规范陈述的书写格式（must/must not + 规范源语义链接）。
- Implement 执行前 pre-flight 检查规范是否已注入。
- Review Standards 轴职责收窄为只查已有规范。

### Out of scope

- 规范文档的版本追踪与漂移检测（卡片只快照执行时点的规范，不记 revision）。
- CLI 解析提取逻辑或提供按路径匹配规范的确定性接口（提取是 skill 的语义工作）。
- 强迫所有项目采用相同的规范组织方式或提取方式（每个项目自定义提取逻辑）。
- 规范提取结果的固定 tier 归属（由 Design 按规范性质决定写入 Constraints 或 Conventions）。

## 场景

1. **标准部署**：项目运行 `flowforge init`，获得内置默认提取说明，可直接使用或按本项目实际补充后使用。
2. **Align 提取规范**：Align 读取提取说明，按其描述的逻辑识别当前需求适用的项目规范，传递给 Solution Design。
3. **Design 转换并合规设计**：Solution Design 接收 Align 识别的规范，在设计决策时对照规范确保合规，并将规范转换成 `must`/`must not` 陈述（含 tier 归属和规范源链接）写入设计 authority。
4. **Plan 转写**：Plan 从设计 authority 读取 `must`/`must not`，机械转写到卡片的 Constraints 或 Conventions；不读提取说明、不做语义判断。
5. **Implement pre-flight**：Implement 开始执行前确认卡片已含 `must`/`must not`；缺失时退回 Plan，不开始编码。
6. **Review 只查已有**：Review 读取卡片中已注入的规范，只检查实现是否符合这些规范，不重新从项目规范源提取、不检查遗漏。
7. **项目自定义提取逻辑**：项目修改提取说明文档，采用自己的组织方式，Align 按修改后的逻辑提取。

## 约束

- 提取逻辑是项目内 Markdown 说明文档，不是 CLI 参数、不是固定 schema、不是 yaml manifest。
- 规范的语义工作（识别哪些适用、转换成 must/must not、决定 tier 归属）在 Align 和 Design 阶段完成；Plan 只做机械转写。
- 规范在三个阶段的一致性由单一来源保证：Design authority 是 must/must not 的唯一权威来源，Plan 从中转写，不独立理解。
- FlowForge 统一的部分只有：must/must not 陈述在 ticket 中的书写格式和 Align/Design/Plan/Review 的流程职责。
- 不引入通过 CLI 传递长文本的接口。
- 不新增 ticket tier；复用现有 Constraints（硬不变量）和 Conventions（约定性）。
- 不在 ticket metadata 中记录规范 revision，不做机器追踪。

## 术语

- **规范提取说明（standards guide）**：一份项目内 Markdown 文档，描述本项目规范文档在哪、如何为一张卡片判断适用规范。结构由项目自定义。Align 读取此文档识别适用规范。
- **规范陈述（standards clause）**：由 Solution Design 转换并写入设计 authority 的单条规范性陈述，使用 `must` / `must not` 形式，附规范源语义链接和 tier 归属。Plan 从设计 authority 机械转写到卡片。
