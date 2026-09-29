# 调研：Agent 产出 UI 的「治理面」社区方案（Skill / 规则条款 / 质量门 / 设计系统分层）

> 建议落库路径：`docs/research/2026-09-29-ui-enforcement-agent-governance-research.md`
> 本文档由 flowforge-research 调研产出（2026-09-29），只做社区方案调研，不做设计裁决。
> 证据标注约定：【已验证·一手】= 直接抓取原文核对；【转述】= 来自搜索结果摘要或二手引用，未逐字核对原文；【未核实】= 无法验证或来源失效。

## 背景

本项目用 pi 类 coding agent + FlowForge skill 体系协作开发，前端技术栈为 antd v5 + @ant-design/x + 自有 `@tangram/ui` 封装（`AppLayout` / `TangramThemeProvider` / `ReadableRegistry`），且已有成文规范 `docs/805-frontend-design-guide.md`（L/T/I/F/S 编号条款、`[Constraints]`/`[Conventions]` 分层、第六章 17 条反例集）。引入 antd 的初衷是「让 agent 产出的页面天然拥有完善的交互与视觉」，但实际产出仍系统性违反规范（自建第二套 Header、暗底暗字、信息密度差、英文直出），说明**成文规范对 agent 的约束力不足**。本报告回答「Agent 治理面」的社区方案：引入什么 SKILL、AGENTS.md 怎么写、工作流加什么质量门，为后续方案裁决提供证据。

---

## 问题一：项目级 SKILL 生态——官方与社区里有什么

### 结论

1. Anthropic 官方 `anthropics/skills` 仓库中**没有**「组件库唯一入口强制」类技能；官方技能的共性写法是「角色设定 + 先计划后实现 + 负面清单（反 AI 味）+ 截图自审」，正文短、命令式、以禁止性约束为主。`frontend-design` 与本项目「严格收敛到既有设计系统」的目标方向相反（它教 agent 摆脱模板味、自由发挥），**不能直接安装使用，但结构与清单形态高度可借鉴**。
2. 官方技能中最可复用的是 `webapp-testing`（Playwright 截图自审的标准件）与 `brand-guidelines`（token 权威源进 skill 的最小实现）。
3. 社区 registry（skills.sh、awesome-claude-skills 等）以「品味提升 / UI/UX 审计」类技能为主，存在零星「design-system-enforcement」类技能但**均未成为高星标准件**；Figma 官方 `mcp-server-guide` 中的 `create-design-system-rules` 是目前最接近「为项目生成设计系统强制规则」的**官方模板**，其规则写法建议（IMPORTANT 前缀、具体化、可执行化、失效迭代）对 AGENTS.md 条款维护直接适用。

### 证据

**（1）anthropics/skills 仓库结构与写法【已验证·一手】**（抓取仓库文件树与 SKILL.md 原文）：

- 仓库结构：每个技能为 `skills/<name>/SKILL.md`（YAML frontmatter：`name` / `description` / `license`）+ 可选 `scripts/`（确定性脚本）、`references/`（按需加载文档）、`assets/`（字体/模板等资源）。例：`webapp-testing` 带 `scripts/with_server.py` 与 `examples/`；`canvas-design` 带 `canvas-fonts/`（约 30 个字体文件）；`claude-api` 带 `references/` 大量分层文档。仓库地址：https://github.com/anthropics/skills
- 触发机制（Anthropic 官方说明，steering 博客一手 + 工程博客转述）：skill 只有 `name`+`description` 在会话开始时预加载，正文在被调用（斜杠命令或任务自动匹配 description）时才进入上下文；「过程性指令放 skill 而非 CLAUDE.md」。见 https://claude.com/blog/steering-claude-code-skills-hooks-rules-subagents-and-more 与 https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills （后者为搜索摘要转述要点：渐进式披露、frontmatter 决定触发、确定性工作优先写成脚本、从评测失败模式迭代技能）。
- **frontend-design**（原文全文已核对：https://github.com/anthropics/skills/blob/main/skills/frontend-design/SKILL.md ）：
  - 角色：把任务当作「设计工作室 design lead」的定制化交付，而非模板。
  - 过程：两遍工作法——先产出「紧凑设计计划」（4-6 个命名色 token、字型与角色、布局概念 + ASCII wireframe、独特性原则），对照 brief 自查「这段计划是否换一个 brief 也会长这样」，修订后才开始写代码。
  - 硬约束清单（负面清单形态）：列出 5 类「AI 生成页面的常见 tell」并要求默认禁用（米白底 + 高对比衬线 + 陶土橘、近黑底 + 荧光绿/朱红点缀、报纸 hairline 分栏、SaaS 卡片套件「同一圆角 + 同一灰阴影 + 渐变装饰」、模板 chrome「ALL-CAPS eyebrow / 中点 meta 串 / 箭头按钮」）；并注明「brief 明确要求时这些才可用」——即**负面约束 + 显式豁免条件**的写法。
  - 自查 checklist（质量底线）：响应式到移动端、键盘焦点可见、reduced-motion、视觉可访问（对比度）；以及「**build 时自评，环境允许就截图自查——a picture is worth 1000 tokens**」。
  - 文案规范：面向最终用户命名（"user manages notifications, not webhook config"）、CTA 用主动语态且全流程同名、错误与空态是「给方向的时刻而非情绪」。
  - ⚠️ 方向性冲突（研究者判断）：该技能的价值观是「反收敛、求独特」，直接装进本项目会鼓励 agent 偏离 antd 既有视觉体系；但其「反 AI 味清单」「两遍工作法」「截图自评」三个结构组件正是本项目反例集与自审循环想要的形态。
