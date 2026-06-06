## 核心约束

### 1. `src/` 是部署边界

**`src/` 下的所有文件都会被部署到目标项目。** 不部署的内容绝不能放进 `src/`。

- `src/agents/` → 部署到 `.agents/skills/`（如有 `.claude/` 则额外同步）
- `src/flowforge/` → 部署到 `.flowforge/`
- `src/wiki-tpl/` → 部署到项目的 `ff-wiki/` 知识库
- `src/AGENTS.md` → 部署到目标项目根目录的 `AGENTS.md`

反之，**不**部署的内容必须放在 `src/` 之外：开发文档在 `docs/`，构建/安装工具在 `scripts/`，测试在 `tests/`。

**添加任何文件前，先问自己："这个文件会部署到目标项目吗？"** 如果不会，就不该放在 `src/` 里。

### 2. 所有开发围绕 SKILL 最佳实践

本项目以 **SKILL** 为核心产出物。每个设计决策都必须对照 SKILL 编写最佳实践来评估：

- **单一职责**——每个 SKILL 只做一件事，并说明何时使用
- **行动导向**——指南描述做什么，而非参考性叙述
- **渐进式披露**——先路由，再按需展开细节
- **薄适配器模式**——SKILL 委托给脚本和规范定义，而不是把所有内容都内联
- **校验闸门**——每个输出格式都有 schema；每个指南都有校验器
- **机制与策略分离**——核心工作流规则在 `flowforge/rules/` 中，项目级覆盖在已安装的项目中
- **自洽命中**——SKILL 靠自己完备的描述让模型准确识别激活时机，不依赖其他 SKILL 的显式委托

编写或审查任何 SKILL、脚本、指南时，请对照这些原则进行验证。

#### 2.1 Description 编写原则（SKILL 命中率关键）

`description` 是 SKILL 被模型发现的**唯一入口**。按 Anthropic 官方最佳实践：

**✅ 必须做：**
- 描述**何时使用**（触发场景），不描述**内部能力**
- 列出具体触发信号词：用户会说什么、Agent 刚完成什么动作
- 关联到现有 SKILL 的路由终态（"当 flowforge-xxx 路由到 xxx 场景时"）
- 明确与相邻 SKILL 的边界——消歧

**❌ 禁止做：**
- 描述实现细节（"通过脚本加载..."、"负责..."）
- 使用抽象术语而不锚定具体场景
- 缺少反例（何时**不**该触发）
- 与其他 SKILL 的 description 重叠或冲突

#### 2.2 Description 审查清单

新增或修改 SKILL 时，必须回答：

1. 能否在 3 秒内说出"用户说了什么话，这个 SKILL 就该激活"？
2. 如果有两个 SKILL 可能同时匹配，它们的 description 是否明确区分了触发条件？
3. 反例（不应该激活的场景）是否明确列出，且与相邻 SKILL 不重叠？
4. description 是为**模型**（不是人）写的——是否包含模型能识别的语义锚点？

审查时检查所有 `src/agents/**/SKILL.md` 的 frontmatter description，确保它们之间互不冲突。

### 3. Agent 工作流驱动设计

**所有设计必须从模拟 Agent 的工作流程出发，而不是直接编写规则文档。**

正确的设计顺序：

1. **先设计 SKILL**——确定 Agent 的触发入口：什么时候激活、收到什么输入
2. **再模拟工作流**——画出 Agent 从触发到完成的完整决策链路：识别场景 → 读取上下文 → 决策动作 → 执行操作 → 写入产物
3. **在流程节点上添加实现**——每个决策节点才需要什么规则、什么 schema、什么校验器
4. **最后才写文档**——文档是对实现的说明，不是实现本身

禁止的行为：

- 在没有 SKILL 入口设计的情况下，直接编写规则文件或指南文档
- 编写"描述性"的规则（告诉人这个系统怎么运作），而不定义"可执行"的规则（告诉 Agent 在什么条件下做什么操作）
- 先定义 artifact 结构再思考 Agent 如何使用它们

## 项目结构

```
FlowForge/
├── docs/              ← 开发文档（不部署）
├── src/               ← 可部署制品（这里的所有内容都会发布）
│   ├── AGENTS.md      ← 目标项目的 AGENTS.md 模板
│   ├── agents/        ← SKILL 及面向 agent 的定义
│   ├── flowforge/     ← .flowforge/ 配置、schema、规则、模板
│   └── wiki-tpl/      ← 知识库目录结构模板
├── scripts/           ← 构建、安装、校验工具（不部署）
├── tests/             ← 测试套件（不部署）
├── package.json
└── README.md
```

## 开发规则

