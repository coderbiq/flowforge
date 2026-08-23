# FlowForge v4

> FlowForge 规范以 [`docs/architecture.md`](docs/architecture.md) 为唯一语义来源；实现状态以 CLI 源码和测试为准。

FlowForge 是面向 AI 协同工程的敏捷工作流中枢。它结合了**多级长期工作记忆（Multi-Tier Working Memory）**、**启发式对话对齐技能（Conversational Agility Skills）**以及**领域活文档持续合流（Living Documentation Synthesis）**，让 Agent 与开发者像资深工程师一样进行需求对齐、极简切片拆解与 TDD 高质量交付。

---

## 核心理念

FlowForge v4 彻底摒弃了“繁杂八股文模板与强门禁状态机”，拥抱以**“人机对话为中心、以工作记忆抗遗忘、以自动化测试为裁判”**的工程范式：

1. **多级长期工作记忆（Multi-Tier Working Memory）**：
   - **Tier 1 (全局记忆 `docs/CONTEXT.md`)**：统一领域语言、全局架构约束与活跃提案索引；
   - **Tier 2 (提案活笔记 `01-workspace/<id>/README.md`)**：长周期（数天）讨论中的事实链、决策与未决问题，支持随时中断并 100% 满血断点续传；
   - **Tier 3 (执行级切片)**：只给 Agent 喂单个切片的极简上下文，杜绝 Context 污染与膨胀。
2. **敏捷工程方法论（Conversational Agility Skills）**：
   - **对话对齐（Grilling）**：在动手写卡片前，通过针对性追问（边界、假设、非目标）对齐真实意图；
   - **曳光弹切片（Tracer Bullets）**：将任务拆解为 15~30 分钟内可打通的端到端极小步；
   - **测试驱动（TDD）**：红-绿-重构循环，代码与自动化测试是唯一的质量裁判；
   - **假设驱动排错与红队审查**：系统性假设驱动 Bug 诊断（Diagnose）与非阻断式对抗性架构审查（Review）。
3. **活文档持续合流（Living Documentation Synthesis）**：
   - 提案交付后，提取通用决策为 ADR，提取领域现状变更作为 Patch 增量合流至 `docs/domains/`，让系统文档随代码演进而常青。

---

## 新手引导

### 1. 安装 FlowForge

#### Linux / macOS

```bash
# 安装最新版本
curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash

# 指定安装目录（实际安装到 <prefix>/bin）
curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash -s -- --prefix "$HOME/.flowforge"
```

#### Windows（PowerShell）

```powershell
irm https://github.com/coderbiq/flowforge/releases/latest/download/install.ps1 | iex
```

### 2. 初始化项目

在项目根目录下执行：

```bash
cd your-project
flowforge init
flowforge project create myproject --wiki-root ff-wiki --src-dir .
```

`flowforge init` 会初始化项目级工作记忆 `docs/CONTEXT.md`、配置目录 `.flowforge/` 并部署原子 Skill 矩阵。

### 3. 开始一个需求（从对齐到交付）

在 Codex / OpenCode 中直接用自然语言与 Agent 对话：

1. **需求对齐与深挖**：
   > "我想在系统里增加一个数据迁移模块，支持历史 ADM 数据导入。"
   - 触发 `flowforge-align`：Agent 会主动追问边界与假设，并创建提案活笔记。
2. **代码探索与事实沉淀**（遇到不确定性时）：
   > "我们先查一下现有数据库 reader 遇到合并单元格是怎么处理的。"
   - 触发 `flowforge-explore`：Agent 探查代码并将文件坐标事实链沉淀至活笔记。
3. **极简切片拆解**：
   > "方案已经明确了，请帮我拆解为可交付的 Tracer Bullet 切片。"
   - 触发 `flowforge-plan`：拆解出 3~5 个端到端切片，每个切片绑定一条自动化测试命令。
4. **TDD 编码交付**：
   > "开始执行 Slice 1。"
   - 触发 `flowforge-implement`：严格遵循红-绿-重构循环，消除认知负荷，运行测试验证。
5. **疑难排错（可选）**：
   > "排查测试偶发失败的原因。"
   - 触发 `flowforge-diagnose`：假设驱动排错，构建最小复现用例，修复根因。
6. **红队审查（可选）**：
   > "审查当前变更是否有架构漂移或隐藏风险。"
   - 触发 `flowforge-review`：非阻断式对抗性审查，提供防御性重构建议。
7. **活文档合流与归档**：
   > "所有切片均已通过测试，请归档本提案并合流系统文档。"
   - 触发 `flowforge-curate`：提取 ADR 和领域 Patch，预览 Diff 确认后合流主干。

---

## SKILL 体系

| SKILL | 阶段 | 职责与触发时机 |
|---|---|---|
| `flowforge-triage` | **分诊 (Triage)** | 任务入口第一道门禁，评估影响半径与不确定性，路由至匹配工作流 |
| `flowforge-align` | **对齐 (Align)** | 深挖需求真实意图，挑战假设，澄清边界，自适应维护 Flat/Hierarchical 活笔记 |
| `flowforge-wayfinder` | **航标 (Wayfinder)** | 应对高迷雾/未知探索，绘制决策图 (`MAP.md`) 并推进决策前沿 |
| `flowforge-explore` | **探索 (Explore)** | 针对具体疑点探查代码库，回填带有 `path:line` 的权威事实链 |
| `flowforge-plan` | **计划 (Plan)** | 多态切片拆解：业务特性拆为 Tracer Bullets，跨模块重构拆为 Expand-Contract 批次 |
| `flowforge-implement` | **实现 (Implement)** | 专注于单切片 TDD 交付（绑定 Seam 接口：编写测试 → 最小实现 → 绿灯重构） |
| `flowforge-diagnose` | **诊断 (Diagnose)** | 假设驱动排错协议（提出假设 → 最小复现 → 状态跃迁追踪 → 根因修复） |
| `flowforge-review` | **审查 (Review)** | 非阻断式红队审查（架构漂移检测 + 认知负荷削减建议 + 并发安全） |
| `flowforge-curate` | **合流 (Curate)** | 提取核心架构决策（ADR），将业务演进 Patch 增量合流至领域活文档 |

---

## 核心设计文档

- [v4 架构设计与演进总览](docs/architecture.md)
- [多级工作记忆体系](docs/memory-system.md)
- [原子 Skill 矩阵与工程方法论](docs/skill-system.md)
- [CLI 设计与轻量辅助工具规范](docs/cli-design.md)
- [知识与活文档系统](docs/knowledge-system.md)
- [实施演进计划与里程碑](docs/implementation-plan.md)

---

## 许可证

MIT License

