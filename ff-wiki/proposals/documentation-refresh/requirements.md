---
flowforge:
  schema: 1
  role: requirement
  id: documentation-refresh-requirements
  revision: 4
---

<a id="documentation-refresh-requirements"></a>
# 文档刷新需求（README 新手引导 + 场景篇 + docs 同步）

## 问题

FlowForge 仓库 20 个 proposal（85/85 ticket）全部交付后，用户可见文档没有跟随：

1. **README 新手引导断裂**。README 仍停留在 2026-08-25 documentation-contract-refinement 之前的表面：全文零次提及 `subagent`、`agents deploy`、`model-set`、PI 宿主（research/2026-09-29-doc-refresh/g2 §Summary）。新用户在 `flowforge init` 之后没有路径指引——不知道可以部署角色、不知道怎么让 agent 建第一个 proposal、不知道 check/frontier 循环怎么转、不知道票可以派发给谁。
2. **已交付 feature 无场景化介绍**。弱模型批量执行、model-set 成本切换、PI 一键批次、前端视觉自审、规范注入、多机部署隔离等 20 个已交付能力只有 proposal 内部记录，用户没有"哪个场景用哪个 feature"的导览。
3. **docs/ 内部事实冲突**。`docs/cli-design.md:7` 写 `docs_dir` 默认 `docs`，真实为 `ff-wiki`（`internal/config/config.go:17`），且与 `docs/architecture.md:66`、`README.md:92` 互相矛盾；subagent 名册规模、front-matter 键集合、宿主清单、逃生阀清单等 12 项 B 类不符与 4 项 C 类失效引用仍未处置（research/2026-09-29-doc-refresh/current-surface-gaps.md §B/§C）。

事实基线：`docs/research/2026-09-29-doc-refresh/`（g1–g4 四份 proposal 事实提取 + current-surface-gaps 差异盘点，全部带 file:line 引用）。

**关键事实修正（rev 2，用户反馈 + 会话取证）**：FlowForge 的真实工作模式是**对话驱动**——用户从不手写任何 Markdown 票、不手工建目录；用户在 agent 会话里用业务语言提出方向/批准/质疑（最高频输入是“继续”），由 agent 激活 skill 链创建 proposal、写票、跑 check/frontier、派发执行。CLI 是 agent 的确定性工具，用户直接敲的只有看进度类命令。取证见 `docs/research/2026-09-29-doc-refresh/real-usage-workflow.md`（内部项目会话与工件分析，敏感细节已脱敏）。任何把“手写票”当作用户路径的文档表述都是错的。

**脱敏约束（rev 3，用户要求）**：对外文档（README、docs/scenarios.md）不得出现任何真实项目信息——本地路径、项目名、业务域名、账号、真实业务模块名一律禁止；案例叙事参考真实实践的工作流骨架，但用虚构的通用应用开发场景（通知系统重构案）呈现。仓内 research 取证文档同步脱敏：只保留方法论结论与泛化的话语类型，不保留原文引语与具体存储路径。

**关键 feature 补录（rev 4，用户反馈）**：安装后初始化包含一个必做环节——`init` 部署的 `<docs_dir>/agents/standards.md` 默认是通用模板，需提示用户让 AI 按项目实况重写它（`standards.guide` 默认指向 `agents/standards.md`，消费方为 align/plan/setup skill 链）。且“项目可定制化、尊重项目差异”是重要 feature：`docs_dir`、`standards.guide`、`agents.hosts`、`agents.models_*`、`agents.test_file_globs` 等按项目/按机器配置，不假设仓库形态。README 快速开始与场景 7 必须呈现这一点。

## 可观察结果

1. 新用户仅按 README 即可从安装走到跑通第一个循环：install → `init` → （可选）`agents deploy` → **在 agent 会话里用一句话说出一个需求** → agent 建立 proposal 并生成带 DAG 的票 → `check`/`frontier` 显示就绪批次 → 用“继续”推进到收票；全程用户不需要手写任何 Markdown；每条命令与 flag 与当前二进制 `--help` 输出一致。
2. `docs/scenarios.md` 提供 8 个精选场景，每个场景含痛点背景、操作路径（真实命令/skill/角色）、涉及 feature 与来源 proposal；README 有指向它的场景速览索引。
3. README 与 docs/*.md 中 current-surface-gaps 清单里指向这些文件的全部 A/B/C 项被消除或显式裁决。
4. `flowforge check --dir docs/proposals/documentation-refresh --strict` 通过。

## 范围

**In**：`README.md`、`docs/scenarios.md`（新增）、`docs/cli-design.md`、`docs/skill-system.md`、`docs/architecture.md`、`docs/CONTEXT.md`、本 proposal 的 tickets。

**Out**：任何代码 / skill 正文（`assets/skills/`）/ subagent 定义（`assets/subagents/`）改动；competitor-reports 与 research 归档；CLI 行为变更。

## 验收场景

- S-验收-1：一名未接触过 FlowForge 的工程师只读 README，在 15 分钟内于临时目录完成 install→init→在会话里说出一个需求→看到 agent 生成 proposal 与票→frontier 显示 ready→check 通过。
- S-验收-2：读者在 README 场景速览中点进任一场景，`docs/scenarios.md` 能完整回答"这解决什么痛点、我该敲什么命令、feature 从哪个 proposal 来"。
- S-验收-3：`grep -n 'docs_dir' docs/cli-design.md` 等历史冲突点复查，全部与 `internal/config/config.go` 一致。

## 约束

- 正文语言保持中文；代码标识与既有术语稳定（`docs_dir`、frontier、ticket 等）。
- 不虚构未交付功能：所有 feature 表述以 20 个 proposal 的 issues 状态与代码事实为准（g1–g4 文档为事实来源）。
- 脱敏：对外文档不出现真实项目信息；案例用虚构通用场景，叙事骨架可取自真实实践（rev 3）。
- README 是入口文档：场景详文放 `docs/scenarios.md`，README 只放索引（决策见 design）。
- 命令示例必须真实可跑，以 `--help` 输出为准，不写计划中的命令。