- **brand-guidelines**（原文全文已核对：https://github.com/anthropics/skills/blob/main/skills/brand-guidelines/SKILL.md ）：极简形态——品牌色（Dark `#141413`、Accent `#d97757` 等 7 个 hex）、字体及应用规则（24pt+ 标题用 Poppins、正文 Lora、缺字体时的回退链）、非文本形状轮换用 accent 色。这就是「**token 权威源 + 应用规则直接写进 skill 正文**」的最小实现，与本项目 `tangramThemeTokens` 的治理诉求同构（改造点：把 antd Design Token 表与 `TangramThemeProvider` 装配规则写进去，替换 Anthropic 品牌值）。
- **webapp-testing**（原文全文已核对：https://github.com/anthropics/skills/blob/main/skills/webapp-testing/SKILL.md ）：官方「让 agent 亲眼看渲染结果」的标准件——原生 Python Playwright；决策树（静态 HTML 直接读 vs 动态应用先起服务）；「侦察-行动」模式（navigate → wait `networkidle` → screenshot/DOM 检查 → 发现 selector → 执行动作）；明确的反面警示「不要在等待 networkidle 前检查 DOM」；helper 脚本当黑盒用（"DO NOT read the source…they pollute your context window"）。
- **canvas-design**（原文全文已核对，经本地克隆读取）：两阶段「设计哲学宣言（.md）→ 视觉表达（.pdf/.png）」+ 收尾强制第二遍打磨（"take a second pass…refine/polish"）。与前端治理关联弱，仅证明官方技能普遍采用「先产出中间制品再实现」的结构。

**（2）社区 skills 集合【转述，仓库清单部分一手验证】**：

- skills.sh（registry）：官方 `frontend-design` 在列，页面显示约 930K 安装（页面展示数据，未复核）。https://www.skills.sh/anthropics/skills/frontend-design
- awesome-claude-skills（BehiSecc，仓库克隆一手验证）：https://github.com/BehiSecc/awesome-claude-skills —— 按类别编排（Document / Development & Code Tools / Security & Web Testing / Collections 等），收录官方 docx/pdf/pptx/xlsx、`web-artifacts-builder`（React+Tailwind+shadcn）、obra/superpowers 的 TDD 技能、`oiloil-ui-ux-guide`（guide/review 双模式，CRAP/HCI 法则/现代极简）、`Design Auditor`（17 条专业规则审计设计，含排版/WCAG/间距/token，打分制）、`agnix`（AI agent 配置的 linter，校验 SKILL.md/CLAUDE.md/hooks/MCP 配置，156 规则）。**未发现现成的「antd/组件库唯一入口」强制类高星技能**。另有 ComposioHQ、travisvn 等同名 awesome 列表（未验证内容）。
- 「设计系统强制」类零散技能（**均为搜索摘要转述，未逐仓库验证内容**）：
  - chrometaphore/design-system-skill（`frontend-design-system`）：扫描项目真实依赖/组件/token，引导 agent 使用而非硬编码发明。https://github.com/chrometaphore/design-system-skill
  - majiayu000/claude-skill-registry 的 `design-system-enforcement`（openskillindex 条目）。
  - rampstackco/claude-skills 的 `design-system` / `design-standards`。
  - Vanszs/Anti-AI-UI（Anti-AI-Slop SKILL.md）、miqdadbadjuber/anti-slop：与官方 frontend-design 的「tells 清单」同源思路的社区防 slop 技能。
