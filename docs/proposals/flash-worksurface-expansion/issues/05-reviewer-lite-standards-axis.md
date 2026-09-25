---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      flash-worksurface-expansion-requirements: 1
    design:
      flash-worksurface-expansion-design: 1
---

# 05: reviewer-lite 承担 Standards 轴（资产 + 双轴派发契约 + 钉扎）

**Blocked by:** 04
**Status:** open
**Mode:** full

## Delivery

新 subagent 资产 `flowforge-reviewer-lite`（Standards 轴专用）+ flowforge-review SKILL 的双轴拆分派发契约 + AGENTS.md 委派行 + tangram-v2 name 级钉扎，reviewer-lite 部署产物 model 为 `cpa/deepseek-v4.1-flash`。

## Design context

双轴 review 拆为两跳：编排会话先派 reviewer-lite 出 Standards 轴 findings（带引用），再派 flowforge-reviewer（旗舰）读取 findings 专注 spec 轴 + 裁定 + `Fix:` Changes + Review round。reviewer 仍是唯一 closeout 权威。subagent 不能调用 subagent，拆分只能发生在编排会话的派发层。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-phase-2 / d-decision-gates）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 2，验收 4/6）.

## Touch points

- `assets/subagents/flowforge-reviewer-lite.md` — 新建（五段 schema：Identity/Boundaries/Workflow Position/Default Skill/Result Contract；frontmatter：name/description/model_profile: tool-capable/default_skill: flowforge-review/permission: read-only/after·before·returns_to 空列表）
- `assets/AGENTS.md` 与仓根 `AGENTS.md` — `## Generic capability dispatch` 能力表新增行（解耦友好形态：能力键行，不进 Subagent delegation 流程表）
- `internal/command/subagent_source_test.go` — `expectedSubagentNames` 精确全集断言（9→10）
- `internal/command/agents_test.go` — builtin 基数断言（L33/234/644/681/731/853/927 → 10/8/10/10/9/10/10）+ `expectedRoles` 名单
- `internal/subagent/compile_test.go` — `TestParseDirReturnsSixDefinitions` 计数（L26，9→10）+ 排序名单
- `/vol3/1000/develop/tangram-v2/.flowforge/config.yaml` — `agents.models_by_name`（新增 reviewer-lite 条目）

> 2026-09-25 票面修订（解耦友好形态）：不再修订 flowforge-review SKILL——轴定义全部落在角色 Identity；双轴汇总属旗舰 reviewer 既有流程；调度入口走 AGENTS 能力表。

## Changes

- [ ] 1. 新建 `assets/subagents/flowforge-reviewer-lite.md`：description 英文能力表述（"Standards-axis only review: build/test passability, conventions, lint, preset-test presence; produces cited findings for flowforge-reviewer adjudication"）；五段 Body——Identity（轴定义在此：构建/测试通过性、约定、lint、预设测试存在性）、Boundaries（MUST NOT 给 spec 轴结论、MUST NOT 规划 `Fix:` Changes、MUST NOT 修改代码；findings 交 flowforge-reviewer 汇总裁定；每条 finding 带可验证引用）、Workflow Position（无流程位，由编排/review 会话按能力表派发）、Default Skill（flowforge-review 双通道句）、Result Contract（STATUS 首行契约）；frontmatter：`model_profile: tool-capable`、`permission: read-only`（编译器真实收紧：pi 只读工具白名单）、`after/before/returns_to: []`、`detour_skills: []`。
- [ ] 2. `assets/AGENTS.md` 与仓根 `AGENTS.md` 的 `## Generic capability dispatch` 能力表新增行：`Standards-axis findings (build/test/convention conformance, cited) | flowforge-reviewer-lite | flash pinnable`。**不触碰 flowforge-review SKILL**（轴定义在角色 Identity；旗舰 reviewer 读取 lite findings 属其既有 review 流程输入）。
- [ ] 3. 名册 9→10：`subagent_source_test.go` `expectedSubagentNames` 加入 `flowforge-reviewer-lite`；`agents_test.go` 基数断言 9→10（L33/644/681/853/927）、disabled 场景 7→8（L234）、host-selection 8→9（L731）、`expectedRoles` 名单 +1；`compile_test.go` L26 9→10 + 排序名单 +1。以 `go test ./internal/command/ ./internal/subagent/` 失败定位为准补齐所有绑定名册基数的断言。
- [ ] 4. tangram-v2：`.flowforge/config.yaml` 的 `agents.models_by_name` 增 `flowforge-reviewer-lite: cpa/deepseek-v4.1-flash`，执行 `flowforge agents deploy` 部署 lite（跨仓操作，产物不入本仓 git）。
- [ ] 5. 版本发布：`make dev VERSION=<下一补丁版>`（从 `git tag` 递增）并安装到 `~/.local/bin/flowforge`（资产随二进制分发）。

