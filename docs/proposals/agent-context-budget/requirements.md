---
flowforge:
  schema: 1
  role: requirement
  id: agent-context-budget-requirements
  revision: 2
---

<a id="agent-context-budget-requirements"></a>
# Agent 上下文预算需求（frontier/check 输出收敛）

来源：GIIS 项目 PI 会话 "migrated data query"（2026-10-15 ~ 10-16）子代理死亡事件复盘 + 2026-10-16 先行排查会话的输出预算机制实测 + 用户约束裁定（400k contextWindow 为刻意的防长上下文劣化设置，禁改）。

## 一、可观察结果

1. **输出有界**：`flowforge frontier` 与 `flowforge check` 的 agent-facing 默认输出（含 `--json`）存在硬性字节预算（上限数值由 design 裁定，量级 ~16KB），不随 notebook 的 issue 数量线性增长。ready 队列保持完整可执行信息（id / feature / title / file_path / 依赖）。
2. **128k 窗口可入场**：窗口 128k 级模型的子代理按技能规定完成入场 frontier/check 调用序列后，会话上下文增量远低于窗口上限，不再出现 `stopReason: "length"` 入场死亡。
3. **票面不经 CLI**：frontier 输出不含任何 `Issue.Body`（票面全文）；票面内容一律以 `file_path` 引用，由消费者按文件读取（"文件负责内容"）。
4. **细节按需可取**：人类与脚本可用 `--full`（或等价机制）获得与现状等价的全量输出；blocker 级诊断在默认输出中逐条可见（门禁可诊断性不降级）。
5. **行为兼容**：`check --strict` 的 exit code 语义、`frontier --quiet`、`frontier --pi-workflow` 行为不变；MCP 工具 `flowforge_frontier` / `flowforge_check` 签名不变且自动获得有界输出。
6. **作用域视图**：frontier / check 支持按 proposal（feature slug）限定范围。子代理实施某 proposal 的票时，入场诊断只含该 proposal 的票；跨 proposal 阻塞引用保留但只带 id + feature（不带正文）。不带范围参数时默认仍是全局队列，供编排会话选派。

## 二、范围

### In scope

- `internal/command/frontier.go`：`--json` schema 收敛 + `--full` + `--feature` 作用域。
- `internal/command/check.go`：文本与 `--json` 输出收敛 + `--full` + `--feature` 作用域。
- `assets/pi/flowforge.ts`：两个 MCP 工具增加可选 `feature` 参数与 description 更新（有界摘要 + 作用域语义）。
- `assets/skills/flowforge-implement/SKILL.md` 入场行（:14）补作用域提示（一行，见 design D6）。
- 输出预算回归测试（合成 fixture notebook，含作用域通道）。
- README / docs 中 `--json` 与作用域行为变更说明。

### Out of scope

- 宿主模型配置（PI `models.json` 的 `contextWindow` / `maxTokens` 等）——用户裁定为刻意的防劣化约束，禁改。
- PI / CPA 宿主侧任何修复（输出截断、子代理 compaction、stopReason 语义消歧）。
- `status` / `inspect` 等其它命令的输出收敛（独立 open item，见 design）。
- skill 文本的大规模重写——预算修复不需要改技能文本；仅 flowforge-implement 入场行补一句作用域提示（见 In scope）。

## 三、场景

### S1 GIIS 复盘场景（必须收敛）

235-issue notebook（giis saas-service-ffv5）。实测序列：implementer 子代理入场（首轮 input 14,645 tokens），第 6 回合调用 `flowforge_check`（402.0KB）+ `flowforge_frontier`（1,385.4KB），下回合 input 595,107 tokens、output 预算被压至 1、`stopReason: "length"`，会话死亡且零落盘；`context: "fresh"` 重试在第 36 回合同 signature 死亡（595,761 tokens）。目标：同样的入场调用序列后，单工具输出 ≤ 预算上限，会话继续。

### S2 小 notebook 无损

<20 票的仓库（如本仓库 flowforge-v2）默认输出语义信息不缺失：ready 队列完整、blocker 可见，人类可读性不劣化。

### S3 CI 门禁不变

gate 脚本以 `flowforge check --strict` / `check` 的 exit code 判定（`errPolicyViolation`），行为与现在完全一致。

### S4 人类排障

`flowforge frontier --full --json` / `flowforge check --full` 输出与旧版等价（含票面 Body 与全量 diagnostics）。

### S5 作用域噪音场景（必须收敛）

GIIS 实测：notebook 含 39 个 proposal / 197 张票（中位数 proposal 仅 4 张），而当日活跃 proposal（policy-migration-traceability）只有 8 张——implementer 入场被强喂全局视图，**约 96% 的图状态与当前工作无关**。噪音不只是体积：其它 proposal 的 waiting-on 列表会诱导模型误判自己被卡。目标：作用域视图只含本 proposal 的票 + 外部阻塞引用（id + feature）。

## 四、约束

1. **禁改模型窗口配置**（用户裁定：400k 是防长上下文劣化的刻意设置，是硬边界而非待修缺陷）。
2. **不依赖宿主侧修复**：不假设 PI 输出截断、子代理 compaction、MCP 层保护的存在。
3. **CLI 是图计算唯一权威**：不引入 CLI 写缓存文件等副作用，frontier/check 保持只读命令语义。
4. **向后兼容**：`--full` 恢复旧形态；存量 `--json` 消费者（脚本/测试/ui）有明确迁移说明（Body 消费者改读 `file_path` 文件——本就符合"文件负责内容"）。

## 五、术语

- **上下文预算**：agent-facing 单命令默认输出的字节上限契约。
- **收敛输出**：计数聚合 + 文件路径引用 + blocker 逐条 + 样例截断的摘要形态，与全量输出（`--full`）相对。
- **作用域视图**：以单一 proposal（feature slug）为边界的 frontier/check 投影；跨 proposal 依赖以外部引用（id + feature）呈现。
- **agent-facing**：被 LLM agent（主会话或子代理）经 CLI 或 MCP 工具直接消费的输出通道。