- **Figma 官方 mcp-server-guide**（仓库克隆一手验证，含 `create-design-system-rules` 全文）：https://github.com/figma/mcp-server-guide —— 这是目前最接近「项目级设计系统规则生成器」的官方材料：
  - 定位：设计系统规则 = 把「资深开发传给新人的隐性知识」（用哪些布局原语/组件、组件文件放哪、命名结构、什么绝不能硬编码、token 与样式怎么处理）编码为项目级指令，供 agent 在每次实现 Figma 设计时遵守。
  - 生成流程（Required Workflow，明令「按序执行，不得跳步」）：调 `create_design_system_rules` 工具拿模板 → 分析代码库（组件目录、token 位置、命名约定、样式方案、架构决策）→ 按模板生成规则 → 存入 `CLAUDE.md` → 用简单组件实现验证、迭代不生效的规则。
  - 规则写法建议（对 AGENTS.md 维护直接适用）：关键规则加 `IMPORTANT:` 前缀；**具体化**（"Always use Button components from `src/components/ui/Button.tsx` with variant prop ('primary'|'secondary')" 而非 "Use the design system"）；**可执行化**（"Colors are defined in `src/theme/colors.ts` – import and use these constants" 而非 "Don't hardcode colors"）；解释 why；规则过多会拖慢 agent，抓 20% 规则解决 80% 一致性问题；规则会过时，需定期 review。
  - 内含「实现类任务的强制流程」：`get_design_context` + `get_screenshot` 都拿到后才开始实现，**完成前必须对照截图做 1:1 视觉校验**——这是「截图基线验收」的官方化写法。

### 对 react+antd+自有组件库项目的改造点（研究者综合推断，非来源直接陈述）

- 新建项目级 skill（如 `tangram-frontend` 或改造 `webapp-testing`）：正文 = UI 任务工作流（先读 805 相关条款 → `@tangram/ui` 组件清单即唯一入口 → 实现 → Playwright 截图自审 → 对照第六章反例集自查）；`references/` 放 805 全文与 antd 官方链接表（渐进式披露，避免常驻上下文）。
- `brand-guidelines` 模式迁移：`tangramThemeTokens` 值 + `TangramThemeProvider` 装配规则写进 skill 正文（短、权威、单一来源）。
- `frontend-design` 的「负面清单 + 豁免条件」结构迁移为第六章反例集的 skill 化表达；其「两遍工作法」对应 FlowForge 的 Plan→Implement 分离。
- Figma `create-design-system-rules` 的「规则不生效→更具体 + IMPORTANT」迭代法迁移为本项目 AGENTS.md/805 条款的维护规程。

---

## 问题二：AGENTS.md / CLAUDE.md / .cursorrules 强制条款的写法与有效性

### 结论

1. **Anthropic 官方明确划界**：CLAUDE.md/rules/skills 是「提示（prompting）」，不是「强制（enforcement）」；「绝对不能发生」的事必须用确定性机制（hooks、permissions、managed settings）兜底。提示式规则在长会话、压力、歧义、prompt injection 下会失效。
2. **唯一找到的实证研究**（5000+ 次 agent 运行）表明：规则文件整体有正收益，但逐条看，**全部受益的规则都是负向约束（"do not X"），全部有害的规则都是正向指令（"follow X style"）**；且随机规则与专家规则收益相当（增益内容无关，疑为 context priming）。这为「短、负向、可验证」条款写法提供了最强的公开证据，但也警示「写了规则就有用」的归因不可靠。
3. 社区共识写法：条款四要素（触发条件 + 单一可观察行为 + 验证命令 rg/git diff + 失败上报义务）；「must happen / must never happen」下沉到 hooks/permissions/tests/CI；30 行以上流程移入 skill。
4. 本项目 805 的形态（编号条款 + Constraints/Conventions 分层 + 反例集）方向正确；**缺口不在条款写法而在三处**（推断）：① 805 是文档不是任务时上下文，agent 执行 ticket 时无强制阅读路由；② 条款没有配验证命令/lint 对应，agent 无法自证合规；③ 「英文直出/信息密度」属正向品味类目标，恰是提示式条款最弱的一类，需靠类型约束（ReadableRegistry）+ review checklist + 截图自审补位。

### 证据

**（1）官方定位与分流【已验证·一手】**（steering 博客全文已核对：https://claude.com/blog/steering-claude-code-skills-hooks-rules-subagents-and-more ）：

- 原文要点（直接摘录意译）：「**"Never do this" 写在 CLAUDE.md 里是错误的工具**。Claude 大多数时候会遵守，但在压力下、长会话中、歧义情境下、或因任务文件里的 prompt injection，模型可能不遵守被提示的规则。真正的护栏必须是确定性的，enforcement 手段是 hooks 和 permissions。PreToolUse hook 可以检查一次调用并以 exit code 2 阻断；managed settings 是管理员下发、用户本地配置不可覆盖的组织级唯一确定性护栏。」
- 分流表（原文）：「"Every time X, always do Y" → hook；30 行流程 → skill；路径相关规则 → `.claude/rules/` 的 `paths:` 作用域（如 `src/api/**` 才加载）；append-system-prompt 指令越多遵从越差，尤其互相矛盾时。」
- skills 触发与成本：「只有 name 与 description 在会话开始加载；正文被调用才载入；compaction 时按预算重注入已调用技能。」

