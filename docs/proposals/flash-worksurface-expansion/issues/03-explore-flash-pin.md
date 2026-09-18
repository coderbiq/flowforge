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

# 03: explore 迁移 flash（项目 opencode.json 钉扎）

**Blocked by:** 01
**Status:** open
**Mode:** lightweight

## Delivery

tangram-v2 项目根新建 `opencode.json`，将 opencode 内建 explore agent 的 model 钉扎为 `cpa/deepseek-v4.1-flash`。

## Design context

explore 不是 flowforge 部署的 subagent（不经编译管线），其模型由项目根 `opencode.json` 的 `agent.explore.model` 钉扎。项目级配置（而非用户级 `~/.config/opencode/opencode.jsonc`）控制影响半径在本项目内。当前 explore 会话多为 glm-5.3（继承宿主默认），已有 1 个 deepseek 会话（Trace sessionId usages，0.8M）无异常。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-phase-1 explore 迁移节）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 1，验收 2）.

## Touch points

- `/vol3/1000/develop/tangram-v2/opencode.json` — 新建（项目根，非 `.opencode/` 目录）

## Changes

- [ ] 1. 新建 `/vol3/1000/develop/tangram-v2/opencode.json`，内容 `{"agent": {"explore": {"model": "cpa/deepseek-v4.1-flash"}}}`。
- [ ] 2. 验证 JSON 合法且结构正确（见 Done）。

## Constraints

- MUST NOT 修改用户级 `~/.config/opencode/opencode.jsonc`（影响半径控制）。
- MUST NOT 修改 implementer 的模型配置或部署产物（Standards clause）。
- Write set: `/vol3/1000/develop/tangram-v2/opencode.json`

## Done and verify

- JSON 合法且键正确: `python3 -c "import json; c=json.load(open('/vol3/1000/develop/tangram-v2/opencode.json')); assert c['agent']['explore']['model']=='cpa/deepseek-v4.1-flash'; print('ok')"` — 输出 ok。
- flowforge 部署不受干扰: `cd /vol3/1000/develop/tangram-v2 && flowforge agents deploy flowforge-implementer` — 成功且 implementer frontmatter 不变。
- 生效确认为观察后置条件（下一 explore 会话 DB 中 model 为 deepseek，由每日 extract --agents explore 覆盖），本 ticket 不阻塞等待。

---

## Execution detail

### Verified contracts

- <filled by flowforge-refine-ticket>

### Execution scenarios

- <filled by flowforge-refine-ticket>

### Expected tests

- <filled by flowforge-refine-ticket>

### Generated artifacts

- <filled by flowforge-refine-ticket>

### Conventions

- <filled by flowforge-refine-ticket>