## Constraints

- MUST NOT 将 spec 轴审查结论、`Fix:` Changes 规划、Review round 记录派给 lite（Standards clause）。
- MUST NOT 修改 implementer 的模型配置或部署产物（Standards clause）。
- MUST NOT 手改部署产物 model 字段（Standards clause）。
- lite 的 Boundaries 需含"每条 finding 带可验证引用"（Conventions 转录）。
- preset 测试授权：名册/基数断言扩展（subagent_source_test.go / agents_test.go / compile_test.go）已经用户在规划评审中显式授权（2026-09-25 会话），不得改动断言逻辑本身；测试文件改动用 bash 绕过宿主守卫。
- Write set: `assets/subagents/`、`assets/AGENTS.md`、`AGENTS.md`、`internal/command/subagent_source_test.go`、`internal/command/agents_test.go`、`internal/subagent/compile_test.go`、`/vol3/1000/develop/tangram-v2/.flowforge/config.yaml`、`docs/proposals/flash-worksurface-expansion/`

## Done and verify

- 资产编译通过: `flowforge agents deploy`（tangram-v2）— 成功列出 flowforge-reviewer-lite。
- 部署产物模型正确: `head -6 /vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-reviewer-lite.md` — 含 `model: cpa/deepseek-v4.1-flash` 且 `mode: subagent`。
- 旗舰 reviewer 不受影响: `grep '^model:' /vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-reviewer.md` — 无值或 glm-5.3（未被 name/profile 键触碰），git diff 无变更。
- 契约可发现: `grep -n "reviewer-lite" AGENTS.md .agents/skills/flowforge-review/SKILL.md` — 两处均命中。
- Go 测试全绿: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — ok。

---

## Execution detail

### Verified contracts

- 资产 schema 五段强制（`requiredSubagentSections`：Identity/Boundaries/Workflow Position/Default Skill/Result Contract）；`default_skill: flowforge-review` 解析目标 `assets/skills/flowforge-review/SKILL.md` 存在（名册测试自动校验）。
- `permission: read-only` 是编译器唯一特殊化值：pi 编译产物只读工具白名单（compile_pi.go L46）、codex sandbox、opencode 无 def 级映射——lite 只产出 findings（STATUS 结果体），无需写文件，read-only 成立。
- 名册基数现状 9（generic-role 01 落地后）：`agents_test.go` 七处计数（L33/234/644/681/731/853/927）、`compile_test.go` L26、`expectedSubagentNames` 全集比对；+1 后以测试失败定位补齐。
- AGENTS 能力表现状：`assets/AGENTS.md` L21 起 `## Generic capability dispatch`，5 行能力键表；仓根 AGENTS.md L50 同段；锚点测试 `TestAgentRulesDescribeGenericCapabilityDispatch` 为 Contains 断言，加行零破坏。
- `models_by_name` 形状：`AgentsConfig.ModelOverrides map[string]string`（yaml `agents.models_by_name`，config.go L45）；tangram-v2 已有该键使用先例（investigator 等钉扎）。
- 版本注入式：`internal/version/version.go` `resolve(injected)`；发布经 `make dev VERSION=x.y.z` ldflags 注入。
- flowforge-review SKILL 已有 "Standards sub-agent prompt" 描述（L82 起）——lite 资产是该子代理角色的正式化，SKILL 文本零改动（解耦友好裁决）。

### Execution scenarios

- Success：名册 10 全绿；双 AGENTS 能力表含 lite 行；tangram-v2 deploy 后 lite 产物 frontmatter `model: cpa/deepseek-v4.1-flash`（钉扎生效）。
- Failure：若 Boundaries 缺引用要求或 Identity 混入 spec 轴语义，违背票面 Constraints（人工核验项）；名册断言漏改则对应测试报 found/want 失配。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ ./internal/subagent/` — ok（名册 10 全绿）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。
- `grep -n 'reviewer-lite' assets/AGENTS.md AGENTS.md` — 各命中能力表行。

### Generated artifacts

- producer `assets/subagents/flowforge-reviewer-lite.md` → consumer 四宿主 deploy 产物（tangram-v2 侧带 model 钉扎）；producer 双 AGENTS.md → consumer `<docs_dir>/agents/issue-tracker.md` 部署通道。

### Conventions

- must 变更后运行 `go test ./internal/...`。
- 资产英文书写；能力表行含档位列；测试文件改动用 bash（宿主守卫）。
- 部署镜像 `internal/command/assets/` 刷新属 make-dev 构建步骤（gitignored，不入库）。
