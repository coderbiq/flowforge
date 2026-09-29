# FlowForge

FlowForge 是面向 AI 协同工程的本地优先工作流：Markdown 保存需求、设计、执行票据和证据，工程 Skill 与 subagent 角色分工推进决策，Go CLI 只负责确定性的工件检查与 DAG 计算。它可以在 Claude Code、OpenCode、Codex 和 PI 四种 agent 宿主上运行同一套方法论。

## 整体设计理论

FlowForge 把工程工作分成四种权威内容：

- **Requirement** 说明为什么做、可观察结果、范围、场景、约束和仍未确定的需求事实。
- **Solution design** 说明如何实现：模块责任、接口与 seam、信息流、迁移顺序、替代方案和验证策略。
- **Ticket** 说明一个可独立验证的执行增量：交付结果、局部设计上下文、触点、有序动作、约束及验证方法。
- **Evidence** 说明实际交付了什么、运行了哪些验证、双轴审查如何处置以及对应实现引用。

人类可读正文是语义权威；ID、revision 和 `consumes` 只为机器追踪服务。流程不保存 `requirements-ready`、`design-ready` 一类阶段状态——能否推进由当前文件事实推导：DAG blocker 必须先解决，设计 gap 只影响明确关联的区域，warning 保持可见但默认允许继续。简单需求可以把四种角色压缩在一张 ticket 中，复杂度决定工件形态，不决定一套固定流程状态。

## 快速开始

### 1. 安装

Linux / macOS：

```bash
curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash
```

Windows PowerShell：

```powershell
irm https://github.com/coderbiq/flowforge/releases/latest/download/install.ps1 | iex
```

### 2. 初始化项目

在项目根目录运行：

```bash
flowforge init
```

`init`（别名 `sync`）会创建 `.flowforge/config.yaml`、`ff-wiki/`（含 `CONTEXT.md`、`adr/`、`proposals/`、`agents/`）、部署工程 Skill 到 `.agents/skills/`。文档根目录默认 `ff-wiki`，可换：

```bash
flowforge config set docs_dir my-wiki
flowforge init --force
```

`init --force` 同步受管资产但保留已有项目配置；部署产物与本机配置会被自动加入 `.gitignore`（每机器文件，见场景 8）。所有命令从子目录运行时会向上查找 `.flowforge/config.yaml`。

**初始化后有一个关键一步：让 AI 重写 `standards.md`。** `init` 部署的 `<docs_dir>/agents/standards.md` 是规范提取说明，默认内容只是通用模板（里面标着“替换为你项目的真实规范”）。在第一次正式提需求前，让 AI 按项目实况重写它：

> “读一下这个项目的分层约定、测试规范和提交规范，把 ff-wiki/agents/standards.md 重写成能指导后续所有 ticket 的规范提取说明。”

重写后，Align 会按这份说明为每个需求挑出适用规范，Design 转写成 `must` / `must not`，Plan 机械转录进每张票（见场景 7）。

尊重项目差异是 FlowForge 的底层设计：文档根（`docs_dir`）、规范说明位置（`standards.guide`）、启用的宿主（`agents.hosts`）、每个角色的模型（`agents.models_by_name` / `models_by_host` / `model_sets`）、测试文件保护范围（`agents.test_file_globs`）等全部按项目、按机器配置，不假设你的仓库长什么样。

### 3. 部署 subagent 角色（可选，推荐）

```bash
flowforge agents deploy
```

把 12 个内置角色（analyst、architect、planner、implementer、reviewer 等流程角色，加 batch-analyst、scribe、executor、investigator 等通用角色）编译成宿主原生子代理文件：`.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/`。只想启用部分宿主时配置 `agents.hosts`；`flowforge agents status` 检查部署漂移。

PI 宿主需要先安装子代理扩展：

```bash
pi install npm:pi-subagents
```

PI 宿主还会得到项目级扩展 `.pi/extensions/flowforge.ts`：拦截对受保护测试文件的写入、把 `flowforge_frontier` / `flowforge_check` 注册为 LLM 原生工具。