- 优先编辑已有文件，而非创建新文件
- 先匹配现有模式，再引入新约定
- 修改 schema 或指南后运行校验器
- 绝不抑制类型错误
- 保持 SKILL 文件聚焦——如果 SKILL.md 开始读起来像参考手册，把那些内容提取出去
- 修改或新增 SKILL 的 description 时，必须通过 2.2 审查清单，并检查与所有其他 SKILL 的 description 是否冲突
- **变更完成后必须执行 `npm test`（或 `node tests/run.js`），确保全部检查通过，不引入回归**
- **变更开始前先检查是否需要调整现有测试或新增测试覆盖**——新增 SKILL、修改 CLI 命令参数、改变 context 脚本输出格式、调整 Backend 接口签名时，对应测试必须同步更新

### bd 操作超时问题

**症状**：`bd create/update/close` 等写操作超时 30-60s，`flowforge task` CLI 报 `ETIMEDOUT`。

**原因**：`bd` 每次写操作后触发 `dolt auto-push` 到远端 GitLab（`git+ssh://git@gitlab.bytesforce.com:7001`）。SSH 端口 7001 在外网不可达时，push 挂起直到超时。**本地数据库写入本身是成功的**。

**解决方案**：所有 `bd` 写操作加 `--sandbox` 标志，禁用 auto-sync：

```bash
bd --sandbox update <id> --claim
bd --sandbox close <id> --reason "..."
bd --sandbox create "title" --type task --parent <epic-id> --labels "..."
```

**注意**：`flowforge task` CLI 内部调 `bd` 时未传 `--sandbox`，在网络不通时会超时。遇到此情况改用 `bd --sandbox` 直接操作，操作完成后用 `flowforge task status` 验证（只读命令不受影响）。

会话收尾时运行 `bd dolt push` 手动同步远端——仅在网络可达时执行，不可达则跳过。

<!-- BEGIN FLOWFORGE v:0.12.0 profile:default -->

## FlowForge SKILL 路由

- 新需求、分析、设计、拆分任务 → `flowforge-design`
- 执行任务、继续推进 → `flowforge-implement`
- 归档、沉淀到 library → `flowforge-archive`
- 实施中发现问题、新认知 → `flowforge-feedback`
- 创建/修改 wiki 文档 → `flowforge-docs`

## 任务操作规则

**所有任务操作必须通过 `flowforge task` CLI，严禁直接读写任务文件。**

- ❌ **禁止** 读取 `tasks.snapshot.md` —— 这是自动生成的只读快照，供人类 git diff 审查
- ❌ **禁止** 读取 `task-map.yaml` —— v0.9 已废弃，任务数据在 beads 后端
- ❌ **禁止** 直接用 `bd create/update/close` 操作 proposal 任务
- ✅ **必须** 使用 `flowforge task status/ready/claim/done` 等命令
- ✅ `bd create/update/close` 仅限与任何 proposal 无关的独立事务
- ✅ 知识持久化用 `bd remember`
- ⚠️ `flowforge task` CLI 超时时，可用 `bd --sandbox` 直接操作（见上方 bd 操作超时问题）

### 任务层级

每个 proposal 的任务空间为 4 层结构（详见 `.flowforge/guides/task-hierarchy.md`）：

```
Main Epic → Type Sub-Epic (分析/设计/实施) → Task → Child Task
```

- 大任务通过 `--parent <parentTaskId>` 拆为子任务，最多 4 层
- 独立任务直接挂在类型子 epic 下
- `tasks.snapshot.md` 按类型分组，父子任务缩进展示

任务查询命令：

```bash
flowforge task status --proposal <id>      # 全部任务状态（含 byType 分组）
flowforge task ready --proposal <id>       # 就绪任务列表
flowforge task blocked --proposal <id>     # 阻塞任务列表
```

## CLI 入口

项目根目录 `flowforge` 是统一入口。常用命令：

```bash
flowforge task ready --proposal <CR-id>     # 就绪任务
flowforge task claim --proposal <CR-id> <id> # 认领任务
flowforge task done --proposal <CR-id> <id>  # 完成任务
flowforge task status --proposal <CR-id>     # 状态概览
flowforge implement-context [CR-id]           # 加载实施上下文
flowforge design-context [CR-id]              # 加载设计上下文
flowforge task --help                         # 任务管理帮助
```

---

以下动作后**必须**激活 `flowforge-progress`：

- 修改 proposal 的 `meta.yaml` status
- 通过 `flowforge task` 完成任务操作
- 在 notes.md 中追加日志
- 创建、归档或移动 proposal 目录

### 会话收尾

1. 质量门禁通过（测试、lint、构建）
2. `git pull --rebase && bd dolt push && git push`

<!-- END FLOWFORGE -->
