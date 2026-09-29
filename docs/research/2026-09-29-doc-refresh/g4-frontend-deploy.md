# g4 — 前端生命周期与部署/阻断事实提取（文档更新消费用）

> 分组：`frontend-standards-lifecycle`、`frontend-vision-loop-poc`、`flash-worksurface-expansion`、
> `deploy-artifact-localization`、`blocked-evidence-persistence`
> 方法：逐文件引用事实（文件路径 + 行号 或命令输出），不做跨组综合，不提改进建议。
> 全部"表面变化"均已在当前仓库实物中复核（下方 `assets/…` / `internal/…` 行号为实测）。

---

## 1. frontend-standards-lifecycle

### 1.1 Feature 一句话
前端规范从此作为一等公民贯穿 align → design → plan → implement → review 全链：编排侧把前端条款转录进票，执行与复核交给视觉模型角色，并可用截图 + 量化输出做条款级核对。

### 1.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| 新 skill（v2） | `flowforge-frontend-implement` v2：槽位契约化（design-tokens / component-entry / counter-examples / verification-entry）、条款优先、前置自检 | `docs/proposals/frontend-standards-lifecycle/design.md:96`；实物 `assets/skills/flowforge-frontend-implement/SKILL.md:1-25`（description 含"follow the ticket's transcribed clauses first, consult slot slices by anchor only"+ 三级权威序 L11-18） |
| 新 agent 角色 | `flowforge-frontend-reviewer`（视觉模型 pin）新增 | `docs/proposals/frontend-standards-lifecycle/design.md:97`、`design.md:44`（Review·Spec+视觉·前端 一行）；实物 `assets/subagents/flowforge-frontend-reviewer.md:1-16`（`model_profile: tool-capable-read-only`、`permission: review-read-only`、`returns_to: [flowforge-reviewer]`） |
| 新 agent 角色 | `flowforge-frontend-implementer`（视觉模型，执行前端票） | `docs/proposals/frontend-standards-lifecycle/design.md:42`、`design.md:24`；实物 `assets/subagents/flowforge-frontend-implementer.md:1-16` |
| 既有 skill 增补 | `flowforge-align` 增补 frontend scope 识别 + manifest 引用 | `docs/proposals/frontend-standards-lifecycle/design.md:98` |
| 既有 skill 增补 | `flowforge-solution-design` 增补前端条款转录义务 + verification-entry 声明 | `docs/proposals/frontend-standards-lifecycle/design.md:99` |
| 既有 skill 增补 | `flowforge-review` 增补前端票三轴编排（视觉轴指向 frontend-reviewer） | `docs/proposals/frontend-standards-lifecycle/design.md:100` |
| AGENTS.md 模板 | 路由表 + 生命周期参与说明 | `docs/proposals/frontend-standards-lifecycle/design.md:101` |
| 槽位 manifest 入口（项目侧） | `.agents/references/frontend/manifest.md` 为唯一入口，声明四个槽位（design-tokens / component-entry / counter-examples / verification-entry） | `docs/proposals/frontend-standards-lifecycle/design.md:60`（代码块 L59-66） |
| Model matrix（新角色档位） | 执行·前端 = flash(视觉) gemini 系；Review·Standards 轴 = reviewer-lite flash(文本)；Review·Spec 轴·后端 = 旗舰；Review·Spec+视觉·前端 = frontend-reviewer flash(视觉) | `docs/proposals/frontend-standards-lifecycle/design.md:36-43`（矩阵表） |
| config 键 | 无新 config 键；沿用既有具名 pin（`hosts/models_by_name` 不配即不激活该角色） | `docs/proposals/frontend-standards-lifecycle/design.md:104`（"无 Go 变更"行）、`design.md:107`（风险 3） |
| 流程契约 | 混合需求拆票：Write set 命中前端目录→前端票，其余→后端票；DAG 边连接；单一前端票无后端耦合不拆 | `docs/proposals/frontend-standards-lifecycle/design.md:81-84` |
| 降级语义 | 缺槽位/缺图像能力/缺 verification-entry → 第 1 步显式 BLOCKED，不静默降级 | `docs/proposals/frontend-standards-lifecycle/design.md:73-74`、`requirements.md:52` |

### 1.3 实施状态
**部分完成 —— 机制侧已部署（实物存在），proposal 自身未收口。**

