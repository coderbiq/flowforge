---
flowforge:
  schema: 1
  role: design
  id: flash-worksurface-expansion-design
  revision: 1
  consumes:
    requirements:
      flash-worksurface-expansion-requirements: 1
  areas:
    mechanism:
      revision: 1
      anchor: d-mechanism
    phase-1:
      revision: 1
      anchor: d-phase-1
    phase-2:
      revision: 1
      anchor: d-phase-2
    phase-3:
      revision: 1
      anchor: d-phase-3
    decision-gates:
      revision: 1
      anchor: d-decision-gates
---

<a id="flash-worksurface-expansion-design"></a>
# Flash 工作面扩大：方案

需求 authority：[Flash 工作面扩大需求](requirements.md#flash-worksurface-expansion-requirements)，修订版 1。

## d-mechanism：模型钉扎机制 <a id="d-mechanism"></a>

现状（已核实代码）：

- `AgentsConfig.Models map[string]string` 按 **ModelProfile** 键钉扎（合法键：`tool-capable`、`tool-capable-read-only`）。
- `resolveCompileOptions`（internal/command/agents.go:251）取 `cfg.Agents.Models[string(def.ModelProfile)]`。
- 模型优先级（internal/subagent/compile_opencode.go:27 `resolveModel`）：config 钉扎 > 既有部署产物的 `model:`（preserve-merge 回填）> 空（继承宿主默认）。
- 各 agent 的 profile：analyst/architect/reviewer = `high-capability`；implementer/planner = `tool-capable`；investigator = `tool-capable-read-only`。

扩展（P2 引入）：`agents.models_by_name: {<agent-name>: <model>}` **按 agent 名**钉扎，优先级高于 profile 键；键名必须命中已知的 subagent 定义名，未知名是 config 错误（与 profile 键校验一致）。理由：reviewer-lite 与 planner 同为 `tool-capable`，profile 键无法单独移动 lite；name 级机制也是后续任意单角色扩面的通用开关。

explore 不是 flowforge 部署的 subagent（opencode 内建 agent），其模型由项目根 `opencode.json` 的 `agent.<name>.model` 钉扎，不经 flowforge 编译管线。

## d-phase-1：investigator + explore 迁移 flash <a id="d-phase-1"></a>

**度量先行**：`executor_metrics.py` 增加 `--agents`（逗号分隔，默认 `flowforge-implementer`）：

- extract：SQL 过滤改为 `agent IN (...)`；observations 表在 `session` 列后增加 `agent` 列；幂等键仍为 session id。
- 本 proposal 使用独立观测文件 `docs/proposals/flash-worksurface-expansion/observations.md`（含 agent 列的新表头），**不迁移** executor-value-measurement 的既有文件（其表无 agent 列且观察窗运行中）。
- report：按 agent 分组聚合；agent 列缺失的旧行按 `flowforge-implementer` 兼容处理。

**investigator 迁移**：tangram-v2 `.flowforge/config.yaml` 增加：

```yaml
agents:
  models:
    tool-capable-read-only: cpa/deepseek-v4.1-flash
```

redeploy 后部署产物 `model: cpa/deepseek-v4.1-flash`（config 钉扎优先于 preserve-merge 回填的 glm-5.3）。该 profile 仅 investigator 一个 agent，无附带影响。

**explore 迁移**：tangram-v2 项目根新建 `opencode.json`：

```json
{ "agent": { "explore": { "model": "cpa/deepseek-v4.1-flash" } } }
```

**不污染运行中的观察窗**：P1 不触碰 implementer 的任何配置；executor-value-measurement 的 implementer 纪元不受影响。P1 部署时间戳即新观测纪元边界（`flagship-baseline:<ts` / `flash-migrated:>ts`，report 时 CLI 传入）。

## d-phase-2：reviewer-lite 承担 Standards 轴 <a id="d-phase-2"></a>

前提机制：`agents.models_by_name`（d-mechanism 扩展）。

**新资产** `assets/subagents/flowforge-reviewer-lite.md`：

- profile `tool-capable`，default_skill `flowforge-review`。
- Identity：只做 Standards 轴（规范符合性：构建/测试通过性、约定、lint、预设测试存在性），产出带引用的 findings 清单。
- Boundaries：MUST NOT 给 spec 轴结论、MUST NOT 计划 `Fix:` Changes（那是 reviewer 的职责）、MUST NOT 修改代码。

**派发契约**（flowforge-review SKILL 修订）：双轴 review 的派发拆为两跳——

1. 编排会话先派 reviewer-lite：Standards 轴 findings（带引用）。
2. 再派 flowforge-reviewer（旗舰）：读取 lite 的 findings 作为输入，专注 spec 轴 + 裁定 + `Fix:` Changes 规划 + Review round 记录。

reviewer 仍是唯一 closeout 权威；lite 的 findings 只是它的输入之一。

**tangram-v2 钉扎**：`agents.models_by_name.flowforge-reviewer-lite: cpa/deepseek-v4.1-flash` + redeploy。

## d-phase-3：设计事实简报（architect 预调研拆分）<a id="d-phase-3"></a>

investigator（P1 后已在 flash）承担 architect 裁定前的事实收集：

**简报契约**（flowforge-research SKILL 增补输出格式）——设计事实简报结构：

- 问题：本次裁定要回答的问题（一句话）
- 事实：每条带可验证引用（文件路径:行号 / javap / 测试输出摘录）
- 约束面：接口签名、现有数据模型、兼容性边界（各带引用）
- 开放项：简报无法回答、必须由 architect 判断的问题清单

**消费步骤**（flowforge-solution-design SKILL 增补）：architect 会话的入口 reading list 允许/鼓励包含一份设计事实简报；architect 基于简报 + 必要抽查做裁定，避免整段读代码库。

**委派行**（AGENTS.md）：设计裁定前的事实收集 → investigator（flash），产出设计事实简报；裁定本身留 architect（旗舰）。

## d-decision-gates：阶段观察门与回滚 <a id="d-decision-gates"></a>

每阶段部署后进入该角色的 flash 纪元；每日 extract/report 覆盖新观测文件。门在部署前预注册如下：

| 门 | 判定 | 后果 |
|---|---|---|
| G1 角色失控（per role） | 该角色 flash 纪元内 `steps>150 或 rep_max>10 或 dur>30min` 的会话 ≥2 | 回滚该角色钉扎（移除 config 条目 + redeploy），开 Fix ticket |
| G1 单次越线 | 1 个会话越线但会话 COMPLETED | 记录观察，不回滚；连续两天出现则回滚 |
| G2 引用质量（investigator/reviewer-lite） | 抽查 1 份产物的引用不可验证 | 该产物重做；同一角色一周内 2 次 → 回滚 |
| G3 实施纪元污染 | executor-value-measurement 的 implementer 纪元出现 config 源污染 | 立即修正 config，重新划定纪元边界 |

阈值沿用 implementer 的 G1 形状（首轮观察后可按角色实际分布调整，调整需修订本 design）。

阶段推进是人工决策：P1 部署 + 观察 ≥2 天且无回滚 → 派发 P2；P2 观察 reviewer-lite ≥2 个 proposal 的 review 无回滚 → 派发 P3。

## Standards clauses

- [Constraints] MUST NOT 在 P1/P2/P3 任何变更中修改 implementer 的模型配置或部署产物。
- [Constraints] MUST NOT 将 spec 轴审查结论、设计裁定、requirements 判断、ticket 切片派发给 flash 角色。
- [Constraints] MUST 在每阶段部署前保持该阶段回滚门已写入本文件 d-decision-gates。
- [Constraints] MUST NOT 手改 tangram-v2 部署产物（`.opencode/agent/*.md`）的 model 字段——模型变更只经 config 钉扎 + redeploy。
- [Conventions] flash 角色的调研/审查产物 SHOULD 每条事实带可验证引用（路径:行号或命令输出）。
- [Conventions] flowforge 仓的 Go 变更 SHOULD 附带 `internal/command` 或对应包的单测，命令 `GOPROXY=https://goproxy.cn,direct go test ./internal/...` 全绿。
