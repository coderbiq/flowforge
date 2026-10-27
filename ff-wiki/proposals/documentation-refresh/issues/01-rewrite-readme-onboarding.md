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

# 01: README 新手引导重写

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

README.md 按设计 d-readme-structure 钉定的七节结构重写：新用户仅按 README 从安装走到第一个 check/frontier 循环，每条命令与当前二进制一致。

## Design context

README 现状（115 行）停在 v5 契约叙事：零提及 subagent / agents deploy / model-set / PI / refine-ticket / 轻量模式（research/2026-09-29-doc-refresh/g2-agents-models.md §Summary；current-surface-gaps.md §A）。新结构七节与新手六步主线由 design d-readme-structure 钉定；"一个需求如何实际推进"叙事仍准确，保留并小幅增补。

See the design authority at [文档刷新方案](../design.md#documentation-refresh-design)（d-readme-structure 节）. Requirement authority: [文档刷新需求](../requirements.md#documentation-refresh-requirements)（问题 1、可观察结果 1、验收 S-验收-1）.

## Touch points

- `README.md` — 全文重写（唯一写集）

## Changes

- [x] 1. 重写定位与“整体设计理论”节：三件套定位（Markdown 工件 + skill/角色体系 + Go CLI），四权威模型压缩保留。
  - cmd: `sed -n '1,17p' README.md`
  - exit: 0
  - output: "定位句含“工程 Skill 与 subagent 角色分工推进决策”“四种 agent 宿主”；理论节压缩为两段（四权威清单 + 文件事实推导段），保留原有语义"
  - artifact: README.md
- [x] 2. 新写“快速开始”六步主线：安装（保留现有脚本段）→ `flowforge init`（副作用 + `sync` 别名）→ `flowforge agents deploy`（可选，四宿主，PI 前置 `pi install npm:pi-subagents`）→ 第一张 compact ticket（最小模板）→ `check`/`frontier`/`status` 循环（真实输出语义）→ 推进与收票（轻量/完整双模式、`frontier --pi-workflow`、Completion evidence 必填）。
  - cmd: `grep -n '^### ' README.md | head -8`
  - exit: 0
  - output: "快速开始六小节：1. 安装 / 2. 初始化项目 / 3. 部署 subagent 角色（可选，推荐）/ 4. 写第一张票 / 5. 校验与查看就绪队列 / 6. 推进、执行与收票；model-set 一句在步骤 3 末尾"
  - artifact: README.md
- [x] 3. 增补“一个需求如何实际推进”步骤 4：轻量/完整双模式与 `Fix:` Changes / repair ticket 分流一句话；步骤 5 增补 refine-ticket 前置。
  - cmd: `grep -nE 'refine-ticket|needs-repair|Fix:|轻量' README.md`
  - exit: 0
  - output: "步骤 4 收票闭环句：Fix: Changes 追加回票 + repair ticket/needs-repair 分流 + Completion evidence/--strict 拒绝；步骤 6 推进段含 refine-ticket 契约前置与轻量/完整双模式（L113-119 区域）"
  - artifact: README.md
- [x] 4. 预留“场景速览”节占位（由票 02 回填索引表）。
  - cmd: `grep -n 'scenarios-index' README.md`
  - exit: 0
  - output: "L123: <!-- scenarios-index: 由 documentation-refresh 票 02 回填 -->"
  - artifact: README.md
- [x] 5. 更新“核心文档”链接清单（docs/scenarios.md 由票 02 创建后再验链，本票先按现有文件写全）。
  - cmd: `sed -n '/## 核心文档/,/^## 开发/p' README.md`
  - exit: 0
  - output: "五链接：architecture/skill-system/cli-design/CONTEXT（现存文件）——scenarios.md 链接按票面约定留给票 02，避免死链"
  - artifact: README.md
- [x] 6. 复核：README 中每条命令/flag 与 `go run ./cmd/flowforge <cmd> --help` 逐字一致；`grep -inE 'subagent|agents deploy|model-set|pi-workflow' README.md` 非零命中。
  - cmd: `for kw in subagent 'agents deploy' 'model-set' pi-workflow; do echo "$kw: $(grep -ic "$kw" README.md)"; done`
  - exit: 0
  - output: "subagent: 4 / agents deploy: 1 / model-set: 1 / pi-workflow: 2；命令清单（flowforge agents/check/config/frontier/status/init + flags --include-gaps/--pi-workflow/--strict/--force）均在本会话 --help 实测输出中逐字存在"
  - artifact: README.md

## Constraints

- 不虚构未交付功能；feature 表述以 g1–g4 事实文档为准。
- 中文正文；既有术语稳定（docs_dir、frontier、ticket、authority）。
- 不动 docs/*.md（票 03–05 负责）；不创建 scenarios.md（票 02 负责）。

## Done and verify

`go run ./cmd/flowforge init --help`、`agents deploy --help`、`frontier --help` 与 README 六步中的命令逐字核对；Change 6 的两个 grep 通过；人类通读：新用户 15 分钟路径无断点（S-验收-1）。

## Implementation note

- 七节结构落笔：定位（1 句）→ 理论（2 段压缩）→ 快速开始（六步）→ 场景速览（占位）→ 推进叙事（保留 + 步骤 4/6 增补）→ 核心文档 → 开发。186 行。
- 新增事实全部来自本会话实测或 g1–g4 引用：init 副作用（`init --help`）；极简票进 frontier READY 带 legacy warning、带 metadata 缺五段契约被 gap 排除（/tmp/ff-readme-demo 双向实测）；model-set 一句来自 g2 §4.2。
- 越界回退：初稿把 scenarios.md 链接写进了核心文档清单，按票面约定（Change 5）移除，留给票 02——避免 01 交付态出现死链。
- 未触碰 docs/*.md 与 scenarios.md（Constraints 遵守）。

## Completion evidence

- 交付行为：README 重写为七节入口文档，快速开始六步主线从安装走到收票闭环；subagent/agents deploy/model-set/pi-workflow 四关键词全部进入正文（Change 6 grep 计数 4/1/1/2）。
- 验证方法与观测：Change 1–6 六条证据四元组逐条落盘（cmd/exit/output/artifact）；`./bin/flowforge check --dir docs/proposals/documentation-refresh` exit 0 无 warning；命令与 --help 实测逐字一致。
- 双轴与发现处置：主会话自审（文档票，无代码 diff）——Standards 轴：中文/术语/约束遵守、无虚构功能；Spec 轴：七节结构与 d-readme-structure 一致、S-验收-1 路径完整。发现 2 项并当场修复（model-set 缺失、scenarios.md 死链预防）。
- 偏差：无。
- 实现参考：本 commit（docs(readme): rewrite onboarding for current surface）。

## Review rounds

### Round 1（用户反馈，2026-09-29）

- 发现：快速开始第 4 步起教用户“手写一张票”，与 FlowForge 真实工作模式相反——用户从不手写票；“一个需求如何实际推进”抽象叙事读不檀。证据：real-usage-workflow.md（内部项目会话取证，已脱敏）。
- Fix changes（本轮追加）：
  - [x] F1. 快速开始改为对话驱动主线：说出需求（真实对话示例）→ agent 建 proposal/写票/跑 check+frontier → “继续”推进批次 → 收口；删掉“写第一张票”小节与极简票模板，票的 Markdown 形态降为一句展示（由 flowforge-plan 生成）。
    - cmd: `grep -n '^### ' README.md | sed -n '4,9p' && grep -c '写第一张票\|极简票' README.md`
    - exit: 0
    - output: "步骤 4/5/6 = 在 agent 会话里说出你的需求 / 用“继续”推进批次 / 随时看进度；手写票与极简票模板零命中；步骤 5 含三种真实动作（方向/批准/质疑）与真实话语摘录"
    - artifact: README.md
  - [x] F2. “一个需求如何实际推进”改为案例叙事（虚构通用场景——通知系统重构案，叙事骨架取自真实实践）：用户话语为叙事主体，方法论角色（import/align/solution-design/plan/implement/review）作为步骤注解；保留 DAG/frontier 循环事实。
    - cmd: `grep -n '^### ' README.md | tail -5 && grep -c '用户：' README.md`
    - exit: 0
    - output: "五步 = 拆分立项/澄清与设计/生成执行票/批次执行与审查/收口；虚构案例引语（工作流骨架同真实实践）；每步括号注解 skill 角色；收口含补救票案例（泛化表述）；开头声明“全程用户没有写过一张票、建过一个目录”"
    - artifact: README.md
  - [x] F3. 场景篇场景 1 同步修正（票 02 范围，由票 02 Fix 轮交付）。
    - cmd: `grep -n '^## 场景 1' docs/scenarios.md`
    - exit: 0
    - output: "场景 1 标题改为“说出一个需求，看它变成 proposal”，正文对话驱动（由票 02 Round 1 Fix F1 落地，见该票 Review rounds）"
    - artifact: docs/scenarios.md

- 处置：F1/F2 本票交付；F3 归属票 02（已同步落地）。真实工作流取证固化于 docs/research/2026-09-29-doc-refresh/real-usage-workflow.md（证据源均为脱敏泛称，见该文档）。
- 偏差补充（Rev 1 evidence 中失效表述）：原 Change 2 的“极简票进 frontier”路径已被本轮否定——它量然是 CLI 真实行为，但不是用户路径；证据四元组保留作为 CLI 行为事实，叙事改为 agent 生成票。

### Round 2（用户反馈：脱敏，2026-09-29）

- 发现：对外文档与票面含真实项目信息（业务模块名、项目名、内部工具名、具体取证数字）。
- Fix changes（本轮追加）：
  - [x] F4. README 全部对话引语与案例叙事替换为虚构通用场景（通知系统重构案）；工作流骨架不变；场景速览第 1 行引语同步。
  - [x] F5. real-usage-workflow.md 重写为脱敏版：证据源泛称化、删原文引语与本地存储路径、保留方法论结论与泛化话语类型表。
  - [x] F6. requirements 升 rev 3（脱敏入约束）、design 升 rev 2（快速开始 4-5 步与案例叙事结构改对话驱动钉定）、六票 consumes 同步 3/2。
- 处置：对外文档（README/scenarios.md）与仓内取证文档敏感词零命中（验证见票 02 Round 2 F3 同款扫描）。

### Round 3（用户反馈：漏 standards 定制化 feature，2026-09-29）

- 发现：README 快速开始漏掉初始化后的必做环节——重写 `<docs_dir>/agents/standards.md`（默认为通用模板）；“项目可定制化、尊重项目差异”的 feature 未呈现。
- Fix changes（本轮追加）：
  - [x] F7. README 第 2 步（初始化项目）追加“让 AI 重写 standards.md”关键一步（含对话示例与消费链说明：Align 挑选→Design 转写→Plan 转录）+ 项目可定制化一段（docs_dir/standards.guide/agents.hosts/models_*/test_file_globs 按项目按机器配置）。
    - cmd: `sed -n '/初始化后有一个关键一步/,/尊重项目差异/p' README.md | head -4 && grep -c 'standards' README.md`
    - exit: 0
    - output: "第 2 步含加粗提示句（L49）+ 对话示例（L51，虚构通用场景）+ 可定制化段（L55）；standards 命中 3（standards.md ×2 + standards.guide ×1）另含场景 7 指引"
    - artifact: README.md
  - [x] F8. 场景 7 增补“起步（一次性）”小节（同一对话示例）并把 standards 定制化列入涉及 feature；消费链第 1 条改为指向重写后的说明。
    - cmd: `grep -n '起步（一次性）\|standards.guide' docs/scenarios.md | head -4`
    - exit: 0
    - output: "场景 7 含起步小节与 standards.guide 默认值说明；涉及 feature 首项 = 规范提取说明项目定制化"
    - artifact: docs/scenarios.md
  - [x] F9. requirements 升 rev 4（feature 补录）、六票 consumes 同步 4，check 零 warning。
    - cmd: `./bin/flowforge check --dir docs/proposals/documentation-refresh`
    - exit: 0
    - output: "graph healthy，零 warning"
    - artifact: docs/proposals/documentation-refresh/

---

## Execution detail

### Verified contracts

- README 现状 115 行，结构 7 节（标题/理论/推进/安装/核心文档/开发/许可证）；安装脚本段（curl install.sh / irm install.ps1）仍有效，保留（research/2026-09-29-doc-refresh/g1 §1.4 引 README:72-99）。
- `flowforge init` 真实副作用：部署 `.agents/skills/`、`<docs_dir>/agents/` 规则、创建 `.flowforge/` 配置与 wiki 层级；别名 `sync`（`go run ./cmd/flowforge init --help` 实测，g1 §4.2）。
- `flowforge agents deploy [name]`：部署到 `.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/`；无名参=全部非停用角色（`agents deploy --help` 实测）。
- PI 宿主前置：`pi install npm:pi-subagents`；扩展 `.pi/extensions/flowforge.ts` 随 deploy 写出（g2 §5.2；docs/cli-design.md PI 宿主说明节）。
- `flowforge frontier` flags：`--include-gaps`、`--json`、`--pi-workflow`、`--quiet`、`--strict`（`frontier --help` 实测）。
- `flowforge status` 输出形态：`FlowForge Local Tracker Status: N/N resolved` + per-feature ✅ 列表（本会话实测）。
- 轻量/完整双模式与 `Fix:` / repair 分流、refine-ticket 五段契约前置：g3 §4.2/§5.2（lightweight-execution-contract、ticket-refinement-contract 的 skill 文本引用）。

### Execution scenarios

- Success：新用户依"快速开始"六步在空仓库跑通；每条命令复制即跑；场景速览节为占位（票 02 回填）。
- Failure：README 出现二进制不存在的命令/flag → 违反 Change 6；出现 docs/*.md 或 scenarios.md 的实质修改 → 违反 Constraints。
- Failure：删除仍准确的"一个需求如何实际推进"叙事或安装脚本段 → 违反 Change 1/3 保留要求。

### Expected tests

- `go run ./cmd/flowforge init --help && go run ./cmd/flowforge agents deploy --help && go run ./cmd/flowforge frontier --help` — 输出与 README 命令逐字一致（人工比对，无自动化测试，文档票）。
- `grep -inE 'subagent|agents deploy|model-set|pi-workflow' README.md` — exit 0 且覆盖四关键词。
- `./bin/flowforge check --dir docs/proposals/documentation-refresh` — 无新增 warning（本票不改图结构）。

### Generated artifacts

- `README.md`（仓库入口文档，直接 git 跟踪；无构建产物）。

### Conventions

- must 命令示例以 `--help` 实测为准（源：requirements 约束 4）。
- must 正文中文、术语稳定（源：requirements 约束 1）。
- 变更后运行 `./bin/flowforge check --dir docs/proposals/documentation-refresh`（源：AGENTS.md）。