证据：
- 无 `issues/` 目录：`ls -la docs/proposals/frontend-standards-lifecycle/` → 仅 `design.md` / `observations.md` / `requirements.md`（无 ticket 可判 closed）。
- 载体演练：端到端五阶段验收在一个内部项目的前端强制工具链 proposal 上跑通（6 票，06:41 立项 → 09:14 全闭）：`docs/proposals/frontend-standards-lifecycle/observations.md:4`、`observations.md:8`。
- 验收通过项："条款四级可追溯 ✓ requirement → design → ticket → review 报告"、`observations.md:9`；"模型矩阵落地 ✓ 13+ runs 全部与矩阵一致"、`observations.md:10`；视觉轴"04/05 双票触发，gemini 判读 16+16 项全 PASS + 量化对照"、`observations.md:12`。
- 未完成项：普适性（验收 2）标为 ⏳ 未做（需非本仓之外的项目实测前端场景）、`observations.md:13`；结论"验收 2（普适性）待其它项目实测后关闭本 proposal"、`observations.md:38`（此处为泛化转述，原文含内部项目名，见票 02 Round 2 脱敏处置）。
- 需求文件自身状态：`docs/proposals/frontend-standards-lifecycle/requirements.md:4` → `> Status: draft`。
- 已知偏差（演练记录，未改文档）：Standards 轴未用 flash lite（子代理不可嵌套）、一次视觉判读误派 frontend-implementer 角色：`observations.md:31-32`。

### 1.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `README.md` §`## 一个需求如何实际推进`（L20）与 §`### 4. 实现、审查并留下证据`（L57） | 只描述文本双轴（Standards / Specification），未提前端票的视觉轴、视觉模型角色与条款四级追溯 | `README.md:20`、`README.md:57-66`（该节两条 review 轴） |
| `docs/architecture.md` §`## 三个组成部分`（L5，agent 列表示意） | 未含前端角色 | `docs/architecture.md:5-27` |
| `docs/architecture.md` §`## 实现边界`（L55）内 §Subagent 名册句（L64） | 名册写"6 个流程角色 + 3 个通用能力角色"，未含 `flowforge-frontend-implementer` / `flowforge-frontend-reviewer` | `docs/architecture.md:64`（名册 juz 句） |
| `docs/skill-system.md` §`## 主交付链`（L5，表 L11-18） | 无前端专项 skill 行，`flowforge-frontend-implement` 未列 | `docs/skill-system.md:11-18` |
| `docs/skill-system.md` §`## Subagent 委派与协作`（L53，角色表 L57-64） | 表中无 frontend-implementer / frontend-reviewer 行 | `docs/skill-system.md:57-64` |
| `docs/skill-system.md` §`## 支持与特殊路径`（L40） | 未提 `flowforge-frontend-implement` 及其槽位/manifest 约定 | `docs/skill-system.md:40-51` |

### 1.5 候选用户场景
> 需求里既要改后端接口又要调前端页面，以前一张票里前端规范全靠实现者自己翻文档，现在设计阶段就把设计系统条款转录进票、拆成前后端两张票，前端票由视觉模型执行并截图核对条款编号。

---

## 2. frontend-vision-loop-poc

### 2.1 Feature 一句话
前端/UI 票可由一个具备图像输入的前端角色执行，交付前强制走"起 dev server → 截图 → 对照反例清单批判 → 修复"的自审循环，截图与批判记录随票交付。

### 2.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| 新 skill | `assets/skills/flowforge-frontend-implement/SKILL.md`：通用 UI 任务工作流（读设计规范 → 组件唯一入口 → 实现 → 截图自审 → 反例自查）；项目专属内容留槽位引用 | `docs/proposals/frontend-vision-loop-poc/requirements.md:24`、`design.md:30-38`（六步骨架） |
| 新 agent 角色 | `assets/subagents/flowforge-frontend-implementer.md`，`model_profile: tool-capable` | `docs/proposals/frontend-vision-loop-poc/requirements.md:25`；实物 `assets/subagents/flowforge-frontend-implementer.md:5` |
| AGENTS.md 模板 | 技能路由表加一行（前端票→frontend-implementer，其余→implementer） | `docs/proposals/frontend-vision-loop-poc/requirements.md:26`、`design.md:20-22` |
| config 键 | `models_by_name: { flowforge-frontend-implementer: cpa/gemini-3.8-flash-high }`（各项目 config 钉扎） | `docs/proposals/frontend-vision-loop-poc/design.md:50` |
| 模型档位 | 默认 `cpa/gemini-3.8-flash-high`，备选 `cpa/deepseek-v4.1-flash` / `cpa/glm-5.3-flash`（均具 image 输入） | `docs/proposals/frontend-vision-loop-poc/requirements.md:27`、`design.md:42-46` |
| CLI | 零 Go 改动，复用 `flowforge agents deploy` 现有机制 | `docs/proposals/frontend-vision-loop-poc/design.md:49` |
| 验收可观测面 | 部署后 `subagent list` 出现 `flowforge-frontend-implementer`（model: gemini-3.8-flash-high） | `docs/proposals/frontend-vision-loop-poc/requirements.md:34` |
| 交付物面 | 截图落 `e2e/__screenshots__/`，批判记录附进 ticket | `docs/proposals/frontend-vision-loop-poc/observations.md:18` |

