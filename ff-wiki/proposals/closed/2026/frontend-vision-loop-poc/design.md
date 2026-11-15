# frontend-vision-loop-poc — Design

## 分工边界(单一真相源)

- **flowforge 仓(本 proposal)**:通用层 — skill 流程骨架、角色定义、模型 pin、AGENTS.md 路由
- **tangram-v2**:项目层 — 805 规范、@tangram/ui 组件唯一入口清单、17 条反例集、PoC 验证场与 fixture

两份调研文档已在 flowforge 仓 `docs/research/`,tangram-v2 侧不再重复设计 flowforge 层机制。

## 交付物设计

### 1. flowforge-frontend-implementer(角色)

- frontmatter:`model_profile: tool-capable`、`thinking: high`(视觉判读需要)
- 描述:Delivers frontend/UI tickets with design-system context and a mandatory
  screenshot self-review loop. Pinned to a vision-capable model.
- 与 flowforge-implementer 的边界:ticket 涉及 JSX/样式/布局/组件层 → frontend-implementer;
  其余 → implementer。由 AGENTS.md 路由表表达。

### 2. flowforge-frontend-implement(SKILL)

流程骨架(每步负向约束优先,依据治理调研问题二结论):

1. 读项目设计规范槽位(路径由项目 AGENTS.md 声明;未声明 → STATUS: BLOCKED)
2. 组件唯一入口清单自查(禁 app 层直 import UI 底层组件)
3. 实现(负向清单:禁手拼 Layout/Header、禁内联硬编码色值、禁绕过 ThemeProvider)
4. 起 dev server → Playwright 截图(desktop+mobile)
5. 对照项目反例集逐条批判截图(需要 image 输入 — 这就是模型 pin 的依据)
6. 修复 → 重复 4-5,稳定后交付,截图与批判记录附进 ticket

### 3. 模型选型(决策已定,可改)

| 候选 | 依据 |
|---|---|
| `cpa/gemini-3.8-flash-high`(默认) | gemini 系视觉判读强,high 档推理 |
| `cpa/deepseek-v4.1-flash` | 已验证工具调用稳定,同为 image 输入 |
| `cpa/glm-5.3-flash` | 生态同源备选 |

注:当前旗舰 `glm-5.3` 纯文本(pi models.json `input: ["text"]`),视觉轴不可用 —
旗舰侧若未来要参与视觉 review,需换或加 vision-capable 旗舰,另行决策。

### 4. 部署与配置

- `flowforge agents deploy` 现有机制,零 Go 改动
- 各项目 config:`models_by_name: { flowforge-frontend-implementer: cpa/gemini-3.8-flash-high }`

## 风险

- gemini-3.8-flash-high 工具调用可靠性未在本项目验证(PoC 验收项 3 覆盖)
- 截图自审的 dev server 生命周期管理(端口/启动等待)由 SKILL 指令承载,PoC 观察
- 治理调研警示:提示式条款最弱区是"正向品味"目标 — PoC 只验证反例拦截,不承诺品味提升
