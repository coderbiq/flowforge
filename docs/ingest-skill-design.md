# 知识策展 SKILL 设计（外部导入 + 归档）

> 状态：draft
> 目标：定义 `flowforge-curate` SKILL 的工作流，将外部资料或 proposal 中的可复用知识沉淀到 library。

## 1. 定位

`flowforge-curate` 是统一的知识策展 SKILL。它处理两种知识来源的入库：

- **Mode A：外部资料导入** — 将长文参考资料（SKILL 文件、工程指南、遗留文档、API 文档等）拆解为 library 原子知识卡片
- **Mode B：proposal 归档** — 从已完成的 proposal 卡片网络中筛选可复用知识，合成到 library

两种模式在"提取知识单元"之后汇入完全相同的共享流程：聚类 → 去重 → 审查计划 → CLI 写入。两者使用**完全相同的 CLI 命令**。

## 2. Description（激活条件）

```
当用户要求将知识沉淀到 library 知识库时激活。包括两种场景：

Mode A — 外部资料导入：
- "把这个文档导入知识库"
- "分析 docs/references/xxx.md 并提取知识"
- "把这份 API 设计规范拆成知识卡片"
- "从这份工程指南中提取可复用的规则"

Mode B — proposal 归档：
- "归档这个 proposal"
- "把 CR26061401 的结论沉淀到 library"
- "这个 proposal 做完了，知识入库"
- "关闭提案并提取可复用知识"

不应激活：
- 用户只是创建单张卡片（使用 card create 命令）
- 用户讨论 proposal 设计（使用 flowforge-design）
- 用户执行实现任务（使用 flowforge-implement）
- 用户反馈问题（使用 flowforge-feedback）
```

## 3. 核心流程

```
模式判断：用户提供文件路径？→ Mode A / 提供 proposal ID？→ Mode B

  Mode A: 外部资料导入            Mode B: proposal 归档
  ┌────────────────────┐        ┌──────────────────────┐
  │ 读取长文文件        │        │ 扫描 proposal 卡片    │
  │ 提取知识单元        │        │ 筛选可复用候选        │
  │ 用自己的话重述      │        │ 评估知识类型          │
  │ 标注来源出处        │        │ 标注来源卡片          │
  └────────┬───────────┘        └──────────┬───────────┘
           │                               │
           └───────────┬───────────────────┘
                       │
                       v
              共享流程（两种模式在此汇合）
  ┌──────────────────────────────────────────────────────┐
  │ 阶段 3：聚类                                         │
  │   - 将知识单元按概念聚类（不按来源结构）               │
  │   - 每个聚类对应一个主题                             │
  │   - 识别知识单元间的依赖关系                         │
  │                                                      │
  │ 阶段 4：生成审查计划                                 │
  │   - 拟议 STR 索引卡结构                              │
  │   - 拟议原子知识卡列表（标题、类型、摘要、所属 STR）   │
  │   - 标注重复/合并候选/警告                           │
  │   - 用户审查确认                                     │
  └──────────────────────────────────────────────────────┘
                       │
                       v
  ┌──────────────────────────────────────────────────────┐
  │ 阶段 5：生成执行计划并写入计划卡                      │
  │   - 将审查计划转化为可分批执行的条目列表               │
  │   - 搜索已有 library（去重），标记 create/merge/skip  │
  │   - 为每个条目分配批次号（每批 5-10 个条目）           │
  │   - 写入计划卡到 library（持久化状态）                │
  │   - 输出第一批的执行清单                             │
  └──────────────────────────────────────────────────────┘
                       │
                       v
            ┌─────────────────────┐
            │  阶段 6：分批执行    │  ← 每轮激活执行一批
            │                     │
            │  ┌───────────────┐  │
            │  │ 6a. 读取计划卡 │  │
            │  │ 找出下一批待处 │  │
            │  │ 理条目        │  │
            │  └───────┬───────┘  │
            │          v          │
            │  ┌───────────────┐  │
            │  │ 6b. 执行本批  │  │
            │  │ - card create │  │
            │  │ - card link   │  │
            │  │ - structure   │  │
            │  │   add         │  │
            │  └───────┬───────┘  │
            │          v          │
            │  ┌───────────────┐  │
            │  │ 6c. 更新计划卡│  │
            │  │ 标记本批条目  │  │
            │  │ 为 done       │  │
            │  └───────┬───────┘  │
            │          v          │
            │  ┌───────────────┐  │
            │  │ 6d. 汇报进度  │  │
            │  │ 剩余 N 条待处 │  │
            │  │ 理，继续？    │  │
            │  └───────────────┘  │
            │                     │
            │  有剩余 → 等待用户  │
            │  重新激活 SKILL     │
            │  全部完成 → 阶段 7  │
            └─────────────────────┘
                       │
                       v
              Mode B only：提案收尾
  ┌──────────────────────────────────────────────┐
  │ 阶段 7：移动 proposal 到 completed            │
  │   - proposal archive 移动目录                │
  │   - 更新 proposal root card 状态             │
  │   - 保留 proposal 原始卡片作为追溯证据        │
  └──────────────────────────────────────────────┘
```