### 2.3 实施状态
**已完成（PoC 三项验收全通过）。**

证据：
- 无 `issues/` 目录（PoC 最小切片）：`ls -la docs/proposals/frontend-vision-loop-poc/` → 仅 `design.md` / `observations.md` / `requirements.md`。
- `docs/proposals/frontend-vision-loop-poc/observations.md:6` → `## 验收结果(全部通过)`。
- 验收 1 角色部署 + pin：`observations.md:10` → "✓ 三项目 `flowforge-frontend-implementer` → `cpa/gemini-3.8-flash-high`"。
- 验收 2b 视觉模型执行：`observations.md:12` → "run 全程 gemini-3.8-flash-high（54 条消息）"。
- 验收 2c 截图自审循环：`observations.md:13`。
- 验收 2d 修复效果：`observations.md:14` → "AP-19 暗底暗字：店铺名对比度 **1.12:1 → 18.43:1**（4 点全 PASS，最低 5.64:1）"。
- 结论：`observations.md:22` → "三件套（具名角色 + 视觉模型 pin + 截图自审 skill）闭环成立，值得完整立项。"
- 需求文件自身状态仍为 `> Status: draft`：`docs/proposals/frontend-vision-loop-poc/requirements.md:4`。
- 残留项（列入后续完整立项）：lint 链、CI 视觉基线、内容层缺口（AP-21 英文直出）、旗舰侧视觉 review 通道（glm-5.3 纯文本）：`observations.md:26`。

### 2.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `README.md` §`### 4. 实现、审查并留下证据`（L57-66） | 只讲 TDD + 文本双轴，未提前端实现的截图自审循环与视觉模型 pin | `README.md:57-66` |
| `docs/skill-system.md` §`## Subagent 委派与协作`（L53，表 L57-64） | `flowforge-implementer` 一行为唯一实现角色行，无前端专用行 | `docs/skill-system.md:57-64` |
| `docs/skill-system.md` §`## 支持与特殊路径`（L40-51） | 无 `flowforge-frontend-implement` 条目 | `docs/skill-system.md:40-51` |
| `docs/cli-design.md` §`## 其他命令`（L47）→ `flowforge agents deploy`（L50） | 未提"前端角色需项目 config 钉视觉模型才激活"这一部署前提 | `docs/cli-design.md:50` |

### 2.5 候选用户场景
> 改一处店铺名的样式，以前只能靠读代码推理颜色对不对，旗舰模型还看不了图；现在前端角色跑起 dev server 截图，自己对着反例清单挑出"暗底暗字"，改完再截一张，对比度从 1.12:1 拉到 18.43:1 才交付。

---

## 3. flash-worksurface-expansion

### 3.1 Feature 一句话
除了实现者，调研者（investigator）、内建 explore、以及 review 的 Standards 轴也可以跑在便宜快速的 flash 模型上；architect 裁定前的事实收集改为由 flash 调研者产出"设计事实简报"交接，每个阶段都带可观测的回滚门。

