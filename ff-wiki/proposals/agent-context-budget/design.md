---
flowforge:
  schema: 1
  role: design
  id: agent-context-budget-design
  revision: 5
  consumes:
    requirements:
      agent-context-budget-requirements: 3
  areas:
    frontier-output:
      revision: 1
      anchor: d1-frontier-schema
    check-output:
      revision: 1
      anchor: d2-check-convergence
    budget-contract:
      revision: 1
      anchor: d3-budget-contract
    scope-selection:
      revision: 1
      anchor: d6-scope-selection
---

<a id="agent-context-budget-design"></a>
# Agent 上下文预算方案（frontier/check 输出收敛）

需求 authority：[Agent 上下文预算需求](requirements.md#agent-context-budget-requirements)，修订版 1。

## 一、决策前沿与已解析事实

### 已决（来自需求）

- 输出有界、Body 不经 CLI、`--full` 保留全量、门禁 exit code 不变、禁改模型配置、不依赖宿主侧修复。

### 已解析事实（代码勘察，含引用）

- `tracker.Issue.Body` 装载每张票全文 Markdown（`internal/tracker/model.go:44`）；`frontier --json` 将 ready / ready_with_warnings / gaps / claimed / **blocked 全量**（含 Body）与**全部 diagnostics** 原样序列化（`internal/command/frontier.go:63`）。235-issue notebook 实测 1,385.4KB。
- `check` 文本模式经 `printDiagnostics` 无界打印每条 diagnostic（`internal/command/check.go:82`）。同一 notebook 实测 402.0KB。
- MCP 封装强制全量：`flowforge_frontier` 固定 `["frontier","--json"]`（`assets/pi/flowforge.ts:261`）；`flowforge_check` 为 check 文本模式。
- **每个执行类技能都会跑这条路径**：implementer 入场 "Run Catalog/frontier diagnostics"（`assets/skills/flowforge-implement/SKILL.md:14`）与收尾（:122）、plan（`flowforge-plan/SKILL.md:139`）、review（`flowforge-review/SKILL.md:162`）→ 预算修复在 CLI 层收敛即覆盖全部技能与宿主路径，技能文本几乎零改动（仅作用域提示一行，见 D6）。
- `frontier` 默认文本模式已是单行/票（`internal/command/frontier.go:97` 起）；膨胀集中在 `--json` 的 Body/blocked/diagnostics。`--quiet`（仅路径）已存在。
- GIIS 实测死亡机制：两工具结果单回合注入约 50 万 tokens，PI 输出预算公式 `floor(本地窗口 − 估算上下文 − 4096, 1)` 把输出压到 1，上游返回 `finish_reason: length`（600k 输入直测 HTTP 200 正常——输入并未超服务端能力）。**消费侧必须适配有限的本地窗口预算，这是硬前提。**
- Body 冗余是纯重复：implementer 派发协议本就要求"读票面文件本身"（`piWorkflowTaskPrefix`，`internal/command/frontier.go:172` 区域），`file_path` 已随输出携带。
- **作用域缺口（rev 2 新增事实）**：GIIS notebook 含 39 个 proposal / 197 张票（中位数 proposal 4 张），当日活跃 proposal 仅 8 张——全局视图下 implementer 吞下的图状态约 96% 与当前工作无关。CLI 已支持目录级作用域（`--dir <proposal 目录>` 实测可用），但 MCP 封装零参数（`assets/pi/flowforge.ts:260,269` `parameters: Type.Object({})`），作用域无暴露入口；技能文本只说 "for this ticket" 却无机制支撑。跨 proposal 依赖用 `feature/ID` 复合键表达（`internal/tracker/model.go:50-51`），作用域视图必须保留跨界阻塞的外部引用。
- **宿主矩阵（rev 3 新增事实，全宿主同病）**：四宿主都部署 agent 定义（claude `.md` / opencode `.md` / codex `.toml` / pi `.md`，`agents.go:28`），仅 PI 额外获得 flowforge.ts 扩展。两个 Bug 在全部宿主成立，只是症状不同：

  | 宿主 | 调用路径 | 输出保护 | GIIS 体量下的症状 |
  |:---|:---|:---|:---|
  | PI | 扩展原生工具（runCli 全量返回） | **无**——原生工具绕过 bash 截断 | 最重：59.5 万 token 直入上下文，输出预算压到 1，`stopReason: length` 死亡 |
  | PI | bash 直跑 CLI | 2000 行 / 50KB | 同 Claude Code |
  | Claude Code | bash 直跑 | 默认 30K 字符（`BASH_MAX_OUTPUT_LENGTH` 可调至 150K） | 不死；JSON 中途截断不可解析，ready 清单被 diagnostics 挤掉，重试每次再浪费 30K |
  | OpenCode | bash 直跑 | ~2000 行 / 50KB | 同上 |
  | Codex | shell 直跑 | 1MiB 捕获，选给模型看的更少（<16KB，需实现期核实） | 同上，更紧 |

  共同点：作用域 Bug 在 CLI 层，四宿主等价受害；所有宿主都遵循"一票一新鲜上下文"，噪音成本每张票重复支付；截断 JSON 触发 fail-fast（至多 2 次重试）后可级联成 `STATUS: BLOCKED`。结论：宿主层修补（如 PI 扩展内截断）救不了其它宿主，CLI 层收敛是唯一全宿主修复点——验证了方案 A 的选择。

## 二、对比设计

| 方案 | 描述 | 评价 |
|:---|:---|:---|
| **A. slim-by-default + `--full`**（选定） | 默认输出去 Body、diagnostics 聚合、blocked 摘要化；`--full` 恢复全量 | 单点治本，同时覆盖 CLI 直跑与 MCP 封装两条路径；需审计存量 `--json` 消费者（见五、验证 3） |
| B. 摘要 + 全量写缓存文件 | 全量 JSON 写 `<docs>/.flowforge/cache/`，stdout 给指针 | 读命令产生写副作用（违反需求约束 3）、缓存失效复杂、污染 git status。弃 |
| C. MCP 层截断 | `flowforge.ts` 截断输出至 N KB | 只修 PI 宿主不修 CLI 直跑；截断破坏 JSON 可解析性；症状修补。弃 |
| D. 查询式子命令 `frontier --issue <id>` | 按票取详情 | detail-on-demand 已由 `context` / `inspect` 承担，无需新增；仅作 A 的注脚 |

弃 B/C 的共同原因：它们把"输出体积是 CLI 的责任"转嫁给消费者或副作用，而仓库第一原则（文件负责内容、CLI 负责图计算）本就要求 CLI 自身输出收敛。

**作用域 seam（rev 2）**：| 方案 | 评价 |
|:---|:---|
| **`--feature` flag + MCP 可选参数**（选定） | 显式、可发现；默认全局不变，编排会话与子代理各取所需 |
| MCP 无参自动猜"当前 proposal" | CLI 无会话状态，只能靠隐式魔法猜测；不可审计 |
| 只改技能文本让 agent 传 `--dir` 路径 | 把责任推给提示词；MCP 路径仍无入口，漏一道 |

## 三、接口决策

### D1 frontier `--json` 新 schema <a id="d1-frontier-schema"></a>

- `ready[]` / `ready_with_warnings[]` / `gaps[]` / `claimed[]`：`{id, feature, file_path, blocked_by[]}`（`TicketRef`——r5 修订：实施期为满足预算省略 `title`，slug 已含于 file_path；`blocked[]` 保留 title 供叙事）。删除 `body`、`mod_time` 等正文性字段。
- `blocked[]`：`{id, feature, title, waiting_on[]}`。
- `diagnostics`：`{counts: {"<severity>:<code>": n}, blockers: [{severity, code, artifact, message(≤160 chars)}], sample_warnings: [{code, artifact}] × ≤5}`——blocker 逐条保留（门禁可诊断性），warning/gap 聚合计数 + 样例。
- 顶层提示字段 `"detail": "pass --full for complete bodies and diagnostics"`。
- `--full`：恢复现状 schema（含 Body 与全量 diagnostics），人类/脚本逃生口。
- 默认文本模式、`--quiet`、`--pi-workflow`、`--strict`/`--include-gaps` 的 ready 计算与 exit code 逻辑不动。

### D2 check 输出收敛 <a id="d2-check-convergence"></a>

- 文本默认：`Checked N issues` + 按 `(severity, code)` 计数行 + blockers 逐条 + warnings/gaps 前 5 条样例 + `--full for all diagnostics` 提示；cycles / dangling / self_blocked 结构性错误原样输出（天然少量）。
- `--json` 默认采用 D1 同款 diagnostics 形态；`--full` 全量。
- `isValid` 计算与 `errPolicyViolation` exit code 语义不变（收敛只动打印层，不动判定层）。

### D3 预算契约与回归测试 <a id="d3-budget-contract"></a>

- 契约常量 `MaxAgentOutputBytes = 16 * 1024`：frontier / check 的 agent-facing 默认输出（frontier 文本、frontier --json、check 文本、check --json，含各自 `--feature` 作用域形态）在任何 notebook 规模下 ≤ 该常量。**实现期必须核对 Codex 展示给模型的上限（宿主矩阵显示其 < 16KB）；若更紧则下调常量，预算以最紧宿主为准**。
- 回归测试：合成 fixture（生成 ≥200 tickets + 多类 diagnostics 的临时 notebook）断言上述通道（含作用域通道）输出字节 ≤ 常量。测试预置于 `internal/command/`，由 Plan 落 ticket（执行者不得改预置测试）。

### D4 MCP 封装与部署

- `flowforge_frontier` / `flowforge_check` 不传 `--full`（自动受益）；新增可选 `feature` 参数（proposal slug，见 D6）；description 更新为有界摘要 + 作用域语义，并提示全量需直跑 CLI。源在 `assets/pi/flowforge.ts`，部署副本（`.pi/extensions/`）经 upgrade 同步，本次 dogfood 手动同步一次。

### D5 兼容策略

- `--full` 为新增 flag；旧 CLI 遇 `--full` 报 unknown flag 属可接受版本耦合（封装不传它，不触发该路径）。
- README 变更说明：Body 消费者改读 `file_path` 文件——与"文件负责内容"原则一致，属纠偏而非破坏。

### D6 作用域选择 <a id="d6-scope-selection"></a>

- CLI：frontier / check 新增 `--feature <slug>`，等价于把扫描范围限定到 `<docs_dir>/proposals/<slug>`（复用现有 `--dir` 机制，slug 优先的 UX）；`--dir` 行为不变。
- 语义：作用域视图内只列本 proposal 的票；本 proposal 票被外部 proposal 的票阻塞时，blocked 条目保留 `waiting_on: [feature/ID]` 外部引用（不拉入外部票正文）。作用域不改变 DAG 计算本身，只是投影过滤。
- MCP：两工具新增可选 `feature` 参数（string，省略 = 全局，供编排会话选派）；description 写明"实施单个 proposal 时传 feature"。
- 技能提示：`assets/skills/flowforge-implement/SKILL.md:14` 入场行补一句 scoped 提示（提示词与机制对齐，避免 agent 继续全局扫）——预算修复不需改技能文本，作用域修复需要这一行。
- 默认全局保持不变：编排会话的选派队列语义不受影响。

## 四、Standards clauses

- [Constraints] "严禁引入通过 CLI 传长文本的接口"（`AGENTS.md` 核心设计原则，CLI/文件分界）→ 本设计是该原则在 frontier/check 上的执行，语义链接：仓库根 `AGENTS.md` 核心设计原则 1。
- [Constraints] "在 `assets/` 中放不部署的内容"禁止（`AGENTS.md` boundaries Never）→ flowforge.ts 改动必须落 `assets/` 并同步部署副本。
- [Conventions] "变更后运行 `go test ./internal/...`"（`AGENTS.md` Commands；GOPROXY 见同文档）。
- [Conventions] 测试分离（`AGENTS.md` Execution unit policy：预置验收测试不由执行者修改）→ 预算回归测试由 Plan 预置。

## 五、验证策略

1. **单元与回归**：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` 全绿，含新增预算回归测试（200-ticket fixture 四通道 ≤ 16KB）。
2. **实测收敛**：对 giis notebook（`/Users/qiangbi/develop/projects/Bytesforce/giis/notebook/200-flowforge/saas-service-ffv5/proposals`，39 proposals / 197 tickets）执行 `frontier --json | wc -c` 与 `check | wc -c`，均 ≤ 16KB；`frontier --feature policy-migration-traceability --json` 仅含 8 张本 proposal 票 + 外部阻塞引用（人工验证入口）。
3. **行为等价**：`--full` 输出与旧版 golden 对比（存量消费者审计：本仓库内 `--json` 消费点为测试与 ui，逐个确认迁移）。
4. **端到端复演**（可选）：128k 窗口模型子代理入场复演 S1 序列，不再触发 `stopReason: "length"`。主验为 1–3。

## 六、后续备注（非机器 open item）

<a id="open-other-commands-audit"></a>- `status` / `inspect` / `agents` / `frontier --quiet` 等通道的输出体积未审计（--quiet 保持全局无界，票 01 明确不动）；预算契约是否泛化、以及超大队列（>~190 ready 条目）的分页/采样，由后续独立 proposal 裁定（不阻塞本设计区域与票据）。
