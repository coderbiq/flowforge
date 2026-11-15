---
flowforge:
  schema: 1
  role: requirement
  id: deploy-artifact-localization-requirements
  revision: 1
---

<a id="deploy-artifact-localization-requirements"></a>
# 部署产物本地化与模型注入需求

## 问题

1. **pi 宿主无 per-subagent 模型通道**：`compile_pi.go` 按 PI 宿主集成方案（`../pi-host-integration/design.md`）§一明文"不写 `model`（省略=继承父会话模型）"，子代理能力分层只剩 `thinking` 字段，换机部署后无法重现角色模型分配（用户在新电脑部署时发现）。
2. **多人协作可移植性**：若将 `.opencode/`、`.pi/`、`.claude/` 等部署产物提交入库，产物 frontmatter 中的 `model:` 会被应用到所有协作机器；不提供该模型/provider 的机器运行即错。[宿主模型外置能力调查](../../research/2026-09-20-host-model-portability.md)证实各宿主"定义文件外置模型"通道强度差异大（pi 可覆盖 frontmatter、opencode 仅在 frontmatter 省略时生效、claude 无 per-role 外置），外置路线不可统一。
3. **本仓库现状已分化**：`.gitignore` 已忽略 `.agents`、`.opencode`，但 `.claude/agents/`、`.codex/agents/`、`.pi/agents/`、`.pi/extensions/flowforge.ts` 仍被 git 跟踪；`.flowforge/config.yaml` 未被跟踪但也未被忽略。
4. **附带 bug**：preserve-merge 对 pi 宿主虚假承诺——手编 `.pi/agents/<name>.md` 加 `model:` 后重新部署，stderr 打印 preserved 提示，但 `CompilePiWithOptions` 忽略 `FallbackModel`，值实际被丢弃（agent-model-preservation 的 preserve-merge 循环宿主无关，未按宿主门控；见 `../agent-model-preservation/design.md`）。

## 目标

采用「模型注入产物 + 产物本地化」路线（用户裁决；否决各宿主用户级外置通道与项目级 `.pi/settings.json` 生成方案）：

1. 模型配置留在每机器的 `.flowforge/config.yaml`，`agents deploy` 时按配置注入各宿主 subagent 定义（含 pi frontmatter `model:`，修订 PI 宿主集成方案 §一）。
2. 部署产物与 `.flowforge/config.yaml` 定性为每机器文件不入库；init 自动配置 `.gitignore`。
3. 修复 preserve-merge 对 pi 的虚假保留。

## 范围与约束

- **config.yaml 语义**：每机器文件（不入库）。团队共享知识载体是 `docs/`（入库）；config 是本机运行时状态，`agents.hosts` 本就应每机器不同（每人使用不同 Coding Agent）。
- **.gitignore 自动化**（init 幂等写入，已存在条目不重复）：`.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/`、`.pi/extensions/`、`.agents/` 整目录，加 `.flowforge/config.yaml` 单文件（`.flowforge/` 其余内容如 `subagents/` 自定义源仍入库）。用户自写且想入库的 agent 以 `git add -f` 例外，AGENTS 模板写明此约定。
- **存量迁移**：upgrade/init 检测到已跟踪的受管产物时不静默改动 git 索引，输出应执行的 `git rm --cached` 指引（是否提供显式 flag 由 design 定）；本仓库存量一次性手工清理。
- **模型值校验**：deploy 时对启用宿主做格式级校验，非法即报错、不写坏配置（fail-fast）；不做存在性校验（纯本地无网络原则）。
- **per-host 模型覆盖层**：`agents.models`（profile 键）为全宿主默认，新增仅对单一宿主生效的覆盖层（键命名归 design）；优先级 per-host 覆盖 > `agents.models`/`models_by_name` > 未配置。
- **未配置时行为不变**：claude 保持账号可移植别名档位（opus/sonnet），opencode/pi 省略 `model`，codex 无 model 概念（`model_reasoning_effort` 映射维持）。
- **不新增 CLI 命令**；不生成项目级 `.pi/settings.json`；不管理各宿主用户级外置通道。
- `models_by_name`（per-agent）语义继续；其与 per-host 层的组合优先级排序归 design。

## 可观察验收

1. pi 宿主：config 任一层配置了某 agent 的模型 → `agents deploy` 后 `.pi/agents/<name>.md` frontmatter 含 `model: <值>`；未配置 → 不含该键。
2. 多宿主同开（如 pi + claude）且仅配置 per-host 覆盖时，两宿主产物分别得到各自宿主的值。
3. 格式非法的模型值（对启用宿主）→ deploy 报错退出、可见失败，不产出坏配置文件。
4. 全新项目 `flowforge init` 后 `.gitignore` 含上述全部条目；重复执行不产生重复行。
5. 手编 `.pi/agents/<name>.md` 加 `model:` 后重新 deploy → 值被真实保留且 stderr 提示正确；config 显式配置该模型时 config 胜出（对齐 agent-model-preservation 的显式通道优先级）。
6. 存在已跟踪受管产物的仓库执行 upgrade/init → 得到 `git rm --cached` 指引输出，git 索引不被自动改动。

## 术语

- **部署产物（deploy artifacts）**：由 `flowforge agents deploy`/assets 部署编译写出的宿主文件（四宿主 agent 目录、`.pi/extensions/flowforge.ts`、`.agents/skills/`）。
- **每机器文件（per-machine file）**：不随 git 旅行、由本机 init/deploy 生成或维护的文件（部署产物与 `.flowforge/config.yaml`）。
- **per-host 模型覆盖（per-host model override）**：config 中仅对单一宿主生效的模型配置层。

## Standards 识别清单（转 solution-design 转换为 must 条款）

- [AGENTS.md] 核心设计原则 2：纯本地确定性文件操作，无网络、无 LLM 调用（约束模型校验上限与 gitignore/init 行为）。
- [AGENTS.md] boundaries：变更后运行 `go test ./internal/...`。
- [AGENTS.md] Ask first 边界：修改 Issue Schema 头规范、变更 CLI 接口签名（本变更新增 config 键与 init/upgrade 行为，不改命令签名；若涉及 schema 头变更需回问）。
- [AGENTS.md] 🚫 Never：`assets/` 只放部署内容（AGENTS 模板新增 gitignore 约定属部署内容，合规）。
- [agent-model-preservation 需求] 显式通道优先级：config 已设值时 config 编译结果为准，不做保留合并（本需求验收 5 与其一致性约束）。
