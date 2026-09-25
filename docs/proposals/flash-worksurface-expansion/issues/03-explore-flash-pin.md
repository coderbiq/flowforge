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

## Completion evidence

闭环证据见各 Change 的 cmd/exit/output/artifact 四元组与 Implementation note（执行记录、验证命令与观测结果）。

# 03: explore 迁移 flash（项目 opencode.json 钉扎）

**Blocked by:** 01
**Status:** closed
**Mode:** lightweight

## Delivery

tangram-v2 项目根新建 `opencode.json`，将 opencode 内建 explore agent 的 model 钉扎为 `cpa/deepseek-v4.1-flash`。

## Design context

explore 不是 flowforge 部署的 subagent（不经编译管线），其模型由项目根 `opencode.json` 的 `agent.explore.model` 钉扎。项目级配置（而非用户级 `~/.config/opencode/opencode.jsonc`）控制影响半径在本项目内。当前 explore 会话多为 glm-5.3（继承宿主默认），已有 1 个 deepseek 会话（Trace sessionId usages，0.8M）无异常。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-phase-1 explore 迁移节）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 1，验收 2）.

## Touch points

- `/vol3/1000/develop/tangram-v2/opencode.json` — 新建（项目根，非 `.opencode/` 目录）

## Changes

- [x] 1. 新建 `/vol3/1000/develop/tangram-v2/opencode.json`，内容 `{"agent": {"explore": {"model": "cpa/deepseek-v4.1-flash"}}}`。
  - cmd: `python3 -c "import json; c=json.load(open('/vol3/1000/develop/tangram-v2/opencode.json')); assert c['agent']['explore']['model']=='cpa/deepseek-v4.1-flash'; print('ok')"`
  - exit: 0
  - output: `ok`
  - artifact: docs/proposals/flash-worksurface-expansion/issues/03-explore-flash-pin.md
- [x] 2. 验证 JSON 合法且结构正确（见 Done）。
  - cmd: `python3 -c "import json; c=json.load(open('/vol3/1000/develop/tangram-v2/opencode.json')); assert c['agent']['explore']['model']=='cpa/deepseek-v4.1-flash'; print('ok')"`
  - exit: 0
  - output: `ok`
  - artifact: docs/proposals/flash-worksurface-expansion/issues/03-explore-flash-pin.md

## Constraints

- MUST NOT 修改用户级 `~/.config/opencode/opencode.jsonc`（影响半径控制）。
- MUST NOT 修改 implementer 的模型配置或部署产物（Standards clause）。
- Write set: `/vol3/1000/develop/tangram-v2/opencode.json`

## Done and verify

- JSON 合法且键正确: `python3 -c "import json; c=json.load(open('/vol3/1000/develop/tangram-v2/opencode.json')); assert c['agent']['explore']['model']=='cpa/deepseek-v4.1-flash'; print('ok')"` — 输出 ok。
- flowforge 部署不受干扰: `cd /vol3/1000/develop/tangram-v2 && flowforge agents deploy flowforge-implementer` — 成功且 implementer frontmatter 不变。
- 生效确认为观察后置条件（下一 explore 会话 DB 中 model 为 deepseek，由每日 extract --agents explore 覆盖），本 ticket 不阻塞等待。

---

## Implementation note

Phase 0 restatement：
- Change 1：新建 tangram-v2 项目根 `opencode.json` 钉扎 explore 模型 → 验收：Done 中 python3 JSON 断言输出 `ok`。
- Change 2：验证 JSON 合法且结构正确 → 验收：同一 python3 断言。

结果：
- Completed：Change 1、Change 2（全部）。
- 命令与结果：
  - python3 JSON 断言（Done/Expected tests 原文命令）→ 输出 `ok`，exit 0。
  - `cd /vol3/1000/develop/tangram-v2 && /vol3/1000/develop/flowforge/bin/flowforge agents deploy flowforge-implementer` → `✓ Deployed 1 subagent(s) to .opencode/agent/`，且 `git diff` / `git status` 对 `.opencode/agent/flowforge-implementer.md` 均为空（无 diff，部署不受干扰）；deploy 日志 `info: preserved local model "cpa/deepseek-v4.1-flash" ...` 表明未改动 implementer 模型配置（Standards clause 符合）。
- 修改文件：`/vol3/1000/develop/tangram-v2/opencode.json`（新建，tangram-v2 commit 84e0f86，仅单文件 `git add opencode.json` 提交，未裹挟其他变更）、本 ticket 文件。
- Write-set compliance：All modifications within write set（tangram-v2/opencode.json + 本 ticket；未触碰 `~/.config/opencode/` 用户级配置、`.flowforge/config.yaml` 及任何代码）。
- 生效确认为观察后置条件（下一 explore 会话 DB 中 model 为 deepseek，由每日 extract --agents 覆盖），本 ticket 不阻塞等待（Done 第 3 条）。
- 派发指令明确要求勾选 Changes、写 Implementation note 并置 Status: closed 后提交；关闭由派发方（工作流所有者）决定。

## Execution detail

### Verified contracts

- `/vol3/1000/develop/tangram-v2/opencode.json` — 不存在（项目根 `ls` 核实；`.opencode/` 目录仅有 agent/、node_modules、package.json 等）。
- opencode 1.18.x 项目级配置：项目根 `opencode.json` 的 `agent.<name>.model` 钉扎内建 agent 模型（opencode 配置 schema；本机会话 DB 中 explore 会话多为 glm-5.3 继承宿主默认、1 个 deepseek 会话无异常，佐证模型可变且无兼容问题）。
- 用户级配置 `~/.config/opencode/opencode.jsonc` 存在但不含 agent 节——本 ticket 不触碰它（影响半径控制）。
- 生效验证为观察后置：下一次 explore 派发会话在 DB 中 model 为 deepseek，由每日 `extract --agents flowforge-investigator,explore` 覆盖（通道已就位）。

### Execution scenarios

- Success：新建 opencode.json 三行 → JSON 校验通过 → flowforge deploy 不受干扰 → 下一次 explore 会话跑在 deepseek 上。
- Failure：JSON 语法错误 → opencode 启动时报配置错误（本 ticket 的 Done 命令在写完后立即校验，拦截于此）。

### Expected tests

- `python3 -c "import json; c=json.load(open('/vol3/1000/develop/tangram-v2/opencode.json')); assert c['agent']['explore']['model']=='cpa/deepseek-v4.1-flash'; print('ok')"` — 输出 ok。
- `cd /vol3/1000/develop/tangram-v2 && /vol3/1000/develop/flowforge/bin/flowforge agents deploy flowforge-implementer` — 成功且 `.opencode/agent/flowforge-implementer.md` 无 diff。

### Generated artifacts

- `/vol3/1000/develop/tangram-v2/opencode.json` — producer：本 ticket；consumer：opencode 运行时。同步断言：文件存在 + JSON 结构断言通过。

### Conventions

- MUST NOT 修改用户级 `~/.config/opencode/opencode.jsonc`（Constraints 转录）。
- 文件为项目级新文件，tangram-v2 仓 commit 一次（与 #02 的 config 变更可同批或分批）。
