---
flowforge:
  schema: 1
  role: design
  id: documentation-refresh-design
  revision: 2
---

<a id="documentation-refresh-design"></a>
# 文档刷新方案（结构决策与分票）

## d-readme-structure：README 新结构

README 重写为"新手可跟跑的入口"，结构钉定：

1. **定位一句话** — 更新为当前事实：Markdown 工件 + 工程 skill 与 subagent 角色 + Go CLI 确定性检查的三件套。
2. **整体设计理论** — 保留现有四权威（Requirement / Solution design / Ticket / Evidence）叙事，压缩为精简版（现有 3 段保留核心）。
3. **快速开始（新手主线，新写，对话驱动，rev 2 修订）** — 六步：
   - 安装（保留现有 install 脚本段）；
   - `flowforge init`（写清副作用：`.flowforge/config.yaml`、`<docs_dir>/` wiki、`.agents/skills/`、`<docs_dir>/agents/`；别名 `sync`）；
   - `flowforge agents deploy`（可选步骤：12 个内置角色编译进 Claude Code / OpenCode / Codex / PI 四宿主；PI 需先 `pi install npm:pi-subagents`）；
   - 在 agent 会话里说出需求（虚构通用场景对话示例；明确“不写任何 Markdown、不建任何目录”；说明 agent 建 proposal、写票、跑 check/frontier）；
   - 用“继续”推进批次（frontier 批次、派发执行角色、四元组证据、双轴审查；三种用户动作：方向/批准/质疑，引语用虚构场景）；
   - 随时看进度（status/check/frontier 为用户可直接敲的看进度命令；PI 原生工具与 `--pi-workflow`）。
4. **场景速览** — 8 行索引表（痛点一句话 + 链接 `docs/scenarios.md`），由票 02 回填。
5. **一个需求如何实际推进** — 案例叙事（虚构通用场景：通知系统重构案），用户话语为叙事主体、方法论角色括号注解，声明“全程用户没有写过一张票”；骨架取自真实实践（rev 2 修订：脱敏，原钉定的真实案例改为虚构）。
6. **核心文档** — 更新链接清单（加 scenarios.md）。
7. **开发 / 许可证** — 保留。

## d-scenarios-landing：场景篇归属决策

场景详文放独立 `docs/scenarios.md`，README 只放索引表。理由：README 已 115 行，8 个场景完整展开（每个 60–120 行）会把入口文档变成手册，与"README 是入口"约束冲突；独立文件允许每个场景带命令块与来源 proposal 链接。

8 个场景与来源 proposal 映射（候选文案已在 g1–g4 工作台文档逐个备好）：

| # | 场景 | 来源 proposal |
|---|---|---|
| 1 | 新手第一张票：15 分钟跑通 init→ticket→frontier | documentation-contract-refinement |
| 2 | 从一份旧 PRD 起步：import 分类交接 + assets verify 查资产漂移 | external-material-intake |
| 3 | 用便宜模型批量执行：轻量模式、evidence 四元组、阻断失败史与 repair 分流 | lightweight-execution-contract、fast-executor-reliability、ticket-refinement-contract、executor-loop-hardening、blocked-evidence-persistence |
| 4 | 角色工厂一次派一批：agents deploy、通用角色、frontier --pi-workflow | subagent-lifecycle、generic-role-orchestration、pi-host-integration |
| 5 | 限额时段换模型不换流程：model-set use、flash 调研/拆轴、价值度量 | model-sets-switching、flash-worksurface-expansion、executor-value-measurement |
| 6 | 前端票自证好看：截图自审循环 + 视觉条款级复核 | frontend-standards-lifecycle、frontend-vision-loop-poc |
| 7 | 规范不靠人肉抄：Align 挑规范、Design 写 must、缺规范的票被退回 | standards-injection、skill-routing-simplification |
| 8 | 多机协作不互踩：模型留在本机 config、部署产物 gitignore、手写 model 保留 | agent-model-preservation、deploy-artifact-localization、wiki-config-single-track |

## d-docs-sync：docs/*.md 修订分配与裁决原则

裁决原则：**代码是唯一事实源**。B/C 类冲突一律以 `internal/` 实现为准改文档；不改代码迁就文档。

- **cli-design.md**（票 03）：A1 completion、A2 model-set、A3 assets verify 入命令清单、A4 check 第 5/6/7 类诊断、A5 未记录 flags、A6 init 真实副作用、A7 sync 别名、A8 config 键全集、A9 raw-script 逃生阀；B1/B2/B9/B11；C1。
- **skill-system.md**（票 04）：A10 六个未记录角色、A11 七个未记录 skill、A15 PI 委派路径；B3（12 角色名册）、B5（front-matter 第三键）、B7（四宿主）；C3。
- **architecture.md + CONTEXT.md**（票 05）：B4（通用角色 4 个）、B10/B12；C2；A13（assets/agents 文件集）、A14（权限标签与 model_profile 词汇表）。

## d-dag：分票与依赖

- 01 README 新手引导重写（root，无依赖）
- 02 场景篇 scenarios.md + README 场景速览节（Blocked by: 01 — 同文件，索引节必须建立在 01 钉定的新结构上）
- 03/04/05 docs/*.md 同步（相互独立，可与 01/02 并行）
- 06 终检收口（Blocked by: 01–05）：A/B/C 清单复核清零、`check --strict`、跨文档命令一致性抽检

执行顺序：本会话先交付 01、02（用户指定优先级）；03–05 经 frontier 后续派发。