### 3.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| CLI（脚本面） | `executor_metrics.py extract --agents <a,b>` 支持多 agent 提取；observations 表在 `session` 后新增 `agent` 列；report 按 agent 分组并逐角色出失控门 | `docs/proposals/flash-worksurface-expansion/requirements.md:50`、`design.md:49`；实物表头 `docs/proposals/flash-worksurface-expansion/observations.md:9`（含 `agent` 列） |
| CLI（命令面） | 不新增 flowforge 子命令；`report` 旧行缺 agent 列按 `flowforge-implementer` 兼容 | `docs/proposals/flash-worksurface-expansion/design.md:50` |
| config 键 | 新增 `agents.models_by_name: {<agent-name>: <model>}`（per-agent 名级钉扎，优先于 profile 键，未知名是 config 错误） | `docs/proposals/flash-worksurface-expansion/requirements.md:52`、`design.md:35-37`；实物 `internal/config/config.go:45`（`ModelOverrides map[string]string \`yaml:"models_by_name,omitempty"\``） |
| config 键（既有 profile 键用法） | `agents.models.tool-capable-read-only` 用于 investigator 钉 `cpa/deepseek-v4.1-flash` | `docs/proposals/flash-worksurface-expansion/design.md:58-65`（d-phase-1 investigator 迁移节） |
| 项目级 config | 目标项目根新增 `opencode.json`：`{"agent": {"explore": {"model": "cpa/deepseek-v4.1-flash"}}}`（explore 非 flowforge 部署 subagent，走 opencode 配置） | `docs/proposals/flash-worksurface-expansion/design.md:68-72`；分类句 `design.md:44` |
| 新 agent 角色 | `assets/subagents/flowforge-reviewer-lite.md`（仅 Standards 轴，`model_profile: tool-capable`，`permission: read-only`） | `docs/proposals/flash-worksurface-expansion/requirements.md:53`、`design.md:78-84`；实物 `assets/subagents/flowforge-reviewer-lite.md:1-11` |
| 流程契约 | `flowforge-review` 双轴拆分派发：编排会话先派 lite 出 Standards findings（带引用），再派旗舰 reviewer 读 findings 专注 spec 轴 + 裁定 + `Fix:` Changes；reviewer 仍是唯一 closeout 权威 | `docs/proposals/flash-worksurface-expansion/design.md:86-91` |
| 简报契约（skill 增补） | `flowforge-research` 增补"设计事实简报"输出格式（问题 / 事实带引用 / 约束面 / 开放项） | `docs/proposals/flash-worksurface-expansion/design.md:97-103`；实物 `.agents/skills/flowforge-research/SKILL.md`（"## Design fact brief" 节，见 L17-32） |
| 简报消费（skill 增补） | `flowforge-solution-design` 增补简报消费步骤：architect 入口 reading list 可含一份设计事实简报 + 必要抽查 | `docs/proposals/flash-worksurface-expansion/design.md:105` |
| AGENTS.md 委派行 | 设计裁定前的事实收集 → investigator(flash) 产出简报，裁定留 architect(旗舰) | `docs/proposals/flash-worksurface-expansion/design.md:107` |
| 回滚门（流程契约） | 四道门预注册：G1 角色失控（per role）/ G1 单次越线 / G2 引用质量 / G3 实施纪元污染；stage 推进为人工决策（P1 观察 ≥2 天无回滚→P2；P2 观察 ≥2 proposal→P3） | `docs/proposals/flash-worksurface-expansion/design.md:111-124` |

### 3.3 实施状态
**已完成（6 票全部 `closed`）。**

证据（每票标题下的 `**Status:** closed` 行）：
- `docs/proposals/flash-worksurface-expansion/issues/01-metrics-multi-agent.md:15` → `**Status:** closed`
- `docs/proposals/flash-worksurface-expansion/issues/02-investigator-flash-pin.md:15` → `**Status:** closed`
- `docs/proposals/flash-worksurface-expansion/issues/03-explore-flash-pin.md:15` → `**Status:** closed`
- `docs/proposals/flash-worksurface-expansion/issues/04-agent-model-override.md:15` → `**Status:** closed`
- `docs/proposals/flash-worksurface-expansion/issues/05-reviewer-lite-standards-axis.md:15` → `**Status:** closed`
- `docs/proposals/flash-worksurface-expansion/issues/06-investigator-design-brief.md:15` → `**Status:** closed`
- 完成证据样例：票 04 "deploy 后产物 frontmatter `model: cpa/deepseek-v4.1-flash`… 未知键 deploy 与 status 均 exit=1 报 `agents.models_by_name: unknown agent "no-such-agent"`"、`issues/04-agent-model-override.md:108`；票 04 "双轴 review 与 disposition：Standards 轴零 findings；Spec 轴零 findings"、`issues/04-agent-model-override.md:117`。
- 观测文件存在且含新表头（`agent` 列）：`docs/proposals/flash-worksurface-expansion/observations.md:9`，文件首行标记 `<!-- flash-migrated epoch start: 1789730224176 -->`、`observations.md:1`。

- 需求文件自身状态：无 `Status` 字段（`docs/proposals/flash-worksurface-expansion/requirements.md:60` 仅 `> Type:` 行），proposal 状态只能由 6 票全 closed 读取。