**（2）实证研究【已验证·一手摘要】**：arXiv:2604.11088《Guardrails Beat Guidance: A Large-Scale Study of Rules, Skills, and Persistent Configuration for Coding Agents》（v1 2026-04-13；抓取 GitHub 上 679 个规则文件 / 25,532 条规则，Claude Code + Opus 4.6 在 SWE-bench Verified 上 5000+ 次运行）https://arxiv.org/abs/2604.11088

- 摘要原文要点：「随机规则与专家精选规则对任务表现的提升相当（判别子集上均为 +13.8pp）；我们数据中**每一条单独有益的规则都是负向约束**（如 "do not refactor unrelated code"），**每一条单独有害的规则都是正向指令**（如 "follow code style"）」；增益基本与内容无关（random/shuffled/错域/换格式均与精选相当），指向 context priming 机制；单条常显有害但在集合中不累积可见伤害（0-50 条规则通过率稳定）。结论原则：「**约束 agent 不能做什么，而不是规定它应该做什么**」。
- ⚠️ 适用性边界（研究者标注）：该研究为通用 SWE 任务、单一模型族、单一 agent 框架；它支持「负向约束优于正向指令」，但「随机规则同样有效」并不意味着规则内容不重要——对 UI 合规这类需要精确拦截的场景，可验证的负向约束 + 确定性检查仍是社区主流做法（见矛盾记录）。

**（3）社区模板与条款样式**：

- Figma 官方 steering/技能文档引用的 AGENTS.md 模板段落【转述（搜索摘要引用；objectstack-ai/objectui 的 AGENTS.md 原文已抓取 61,741 字符但未逐条核对）】https://github.com/objectstack-ai/objectui/blob/main/AGENTS.md ：「UI 前先检查既有设计系统组件/token/文档；能用既有组件就用，wrap/compose 而非重造；禁止在有 DS 等价物时使用原生 HTML 控件（button/input/select）；禁止硬编码颜色/间距/字型/圆角/阴影/断点；**禁止内联样式**（`style={{...}}`、HTML `style=""`）；系统没有所需组件/变体时**先问再发明**；完成前跑 lint/typecheck/UI 测试/Storybook 并修复而非记录例外。」另见 Figma 官方 `mcp-server-guide` 的 `figma-power/steering/` 目录【已验证·一手，仓库文件树】——Figma 把此类内容称为「steering（转向）」文件。
- 「grep-able rule 四要素」【转述】：触发条件（"when editing src/api/**"）+ 可观察行为 + 精确验证命令（`git diff --name-only | grep -E '^db/migrations/' && exit 1`、`rg`）+ 失败上报义务（"report the actual exit status; do not say tests passed otherwise"）。来源如 https://www.sharmaprakash.com.np/ai/rules-agents-actually-follow/ 、https://dev.to/moonrunnerkc/how-to-write-a-claudemd-rule-that-actually-gets-enforced-3npa （未逐篇核对原文）。
- antd 官方仓库自身有 `.github/copilot-instructions.md` 与 `DESIGN.md`【存在性来自搜索结果，内容未读，未核实】——antd 官方也在为 agent 准备项目级指令，可作为「官方库自己怎么写 agent 规则」的后续研究对象。
- 本项目 805 已实践的「引用 + 裁剪」（官方细节给链接不内化）与上述模板的「指向权威源 + 只写裁剪决策」一致【已验证·一手，项目文档】。

---

## 问题三：Agent UI 质量门工作流（2025-2026 主流实践）

### 结论

1. **截图自审循环**是 2025-2026 最主流的 agent 侧质量门：官方（Anthropic webapp-testing skill、frontend-design 的「截图自评」）与社区（Playwright MCP 工作流博客）口径一致——「实现 → 起 dev server → 截图（desktop+mobile）→ 对照验收清单批判 → 修复 → 重复 N 轮 → 稳定后固化为回归测试」。Playwright 官方定位：交互走 accessibility snapshot，**截图专用于视觉验证**。
2. **AI 视觉 review 服务**：Chromatic UI Review 本质是「PR 级人工评审工作流」（reviewers/评论/request changes）+ 快照视觉 diff，其 AI 叙事定位是「agent 生成、Chromatic 验证、人审意图」；Argos 是 MIT 开源的同类（Playwright SDK `argosScreenshot` + PR check）。两者都要求先有 Storybook stories / 稳定截图入口。
3. **CI 侧**：Lighthouse CI `assertions` 设为 `error`（a11y/perf 分数与资源预算）+ axe-core（`@axe-core/playwright` 组件级断言 / `@axe-core/cli --exit`）是标准组合。
4. 层次位置（社区共识综合）：agent 会话内自审（截图 + 即时 axe，最便宜的修复点）→ PR/CI 阻断（视觉基线 + 预算断言）→ 人审（意图性变更）。对「暗底暗字」这类问题 axe 对比度检查能拦；「自建第二套 Header」靠整页截图基线能拦；「英文直出/信息密度差」**没有现成自动化门**，靠 skill checklist + 人审。