角色的模型分配留在本机 `.flowforge/config.yaml`；`flowforge model-set list / show / use` 可在多套命名模型方案间一键切换（自动重部署、失败回滚，见场景 5）。

### 4. 在 agent 会话里说出你的需求

你不需要写任何 Markdown、不需要建任何目录。在宿主会话（PI / Claude Code / OpenCode / Codex）里用业务语言描述方向，比如：

> “通知模块上次重构得比较乱，我整理了一版旧设计放在 docs/old-design.md。把消息模板管理拆出去单独立项，再开一个专题修复发送重试的 proposal。”

Agent 会接着做：先查代码和现有文档，只向你询问仓库回答不了的产品取舍；然后在 `<docs_dir>/proposals/<feature>/` 下建立 proposal，把需求与设计事实写成 Markdown（`requirements.md`、`design.md`）；再用 `flowforge-plan` 生成带 DAG 依赖的执行票；最后跑 `flowforge check` / `frontier` 给你看就绪队列。从旧 PRD / 旧提案起步时，`flowforge-import` 会先把来源材料分类（事实 / 需求候选 / 设计决定 / 证据 / 未知），不会机械照搬。

生成的票长这样（全部由 agent 写，你只需要能看懂）：标题、交付结果、有序变更清单、约束与验证命令，加上机器可核对的五段执行契约；真正阻止开工的依赖写成 `Blocked by`，由 CLI 计算拓扑顺序。

### 5. 用“继续”推进批次

每一轮推进你只需要说话，最高频的一个词就是**“继续”**。Agent 收到后：跑 `flowforge frontier` 取无阻塞批次 → 能并行的票派给干净上下文的执行角色（每票一个）→ 执行者按轻量模式机械执行变更并留下证据四元组（`cmd`/`exit`/`output`/`artifact`）→ `flowforge-review` 双轴审查（Standards / Specification）把发现翻译成 `Fix:` 变更追加回票，实质发现则开 repair 票 → 零发现才写 Completion evidence 并关票。

你的三种真实动作（来自真实使用记录）：

- **方向**：“我的理解是发送前不应该每条消息都查一次数据库，应该看汇总模块是怎么批量取数的”——纠偏设计；
- **批准**：“同意，继续推进先落盘在设计”；
- **质疑**：“不太看得明白当前这些票在做什么”、“批次 5 是什么？”——agent 解释或调整拆票。

### 6. 随时看进度

看进度类命令是你唯一可能直接敲的（推进中的 check / frontier 由 agent 自动跑）：

```bash
flowforge status     # 按 feature 汇总进度
flowforge check      # DAG、诊断与证据检查（推进中 agent 会自动跑）
flowforge frontier   # 无阻塞就绪票投影（clean/warning/gap/blocked 分档）
```

PI 宿主里连命令都不用敲：`flowforge_frontier` / `flowforge_check` 已注册为 LLM 原生工具，直接问“下一批该干什么”；要跑整批时用 `flowforge frontier --pi-workflow` 生成 pi-subagents 执行脚本。

## 场景速览