## 4. 阶段详解

### 4.1 模式判断

根据用户输入判断进入哪个模式：

- 用户提供了文件路径 → **Mode A**（外部资料导入）
- 用户提供了 proposal ID → **Mode B**（proposal 归档）

### 4.2 Mode A：提取知识单元（外部资料）

**Step 1：读取资料，登记来源**

```bash
# SKILL 直接读取文件（这是 SKILL 的职责，不是 CLI 的）
```

应记录：
- 来源文件路径
- 文档标题
- 文档类型（工程指南 / API 文档 / 设计规范 / SKILL 文件 / 其他）

**Step 2：提取知识单元**

识别标准：文档中一个可独立存在的知识断言。它满足：
- 脱离文档上下文仍然有意义
- 表达一个明确的观点、规则、事实或设计
- 可以被其他卡片独立引用

重述原则：
- 必须用自己的话重新表述，不能复制粘贴原文
- 保留来源出处（章节/段落引用）
- 对于 convention 类型，需要提取为可执行的规则形式

### 4.3 Mode B：提取知识单元（proposal 归档）

**Step 1：扫描 proposal 状态**

```bash
flowforge proposal inspect <proposal-id>
flowforge context proposal --proposal <proposal-id>
```

了解 proposal 的完整卡片网络：有多少 requirement、design、decision、finding、log，以及它们之间的链接关系。

**Step 2：筛选可复用候选**

遍历 proposal 卡片，按以下规则筛选：

| 卡片类型 | 筛选规则 | 典型操作 |
|----------|----------|----------|
| `finding` | 是否表达可复用的事实/风险/注意事项 | create 或 merge |
| `decision` | 是否是仍有效的架构决策 | create 或 merge |
| `design` | 是否包含可复用的设计模式/约束 | 提炼为 convention/pattern/decision |
| `convention` | 提案中形成的临时约定是否需要持久化 | create 为 library convention |
| `log` | 默认不进入 library，只作为证据链 | skip（除非包含关键知识） |
| `requirement` | 默认不进入 library | skip |
| `task` | 默认不进入 library | skip |
| `ROOT` | 不进入 library | skip |
| `STR` | 不进入 library | skip |

**Step 3：评估知识类型**

对每个筛选出的候选卡片，判定其最适合的 library 知识类型。提案中的 `design` 卡可能被合成为 `convention`、`decision` 或 `pattern`，取决于其内容特征。

### 4.4 共享：聚类（两种模式通用）

**聚类策略**：
- 不按原文的章节结构或 proposal 的卡片组织方式聚类
- 按概念/主题聚类：相同概念在不同来源中出现，聚类把它们聚合
- 对于工程文档，按 layer / domain / module 等维度聚类
- 每个聚类对应一个 STR 索引卡，含 3-15 个知识单元

**聚类输出**：
```
Cluster: "分页查询规范"
  - convention: 分页查询必须传 pageSize
  - convention: pageSize 上限 100
  - finding: 旧版 API 的 pageSize 默认值是 20
  - example: 分页查询的标准响应格式

Cluster: "错误处理"
  - convention: 所有 API 返回统一错误信封
  - decision: 错误码使用 5 位数字编码
  - fact: 错误信封格式为 {code, message, details}
```