### 3.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `README.md` §`## 安装与初始化`（L72） | 只讲 `flowforge init` / `config set docs_dir`，未提模型分层配置（`agents.models` / `models_by_name`）这个新手必配项 | `README.md:72-100` |
| `docs/architecture.md` §`## 实现边界`（L55）→ `internal/config` 项（L58） | 配置键清单只写"项目根、`docs_dir`、`standards.guide`、`agents.disabled` 与兼容配置解析"，未提 `agents.models` / `models_by_name` / `models_by_host` | `docs/architecture.md:58` |
| `docs/architecture.md` §Subagent 名册句（L64） | 写"6 个流程角色 + 3 个通用能力角色"，未含 `flowforge-reviewer-lite` | `docs/architecture.md:64` |
| `docs/skill-system.md` §`## Subagent 委派与协作`（L53，表 L57-64） | 无 reviewer-lite 行；`flowforge-reviewer` 行写"针对 Standards 与 Spec 双轴进行代码审查"，与"双轴拆两跳派发"不一致 | `docs/skill-system.md:63` |
| `docs/skill-system.md` §`## 支持与特殊路径`（L40）→ `flowforge-research` 行（L44） | 只写"针对一个缺失的一手资料事实进行调查并留下 Markdown 结果"，未提设计事实简报四段结构 | `docs/skill-system.md:44` |
| `docs/cli-design.md` §`## 其他命令`（L47） 全部小节 | 无 `executor_metrics.py extract --agents` / report 按 agent 分组这一观测面的记载 | `docs/cli-design.md:47-55` |
| `docs/cli-design.md` §`## 稳定边界`（L65） | 未提 model 优先级链（name 键 > profile 键 > preserve-merge 回填 > 宿主默认） | `docs/cli-design.md:65-70` |

### 3.5 候选用户场景
> 每次调研第三方库都走旗舰模型，一次会话烧掉上千万 token 却只产出一份带引用的事实清单；现在事实收集交给 flash 调研者按"设计事实简报"格式交接，旗舰只在简报上做裁定，同时 review 的规范符合性轴也拆给 flash 角色跑。

---

## 4. deploy-artifact-localization

### 4.1 Feature 一句话
模型分配留在每台机器自己的 `.flowforge/config.yaml`，部署时按配置注入各宿主的 subagent 定义（pi 也会写出 `model:`）；部署产物与 config 被 `init`/`upgrade` 自动加入 `.gitignore` 定性为每机器文件，已入库的存量产物只给出 `git rm --cached` 指引而不自动改索引。

### 4.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| config 键（新） | `agents.models_by_host`：外层宿主键（`opencode`/`claude`/`pi`；`codex` 是配置错误），内层 agent 名**或** profile 键，值 = 模型字符串 | `docs/proposals/deploy-artifact-localization/design.md:36-45`；实物 `internal/config/config.go:46`（`ModelHostOverrides map[string]map[string]string \`yaml:"models_by_host,omitempty"\``） |
| 优先级链（新） | 六级：`models_by_host.<host>[name]` > `models_by_host.<host>[profile]` > `models_by_name.<name>` > `models.<profile>` > preserve-merge 回填 > 宿主默认 | `docs/proposals/deploy-artifact-localization/design.md:47-55` |
| CLI 行为（新增，无新命令） | 不新增 CLI 命令；`resolveCompileOptions(cfg, def, hostKey)` 增宿主参数属包内缝，CLI 接口签名不变 | `docs/proposals/deploy-artifact-localization/design.md:57`、`requirements.md:39` |
| CLI 行为（新增） | `agents deploy` 对启用的 model 承载宿主做格式级校验（非空、无空白/控制字符；opencode/pi 需 `provider/model`；claude 单 token 合法），fail-fast 且在任何写入前 | `docs/proposals/deploy-artifact-localization/design.md:61-72`、`design.md:69-71` |
| 错误面 | 未知名沿用 `agents.models_by_name: unknown agent %q` 风格；未知宿主/未知键为 `unknown host %q` / `unknown key %q`；`models_by_host.codex` 存在即配置错误 | `docs/proposals/deploy-artifact-localization/design.md:63`、`design.md:65` |
| `flowforge init` 行为（新） | 幂等追加 `.gitignore`：`.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/`、`.pi/extensions/`、`.agents/`、`.flowforge/config.yaml`；带 marker 注释 `# flowforge: managed deploy artifacts (per-machine)` | `docs/proposals/deploy-artifact-localization/design.md:88`、`requirements.md:32`；实物仓根 `.gitignore` 末段（marker + 7 条） |
| `flowforge upgrade` 行为（新） | `syncProjectAssets` 同通道调用 gitignore 管理 | `docs/proposals/deploy-artifact-localization/design.md:88` |
| 存量迁移输出（新） | `init`/`upgrade` 用本地 `git ls-files --error-unmatch <受管路径>` 检测已跟踪产物，向 stderr 输出逐条 `git rm --cached -r <path>` 指引 + 原因句；绝不自动改 git 索引；不提供 `--untrack` flag | `docs/proposals/deploy-artifact-localization/design.md:90` |
| pi 宿主产物变化 | `.pi/agents/<name>.md` frontmatter 增 `model:`（`omitempty`，位置在 `description` 与 `thinking` 之间），取 `resolveModel(opts)`；空时省略键 | `docs/proposals/deploy-artifact-localization/design.md:76-78`；实物 `assets/subagents/*.md` 源文件本身不含 `model:`（如 `assets/subagents/flowforge-implementer.md:5` 仅 `model_profile`） |
| 既有方案修订 | [PI 宿主集成方案 §一] 的"不写 `model`"一行被覆盖，其 design revision 升至 2 | `docs/proposals/deploy-artifact-localization/design.md:79`；实物 `docs/proposals/pi-host-integration/design.md:26`（"`model`：**有条件写入**（revision 2，由部署产物本地化…修订）"） |
| 修复的既有 bug | pi preserve-merge 虚假保留：`CompilePiWithOptions` 忽略 `FallbackModel` 却打印 preserved 提示 | `docs/proposals/deploy-artifact-localization/requirements.md:23`、`design.md:80` |
| 提示文案变化 | preserved 提示改为 `(set agents.models_by_name/models_by_host in .flowforge/config.yaml to pin explicitly)` | `docs/proposals/deploy-artifact-localization/design.md:81`；实物 `internal/command/agents_test.go:1627`、`agents_test.go:1682` |
| AGENTS 模板约定（新） | `assets/AGENTS.md` 增补：部署产物与 `.flowforge/config.yaml` 是每机器文件（init 自动 gitignore）；用户自写想入库的 agent 用 `git add -f` | `docs/proposals/deploy-artifact-localization/design.md:92` |
| 回归契约 | 未配置任何新键的项目：四宿主产物逐字节不变（pi 仍省略 `model`；claude 仍 opus/sonnet；codex 不变） | `docs/proposals/deploy-artifact-localization/design.md:105` |