| 场景 | 痛点一句话 |
|---|---|
| [1. 说出一个需求](docs/scenarios.md#场景-1说出一个需求看它变成-proposal) | 是不是要学一堆 Markdown 模板？都不用 |
| [2. 从旧 PRD 起步](docs/scenarios.md#场景-2从一份旧-prd-起步) | 几百行旧 PRD 说不清哪段还算数 |
| [3. 用便宜模型批量执行](docs/scenarios.md#场景-3用便宜模型批量执行一个-proposal) | flash 会“勾了却没做”、失控重试 |
| [4. 角色工厂一次派一批](docs/scenarios.md#场景-4角色工厂一次派一批票) | 主会话不该亲自读 20 个文件 |
| [5. 限额时段换模型](docs/scenarios.md#场景-5限额时段换模型不换流程) | 旗舰限速，想整体切便宜模型 |
| [6. 前端票自证好看](docs/scenarios.md#场景-6前端票自证好看) | 读代码推不出颜色对不对 |
| [7. 规范不靠人肉抄](docs/scenarios.md#场景-7规范不靠人肉抄进每张卡) | 漏抄就漏执行 |
| [8. 多机协作不互踩](docs/scenarios.md#场景-8多机协作不互踩) | 同事 commit 覆盖你的模型配置 |

完整场景（含命令路径与来源 proposal）见[场景导览](docs/scenarios.md)。

## 一个需求如何实际推进

以下是一个完整案例的叙事（虚构的通用场景：一个团队协作应用的通知系统重构；对话为示意引语，工作流骨架来自真实使用记录），方法论角色在括号内标注——**全程用户没有写过一张票、建过一个目录**。

### 1. 拆分立项

> 用户：“notification-delivery 这个提案太大了，并且经过多轮讨论其中有一些错误。把消息模板单独创建一个 proposal 后面再来分析设计，当前把其它问题收个尾，然后创建一个新的发送重试 proposal 来专题修复这个环节的问题。”

Agent 读取旧提案的全部工件与源笔记，把仍然有效的事实分类保留（`flowforge-import`：事实 / 需求候选 / 设计决定 / 证据 / 未知五类，不机械照搬），建立两个新 proposal 目录并写下各自的需求事实（`flowforge-align` 只保存会改变方案空间的事实，其余先查代码与现有文档）。

### 2. 澄清与设计

> 用户：“对整个方案做一次 review，避免出现消息发出去了才发现模板没渲染这样的低级错误。”
> 用户：“我不清楚你在说的批次 2 批次 3，不过我猜测当前设计还有大量性能问题——发送前不应该每条消息都查一次数据库，这一点可以看汇总模块是怎么批量取数的。”

Agent 写全景评审文档、比较责任边界与 seam、迭代方案（`flowforge-solution-design`：把每个设计区域的已定/警告/缺口写清楚，项目规范在此时写成 `must` / `must not`），用户对齐后批准：

> 用户：“同意，继续推进先落盘在设计。”

### 3. 生成执行票

Agent 用 `flowforge-plan` 生成带真实 DAG 依赖的执行票（每张含交付结果、有序变更、约束、验证命令与五段机器执行契约），随后跑 `flowforge check` / `frontier` 展示就绪队列。看不懂就问：

> 用户：“不太看得明白当前这些票是在做什么。”

Agent 逐票解释、必要时重拆；缺契约的票会先经 `flowforge-refine-ticket` 补齐仓库事实再放行。

### 4. 批次执行与审查

用户说“继续”，Agent 每轮跑 `flowforge frontier` 取无阻塞批次，把票派给干净上下文的执行角色：轻量模式机械执行变更并留证据四元组（`cmd`/`exit`/`output`/`artifact`），`flowforge-review` 双轴审查把发现翻成 `Fix:` 变更回票、实质发现开 repair 票。用户随时可以纠偏：

> 用户：“重试间隔在两个配置里重复定义了，应该合并成一个。”

执行者被阻断时必须先写下逐字错误与已试命令再返回；这段失败史会被转写为持久契约，重派的新上下文不再踩同一个坑。

### 5. 收口

每关一票就重跑 `check` + `frontier` 取下一批；全部票关闭且诊断处置完毕，proposal 完成。收口后仍发现质量问题时不手改：直接补开新票修复（真实实践：一批票关闭后检出未提交文件与多类收口问题，随即补开补救票并用 `check`/`frontier` 重算依赖）。CLI 计算拓扑顺序，Agent 不在上下文里自行猜测。

## 核心文档

- [场景导览](docs/scenarios.md)
- [架构与权威模型](docs/architecture.md)
- [Skill 职责与协作](docs/skill-system.md)
- [CLI 行为与策略](docs/cli-design.md)
- [域词汇表](docs/CONTEXT.md)

## 开发

```bash
make dev
GOPROXY=https://goproxy.cn,direct go test ./internal/...
```

许可证：MIT。