### 4.5 共享：生成审查计划

审查计划是审查点，**此时不写入任何卡片**。格式：

```markdown
## 审查计划

### 来源
- 文件：docs/references/api-design-guide.md（Mode A）
- 或 Proposal：CR26061401（Mode B）
- 提取知识单元数：23

### 拟议 STR 索引卡

| STR ID | 标题 | 入口卡片数 |
|--------|------|-----------|
| STR-API-PAGINATION | 分页查询规范 | 4 |
| STR-API-ERROR | 错误处理规范 | 3 |
| STR-API-AUTH | 认证与授权 | 5 |

### 拟议知识卡

| 类型 | 标题 | 摘要 | 所属 STR | 操作 |
|------|------|------|----------|------|
| convention | 分页查询必须传 pageSize | 所有列表查询接口必须包含 pageSize 参数 | STR-API-PAGINATION | create |
| convention | pageSize 上限 100 | pageSize 参数最大值为 100 | STR-API-PAGINATION | create |
| finding | 旧版 API pageSize 默认值 | 旧版 API 使用 pageSize=20 作为默认值 | STR-API-PAGINATION | create |
| example | 标准分页响应格式 | {data, total, pageSize, page} | STR-API-PAGINATION | create |

### 重复/合并候选

| 拟议卡片 | 匹配已有卡片 | 建议 |
|----------|-------------|------|
| 统一错误信封 | CONV-003 错误处理规范 | merge：追加到 CONV-003 |

### 警告

- "API 网关限流策略"知识单元过大，建议拆分为 3 张卡片
- "请求参数校验规则"过于模糊，建议补充具体校验规则

### 下一步

审查确认后执行写入阶段，将以上卡片写入 library。
```

### 4.6 共享：生成执行计划并写入计划卡

审查通过后，SKILL 将审查计划转化为可分批执行的条目列表，并写入一张**计划卡**持久化。

#### 为什么需要计划卡

一次导入可能产生 20-50+ 个知识单元。LLM 的上下文窗口无法在一轮对话中完成全部写入，同时需要：
- 状态持久化：中断后可恢复
- 进度可见：用户知道完成了多少、还剩多少
- 可审查：每批执行前可再次确认
- 可回退：出错的批次可以重新处理

计划卡用 FlowForge 自身的卡片系统解决这些需求。

#### 计划卡结构

计划卡是一张特殊内容卡（建议使用 `FIND` 类型，kind: `curation-plan`），写入 library：

```markdown
---
id: FIND-xxx
title: "策展计划：<来源名称>"
type: finding
status: active
tags: [curation-plan]
links:
  - target: <source-card-id>       # Mode A: 来源登记卡 / Mode B: ROOT-<proposal>
    relation: derived-from
created: 2026-06-20
---

# 策展计划：<来源名称>

## 来源

- 文件：docs/references/api-design-guide.md
- 或 Proposal：CR26061401
- 知识单元总数：23
- 已处理：0 / 23
- 已跳过：0

## 计划条目

### 批次 1（条目 1-8）

- [ ] CONV / 分页查询必须传 pageSize / STR-API-PAGINATION / create
- [ ] CONV / pageSize 上限 100 / STR-API-PAGINATION / create
- [ ] FIND / 旧版 API pageSize 默认值 / STR-API-PAGINATION / create
- [ ] EX / 标准分页响应格式 / STR-API-PAGINATION / create
- [ ] CONV / 所有 API 返回统一错误信封 / STR-API-ERROR / merge:CONV-003
- [ ] DEC / 错误码使用 5 位数字编码 / STR-API-ERROR / create
- [ ] FACT / 错误信封格式定义 / STR-API-ERROR / create
- [ ] CONV / Token 刷新策略 / STR-API-AUTH / create

### 批次 2（条目 9-16）

- [ ] CONV / 权限校验中间件 / STR-API-AUTH / create
- [ ] DEC / JWT 过期时间 2 小时 / STR-API-AUTH / create
...

### 批次 3（条目 17-23）

- [ ] ...
```