### 4.3 实施状态
**已完成（3 票全部 `closed`）。**

证据：
- `docs/proposals/deploy-artifact-localization/issues/01-models-by-host.md:15` → `**Status:** closed`
- `docs/proposals/deploy-artifact-localization/issues/02-pi-model-injection.md:15` → `**Status:** closed`
- `docs/proposals/deploy-artifact-localization/issues/03-artifact-localization.md:15` → `**Status:** closed`
- 完成证据样例：票 01 交付"`agents.models_by_host`（外层宿主键/内层 name+profile 双键）落地… 全套测试绿"、`issues/01-models-by-host.md:19`；票 02 交付"pi 编译产物 frontmatter 支持 `model:` 字段…"、`issues/02-pi-model-injection.md:19`；票 03 交付"init/upgrade 自动管理 `.gitignore`… 部署前 `git ls-files` 检测…"、`issues/03-artifact-localization.md:19`。
- design 自身为 `revision: 3`（含 2026-09-25 用户四项口味题确认）：`docs/proposals/deploy-artifact-localization/design.md:5`、`design.md:31`。
- 本仓库存量一次性手工清理项（按 d-artifact-localization 指引人工执行，不入 ticket 自动化）：`docs/proposals/deploy-artifact-localization/design.md:107`。

### 4.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `README.md` §`## 安装与初始化`（L72-100） | 讲 init 创建的目录与 `--force`，未提 init 会维护 `.gitignore`、产物与 config 是每机器文件、新机器要自己配模型 | `README.md:72-100` |
| `README.md` 全文（L1-115） | 无"部署产物 / 每机器文件"概念的任何小节（`grep -n 'gitignore\|部署产物' README.md` → 无命中） | `README.md:1-115`（含 `## 核心文档` L101-107、`## 开发` L108-114） |
| `docs/architecture.md` §`## 实现边界`（L55）→ `internal/config` 项（L58） | 配置键清单未含 `models_by_host` | `docs/architecture.md:58` |
| `docs/architecture.md` §`## 实现边界`（L55）→ `internal/subagent` 项（L59） | 写"四宿主原生格式编译与 Model Profile 映射"，未提 `model` 注入与 per-host 覆盖 | `docs/architecture.md:59` |
| `docs/cli-design.md` §`## 其他命令`（L47）→ `flowforge agents deploy [name]`（L50） | 一行式描述，未提模型注入、校验 fail-fast、每机器文件定位 | `docs/cli-design.md:50` |
| `docs/cli-design.md` §`## 其他命令`（L47）→ `flowforge upgrade`（L54） | 未提 upgrade 也会维护 gitignore 与做已跟踪产物检测 | `docs/cli-design.md:54` |
| `docs/cli-design.md` §`## 目录解析`（L5）→ init 段（L9） | init 行为只写"创建配置、CONTEXT.md、adr/ 和 proposals/，部署 agents/ 与 .agents/skills/"，无 gitignore 管理 | `docs/cli-design.md:9` |
| `docs/cli-design.md` §`## 稳定边界`（L65） | 未提 per-host 覆盖层与优先级链、格式级校验边界（不做存在性校验、纯本地无网络） | `docs/cli-design.md:65-70` |

