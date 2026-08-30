---
flowforge:
  schema: 1
  role: requirement
  id: standards-injection
  revision: 1
---

<a id="standards-injection"></a>
# 项目规范注入任务卡片

## 问题陈述

FlowForge 当前 Align → Solution Design → Plan 流程只解决需求目标如何实现，不解决实现必须遵守哪些项目规范。在一个真实项目中，实现一个需求除需求目标外还必须遵守项目规范（分层依赖规则、接口契约、测试策略等），这些规范根据每张卡片要实现的分层和具体需求场景决定。

tangram-v2 的 `backend-architecture-refactoring-wayfinding` 是真实参照：这批卡片今天靠人工把 `docs/112-backend-dependency-rules.md` 的分层依赖规则手工抄进 Constraints（如 `isTransitive = false`、`AgentRuntime 公开签名不得引用 agentscope 包`）。这是可行的，但依赖执行 agent 恰好知道这些规范、或在 review 阶段才发现违规。需求是：在 Plan 规划卡片时主动把适用规范提取进卡片，让 implementer 直接遵守。

## 可观察结果

1. 部署 FlowForge 到一个项目后，该项目有一份说明文档，描述本项目的规范文档在哪、以及如何为一张卡片判断它适用哪些规范。
2. Plan 发布的每张卡片包含从项目规范提取的适用规范陈述，以 `must be` / `must not be` 形式写入卡片，并链接到规范源。
3. 实现执行前可以检查卡片是否已完成规范提取；缺失时退回 Plan 而非开始执行。
4. Review 的 Standards 轴只检查实现是否符合卡片中已有的规范，不重跑提取逻辑、不检查遗漏。

## 范围

### In scope

- 项目级规范提取说明的内置默认版本与部署。
- Plan 阶段读取提取说明、为每张卡片提取适用规范、按规范性质分流到 Constraints 或 Conventions。
- ticket 契约中规范陈述的书写格式（must/must not + 规范源语义链接）。
- Implement 执行前 pre-flight 检查规范是否已注入。
- Review Standards 轴职责收窄为只查已有规范。

### Out of scope

- 规范文档的版本追踪与漂移检测（卡片只快照执行时点的规范，不记 revision）。
- CLI 解析提取逻辑或提供按路径匹配规范的确定性接口（提取是 skill 的语义工作）。
- 强迫所有项目采用相同的规范组织方式或提取方式（每个项目自定义提取逻辑）。
- 规范提取结果的固定 tier 归属（硬不变量归 Constraints，约定性归 Conventions，由规范性质决定）。

## 场景

1. **标准部署**：项目运行 `flowforge init`，获得内置默认提取说明，可直接使用或按本项目实际补充后使用。
2. **Plan 提取**：Plan 读取提取说明，按其描述的逻辑为本卡片定位适用规范，提取后按每条规范的性质决定写入 Constraints 还是 Conventions。
3. **Implement pre-flight**：Implement 开始执行前确认卡片已含规范提取；缺失时退回 Plan，不开始编码。
4. **Review 只查已有**：Review 读取卡片中已注入的规范，只检查实现是否符合这些规范，不重新从项目规范源提取、不检查 Plan 是否遗漏了本应注入的规范。
5. **项目自定义提取逻辑**：项目修改提取说明文档，采用自己的组织方式（分层、模块、场景或任何结构），Plan 按修改后的逻辑提取。

## 约束

- 提取逻辑是项目内 Markdown 说明文档，不是 CLI 参数、不是固定 schema、不是 yaml manifest。
- FlowForge 统一的部分只有：提取结果在 ticket 中的书写格式（must/must not + 语义链接）和 Plan/Review 的流程职责。
- 不引入通过 CLI 传递长文本的接口。
- 不新增 ticket tier；复用现有 Constraints（硬不变量）和 Conventions（约定性）。
- 不在 ticket metadata 中记录规范 revision，不做机器追踪。

## 术语

- **规范提取说明（standards guide）**：一份项目内 Markdown 文档，描述本项目规范文档在哪、如何为一张卡片判断适用规范。结构由项目自定义。
- **规范陈述（standards clause）**：从项目规范提取后写入卡片的单条规范性陈述，使用 `must be` / `must not be` 形式，附规范源语义链接。
