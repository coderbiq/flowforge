# frontend-vision-loop-poc — Requirements

> Type: PoC(最小验证切片,非完整立项)
> Status: draft
> Sources: `docs/research/2026-09-29-ui-enforcement-agent-governance-research.md` + `2026-09-29-ui-enforcement-toolchain-research.md`(tangram-v2 产出,已导入)

## 问题

tangram-v2 的 UI 执行链(前端实现、视觉自查)当前由纯文本模型或人工承担:

1. 旗舰 `glm-5.3` 无视觉输入 — 无法做"截图自审"(社区主流的 agent UI 质量门,见调研问题三)。
2. flowforge 角色体系没有前端执行角色,UI 任务与后端任务共用 `flowforge-implementer`,设计系统上下文(规范/组件唯一入口/反例集)无路由。
3. 治理调研结论:"skill 承载流程与 checklist + 常驻规则文件只做路由" — 当前缺前端 skill。

## PoC 目标(验证一个闭环)

在 tangram-v2 上验证:**视觉模型 + flowforge 前端角色 + 截图自审 skill** 三件套能跑通
「实现 UI 变更 → 起 dev server → 截图 → 对照反例清单批判 → 修复」循环。

## 范围(In)

- `assets/skills/flowforge-frontend-implement/SKILL.md`:通用 UI 任务工作流(读设计规范 → 组件唯一入口 → 实现 → 截图自审 → 反例自查);项目专属内容留槽位引用
- `assets/subagents/flowforge-frontend-implementer.md`:`model_profile: tool-capable`
- AGENTS.md 模板技能路由表加一行
- 部署产物 + `models_by_name` pin(默认 `cpa/gemini-3.8-flash-high`,备选 `deepseek-v4.1-flash`/`glm-5.3-flash`,均具备 image 输入)

## 范围(Out)

- ESLint/stylelint/axe 工具链(调研推荐组合 1+2+5+6+7+8+9)— 属完整立项,另行 proposal
- Storybook 基建、CI 视觉基线
- 内容层强制(英文直出/内部 ID,调研 17 号缺口)
- Go 侧 tier 体系变更(具名 pin 已够)

## 验收

1. tangram-v2 部署后,`subagent list` 出现 `flowforge-frontend-implementer`(model: gemini-3.8-flash-high)
2. 对一个含反例的页面(建议:store-mate 暗底暗字事故页)跑通截图自审循环,产出 ≥1 张截图 + 对照反例的批判记录
3. PoC 结论落 `observations.md`:视觉模型判读准确率、循环耗时、是否值得完整立项