### 证据

**（1）截图自审循环【已验证·一手（官方 skill 原文）+ 转述（社区流程）】**：

- 官方 `frontend-design` 原文：「Critique your own work as you build, **taking screenshots to review if your environment supports it — a picture is worth 1000 tokens**」【已验证·一手】。
- 官方 `webapp-testing` 原文提供完整工具链：`with_server.py` 管理服务生命周期、`page.screenshot(path, full_page=True)`、`networkidle` 等待、console log 捕获【已验证·一手】。
- Playwright MCP 官方文档【转述（官方文档页摘要）】https://playwright.dev/mcp/introduction 、https://playwright.dev/mcp/tools/screenshots ：设计为 accessibility-tree 结构化交互，截图作为补充的视觉验证工具；社区实践把它组成「build→inspect→screenshot→critique→fix→repeat（最多 3 轮）」的 agent prompt（多篇博客转述：sdet.qa、egghead、qaskills 等，未逐篇核对）。
- Figma `create-design-system-rules` 的 Required Flow 也内置「完成前对照 `get_screenshot` 做 1:1 校验」【已验证·一手】——截图验收正在从社区实践变成官方模板标配。

**（2）AI 视觉 review 服务【转述（官方页面摘要）】**：

- Chromatic：UI Review 为 PR 视觉评审工作流（快照 diff、指派 reviewer、request changes）https://www.chromatic.com/docs/review/ ；AI 工作流页把 Chromatic 定位为验证层：agent 基于「验证过的 UI 上下文（组件/props/stories/patterns）」生成代码，Chromatic 跑 visual/interaction/a11y 检查，人审模糊 diff 与 UX 意图 https://www.chromatic.com/frontend-workflow-for-ai 。
- Argos：MIT 开源视觉回归平台，Playwright quickstart 为测试内 `argosScreenshot(page, "name")` + CI 上传 + PR check accept/reject https://github.com/argos-ci/argos 、https://github.com/argos-ci/docs/blob/main/docs/quickstart/playwright-quickstart.md 。

**（3）CI 断言【转述（官方文档摘要）】**：

- Lighthouse CI：`assertions` 以 `'error'` 级别配置 `categories:accessibility` minScore、`categories:performance`、资源体积预算，违例即非零退出 https://github.com/GoogleChrome/lighthouse-ci/blob/main/docs/configuration.md 。
- axe-core：`@axe-core/playwright` 在每个交互态渲染后 `AxeBuilder.analyze()` 断言 `violations` 为空；`@axe-core/cli --exit` 适合 CI 阻断 https://github.com/dequelabs/axe-core-npm 。

---

## 问题四：公开设计系统的治理分层案例

### 结论

成熟设计系统的治理是**分层分工**：文档层（何时用/不用）→ IDE 实时层（VS Code 扩展诊断与 quick fix）→ lint 层（CI 可阻断的 token/选择器/布局规则）→ codemod 层（存量违规批量修复 + 版本迁移）→ review 层（人审 + checklist）。共同经验是：**单靠「禁止」会被绕过，必须把合规路径做成最省力路径**（quick fix、codemod、顺滑 API）。各系统公开完整度不一：Polaris（lint+IDE 最全）、Primer（codemod+组件生命周期文档最全）、Geist（文档质量+设计原则，未见公开强制 lint）、Ant Design（价值观与规范文档为主，官方未见 agent 专项治理件）。

### 证据

- **Shopify Polaris**【分层一手/规则页转述】：
  - IDE 层：官方 VS Code 扩展「Polaris for VS Code」（内联诊断、quick fix、组件/token 建议）https://marketplace.visualstudio.com/items?itemName=Shopify.polaris-for-vscode 。
  - lint 层：`stylelint` Polaris 插件 40+ 规则；典型防绕过规则【转述，来自官方规则页摘要】：`layout/property-disallowed-list` 拦 `display: grid` 等手写布局并引导用 Polaris `InlineGrid` 布局组件；`conventions/selector-disallowed-list` 拦覆盖 Polaris 私有类名（随时会变，正确路径是组件 API 或上游贡献）；border 类规则强制 Polaris token 而非 legacy Sass helper。规则页：https://polaris-react.shopify.com/tools/stylelint-polaris （注：本次直接抓取该页被重定向至 shopify.dev 通用页，具体规则页内容未逐页复核，npm 包页面 403——整体标转述）。