### 4.5 候选用户场景
> 换一台新电脑部署，pi 宿主的子代理定义里从来没有 `model:` 字段，模型分配无法重现；现在模型留在本机 config、部署时按宿主注入，产物和 config 由 init 自动 gitignore，多个协作机器各配各的也不会互相污染。

---

## 5. blocked-evidence-persistence

### 5.1 Feature 一句话
执行者被阻断时，返回 `STATUS: BLOCKED` 之前必须先在票里写下结构化失败证据（逐字错误、已试命令与退出码、下一步假设）；`check` 会告警、`frontier` 把它标成带 warning 的就绪票，`refine-ticket` 消费完这段证据后把它转成 Verified contracts 并移除。

### 5.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| 新诊断码 | `blocked-evidence-present`（severity warning，Source 指向票文件）；open 票正文含标题 `## Blocked evidence` 即命中 | `docs/proposals/blocked-evidence-persistence/design.md:23-24`；实物 `internal/tracker/catalog.go:53`（`DiagnosticBlockedEvidencePresent DiagnosticCode = "blocked-evidence-present"`） |
| 诊断消息 | 面向派发方："open ticket carries blocked evidence; run flowforge-refine-ticket to consume it before redispatch" | `docs/proposals/blocked-evidence-persistence/design.md:24`；实物 `internal/tracker/catalog.go:305` |
| `flowforge check` 行为 | warning → `--strict` 判失败（与 evidence 族一致） | `docs/proposals/blocked-evidence-persistence/design.md:25`、`requirements.md:38` |
| `flowforge frontier` 行为 | 零改动：`classifyReady` 按 severity 归类，warning → `ready_with_warnings` 既有桶；需求端要求其"不进 clean ready" | `docs/proposals/blocked-evidence-persistence/design.md:27`、`requirements.md:41` |
| 新票内小节约定 | `## Blocked evidence`（正文小节，不属 Issue Schema 头，与 evidence 四元组同层的 Markdown 标记）；位置约定追加在票末尾（Implementation note 之后，解析不依赖位置） | `docs/proposals/blocked-evidence-persistence/requirements.md:30`、`design.md:28` |
| skill 变更（写入契约） | `assets/skills/flowforge-implement/SKILL.md` 六处 BLOCKED 出口各增补同一指令句：返回前先在票末尾追加三要素小节 + 双通道声明（会话返回照旧，工件为权威） | `docs/proposals/blocked-evidence-persistence/design.md:32-33`；实物证据 `docs/proposals/blocked-evidence-persistence/issues/02-block-and-record-contract.md:44`（`grep -c '## Blocked evidence' … → 6`） |
| agent 角色契约变更 | `assets/subagents/flowforge-implementer.md` 的 Non-negotiables 增第四句 `Block and record:`（英文锚点 `## Blocked evidence`），既有三锚点（`at most 2 times` / `STATUS: BLOCKED` / `failed repair rounds`）逐字不动 | `docs/proposals/blocked-evidence-persistence/design.md:32`、`requirements.md:31`；实物 `assets/subagents/flowforge-implementer.md:22-26`（三锚点）、证据 `issues/02-block-and-record-contract.md:44-49`（`- Block and record:` 追加于 Budget closure 之后、`## Boundaries` 之前） |
| skill 变更（消费循环） | `assets/skills/flowforge-refine-ticket/SKILL.md` 新增前置步骤 `1. Consume blocked evidence`：转写 Verified contracts（验证不了的正确调用不猜测）→ 移除小节 → 复跑 `flowforge check` 确认诊断消失 | `docs/proposals/blocked-evidence-persistence/design.md:37-43`、`requirements.md:42-44`；实物证据 `issues/03-refine-consumption-loop.md:45`（六步骤序，消费为步骤 1 L13，顺序声明 L21） |
| CLI 面 | 不新增 CLI 子命令与接口签名；不修改 Issue Schema 头规范 | `docs/proposals/blocked-evidence-persistence/requirements.md:30-31` |
| config 键 | 无新增 | `docs/proposals/blocked-evidence-persistence/design.md:20`（"不新增 CLI 面"） |
| 清除条件 | 小节移除或票转 `closed` 后解析不再命中，无状态迁移代码 | `docs/proposals/blocked-evidence-persistence/design.md:26`、`requirements.md:39` |

### 5.3 实施状态
**已完成（3 票全部 `closed`）。**