#### 批次大小

每批 5-10 个条目。依据：
- 每批创建 5-10 张卡片 + 链接 + STR 操作，在 LLM 单轮上下文内可控
- 批处理完成后需要搜索已有 library 去重，下一批的上下文不会受到前一批的污染
- 用户审查负担可控：每批结果 5-10 条，一眼可扫完

#### 搜索去重与操作标记

每个条目需要标记操作类型：

| 操作 | 含义 | CLI 命令 |
|------|------|----------|
| `create` | 创建新卡片 | `card create` |
| `merge:CONV-xxx` | 合并到已有卡片 | `card read` + `card update` + `card link` |
| `skip:reason` | 跳过（过大/模糊/不适用） | 无，仅记录 |

### 4.7 分批执行机制

`flowforge-curate` 的每次激活执行**一个批次**。这是核心循环：

```
用户激活 SKILL（"继续策展" 或 "flowforge-curate"）
  │
  ├─ 检查是否存在进行中的计划卡
  │   ├─ 没有 → 从阶段 1 开始（新策展任务）
  │   └─ 有 → 进入阶段 6（继续执行）
  │
  v
6a. 读取计划卡
  - 找到下一个未完成的批次
  - 读取该批次的所有条目
  - 如果所有批次已完成 → 进入阶段 7
  │
  v
6b. 执行本批
  - 对每个条目：
    - create: card create --type <type> --title "<title>" --status draft --tags "..."
    - merge: card read <target> --summary，然后 card update 追加内容
    - skip: 记录跳过原因
  - 创建本批涉及的 STR 索引卡（如果尚未创建）
  - 建立卡片间链接：card link
  - 维护 STR 索引：structure add
  │
  v
6c. 更新计划卡
  - 将本批已处理的条目标记为 [x]
  - 更新 "已处理" 计数
  - 记录已创建的卡片 ID（用于追溯）
  │
  v
6d. 汇报进度
  - 输出本批完成摘要
  - 剩余批次数和条目数
  - 如果还有剩余：提示用户说"继续"来激活下一批
  - 如果全部完成：进入阶段 7
```

#### 重新激活 SKILL 的入口

用户可以通过以下方式继续策展：

- "继续策展" — 读取计划卡，执行下一批
- "继续导入 <计划卡 ID>" — 指定计划卡继续
- "flowforge-curate" — SKILL 激活后自动检测进行中的计划卡

#### 中断恢复

- 如果一批执行到一半中断（例如 LLM 上下文溢出），下一轮激活时：
  1. 读取计划卡，找到当前批次
  2. 检查哪些条目已创建卡片（通过 `card search` 按标题搜索）
  3. 只处理尚未创建的条目
  4. 更新计划卡状态

#### 计划卡完成后的清理

全部批次完成后：
- 更新计划卡 status 为 `done`
- 计划卡保留在 library 中作为导入记录
- 输出最终报告（包含所有已创建卡片的 ID 列表）

### 4.8 分批执行的 CLI 命令

每批执行使用的 CLI 命令（与之前设计相同，但按批次执行）：

```bash
# 1. 创建本批涉及的 STR 索引卡（如果尚未创建）
flowforge card create --type structure --title "分页查询规范" --status active

# 2. 创建本批原子知识卡（status: draft）
flowforge card create --type convention --title "分页查询必须传 pageSize" \
  --status draft --tags "layer:api,scenario:page-query"
flowforge card create --type convention --title "pageSize 上限 100" \
  --status draft --tags "layer:api,scenario:page-query"
# ... 本批其余条目

# 3. 建立卡片间链接
flowforge card link CONV-xxx CONV-yyy --relation references

# 4. 将卡片加入 STR 索引
flowforge structure add --index STR-API-PAGINATION --card CONV-xxx
flowforge structure add --index STR-API-PAGINATION --card CONV-yyy

# 5. 合并（如果需要）
flowforge card read CONV-003 --summary
flowforge card update CONV-003          # 追加新内容
flowforge card link CONV-003 FIND-xxx --relation derived-from

# 6. 更新计划卡：标记本批条目为 [x]
flowforge card update FIND-xxx-plan     # 更新计划卡正文

# 7. 重建索引
flowforge index rebuild
```