- **GitHub Primer**【转述（官方 repo 摘要）】：`primer/react-migrate` 提供一次性 codemod（`npx primer-react-migrate src -p v37 --create-commits`；v35/v37 preset；废弃 Button/Dialog/Tooltip/Octicon/TabNav 的定向转换），repo 已于 2025-08-07 归档 https://github.com/primer/react-migrate ；组件生命周期文档管理废弃节奏 https://primer.github.io/contribute/component-lifecycle/ ；v37 迁移讨论（CSS Modules、禁止 wildcard 深导入、`TabNav`→`UnderlineNav`）https://github.com/primer/react/discussions/5165 。底层用 jscodeshift（AST codemod 工具链）https://github.com/facebook/jscodeshift 。
- **Vercel Geist**【转述（官方页面摘要）】：定位「跨团队跨端一致实现的 foundations + tokens + components」https://vercel.com/geist/stack ；每个组件页含「when to use / behavior / content / accessibility」引导（如 Tabs 只用于同 scope 的兄弟视图）https://vercel.com/geist/tabs ；《Web Interface Guidelines》是活文档并明确「审计生成的 UI 代码、把重复反馈固化为更好的默认值/工具/系统」https://vercel.com/design/guidelines 。未见公开强制 lint（未核实）。另有《AI-powered prototyping with design systems》博客未读（列参考）。
- **Ant Design 官方**【一手（项目 805 已引用官方规范页）+ 存在性未核实】：`DESIGN.md`（设计价值观）与 spec 文档体系是文档层治理；仓库存在 `.github/copilot-instructions.md`（搜索结果列出，内容未读）。antd 生态未见官方「组件白名单 lint」。
- **组件白名单 lint 的通用实现**【转述（ESLint 官方规则文档摘要）】：`no-restricted-imports` 的 `paths` + `allowImportNames` 可对 `antd` 做 deny-by-default 白名单（或整体禁直引 + `patterns` 拦 `antd/*` 深导入，仅对 `packages/tangram-ui/**` 豁免），配自定义 message 指向 `@tangram/ui` https://eslint.org/docs/latest/rules/no-restricted-imports 。

---

## 问题五：「防 UI 漂移」组合拳是否有成体系总结

### 结论

**各层都有权威单点（Anthropic skills/steering、Figma steering、Polaris/Primer 分层、Chromatic/Argos/Lighthouse/axe），但「组件白名单 lint + 截图基线 + 项目 skill」三件套的成体系串联目前只存在于厂商博客与社区模板，未见官方或大型开源模板把五层串成一体的「标准答案」。** 最接近成体系的是 designproject.io 的「Agentic design system」四层上下文架构文（厂商博客，质量较高、有研究引用）。组合拳的社区共识形态为：**确定性门（lint/CI）拦可 grep 的违规 + 截图证据拦视觉漂移 + skill 承载流程与 checklist + 常驻规则文件只做路由**。

### 证据

- **designproject.io《Agentic design system: How to stop UI drift in your codebase》**【已验证·一手全文】https://designproject.io/blog/agentic-design-system-context/ ：
  - 核心论点：「组件库存在 ≠ agent 会用它。agent 可能发现不了、误解用途、或认为本地新造一个更快——可用性与正确使用之间的裂缝就是漂移开始的地方。」
  - 四层上下文架构：`AGENTS.md`/`CLAUDE.md` = 常驻路由层（架构边界、硬约束、构建命令、「改 UI 前先读 PRODUCT.md 与 DESIGN.md」）；**skills = 任务层**（详细流程只在匹配时加载）；`PRODUCT.md` = 产品意图层（受众/性格/参考与反参考，防「统计安全的通用 SaaS 脸」）；`DESIGN.md` = 视觉实现层（token 角色、组件清单与状态、明确 do/don't，「系没覆盖就先问再发明，批准后回写 DESIGN.md」）。
  - 五条落地法与维护原则：常驻文件保持稀缺；流程进 skill；规则带语义（"purple 用于 primary/focus/selection，绝不做通用成功色"这种「值 + 决策规则」写法）；产品与设计文档分离维护；源头更新。
  - 引用 ETH Zurich 研究（arXiv:2602.11988，转述）：**LLM 生成的 context 文件使任务成功率降约 3%，人写的提升约 4%**——佐证「常驻文件要小且人写、避免 agent 自动生成规则文件」。
- 社区 checklist 与博客【转述，未逐篇核对原文】：usewalkie《12-Point Checklist for Verifying AI-Built UI》（empty/loading/error/long-text/dark/focus/mobile/disabled/slow-network 等 12 态验证）https://usewalkie.com/blog/verifying-ai-built-ui-checklist/ ；anti-slop 技能群（Vanszs/Anti-AI-UI、miqdadbadjuber/anti-slop）；snapdiff.ai《Catching UI drift…》一文搜索给出的 URL 已 404【未核实，疑似失效或迁移】。
- FlowForge 侧对照（项目事实）：AGENTS.md 已有「Entry reading list（ticket + 关联权威文档）」「Review 收敛」「Expected tests 预设」机制——这三者正是社区三件套在流程层的对应物；缺的是 UI 专用技能绑定与机械验证命令。

---

## 可借鉴方案清单表

