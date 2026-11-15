---
flowforge:
  schema: 1
  role: design
  id: roadmap-management-design
  revision: 1
  consumes:
    requirements:
      roadmap-management-requirements: 1
---
<a id="roadmap-management-design"></a>
# Roadmap 需求树管理：设计

消费[需求 revision 1](requirements.md#roadmap-management-requirements)。本设计新增意图层（roadmap 方向文件）与 proposal 生命周期命令，不触碰 proposal 内部 artifact 契约与 catalog/DAG/frontier 语义。

<a id="physical-layout"></a>
## D1 物理布局

```
<docs_dir>/
├── roadmap/
│   ├── <direction>.md            # 活跃方向，schema v1 role: map
│   └── closed/<direction>.md     # 完结/放弃方向归档（不按年）
└── proposals/
    ├── <slug>/                   # 工作 proposal（现状不动，无 active/）
    └── closed/<YYYY>/<slug>/     # 归档，按年
```

- `featureForPath` 取叶子目录名，`closed/2026/<slug>` 迁移后 feature 名不变，跨 proposal `Blocked by`/`consumes`/exempt 配置零感知。
- `roadmap/` 不在 catalog 扫描根（proposals/）之下，天然无执行语义；归档目录 `roadmap/closed/` 亦然。
- 方向文件 frontmatter：`role: map`、`id: <direction>-map`、`revision`（树结构变化时 +1：增删节点/换父；文字润色不动）。

<a id="node-format"></a>
## D2 节点格式与解析规则

节点 = 树列表项，三部分：

```markdown
- **<标题>**: <一句话概要>
  <可选描述段：意图/背景/雾/孵化条件/兄弟约束；不记验收行为与接口细节>
  → [<slug>](../proposals/<path>/README.md)
```

解析规则（status/reconcile 提取用，刻意最窄语法）：

1. 链接行：节点项下以 `→` 起始的行，恰含一个 markdown 链接 `→ [<slug>](<relpath>)`；relpath 解析出目标 proposal 目录。
2. 节点身份：标题 + 列表位置；无 ID。
3. 一节点至多一链接；多链是退化态，reconcile 报告。
4. 链接固定相对方向文件；close/reopen 命令重写。

<a id="derived-state"></a>
## D3 派生状态（只算不存）

节点状态由链接目标 proposal 目录内文件与票状态计算：

| 条件（按序判定的第一命中） | 状态 |
|:---|:---|
| 无链接 | 概要 |
| 目标目录不存在 | 死链（reconcile 报告） |
| 无 requirements.md | 已孵化·分析中 |
| 有 requirements.md 无 design.md | 分析完成 |
| 有 design.md，票未全终态 | 实施中 `n/m` |
| 票全终态（closed/resolved/wontfix）且 completion evidence 齐备 | 已交付 |

方向汇聚 = 各节点派生状态列表；未孵化节点不进"完成度"分母（聚合管事实，树干管意图）。全部计算在 `status` 内联完成，无持久化。

<a id="commands"></a>
## D4 命令

### `flowforge proposal close <slug>`
1. 定位 live 目录 `proposals/<slug>/`（不存在则报错）。
2. 校验：票全终态才可 close（非终态票存在时拒绝并列出）。
3. `mv` 到 `proposals/closed/<当年年份>/`（目录按需创建）。
4. **重写全库引用**：扫描 `<docs_dir>` 下所有 `.md`（roadmap/ + proposals/ 各 proposal 正文），将指向旧路径的相对链接改写为新路径（按文件相对位置重算 `../` 前缀）。
5. 输出迁移与重写计数。

### `flowforge proposal reopen <slug>`
1. 在 `proposals/closed/*/<slug>/` 中查找唯一匹配（多处匹配报错列出）。
2. live 侧同名已存在则拒绝。
3. `mv` 回 `proposals/<slug>/` + 全库引用重写（同 close 第 4 步）。

### status 扩展
- 检测 `<docs_dir>/roadmap/` 存在时：先输出各方向文件树状汇聚（节点 + 派生状态），再输出原 feature 平铺视图；live 与 closed 计数分开。
- reconcile 检查内联在 status：死链、孤儿 proposal（live 且无任何方向文件链接，仅报告）、多链节点。

### check 扩展
- 新诊断 `roadmap-dead-link`（warning）：方向文件链接指向不存在的 proposal 目录。

<a id="skill"></a>
## D5 flowforge-roadmap skill

新建 `assets/skills/flowforge-roadmap/SKILL.md`（源）→ `flowforge init` 部署到 `.agents/skills/`。四动作：

1. **chart**：松散大想法 → 方向文件。轻量 grilling（≤3 问）定意图+边界，写初始树。反过度建设：一个 proposal 装得下的想法拒绝建图，路由 align。
2. **grow**：加/拆/换父节点、维护概要与描述段。守 D2 格式规则。
3. **hatch**：前置校验（概要+附注像样、用户确认范围）→ 扫描兄弟节点与已链 proposal 附注收集跨节点约束为孵化备注 → 建 `proposals/<slug>/README.md` 存根（标题、一句话、回链方向文件+节点标题、一句"requirement authority 由 align 建立"）→ 节点写链接行 → 交接 flowforge-align（输入=概要+孵化备注+存根位置）。不写 requirements.md。
4. **reconcile**：跑 status 视图，报告死链/孤儿/多链，修复建议路由 chart/grow/hatch/close。

边界（skill 内明写）：不澄清需求、不做设计、不切票、不写节点状态、不建节点依赖。

<a id="migration"></a>
## D6 存量迁移

本仓 25 个已 resolved proposal 迁入 `ff-wiki/proposals/closed/2026/`（验证迁移语义；`documentation-contract-refinement` 等证据豁免名单按 feature 名继续生效）。新建 `roadmap/flowforge-self.md` 方向文件收纳后续本仓工作，节点「roadmap 需求树管理」链接本 proposal——dogfooding。

<a id="standards"></a>
## Standards clauses

- must 构建走 `make dev` 或 `go build -trimpath`，变更后 `go test ./internal/...` 全绿 — AGENTS.md Commands
- must not 引入通过 CLI 传长文本的接口；文件负责内容，CLI 负责计算 — AGENTS.md 核心设计原则
- must not 在 `assets/` 放不部署的内容 — AGENTS.md boundaries
