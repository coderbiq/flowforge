---
flowforge:
  schema: 1
  role: design
  id: standards-injection-design
  revision: 2
  consumes:
    requirements:
      standards-injection: 2
  areas:
    standards-guide:
      revision: 1
      anchor: standards-guide
    align-extraction:
      revision: 2
      anchor: align-extraction
    design-conversion:
      revision: 2
      anchor: design-conversion
    plan-transcription:
      revision: 2
      anchor: plan-transcription
    implement-preflight:
      revision: 1
      anchor: implement-preflight
    review-standards:
      revision: 1
      anchor: review-standards
    setup-deploy:
      revision: 1
      anchor: setup-deploy
---

<a id="standards-injection-design"></a>
# 项目规范注入任务卡片 — 方案设计

需求 authority：[项目规范注入任务卡片需求](requirements.md#standards-injection)，修订版 2。

## 设计区域

### 1. 规范提取说明（standards guide）

<a id="standards-guide"></a>

一份项目内 Markdown 文档，描述本项目的规范文档在哪、以及如何为一张卡片判断它适用哪些规范。结构由项目自定义，FlowForge 不规定其内部结构。

**部署形态**：作为受管资产（managed asset），与 `agents/domain.md`、`agents/issue-tracker.md` 同级，部署到 `<docs_dir>/agents/standards.md`。由 `flowforge init`/`upgrade` 部署内置默认版本。

**默认版本内容**：示范一种通用提取思路（如按文件路径推导分层、按 Touch points 和 Write set 推导涉及的模块与场景），并在正文中以注释引导用户如何根据本项目实际补充。默认版本不假设任何特定项目结构；它提供一个可工作的起点，项目应在其基础上特化。

**配置登记**：新增配置键 `standards.guide`（默认 `agents/standards.md`，相对于 `docs_dir`），指向提取说明文档位置。项目可改指到其他路径。`flowforge config get/set` 支持该键。

**为什么不复用 `KnowledgeSources`**：`KnowledgeSources` 是 list 结构，设计初衷是登记多个知识源元数据。规范提取说明是单一文档引用，用独立标量键 `standards.guide` 更清晰，且避免给 `KnowledgeSources` 增加语义负担。`KnowledgeSources` 保持现状不动。

### 2. Align 提取规范

<a id="align-extraction"></a>

`flowforge-align` 的 step 2「Resolve repository facts first」扩展为同时读取提取说明：

1. 读取 `standards.guide` 配置指向的提取说明文档。
2. 按提取说明中描述的逻辑，识别当前需求适用的项目规范（哪些规范与本需求的范围、场景、约束相关）。
3. 将识别结果传递给 Solution Design——作为 requirement constraint 或作为规范识别清单附在 requirement authority 中。

Align 只做**识别**：判断哪些规范适用，不做 `must`/`must not` 转换，不决定 tier 归属。转换是 Design 的职责。

**为什么 Align 而非 Design 做提取**：规范适用性取决于需求范围和场景，这些是 Align 的职责。如果 Design 同时做识别和转换，可能先入为主地按设计方案筛选规范而遗漏。Align 先识别范围适用性，Design 再在设计决策时对照并转换，形成两层保障。

### 3. Design 转换 must/must not 并合规设计

<a id="design-conversion"></a>

`flowforge-solution-design` 扩展两个职责：

**3a. 设计合规**：step 1「Establish the decision frontier」新增接收 Align 传递的规范识别清单。step 3「Compare credible designs」把规范合规性作为设计替代方案的比较维度——违反规范的设计替代方案被否决或需要显式 waiver。

**3b. 转换写入设计 authority**：step 4「Maintain authority incrementally」新增——将 Align 识别的规范转换成 `must`/`must not` 陈述，含 tier 归属和规范源语义链接，写入设计 authority。设计 authority 是 `must`/`must not` 的唯一权威来源。

**落点由 Design 决定**：Design 按每条规范本身的性质决定 tier 归属：
- 硬不变量（违反即失败、影响完成判定）→ 标注为 Constraints。
- 约定性（implementer 须读但较软、不直接决定完成）→ 标注为 Conventions。

**书写格式**（在设计 authority 中）：

```markdown
## Standards clauses

- must not <forbidden behavior> — <规范源语义链接> [Constraints]
- must <required behavior> — <规范源语义链接> [Conventions]
```

Design authority 中每条 clause 标注 `[Constraints]` 或 `[Conventions]` 指示 Plan 转写时的目标 tier。链接使用规范源文件的相对路径 + anchor，不记 revision。

**为什么 Design 而非 Plan 做转换**：转换是语义工作——需要判断规范在设计语境下的具体含义、与设计决策的关系、以及违反后果的严重程度。Design 同时做设计决策和规范转换，保证两者一致。Plan 不具备设计上下文，无法做这个判断。

**为什么单一来源保证一致性**：三个阶段（Align 识别、Design 转换、Plan 转写）中，`must`/`must not` 的语义工作只在 Design 做一次。Plan 从设计 authority 机械转写，不独立理解规范。这消除了三个 skill 各自理解导致的不一致风险。

### 4. Plan 机械转写

<a id="plan-transcription"></a>

`flowforge-plan` 的「Resolve effective authority」一步改为从设计 authority 读取 `must`/`must not` clauses：

1. 读取设计 authority 中的 Standards clauses 段落。
2. 按每条 clause 标注的 tier 归属（`[Constraints]` 或 `[Conventions]`），机械转写到卡片的 `## Constraints` 或 Tier 3 `### Conventions`。
3. 转写时不修改陈述内容、不重新判断 tier、不读提取说明。

Plan 不再做：读提取说明、识别适用规范、转换 must/must not、判断 tier。这些工作由 Align 和 Design 完成。

**自检**：Plan 确认已从设计 authority 转写所有 Standards clauses 到卡片。如果设计 authority 无 Standards clauses 段落且卡片有 Write set，Plan 在卡片 Constraints 中标注 `standards: pending` 退回 Design。

**compact ticket 情况**：对于跳过独立 Solution Design 的 compact ticket，如果 Align 在 requirement authority 中已识别适用规范，Plan 直接从 requirement authority 的规范识别清单做简单转写（不读提取说明）。如果 Align 未识别（compact ticket 可能跳过独立 Align），Plan 标注 `standards: pending`。

### 5. Implement pre-flight 检查

<a id="implement-preflight"></a>

`flowforge-implement` 的 step 1「Resolve the effective specification」增加 pre-flight：

- 读取卡片 Constraints，检查是否有以下三种状态之一：
  - `must`/`must not` 规范陈述已存在 → 通过，继续执行。
  - `standards: pending` → 规范转写未完成，退回 Plan。
  - `standards: none found` → 设计 authority 中无适用规范，通过，继续执行。
- 缺少以上任何标记且卡片涉及代码变更（有 Changes 和 Write set）时，视为规范转写缺失，退回 Plan。

**判定标准**：卡片只要有 Changes 和 Write set 就必须有规范标记。纯文档类卡片（无 Write set）不要求。

**为什么退回 Plan 而非自行提取**：退回 Plan 让 Plan 从设计 authority 补充转写。如果设计 authority 本身缺少 Standards clauses，Plan 再退回 Design。

### 6. Review Standards 轴职责收窄

<a id="review-standards"></a>

`flowforge-review` 的 step 3「Identify the standards sources」改为：

1. 读取卡片中已注入的规范陈述（Constraints 和 Conventions 中的 `must`/`must not` 条目）。
2. Standards 轴只检查实现是否符合卡片中已有的规范。
3. 不重新从项目规范源提取、不检查 Plan 是否遗漏了本应注入的规范。

**Step 3 原有逻辑保留**：Review 仍会查找仓库中 documented coding standards 文件（如 `CODING_STANDARDS.md`）和 smell baseline——这些是通用代码质量基线，与本需求的「项目规范注入」正交。两者并存：卡片内规范是项目特定规范，通用 standards 查找是跨项目代码质量基线。

### 7. Setup 部署

<a id="setup-deploy"></a>

`flowforge-setup` 新增 Section D「Standards」：

1. 检查 `<docs_dir>/agents/standards.md` 是否存在（已由 `init` 部署受管资产）。
2. 向用户展示默认版本内容摘要。
3. 引导用户根据本项目实际补充：规范文档在哪、如何为卡片判断适用规范。
4. 确认后写入。用户可在后续随时编辑该文件。

Setup 不强制用户完成提取说明的特化；默认版本已可工作。Setup 只展示和引导。

## 实现边界

- **Config** `internal/config`：`Config` 新增 `Standards StandardsConfig` 字段；`StandardsConfig{Guide string}`；`Get/Set/List` 支持 `standards.guide` 键；`Save` 序列化新字段。（v5.4.0 已实现，无需改动）
- **Managed asset** `assets/agents/`：`standards.md` 文件已存在。（v5.4.0 已实现，无需改动）
- **Skills** `assets/skills/`：
  - `flowforge-align/SKILL.md`：step 2 新增读取提取说明、识别适用规范、传递给 Design。
  - `flowforge-solution-design/SKILL.md`：step 1 接收规范清单；step 3 规范合规性作为比较维度；step 4 新增转换 must/must not 写入设计 authority。
  - `flowforge-plan/SKILL.md`：step 1 改为从设计 authority 读取 Standards clauses 机械转写；删除读提取说明和提取逻辑。
  - `flowforge-implement/SKILL.md`：step 1 pre-flight 检查（v5.4.0 已实现，措辞调整 `none found per guide` → `none found`）。
  - `flowforge-review/SKILL.md`：step 3 读卡片内规范（v5.4.0 已实现，无需改动）。
  - `flowforge-setup/SKILL.md`：Section D（v5.4.0 已实现，无需改动）。
  - `_shared/ARTIFACT-CONTRACT.md`：Standards clauses 格式改为「设计 authority 产物，Plan 转写」；新增设计 authority 中 Standards clauses 段落格式。
- **Tracker** `internal/tracker`：不改动。
- **文档** `docs/architecture.md`、`docs/cli-design.md`、`docs/skill-system.md`：更新 Align/Design/Plan 职责描述。

## 验证策略

- Config 层：`standards.guide` get/set/list 往返测试（v5.4.0 已实现）。
- 受管资产：`standards.md` 部署后 `assets verify` 报 `current`（v5.4.0 已实现）。
- Skill 层：手动 walkthrough — Align 提取规范传递 Design、Design 转换写入设计 authority、Plan 从设计 authority 转写到卡片、Implement pre-flight 通过/退回、Review 只查已有。
- 契约一致性：ARTIFACT-CONTRACT 扩展与 Align/Design/Plan skill 改动保持一致。

## 替代方案

- **Plan 做提取和转换（v5.4.0 原方案）**：已否决。规范语义工作在 Plan 做太晚——设计已完成才提取规范，设计可能已违规；且三个 skill 各自理解规范无法保证一致性。
- **固定 schema manifest（layer/scenario → doc+anchor）**：已否决。强迫所有项目按分层/场景组织规范。
- **CLI `standards match` 子命令**：已否决。CLI 解析提取逻辑越界。
- **复用 `KnowledgeSources` 配置**：已否决。list 结构语义不符。
- **新增 ticket tier**：已否决。割裂 Constraints 内约束。
- **机器追踪规范 revision**：已否决。卡片只快照执行时点规范。

Plan 可消费的实现区域为：Align 提取职责、Design 转换+合规设计、Plan 改为转写、ARTIFACT-CONTRACT 更新、文档更新。Config/受管资产/Setup/Implement/Review 已在 v5.4.0 实现，无需改动。每个区域已具备责任、seam、约束和可行验证；ticket 切分及 DAG 仍等待 Plan 阶段。