| # | 方案 | 来源 | 改造成本 | 预期拦截面 | 与本项目 skill 体系的关系 |
|---|------|------|---------|-----------|------------------------|
| 1 | 项目级 `tangram-frontend` skill：UI 任务工作流（读 805 → 组件唯一入口清单 → 实现 → 截图自审 → 反例集自查） | anthropics/skills（webapp-testing/frontend-design 结构）+ Figma steering | 中（新写 SKILL.md + references） | UI 任务流程合规：选型、四态、展示分级、布局壳 | 新增技能进 AGENTS.md 路由表，与 flowforge-implement/review 衔接（implement 的绑定技能或其 references） |
| 2 | AGENTS.md「UI 负向约束 + 验证命令」：可 grep 条款配 `rg`/lint 命令与失败上报义务 | claude.com steering + arXiv:2604.11088 + Figma「IMPORTANT/具体化」法 | 低 | 可机械验证违规：手拼 Layout/Header、antd token 外色值、绕过 Provider | 铁律区追加；805 `[Constraints]` 条款转录为验证命令 |
| 3 | ESLint `no-restricted-imports` 组件白名单（apps 禁 antd 直引，`@tangram/ui` 唯一入口）+ stylelint token 强制 | ESLint 官方规则 + Polaris stylelint 模式 | 中（前端引入 ESLint/stylelint 链） | import 级绕过（第二套 UI、深导入）；裸色值/裸间距 | CI 质量门；FlowForge review 的机械前置，降低 reviewer 负担 |
| 4 | 截图自审循环（Playwright 起服务 → 桌面/移动截图 → 对照清单批判 → 修复 ≤N 轮） | webapp-testing（官方）+ Playwright MCP 实践 | 中 | 布局漂移（第二套 Header）、暗底暗字、密度差、溢出 | flowforge-implement 的执行环节 / tangram-frontend skill 的收尾步骤；宿主能力映射待验证（见开放问题 3） |
| 5 | 视觉基线 CI（Playwright `toHaveScreenshot` 或 Argos/Chromatic） | Argos（MIT 开源）/ Chromatic / Playwright | 中-高（需稳定截图入口，Chromatic 还需 Storybook） | 回归性视觉漂移（共享组件被改坏、级联样式回归） | ticket 完成门 + CI；与 FlowForge「Expected tests」预设衔接 |
| 6 | axe 断言 + Lighthouse 预算（a11y 分数 error 级） | axe-core / lighthouse-ci 官方 | 低-中 | 对比度（暗底暗字）、键盘可达、性能预算 | CI 门；部分覆盖 805 质量底线条款 |
| 7 | 四层上下文架构：AGENTS.md 只做路由 + skill 承流程 + PRODUCT/DESIGN 分层（805 即 DESIGN.md 角色） | designproject.io（引用 ETH 研究） | 低 | agent「找不到系统/上下文过载/统计安全脸」；常驻 token 浪费 | 805 定位为 DESIGN 层；AGENTS.md 加「改 UI 前先读 805」路由块（现已部分具备） |
| 8 | 规则失效迭代法：不生效→更具体 + IMPORTANT 前缀 + 定期 review | Figma create-design-system-rules | 低 | 规则熵增与过时 | 805/AGENTS.md 的维护规程（配合 WIKI-STRUCTURE 清单） |
| 9 | codemod 批量修复存量违规（jscodeshift） | Primer react-migrate / react-codemod | 高（后期才需要） | 存量 17 条反例的批量收敛、库升级迁移 | 一次性治理工具，非日常门 |
| 10 | 「AI 味/反例」负面清单技能化（带豁免条件） | frontend-design tells 清单 + 805 第六章 | 低 | 模板脸、装饰性默认、无信息量的 chrome | tangram-frontend skill 的自查段；与 Review 的 Spec 轴引用同源 |

## 开放问题（留给我方裁决）

1. **lint 链引入时机**：方案 3（ESLint 白名单 + stylelint）是社区拦截力最强的一层，但意味着前端构建链新增工具与规则维护成本——先只配 `rg` 验证命令（方案 2）还是直接上 lint，需要权衡。
2. **技能归属**：UI 工作流是新建 `tangram-frontend` skill，还是作为 `flowforge-implement` 的绑定 references 注入（AGENTS.md 技能路由表已有 implement 行）？两者可并存，边界怎么划。
3. **截图自审的宿主实现**：Playwright 本地脚本（webapp-testing 模式）vs Playwright MCP vs CI 侧 Argos/Chromatic；pi 宿主自身已具备视频/图像读取能力（本次调研即用 fetch 读图），能否直接承担「渲染结果进上下文」的角色需实测。
4. **基线策略**：整页截图基线噪声大（数据驱动页面对比需稳定 fixture），Argos/Chromatic 需要截图入口先工程化——首期先做「agent 会话内自审截图」还是直接建 CI 基线。
5. **正向品味类目标（英文直出、信息密度）的落点**：ReadableRegistry 类型层约束 + skill checklist + 人审三者边界；axe/Lighthouse 均拦不住这类问题，需明确人审保留范围。
6. **规则有效性预期管理**：arXiv:2604.11088 显示规则增益可能主要来自 context priming 且随机规则同样有效——我方是否要为 805/AGENTS.md 条款建立小型评测（拿 store-mate 反例做回归题）来真正度量拦截率，而非依赖「写了规则」。
7. **pi 宿主与 Claude Code 机制的等价性**：本报告引用的 hooks/permissions/`paths:` 作用域/managed settings 均为 Claude Code 机制；pi 宿主的等价物（权限/deny 工具、目录作用域规则）未在本调研中核实，需单独确认。

