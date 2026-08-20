# FlowForge Skills Specification & Engineering Methodology

> FlowForge v4 五大原子 Skill 规范与敏捷工程方法论。

---

## 1. 设计哲学：方法论优于形式主义

过去的版本将所有职责强塞给 `design` 和 `implement`，导致 Agent 陷入“黑盒生成八股文”的泥潭。
FlowForge v4 借鉴 **`mattpocock/skills`** 的敏捷与极简哲学，构建了 5 个职责专一的原子 Skill：

```text
┌─────────────────────────────────────────────────────────────┐
│                      FlowForge Skill 矩阵                   │
└─────────────────────────────────────────────────────────────┘
  1. flowforge-align (Grill & Memory)
     └─ 目标：在对话中深挖需求、澄清边界、维护 Proposal 活笔记
     └─ 借鉴：/grill-me + Working Memory

  2. flowforge-explore (Fact Finding)
     └─ 目标：针对特定疑点深入代码库调研，将事实回填至 Proposal
     └─ 借鉴：/explore + /diagnosing-bugs

  3. flowforge-plan (Tracer Bullet Slice)
     └─ 目标：在共识达成后，将方案拆解为 3~6 个极小可验证切片
     └─ 借鉴：Task Decomposition + Vertical Slice

  4. flowforge-implement (TDD Cycle)
     └─ 目标：专注执行一个切片：写红测试 → 写代码 → 绿测试 → 重构
     └─ 借鉴：/tdd + Red-Green-Refactor

  5. flowforge-curate (Living Doc Synthesis)
     └─ 目标：提案交付后，提取增量合流至系统全局文档，归档提案
     └─ 借鉴：Living Documentation Merge
```

---

## 2. 五大原子 Skill 详细规范

### 2.1 `flowforge-align` (对齐与记忆捕获)
* **触发场景**：用户提出新想法、复杂业务需求、或恢复一个正在讨论的提案。
* **输入契约**：用户的自然语言诉求或现有 Proposal ID。
* **执行准则**：
  1. **禁止先建全套工件**：在共识达成前，绝不盲目生成代码或大量 Feature 卡片；
  2. **Grilling 追问机制**：主动挖掘边界、假设与非目标（Non-goals），每次提问 1~3 个关键架构/业务选择题；
  3. **记忆实时捕获**：将讨论结论立即以极简 Bullet Points 写入 `01-workspace/<proposal_id>/README.md`。

### 2.2 `flowforge-explore` (代码与数据事实探查)
* **触发场景**：讨论过程中遇到不确定的代码逻辑、遗留表结构、接口边界时。
* **执行准则**：
  1. **只查事实，不写业务代码**：快速定位具体函数、数据样本、异常堆栈；
  2. **结构化回填**：必须输出带有文件坐标（`file:line`）的事实链，回填至活笔记的 `[Explored Facts]`。

### 2.3 `flowforge-plan` (Tracer Bullet 任务拆解)
* **触发场景**：用户明确表示“方案已对齐，开始规划”或调用 `/plan`。
* **执行准则**：
  1. **Tracer Bullet 原则**：拆解为 3~6 个纵向极小切片（Vertical Slice），每个切片 15~30 分钟内可完成；
  2. **Test-Driven 绑定**：每个切片必须绑定一条明确的自动化测试验证命令；
  3. **拒绝伪代码**：不在计划中堆砌不可靠的伪代码，保持切片定义极简。

### 2.4 `flowforge-implement` (TDD 增量实现)
* **触发场景**：用户要求执行某个具体 Slice。
* **执行准则**：
  1. **红-绿-重构循环 (Red-Green-Refactor)**：
     - **Red**：先编写/更新失败的单元测试；
     - **Green**：编写最小代码使测试通过；
     - **Refactor (Cognitive Load Reduction)**：优化代码结构，消除死代码，内联浅封装，早返回，保持测试全绿。
  2. **测试作为唯一裁判**：只有绑定的测试命令 100% 运行通过，该切片才算完成。

### 2.5 `flowforge-diagnose` (假设驱动排错)
* **触发场景**：遇到复杂 Bug、偶发测试失败、生产回归或死锁排查。
* **执行准则**：
  1. 提出可证伪的假设；
  2. 构建最小隔离复现用例；
  3. 追踪变量与状态跃迁；
  4. 修复根因而非表面症状。

### 2.6 `flowforge-review` (非阻断式对抗性审查)
* **触发场景**：提案切片交付后、合并前进行红队安全与架构漂移审查。
* **执行准则**：
  1. 架构漂移检测（对比 `docs/CONTEXT.md`）；
  2. 认知负荷削减建议（死代码、浅封装）；
  3. 并发与边界安全检查。

### 2.7 `flowforge-curate` (活文档合流与归档)
* **触发场景**：Proposal 所有切片完成并通过测试，准备归档。
* **执行准则**：
  1. **Delta 提取**：提取关键架构决策为 ADR，提取业务现状变更生成 Domain Patch；
  2. **Diff 审查与合流**：展示对 `docs/domains/<domain>/README.md` 的变更 Diff，用户确认后合并，并将 Proposal 移入 `00-archive/`。