### 4.9 Mode B only：提案收尾

仅当 Mode B（proposal 归档）时，且所有批次完成后执行：

```bash
# 移动 proposal 到 completed
flowforge proposal archive <proposal-id>

# 更新 proposal root card 状态
flowforge card update ROOT-<proposal-id> --status completed
```

归档关键规则：
- proposal 内原始卡片保留在 `03-completed/` 下，作为历史追溯
- library 卡通过 `derived-from` 关系链接回原始 proposal 卡片
- 不直接移动原始 proposal 卡片到 library
- log 卡默认不进入 library，只作为证据链保留

### 4.10 输出报告

#### 每批完成后的进度报告

```markdown
## 批次 1/3 完成

### 本批创建
- STR-API-PAGINATION：分页查询规范（索引卡）
- CONV-xxx：分页查询必须传 pageSize
- CONV-yyy：pageSize 上限 100
- FIND-xxx：旧版 API pageSize 默认值
- EX-xxx：标准分页响应格式
- ...（共 8 张）

### 进度
- 已完成：8 / 23
- 剩余批次：2（共 15 条）

### 下一步
- 说"继续"处理下一批
```

#### 全部完成后的最终报告

```markdown
## 策展完成

### 来源
- 文件：docs/references/api-design-guide.md（Mode A）
- 或 Proposal：CR26061401（Mode B）

### 新增卡片
- STR-API-PAGINATION：分页查询规范（索引卡）
- STR-API-ERROR：错误处理规范（索引卡）
- CONV-xxx：分页查询必须传 pageSize
- CONV-yyy：pageSize 上限 100
- FIND-xxx：旧版 API pageSize 默认值
- EX-xxx：标准分页响应格式
- ...（共 20 张）

### 合并卡片
- CONV-003：追加"统一错误信封"内容

### 跳过
- "API 网关限流策略"：单元过大，建议单独处理

### 提案收尾（Mode B only）
- proposal CR26061401 已移动到 03-completed/

### 下一步
- 审查 draft 卡片，确认后提升为 active
- 处理跳过的过大知识单元
```

## 5. 与 proposal 归档 SKILL 的合并分析

### 5.1 差异清单

| 维度 | ingest | archive |
|------|--------|---------|
| 输入 | 长文文件 | proposal 卡片网络 |
| 提取方式 | 从文本中识别知识断言，用自己的话重述 | 从已有卡片中筛选可复用候选 |
| 聚类输入 | 文本提取出的知识单元 | 已有卡片的 type/tags/links |
| 来源链接 | 指向原文章节/段落 | 指向 proposal 原始卡片 |
| 提案操作 | 无 | 移动 proposal 目录到 completed，更新 proposal 状态 |
| CLI 命令 | 完全相同 | 完全相同 |
| 审查→写入流程 | 相同 | 相同 |

### 5.2 合并可行性分析

**实际上只有两个本质差异**：

1. **提取阶段的输入不同**：ingest 从文本中解析，archive 从卡片中筛选。但两种输入的终点相同——都是一组待分类的"知识单元"。ingest 的提取是"阅读→识别→重述"，archive 的提取是"扫描→筛选→评估"。核心操作都是：**判断一条知识是否可复用 + 判定它的类型**。

2. **archive 多了一个提案收尾操作**：移动 proposal 目录到 completed，更新 proposal 状态。这是纯 CLI 操作，不涉及知识判断。

**其余步骤完全重合**：聚类 → 类型判定 → 去重 → 生成审查计划 → 用户确认 → CLI 写入 → 重建索引。调研的所有方法论（Zettelkasten 文献笔记→永久笔记、Progressive Summarization、Evergreen Notes 的收集→处理→连接）对两种场景同样适用。

**结论：可以合并为一个 SKILL**，差异收敛为入口模式的选择。

