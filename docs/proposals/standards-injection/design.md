---
flowforge:
  schema: 1
  role: design
  id: standards-injection-design
  revision: 1
  consumes:
    requirements:
      standards-injection: 1
  areas:
    standards-guide:
      revision: 1
      anchor: standards-guide
    ticket-standards-clause:
      revision: 1
      anchor: ticket-standards-clause
    plan-extraction:
      revision: 1
      anchor: plan-extraction
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

需求 authority：[项目规范注入任务卡片需求](requirements.md#standards-injection)，修订版 1。

## 设计区域

### 1. 规范提取说明（standards guide）

<a id="standards-guide"></a>

一份项目内 Markdown 文档，描述本项目的规范文档在哪、以及如何为一张卡片判断它适用哪些规范。结构由项目自定义，FlowForge 不规定其内部结构。

**部署形态**：作为受管资产（managed asset），与 `agents/domain.md`、`agents/issue-tracker.md` 同级，部署到 `<docs_dir>/agents/standards.md`。由 `flowforge init`/`upgrade` 部署内置默认版本。

**默认版本内容**：示范一种通用提取思路（如按文件路径推导分层、按 Touch points 和 Write set 推导涉及的模块与场景），并在正文中以注释引导用户如何根据本项目实际补充。默认版本不假设任何特定项目结构；它提供一个可工作的起点，项目应在其基础上特化。

**配置登记**：新增配置键 `standards.guide`（默认 `agents/standards.md`，相对于 `docs_dir`），指向提取说明文档位置。项目可改指到其他路径。`flowforge config get/set` 支持该键。

**为什么不复用 `KnowledgeSources`**：`KnowledgeSources` 是 list 结构，设计初衷是登记多个知识源元数据。规范提取说明是单一文档引用，用独立标量键 `standards.guide` 更清晰，且避免给 `KnowledgeSources` 增加语义负担。`KnowledgeSources` 保持现状不动。

**为什么提取说明是受管资产而非项目自建**：用户定调要求 FlowForge 内置一份普适通用版本在 setup 时注入。受管资产机制（embed + deploy + verify）已成熟，复用 `copyDir`/`compareManagedAssets` 即可部署与校验。用户修改后 `assets verify` 报告 `drifted`（同其他 agents 文档），`--force` 可刷新。

### 2. Ticket 契约扩展 — 规范陈述书写格式

<a id="ticket-standards-clause"></a>

ARTIFACT-CONTRACT 扩展 ticket 的 hand-off 规则：规范提取产出以 `must be` / `must not be` 规范性陈述写入卡片，附规范源语义链接，不复制整段 rationale（复用现有「下游只摘要本地所需含义并链接」原则）。

**落点不固定**：提取后按每条规范本身的性质决定落点：
- 硬不变量（违反即失败、影响完成判定）→ `## Constraints`，与现有 ticket-specific invariants 并列。
- 约定性（implementer 须读但较软、不直接决定完成）→ Tier 3 `### Conventions`，与现有 non-obvious code convention 并列。

**书写格式**：

```markdown
## Constraints

- <ticket-specific invariant>
- must not <forbidden behavior> — <规范源语义链接>
- must <required behavior> — <规范源语义链接>
- Write set: <allowed directories/files only>
```

```markdown
### Conventions

- <non-obvious code convention in the touch area>
- must <required behavior> — <规范源语义链接>
- must not <forbidden behavior> — <规范源语义链接>
```

链接使用规范源文件的相对路径 + anchor，如 `../docs/112-backend-dependency-rules.md#4-spring-boot-集中管理`。不记 revision，不做机器追踪。

**为什么不新增 tier**：现有两层（Constraints 硬、Conventions 软）已天然对应规范的两类性质。新增 tier 会割裂 Constraints 内的「需求约束」与「规范约束」，且 ARTIFACT-CONTRACT 现有原则是「复杂度决定工件形态，不决定固定流程状态」，新增 tier 违背此原则。

### 3. Plan 提取职责

<a id="plan-extraction"></a>

`flowforge-plan` 的「Resolve effective authority」一步扩展为同时 resolve 规范提取说明：

1. 读取 `standards.guide` 配置指向的提取说明文档。
2. 按提取说明中描述的逻辑，为本卡片定位适用规范（逻辑由项目自定义，Plan 只执行）。
3. 提取后按每条规范性质分流到 Constraints 或 Conventions。
4. 自检：确认本卡片的适用规范已提取并写入对应 tier。缺失时在 ticket 中留空并标注 `standards: pending`（见下文 pre-flight 判定）。

Plan 不做规范正确性判断；它只按提取说明执行提取。如果提取说明本身不完整或无法为某卡片定位规范，Plan 在 ticket Constraints 中写明 `standards: none found per guide`，表示已尝试提取但无适用规范——这与「未提取」区分开。

**为什么不加 CLI 诊断**：提取是语义工作，不是确定性计算。CLI 不解析提取说明，不校验提取结果是否「正确」，不提供按路径匹配规范的确定性接口。这符合「CLI 负责图计算（高确定性）」原则。

### 4. Implement pre-flight 检查

<a id="implement-preflight"></a>

`flowforge-implement` 的 step 1「Resolve the effective specification」增加 pre-flight：

- 读取卡片 Constraints，检查是否有以下三种状态之一：
  - `must`/`must not` 规范陈述已存在 → 通过，继续执行。
  - `standards: pending` → 规范提取未完成，退回 Plan。
  - `standards: none found per guide` → 已尝试提取但无适用规范，通过，继续执行。
- 缺少以上任何标记且卡片涉及代码变更（有 Changes 和 Write set）时，视为规范提取缺失，退回 Plan。

**判定标准**：卡片只要有 Changes 和 Write set 就必须有规范提取标记。纯文档类卡片（无 Write set）不要求。这避免对不需要规范的轻量卡片产生噪声。

**为什么退回 Plan 而非自行提取**：Implement 是 lightweight 模式时无法搜索代码库或做语义判断；full 模式虽然能，但提取归属 Plan，Implement 不应越权。退回 Plan 让 Plan 补充提取后重新发布。

### 5. Review Standards 轴职责收窄

<a id="review-standards"></a>

`flowforge-review` 的 step 3「Identify the standards sources」改为：

1. 读取卡片中已注入的规范陈述（Constraints 和 Conventions 中的 `must`/`must not` 条目）。
2. Standards 轴只检查实现是否符合卡片中已有的规范。
3. 不重新从项目规范源提取、不检查 Plan 是否遗漏了本应注入的规范。

**Step 3 原有逻辑保留**：Review 仍会查找仓库中 documented coding standards 文件（如 `CODING_STANDARDS.md`）和 smell baseline——这些是通用代码质量基线，与本需求的「项目规范注入」正交。两者并存：卡片内规范是项目特定规范，通用 standards 查找是跨项目代码质量基线。

### 6. Setup 部署

<a id="setup-deploy"></a>

`flowforge-setup` 新增 Section D「Standards」：

1. 检查 `<docs_dir>/agents/standards.md` 是否存在（已由 `init` 部署受管资产）。
2. 向用户展示默认版本内容摘要。
3. 引导用户根据本项目实际补充：规范文档在哪、如何为卡片判断适用规范。
4. 确认后写入。用户可在后续随时编辑该文件。

Setup 不强制用户完成提取说明的特化；默认版本已可工作。Setup 只展示和引导。

## 实现边界

- **Config** `internal/config`：`Config` 新增 `Standards StandardsConfig` 字段；`StandardsConfig{Guide string}`；`Get/Set/List` 支持 `standards.guide` 键；`Save` 序列化新字段。
- **CLI** `internal/command`：`init.go` 无需改动（受管资产部署已覆盖 `assets/agents/` 全目录）；`config.go` 命令无需改动（`Get/Set` 在 service 层处理新键）。
- **Managed asset** `internal/command/assets/agents/`：新增 `standards.md` 文件。`assets_compare.go` 的 `compareManagedAssets` 已遍历 `agents/` 目录，无需改动。
- **Skills** `.agents/skills/`：
  - `flowforge-setup/SKILL.md`：新增 Section D。
  - `flowforge-plan/SKILL.md`：扩展 step 1 resolve 提取说明；扩展 step 3 写入规范陈述；新增提取自检。
  - `flowforge-implement/SKILL.md`：step 1 增加 pre-flight 检查。
  - `flowforge-review/SKILL.md`：step 3 改为读卡片内规范 + 保留通用 standards 查找。
  - `_shared/ARTIFACT-CONTRACT.md`：扩展 hand-offs 规则，定义规范陈述书写格式。
- **Tracker** `internal/tracker`：不改动。不做机器追踪，不扩展语义诊断。
- **文档** `docs/architecture.md`、`docs/cli-design.md`、`docs/skill-system.md`：记录 `standards.guide` 配置项与提取说明定位。

## 验证策略

- Config 层：`standards.guide` get/set/list 往返测试；默认值 `agents/standards.md`；空值拒绝。
- 受管资产：`standards.md` 部署后 `assets verify` 报 `current`；修改后报 `drifted`；`--force` 刷新。
- Skill 层：手动 walkthrough — 部署默认提取说明、Plan 提取规范写入卡片、Implement pre-flight 通过/退回、Review 只查已有。
- 契约一致性：ARTIFACT-CONTRACT 扩展与 Plan/Review/Implement skill 改动保持一致。

## 替代方案

- **固定 schema manifest（layer/scenario → doc+anchor）**：已否决。强迫所有项目按分层/场景组织规范，与「项目自定义提取逻辑」需求冲突。
- **CLI `standards match` 子命令**：已否决。CLI 解析提取逻辑越界，违反「CLI 不做内容 API」边界。
- **复用 `KnowledgeSources` 配置**：已否决。list 结构语义不符，独立标量键更清晰。
- **新增 ticket tier `### Applicable standards`**：已否决。割裂 Constraints 内约束，违背「复杂度决定工件形态」原则。
- **机器追踪规范 revision**：已否决。卡片只快照执行时点规范，不记 revision。

Plan 可消费的实现区域为：Config 新增字段、受管资产 standards.md、Plan 提取职责、Implement pre-flight、Review 职责收窄、ARTIFACT-CONTRACT 扩展、Setup Section D、文档更新。每个区域已具备责任、seam、约束和可行验证；ticket 切分及 DAG 仍等待 Plan 阶段。
