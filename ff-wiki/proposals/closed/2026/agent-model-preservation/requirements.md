---
flowforge:
  schema: 1
  role: requirement
  id: agent-model-preservation-requirements
  revision: 1
---

<a id="agent-model-preservation-requirements"></a>
# 部署重编译保留项目已设 model 需求

## 问题

部署产物是编译生成物（定义源 + config.yaml），`deploySubagents` 无条件 `os.WriteFile` 覆盖（agents.go L370-378）。项目在已部署文件 `.opencode/agent/*.md` / `.claude/agents/*.md` frontmatter 里设置的 `model:` 在升级（re-exec `init --force` → 重部署）或 `agents deploy` 时被静默抹掉。tangram-v2 v5.7 时代手编文件定 model 的项目升级即失。对照组：skills 部署有漂移保护（"project-customised (preserved)"），agents 是静默覆盖例外。

## 目标

重编译部署既有 subagent 时，若目标文件已存在且其 frontmatter 含 `model:` 而新编译内容不含（config 未设该 profile 的 model），则保留旧值合并进新内容后写回；并在 stderr 提示保留了本地定制（仿 skills 的 preserved 消息风格）。

## 范围与约束

- 优先级规则：config.yaml `agents.models` 是显式通道——config 已设值时以 config 编译结果为准，不做合并（否则删 config 不生效，形成删不掉的旧值）。
- 适用宿主：frontmatter 含 `model:` 字段的 opencode 与 claude；codex 编译产物无 model 字段，不适用。
- 移除路径：项目删除部署文件中的 `model:` 行后，下次部署不再保留（自然回退 host 默认）。
- 不改 Issue Schema、不改 CLI 接口签名；合并发生在 deploySubagents 写文件前。

## 可观察验收

1. 既有部署文件含 `model: X`、config 未设 → `agents deploy` / `init --force` 后文件仍含 `model: X`，其余内容为新版本，stderr 有 preserved 提示。
2. 既有部署文件含 `model: X`、config `agents.models` 已设 `Y` → 部署后为 `model: Y`（config 胜出）。
3. 既有部署文件无 `model:`、config 未设 → 部署后仍无（不凭空引入）。
4. 首次部署（无既有文件）行为不变。
5. codex 宿主部署行为不变（无 model 概念）。

## 术语

- **保留合并（preserve-merge）**：重编译前读取既有部署文件 frontmatter 的 `model:`，在新内容缺该字段时注入旧值的部署行为。
- **显式通道（explicit channel）**：config.yaml `agents.models`——配置在案的值总是优先于文件残留值。
