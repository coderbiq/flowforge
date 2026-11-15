---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      documentation-refresh-requirements: 4
    design:
      documentation-refresh-design: 2
---

# 02: 场景篇 docs/scenarios.md + README 场景速览

**Blocked by:** 01
**Status:** closed
**Mode:** lightweight

## Delivery

新建 `docs/scenarios.md`：8 个精选场景（痛点背景 → 操作路径 → 涉及 feature → 来源 proposal），并在 README"场景速览"节回填 8 行索引表。读者从任一场景能回答"这解决什么痛点、我该敲什么命令、feature 从哪个 proposal 来"。

## Design context

场景归属决策与 8 场景 × 来源 proposal 映射由 design d-scenarios-landing 钉定；候选文案（每场景 30–60 字故事）已在 g1–g4 工作台文档"候选用户场景"节逐个备好，本票扩写为完整场景而非照抄。依赖理由：README"场景速览"节由票 01 建立占位，本票回填索引表——同文件顺序编辑。

See the design authority at [文档刷新方案](../design.md#documentation-refresh-design)（d-scenarios-landing 节）. Requirement authority: [文档刷新需求](../requirements.md#documentation-refresh-requirements)（问题 2、可观察结果 2、验收 S-验收-2）.

## Touch points

- `docs/scenarios.md` — 新建（主交付）
- `README.md` — 仅"场景速览"节占位替换为索引表 + "核心文档"节补 scenarios.md 链接

## Changes

- [x] 1. 新建 `docs/scenarios.md`：标题 + 一段导览；8 个场景各一节，结构统一为【痛点背景】【操作路径】（真实命令/skill/角色名）【涉及 feature】【来源】（proposal 相对链接）。
  - cmd: `grep -c '^## 场景' docs/scenarios.md`
  - exit: 0
  - output: "8（场景 1–8 齐备，每节四段结构统一）"
  - artifact: docs/scenarios.md
- [x] 2. 场景 1–8 按 d-scenarios-landing 映射表扩写：候选故事取自 g1/g2/g3/g4 各 §.5，命令以各工作台文档 §2 表面变更清单的 file:line 引用为准确认。
  - cmd: `grep -oE 'proposals/[a-z-]+/' docs/scenarios.md | sort -u | wc -l`
  - exit: 0
  - output: "21 个 proposal 路径引用（20 个真实 proposal 全部存在；my-first-feature 为场景 1 代码注释内示例路径非链接）；关键命令均有 g 文档 file:line 依据（model-set 子命令 g2§4.2、--pi-workflow g2§5.2、max_steps/四元组 g3§1.2/§3.2、六级链 g4§4.2）"
  - artifact: docs/scenarios.md
- [x] 3. README“场景速览”占位节替换为 8 行索引表（场景一句话 + 链接）；“核心文档”清单补 `docs/scenarios.md`。
  - cmd: `grep -c 'scenarios.md#场景-' README.md && grep -n '场景导览' README.md`
  - exit: 0
  - output: "8 行索引链接；核心文档首行即 [场景导览](docs/scenarios.md)"
  - artifact: README.md
- [x] 4. 校验：scenarios.md 内所有 proposal 相对链接可解析；README 索引表 8 行链接与 scenarios.md 节锚点一一对应。
  - cmd: `python3 anchor-check（见 Implementation note）`
  - exit: 0
  - output: "headings=8 links=8 mismatched=NONE（GitHub 锚点规则自动生成对照）；20 个来源 proposal 目录逐一存在"
  - artifact: docs/scenarios.md

## Constraints

- 每个场景的命令/角色名/skill 名必须真实存在（以 assets/skills/、assets/subagents/ 目录与 CLI --help 为准）。
- 不重复 README 六步主线细节，场景篇按"用户处境"组织，可交叉引用 README。
- 来源 proposal 链接用相对路径。

## Done and verify

Change 4 的链接校验通过；`ls docs/scenarios.md` 存在；README 索引 8 行与 scenarios.md 8 节锚点一一对应；人类抽读 2 个场景（S-验收-2：痛点/命令/来源三问可答）。

## Implementation note

- 8 场景均从 g 文档候选故事扩写为四段结构（不是照抄 30–60 字素材）；多源场景（3/4/8）按 design 映射表归并，操作路径按实际使用顺序重排。
- 锚点校验脚本：按 GitHub 规则（lowercase、去标点、空格转连字符）从 scenarios.md 标题生成锚点集合，与 README 8 条链接逐一比对，mismatched=NONE。
- 命令事实抽检：model-set list/show/use、frontier --pi-workflow、agents status、max_steps 默认 200、六级优先级链均有 g 文档 file:line 依据；场景 1 的极简票行为（READY+legacy warning）为本会话 /tmp 实测。
- README 只改场景速览与核心文档两节（票 01 建立的结构未动，Constraints 遵守）。

## Completion evidence

- 交付行为：docs/scenarios.md 新建（8 场景四段结构，20 个来源 proposal 相对链接全部可解析）；README 场景速览 8 行索引表 + 核心文档场景导览链接。
- 验证方法与观测：Change 1–4 证据四元组逐条落盘；锚点自动对照 8/8 无失配；`./bin/flowforge check --dir docs/proposals/documentation-refresh` exit 0。
- 双轴与发现处置：主会话自审——Standards 轴：命令/角色名实物为准、相对链接、中文正文；Spec 轴：四段结构统一、S-验收-2 三问可答、场景映射与 d-scenarios-landing 一致。发现 1 项并当场处理（my-first-feature 路径为代码注释非链接，确认为误报非死链）。
- 偏差：无（Rev 1）；Rev 2 见 Review rounds Round 1。
- 实现参考：本 commit（docs: add scenarios guide and README index）。

## Review rounds

### Round 1（用户反馈，2026-09-29，同票 01 Round 1）

- 发现：场景 1 原文教用户“写一张极简票进 frontier”，与真实工作模式相反（用户从不手写票；取证见 real-usage-workflow.md）。
- Fix changes（本轮追加）：
  - [x] F1. 场景 1 改为对话驱动：标题“说出一个需求，看它变成 proposal”；痛点改为“是不是要学一堆 Markdown 模板？都不用”；操作路径换成虚构起步对话（通用场景）+ agent 自动跑 check/frontier；推进方式 = “继续/同意/质疑”；删极简票教学与 legacy warning 路径叙述。
    - cmd: `sed -n '/^## 场景 1/,/^\*\*涉及/p' docs/scenarios.md | head -10 && grep -c '极简票' docs/scenarios.md`
    - exit: 0
    - output: "新标题 + “答案：都不用”开篇 + 起步对话引块（虚构通用场景）+ “继续（最高频）”表述 + 涉及 feature 改“对话驱动立项”；极简票零命中"
    - artifact: docs/scenarios.md
  - [x] F2. README 场景速览第 1 行同步：链接锚点与新标题一致，痛点句改写，顺带修正原错别字“从哪于手”。
    - cmd: `python3 anchor-check（README 8 链接 vs scenarios 8 标题）`
    - exit: 0
    - output: "headings=8 links=8 mismatched=NONE；第 1 行 = [1. 说出一个需求](docs/scenarios.md#场景-1说出一个需求看它变成-proposal) | 是不是要学一堆 Markdown 模板？都不用"
    - artifact: README.md
- 处置：Fix 基于用户反馈与 real-usage-workflow.md 取证；scenarios.md 其余 7 场景无手写票表述（“手写”仅出现于场景 8 指手写 model 字段，语义正确）；Implementation note 中“极简票行为实测”保留作为 CLI 行为事实但不再作为用户路径叙述。

### Round 2（用户反馈：脱敏，2026-09-29，同票 01 Round 2）

- 发现：场景 1 起步对话含真实项目业务信息；导览段残留“第一张票”旧表述。
- Fix changes（本轮追加）：
  - [x] F3. 场景 1 对话替换为虚构通用场景（与 README 快速开始同款通知模块示例，标明虚构）；导览段改“安装 → 说出需求 → frontier 循环”并加“场景对话均为虚构的通用示例”。
  - [x] F4. 全文档敏感词扫描（对外两文件 + 本票所属 proposal 与取证文档）：giis/bytesforce/policy/orders/customer/stagingtable/镜像同步/herdr/NUC/本地用户路径 零命中。
    - cmd: `grep -ricE 'giis|bytesforce|policy|orders|customer|stagingtable|镜像同步|herdr|NUC|qiangbi|/Users/' README.md docs/scenarios.md docs/proposals/documentation-refresh/ docs/research/2026-09-29-doc-refresh/real-usage-workflow.md | grep -v ':0'`
    - exit: 1（全部文件零命中）
    - output: "无输出（每个文件计数为 0）"
    - artifact: README.md + docs/scenarios.md

### Round 3（用户反馈：漏 standards 定制化 feature，2026-09-29，同票 01 Round 3）

- 发现：场景 7 只讲了规范进场链，没讲起点——`init` 部署的 standards.md 默认是通用模板，需一次性重写；项目可定制化未呈现。
- Fix changes（本轮追加）：
  - [x] F5. 场景 7 增补“起步（一次性）”小节：重写对话示例（虚构通用场景，与 README 同款）、`standards.guide` 默认值 `agents/standards.md`；涉及 feature 首项补“规范提取说明项目定制化”。
    - cmd: `grep -c 'standards' docs/scenarios.md && grep -n '起步（一次性）' docs/scenarios.md`
    - exit: 0
    - output: "起步小节在痛点之后、进场链之前；standards 命中 ≥ 5"
    - artifact: docs/scenarios.md

### Round 4（终检票 06 回流，2026-09-29）

- 发现（票 06 F2/F3）：20 个"来源"链接用 `../proposals/**`——从 `docs/scenarios.md` 解析到不存在的仓库根 `proposals/`（实际在 `docs/proposals/`）；其中 2 个目标 `requirements.md` 不存在（documentation-contract-refinement 与 lightweight-execution-contract 只有 `spec.md`）。根因：票 02 首轮验证只 ls 了目录存在性，未验证相对路径解析与目标文件。
- Fix changes（本轮追加）：
  - [x] F6. 全部 `../proposals/` → `proposals/`（同级相对路径）；两个不存在的 requirements.md 改指各自 spec.md。
    - cmd: `python3 link-check（scenarios.md 全部 21 条相对链接逐一 os.path.exists）`
    - exit: 0
    - output: "total=21 broken=0；六份文档全量链接复扫均 NONE"
    - artifact: docs/scenarios.md
- 处置：教训落账——链接验证必须验证"从文档自身目录出发的解析结果"而非目录名存在性；此验证方式已用于票 06 复核脚本。

---

## Execution detail

### Verified contracts

- 8 场景 × 来源映射表在 design d-scenarios-landing（含三场景多 proposal 归并：场景 3 五源、场景 4 三源、场景 8 三源）。
- 候选故事位置：g1 §1.5/§2.5/§3.5/§4.5/§5.5（L62/110/158/202/247）、g2 §1.5–§5.5（L74/118/158/198/245）、g3 §1.5–§5.5（L61/109/175/228/285）、g4 §1.5–§5.5（L52/96/147/198/249）。
- 命令事实源：各工作台文档 §2"表面变更清单"逐条 file:line（如 model-set use 见 g2 §4.2 引 internal/command/modelset.go:25/52/88/152；frontier --pi-workflow 见 g2 §5.2 引 internal/command/frontier.go:139/207-208）。
- skill/角色实物：`assets/skills/` 28 个（含 flowforge-refine-ticket、flowforge-frontend-implement）、`assets/subagents/` 12 个（本会话实测 ls）。
- README"场景速览"占位节由票 01 Change 4 建立（DAG 前置）。

### Execution scenarios

- Success：8 节齐备、结构四段统一；README 索引表 8 行可点；S-验收-2 三问可答。
- Failure：场景中出现不存在的命令/角色（如未交付的 `flowforge daemon`）→ 违反 Constraints；README 除场景速览与核心文档两节外被改动 → 越界。
- Failure：照抄 g 文档候选故事原文而未核对 §2 命令引用 → 30–60 字故事是素材不是成品，违反 Change 2"扩写"。

### Expected tests

- `for l in $(grep -oE 'proposals/[a-z-]+/' docs/scenarios.md | sort -u); do ls docs/$l >/dev/null; done` — exit 0（全部来源 proposal 存在）。
- `grep -c '^## 场景' docs/scenarios.md` — 输出 8。
- `./bin/flowforge check --dir docs/proposals/documentation-refresh` — 无新增 warning。

### Generated artifacts

- `docs/scenarios.md`（新增导览文档）；`README.md`（索引表回填）。

### Conventions

- must 场景内命令/角色名以实物为准（源：requirements 约束 2/4）。
- must 相对链接（源：Constraints）。
- 变更后运行 `./bin/flowforge check --dir docs/proposals/documentation-refresh`（源：AGENTS.md）。
