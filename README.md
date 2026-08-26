# FlowForge v5

FlowForge 是面向 AI 协同工程的本地优先工作流：Markdown 保存需求、设计、执行票据和证据，工程 Skill 分工推进决策，Go CLI 只负责确定性的工件检查与 DAG 计算。

## 整体设计理论

FlowForge 把工程工作分成四种权威内容：

- **Requirement** 说明为什么做、可观察结果、范围、场景、约束和仍未确定的需求事实。
- **Solution design** 说明如何实现：模块责任、接口与 seam、信息流、迁移顺序、替代方案和验证策略。
- **Ticket** 说明一个可独立验证的执行增量：交付结果、局部设计上下文、触点、有序动作、约束及验证方法。
- **Evidence** 说明实际交付了什么、运行了哪些验证、双轴审查如何处置以及对应实现引用。

人类可读正文是语义权威；ID、revision 和 `consumes` 只为机器追踪服务。下游工件链接并概括当前任务真正需要的上游含义，不复制整段需求或设计。

流程不保存 `requirements-ready`、`design-ready` 一类阶段状态。是否能推进由当前文件事实推导：DAG blocker 必须先解决；设计 gap 只影响明确关联的区域；warning 保持可见但默认允许继续。需要承担风险时，可显式使用 `--include-gaps` 或精确 waiver，而不是伪造状态让门禁通过。

简单需求可以把四种角色压缩在一张 ticket 中。只有当内容需要独立评审、被多个 ticket 消费、跨会话长期维护或已经妨碍阅读时，才提升为 `requirements.md`、`design.md`、`spec.md` 或独立 evidence。复杂度决定工件形态，不决定一套固定流程状态。

## 一个需求如何实际推进

假设用户提出：

> “让 `flowforge check` 显示 warning 数量，保持现有诊断顺序和默认退出行为。”

### 1. 找到下一位责任人

当入口不明确时，`flowforge-route` 只判断当前尚未解决的问题属于谁，不创建内容。如果结果、范围或约束还不清楚，交给 `flowforge-align`；如果需求清楚但接口、责任或验证 seam 未定，交给 `flowforge-solution-design`。

这个例子结果明确、影响局部，并复用现有命令输出 seam，因此可以跳过独立设计文档，直接进入 Plan；ticket 的 Design context 会说明它复用的既有 seam，而不是记录一个 “design-ready” 状态。

如果需求从 PRD、旧 proposal、会议记录或其它本地材料开始，先使用 `flowforge-import`。它按来源事实、需求候选、设计决定、交付/验证证据、未知/冲突分类，并保留文件与标题定位；它不机械转换文档，也不创建 ticket。需求候选回到 Align，模块/seam/迁移等选择交给 Solution Design。

外部材料的完整推进是：Import 分类并交接 → Align 接受需求候选并发布 requirement → 必要时 Solution Design 决定责任与 seam 并发布 design → 每次 authority 发布后以 `check --strict` 验证 → Plan 展示 title、Delivery 和真实 DAG 边，等待用户接受后才写 issue。这个顺序是工件关系的即时校验，不是一套可卡住的状态流转。

### 2. 保存需求和设计事实

`flowforge-align` 只保存会改变方案空间的事实，例如：warning 数量必须可观察、诊断顺序不能改变、默认模式不能因 warning 失败。它先查代码和现有文档，只向用户询问仓库无法回答的产品取舍。

从来源材料形成 requirement/design 时，正文按目标语言重新表达确认含义，代码标识与既有术语保持稳定。每次发布或修订带 schema metadata 的 authority 后，作者运行 `flowforge check --dir <feature-dir> --strict`；这是对 anchor、revision、链接和 open item 的发布自检，不会写入 readiness 状态或阻塞未受影响工作。

如果需求改成“统一 check/frontier/status 的诊断策略”，就会触发 `flowforge-solution-design`：比较责任边界，确定共享策略接口、调用方、失败语义、兼容顺序和测试 seam，并把每个设计区域的 resolved、warning、gap 或 blocker 写清楚。Plan 只消费已经确定的区域。

### 3. 发布可执行 ticket

`flowforge-plan` 在 `<docs_dir>/proposals/<feature>/issues/` 写 Markdown ticket。它必须包含一项可观察交付、当前任务需要的设计上下文、稳定触点、有序变更、容易违反的约束，以及成对的完成条件和验证命令。真正阻止开工的依赖才写入 `Blocked by`。

发布后运行：

```bash
flowforge check
flowforge frontier
```

`check` 校验 DAG、工件角色、authority revision、语义链接、scoped open item、waiver 和完成证据。`frontier` 只投影 `issues/` 下的可执行 ticket，并分开展示 clean、warning、gap 和 blocker。

### 4. 实现、审查并留下证据

`flowforge-implement` 读取 ticket 及其链接的有效需求和设计，在已选 seam 上用 TDD 交付。实现过程中：路径或命令变化可作为事实修正；责任、接口或 seam 变化返回 Solution Design；可观察行为变化返回 Align。

完成代码后，`flowforge-review` 对同一个固定 diff 分别检查：

- **Standards**：是否符合仓库规范和代码质量基线；
- **Specification**：是否真正实现有效需求与设计。

Implement 处理完两轴发现，记录实际验证结果、处置、偏差和提交引用，最后才把 ticket 设为 `closed`。没有 `Completion evidence` 的 closed ticket 会被 `check` 警告，并被 `check --strict` 拒绝。

### 5. 继续 DAG 前沿

每完成一个 ticket，再运行 `flowforge check` 和 `flowforge frontier`。CLI 计算下一批无阻塞工作；Agent 不在上下文中自行猜测拓扑顺序。全部 ticket 关闭且诊断得到处置时，一个需求才完成。

## 安装与初始化

Linux / macOS：

```bash
curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash
```

Windows PowerShell：

```powershell
irm https://github.com/coderbiq/flowforge/releases/latest/download/install.ps1 | iex
```

在项目根目录初始化：

```bash
flowforge init
```

默认创建 `.flowforge/config.yaml`、`docs/CONTEXT.md`、`docs/adr/`、`docs/proposals/`，并部署 `.agents/skills/` 与 `docs/agents/`。如需其他文档根目录：

```bash
flowforge config set docs_dir ff-wiki-v5
flowforge init --force
```

`init --force` 同步受管资产，但保留已有项目配置。`check`、`frontier` 和 `status` 从子目录运行时也会向上查找 `.flowforge/config.yaml`。

## 核心文档

- [架构与权威模型](docs/architecture.md)
- [Skill 职责与协作](docs/skill-system.md)
- [CLI 行为与策略](docs/cli-design.md)
- [本次文档契约重构规格](docs/proposals/documentation-contract-refinement/spec.md)

## 开发

```bash
make dev
GOPROXY=https://goproxy.cn,direct go test ./internal/...
```

许可证：MIT。