## 矛盾与张力记录

- **「随机规则同样有效」vs 社区「规则要具体可验证」**：arXiv:2604.11088 的整体结论（增益内容无关）与社区主流写作建议（具体化 + 验证命令）存在张力。论文自己的逐条分析（负向约束有益、正向指令有害）与社区建议兼容，但整体效应提示：不要指望条款内容本身产生大部分收益，确定性检查才是下限。两说并存，未消解。
- **官方 frontend-design 的「反收敛」vs 本项目「强收敛」**：直接安装官方 frontend-design 会与 805 目标冲突（它鼓励脱离模板、自造 token 系统）；只能取其结构不复用其价值观。
- **Chromatic 的 AI 叙事 vs 产品本质**：官方页把 Chromatic 定位为「AI 工作流的验证层」，但其 UI Review 核心是人工评审工作流 + 快照 diff，不是自动 AI 评审——采信其「人审意图、机器拦回归」的分工，不采信「AI 自动 review」的想象。

## 未核实与缺失证据

- snapdiff.ai《Catching UI drift from AI coding agents》——搜索给出的 URL 404，内容未核实（不采信其具体主张）。
- Polaris stylelint 具体规则页与 npm 包（`stylelint-plugin-polaris`）——官方文档站重定向、npm 页 403，规则描述来自搜索摘要（引用官方规则页标题），未逐页复核。
- objectstack-ai/objectui `AGENTS.md`——原文已抓取（61,741 字符）但未逐条核对，条款样式按搜索摘要转述。
- antd 官方 `.github/copilot-instructions.md` 与 `DESIGN.md` 内容未读。
- Radix/WorkOS 治理分层未调研（时间盒内未做，Polaris/Primer/Geist/antd 已覆盖问题四）。
- skills.sh 安装数（~930K）、awesome 列表其他同名仓库（ComposioHQ/travisvn 等）未复核。
- arXiv:2602.11988（ETH 上下文文件研究）仅经 designproject.io 转述，未读原文。
- 所有 Reddit 讨论帖（CLAUDE.md 遵从率的民间对照）未采纳入结论——无法验证发帖人与样本。

## 引用清单（保留）

一手验证：
- anthropics/skills 仓库与四个 SKILL.md 原文 — 官方技能结构、写法、截图自审、Playwright 工具链的权威依据（https://github.com/anthropics/skills）
- claude.com《Steering Claude Code》— 规则/技能/hooks/子代理的官方分流与「提示≠强制」论断（https://claude.com/blog/steering-claude-code-skills-hooks-rules-subagents-and-more）
- arXiv:2604.11088 — 规则文件有效性唯一实证研究（负向约束 > 正向指令）（https://arxiv.org/abs/2604.11088）
- figma/mcp-server-guide（含 `figma-power/steering/create-design-system-rules.md` 全文）— 设计系统规则生成的官方模板与迭代法（https://github.com/figma/mcp-server-guide）
- designproject.io《Agentic design system》— 四层上下文架构成体系总结（https://designproject.io/blog/agentic-design-system-context/）
- BehiSecc/awesome-claude-skills — 社区技能生态实貌（https://github.com/BehiSecc/awesome-claude-skills）
- 项目内 docs/805-frontend-design-guide.md、AGENTS.md — 背景对照。

转述采信（官方页面/文档经搜索摘要）：
- Playwright MCP 文档、Chromatic docs（review / frontend-workflow-for-ai）、Argos repo+docs、lighthouse-ci configuration、axe-core-npm、ESLint no-restricted-imports、Polaris（VS Code 扩展 + stylelint 规则页）、Primer（react-migrate + discussion #5165 + lifecycle）、Vercel（geist stack/tabs + design guidelines）、Anthropic engineering（equipping agents / context engineering）、usewalkie checklist。

弃用：
- snapdiff.ai 文章（404）、各镜像站（officialskills.sh/agentics.io/baaderagency 等，非一手）、enbi-dev/nexu-io 等 anthropic skills 镜像仓库、未验证的社区技能仓库细节（chrometaphore/majiayu000/rampstackco，仅作线索保留）、Reddit 民间对照帖。
