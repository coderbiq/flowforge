# 调研:前端强制方案的普适性 — 通用 skill 分层、跨栈锚点、视觉验证收敛点(第三轮)

> 2026-09-29,flowforge 仓。目的:检验 frontend-vision-loop-poc 的流程设计是否过度拟合
> tangram-v2 单项目场景,为"flowforge 通用 proposal + tangram-v2 项目 proposal"的切分
> 提供社区参照。前两轮调研(治理面/工具链)见同目录 2026-09-29-ui-enforcement-*.md。

## Q1:通用 skill 与项目上下文的分层先例

**结论**:官方没有强制的 generic/project 切分指南,但存在三层惯例与一套可依赖的设计原则:

1. 社区惯例:项目级 skill 放 repo 内(`.claude/skills/` 等价 `.agents/skills/`),跨项目
   可复用的提取到用户级或插件/共享仓;**提取时机 = 在多个仓成功运行且改动极小时**
   (而非预先抽象)。
2. 官方 best-practices 的可迁移原则:
   - **自由度分级**:高(文本指令,多解皆可)/ 中(伪代码+参数)/ 低(具体脚本,
     "操作脆弱易错、一致性关键、必须固定序列"时)。确定性检查**用脚本不用 prose**。
   - **渐进披露**:SKILL.md <500 行作目录,references 一层深,长 reference 带目录。
   - **validator 反馈循环**:"跑验证器 → 修 → 重复"显著提升产出;参考文档
     (如 STYLE_GUIDE.md)可充当 validator — 与本仓"反例集槽位"同构。
3. 对本方案的检验:PoC skill 的"三槽位 + 量化脚本"结构与上述原则一致;
   **风险不在分层而在槽位内容的形态假设**(见 Q2/Q3)。

来源:platform.claude.com agent-skills/best-practices(渐进披露、自由度分级、
validator 循环);第一轮检索聚合(社区 skill 组织惯例)。

## Q2:跨技术栈的通用锚点

**结论**:token 管道有跨栈标准锚;组件边界强制是 JS/React 生态内的,跨栈需适配层。

- **W3C Design Tokens (DTCG) 2025.10**:Community Group Report(非正式 Recommendation),
  但已是事实锚点:token 源(DTCG JSON)→ CI 语义校验 → Style Dictionary 生成平台产物
  → stylelint 强制消费侧(`declaration-strict-value` 等)。该管道框架无关
  (CSS/Less/原生/生成物均可消费)。
- 组件唯一入口强制(our component-entry-list)依赖 ESLint `no-restricted-imports` /
  `eslint-plugin-boundaries` / `react/forbid-elements` — **JS 生态内**,Vue/Svelte 有
  等价规则但非同一套;原生/非 JS 栈无直接等价。
- **修订含义**:flowforge 通用层的"组件入口槽位"应声明为**项目自证的契约**
  (项目提供白名单文件 + 自检命令),而不是规定用哪种 lint 实现。

## Q3:视觉验证的宿主抽象

**结论**:Playwright(+MCP)是 web 前端的社区收敛点;自主截图验证的生产采用度
仍有限(共识是"agent 探索/取证用 MCP,稳定断言落 CI `toHaveScreenshot`");
非 web 前端(移动/native)无收敛方案。

- 收敛的工作分工:**交互走 accessibility snapshot,截图专用于视觉验证**(与本仓
  PoC skill 第 3 步一致 — 已验证为社区同构)。
- **修订含义**:范围声明 web-only;截图入口作为**项目声明的验证命令槽位**
  (tangram-v2 用 poc-screenshot.mjs,别的项目可以是 Playwright MCP、Storybook、
  自有脚本),skill 只约定"可重复入口 + 量化输出契约"。

## Q4:多项目复用与版本化

**结论**:社区模式 = 中心化版本 skill 目录 + 本地 overlay;skill 打包为
SKILL.md + scripts + references + assets(+tests),按不可变 SHA 锚定版本;
确定性 CI 与 agent 判断分离,配 ownership/评审/弃用策略。

- 对本仓:flowforge 已具备等价机制(assets/skills 源 + agents deploy 分发 +
  AGENTS.md 路由),**无需引入新分发机制**,补的是槽位协议的稳定性承诺
  (版本化契约文档)。

## Q5:既有 "agentic design system" 实践的跨项目形态

**结论**:可参照实践(含 designproject.io 四层架构)收敛为:版本化 DTCG token +
共享组件库 + 项目级主题 + **bounded named slots(有界命名槽位)**+ 组件元数据
对 agent 可读(文档或 MCP);强制落在 CI(import/token/a11y/视觉),agent 只
在批准的原语内组装。

- "bounded named slots" 与本仓三槽位协议同构 → 分层方向有社区印证。
- 组件元数据 MCP 化(Storybook MCP 等)是比静态白名单文件更进一步的方向,
  可列为完整 proposal 的可选项而非必选。

## 综合修订建议(skill v2)

1. 槽位协议从"三文件"升级为**声明式 frontend context 契约**:design-tokens 源
   (推荐 DTCG)、component-entry(项目自证白名单 + 自检命令)、counter-examples
   (任意形态的反例清单)、verification-entry(**项目声明的可重复视觉验证命令 +
   量化输出**)。skill 不假设 Playwright/antd/React。
2. 验证入口与量化脚本属项目层(低自由度、官方"确定性用脚本"原则);flowforge 只
   约定输出契约(可 grep 的 PASS/FAIL + 数值)。
3. 范围声明:web-only;非 web 宿主为 adapter 扩展点,不在本期。
4. 分发沿用现有 assets + deploy;槽位协议文档版本化。

## 开放问题

1. 槽位契约的版本化位置(AGENTS.md vs 独立 manifest)与向后兼容策略
2. design-tokens 槽位是否首期强制(tangram-v2 现状是 antd token 体系,非 DTCG —
   迁移是项目层决策)
3. 组件元数据 MCP 化何时值得投入(先看静态白名单的拦截率数据)

## 引用

- https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices
- https://www.w3.org/community/reports/design-tokens/CG-FINAL-format-20251028/
- https://www.designtokens.org/TR/2025.10/format/
- https://playwright.dev/mcp/screenshots · https://playwright.dev/docs/next/test-snapshots
- https://github.com/microsoft/playwright-mcp · https://storybook.js.org/docs/ai/mcp/overview
- https://stylelint.io/user-guide/customize/
- agentskills.io skill specification(版本化/portable 结构)