证据：
- `docs/proposals/blocked-evidence-persistence/issues/01-blocked-evidence-diagnostic.md:15` → `**Status:** closed`
- `docs/proposals/blocked-evidence-persistence/issues/02-block-and-record-contract.md:15` → `**Status:** closed`
- `docs/proposals/blocked-evidence-persistence/issues/03-refine-consumption-loop.md:15` → `**Status:** closed`
- 完成证据样例：票 01 "表驱动 8 票：open/heading-only/大小写与尾随空白变体/ready-for-agent → blocked-evidence-present warning"、`issues/01-blocked-evidence-diagnostic.md:31`；负例 4 类零产出（closed+小节、移除小节、`###` 级标题、带后缀标题）、同文件 `:34`。
- 票 02 结构断言三端（SKILL.md 指令句 / 定义源第四句 / 部署产物第四句）+ 既有三锚点回归：`issues/02-block-and-record-contract.md:44`、`:49`。
- 票 03 顺序锁：`strings.Index("Consume blocked evidence") < strings.Index("Fill the five execution-contract sections")`、`issues/03-refine-consumption-loop.md:49`。
- 边界依赖：票 03 `Blocked by: 01`（`issues/03-refine-consumption-loop.md:13`），票 02 无依赖（`issues/02-block-and-record-contract.md:13`）。
- 素材来源：某内部项目一张票在同一命令上重试 290 次的事故（无失败记忆，票号隐去）：`docs/proposals/blocked-evidence-persistence/design.md:9`。

### 5.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `README.md` §`### 4. 实现、审查并留下证据`（L57-66） | BLOCKED 只作为会话返回提及（L60 "返回 Align"、L64 "返回 Solution Design"），未提写票内证据与双通道 | `README.md:57-66` |
| `README.md` §`### 5. 继续 DAG 前沿`（L68-71） | 只讲 clean/gap/blocker 与"全部 ticket 关闭"，未提 `ready_with_warnings` 这一"可派发但带阻断史"的展示 | `README.md:68-71` |
| `docs/skill-system.md` §`## 主交付链`（L5，表 L11-18）→ `flowforge-implement` 行（L17） | 行内未提 BLOCKED 前的 block-and-record 写入契约 | `docs/skill-system.md:17` |
| `docs/skill-system.md` §`## 主交付链`（L5，表 L11-18）→ `flowforge-refine-ticket` 行 | 该表无 `flowforge-refine-ticket` 行（`grep -n 'refine-ticket' docs/skill-system.md` → 无命中），消费循环无承载小节 | `docs/skill-system.md:11-18` |
| `docs/skill-system.md` §`## 继续、跳过与返回`（L81-93） | 列了 blocker 与 review 循环，未提"阻断票先消费 Blocked evidence 再补契约" | `docs/skill-system.md:81-93` |
| `docs/cli-design.md` §`## Artifact Catalog`（L11）→ 确定性诊断覆盖清单（L15-21） | 诊断清单无 `blocked-evidence-present` | `docs/cli-design.md:15-21` |
- `docs/skill-system.md` §`## 支持与特殊路径`（L40-51）→ `flowforge-research` 行 | 未提 `flowforge-refine-ticket` 的 Blocked evidence 消费步骤 | `docs/skill-system.md:40-51` |
| `docs/cli-design.md` §`## `flowforge check``（L23-29） | 未提该新 warning 及 strict 语义 | `docs/cli-design.md:29` |
| `docs/cli-design.md` §`## `flowforge frontier``（L31）→ 投影清单（L39-43） | 只有 clean/warning/gap/blocker 四类表述，未点名 `ready_with_warnings` 桶（`--json` 分组在 L45 提"warning"） | `docs/cli-design.md:39-45` |

### 5.5 候选用户场景
> 一张票在同一个命令上重试了两百多次，换会话重派时新上下文完全不知道上次怎么失败的；现在阻断前必须把逐字错误和试过的命令写进票里，check 会提醒派发方，refine 把这段失败史转成持久契约，重派时不用再踩同一个坑。

---

## 附：本组事实的复核命令（可复验）
```bash
# 角色与 skill 实物
ls assets/subagents/            # 含 flowforge-frontend-implementer.md / flowforge-frontend-reviewer.md / flowforge-reviewer-lite.md
ls assets/skills/               # 含 flowforge-frontend-implement
# config 键
grep -n 'ModelOverrides\|ModelHostOverrides' internal/config/config.go   # L45 models_by_name / L46 models_by_host
# 诊断码
grep -rn 'DiagnosticBlockedEvidencePresent' internal/tracker/catalog.go   # L53 定义 / L305 产出
# 提示文案
grep -rn 'models_by_name/models_by_host' internal/command/agents_test.go   # L1627 / L1682
# init 写入的 gitignore 条目（本仓自身已生成）
tail -12 .gitignore
```