### 5.3 合并后的 SKILL 模型

```
flowforge-curate（知识策展 SKILL）

两个入口模式：
  Mode A: 外部资料导入
    输入：文件路径 → 读取文本 → 提取知识单元 → 用自己的话重述

  Mode B: proposal 归档
    输入：proposal ID → 扫描 proposal 卡片 → 筛选可复用候选

共享流程（两种模式在此汇合）：
   → 聚类：按概念组织知识单元
   → 生成审查计划：拟议 STR 索引 + 原子知识卡 + 合并/跳过/警告
   → 用户审查确认
   → 生成计划卡：分批组织，写入 library 持久化
   → 分批执行：每轮激活处理一批（5-10 条）
     → card create → card link → structure add → 更新计划卡进度
   → 全部批次完成后：index rebuild
   → (Mode B only) 提案收尾：proposal archive 移动目录

输出：导入/归档报告
```

### 5.4 合并的收益

| 收益 | 说明 |
|------|------|
| 减少 SKILL 数量 | 从 2 个减为 1 个，降低 activation 冲突风险 |
| 方法论统一 | 聚类、去重、审查、写入的流程只维护一处 |
| 知识类型一致 | 两种入口产出的卡片类型相同，不会出现两套 taxonomy |
| CLI 命令复用 | 已经在用同一套命令，合并后 SKILL 本体更薄 |
| 用户心智模型统一 | "往知识库里加东西"是一个动作，不管是加文档还是加 proposal |

### 5.5 合并的风险和对策

| 风险 | 对策 |
|------|------|
| 入口模式判断错误 | 根据用户输入判断：给了文件路径→Mode A，给了 proposal ID→Mode B |
| SKILL description 覆盖两种场景可能模糊 | 写清楚两种激活场景，加上反例 |
| 两种模式的处理逻辑混在一起 | SKILL 内部分支清晰，共享流程独立成段 |
| archive 特有的提案收尾操作被遗漏 | 在共享流程结束后明确标注 Mode B only 步骤 |

## 6. SKILL 本体编写建议

真正的 `assets/skills/flowforge-curate/SKILL.md` 应保持短小，作为薄适配器。建议只包含：

1. 激活后先检查是否存在进行中的计划卡
   - 有计划卡 → 进入分批执行模式（阶段 6）
   - 无计划卡 → 进入新建策展流程
2. 根据用户输入判断模式：文件路径→Mode A，proposal ID→Mode B
3. Mode A：读取文件，提取知识单元，用自己的话重述
4. Mode B：扫描 proposal 卡片，筛选可复用候选，评估知识类型
5. 共享：聚类 → 生成审查计划 → 用户审查 → 生成计划卡
6. 分批执行：每轮激活处理一批（5-10 条），更新计划卡进度
7. Mode B only：全部完成后移动 proposal 到 completed
8. 只使用 CLI 原子操作（card create / card link / structure add / index rebuild）
9. 不直接操作文件，不自行构建索引
10. 默认创建 `status: draft` 卡片
11. 每轮必须输出审查报告或进度报告，不直接写入

详细方法论和卡片模板应放在 docs 中，不塞进 SKILL 本体。

## 7. 反推 CLI 能力需求

curate SKILL 依赖的 CLI 原子操作已全部存在或已在 MVP 计划中：

| 命令 | 用途 | 状态 |
|------|------|------|
| `card create --type --status --tags` | 创建卡片 | MVP |
| `card read --summary/--section` | 读取已有卡片 | MVP |
| `card update` | 更新卡片内容/状态 | MVP |
| `card link --relation` | 建立类型化链接 | MVP |
| `card search --scope library` | 搜索已有卡片 | MVP |
| `structure add --index --card` | 维护 STR 索引 | MVP |
| `index rebuild` | 重建 sqlite 索引 | MVP |
| `proposal archive <id>` | 移动 proposal 到 completed | MVP |
| `proposal inspect <id>` | 扫描 proposal 状态 | MVP |

**不需要新增任何 CLI 命令**。curate 的"智能"部分完全由 SKILL 承担。