---
flowforge:
  schema: 1
  role: design
  id: agent-model-preservation-design
  revision: 2
  consumes:
    requirements:
      agent-model-preservation-requirements: 1
---

<a id="agent-model-preservation-design"></a>
# 部署重编译保留项目已设 model 方案

依据：[部署重编译保留项目已设 model 需求](requirements.md#agent-model-preservation-requirements)。用户裁决的合并方案：升级重编译时提取既有部署文件的 `model:`，新内容缺失时合并写回。

## <a id="d-preserve-merge"></a>d-preserve-merge：合并点与优先级

- 合并点在 `deploySubagents`（agents.go）每宿主写文件前：读取目标路径既有内容，提取 frontmatter `model:` 值；传给 compile 作为 fallback。
- 实现缝：`CompileOptions` 增 `FallbackModel string`；opencode/claude 编译器在 `Model` 为空且 `FallbackModel` 非空时采用 fallback（`model` 字段产出逻辑集中各编译器，frontmatter 序列化不复造）。宿主尾部默认各守现状：opencode 省略字段，claude 输出 profile 默认（`ClaudeModel()`）。
- claude 适配器今日丢弃 `CompileOptions`（agents.go L166-168）；本变更为其补 opts 布线（纯包内缝，不改 CLI 接口签名），优先级与 opencode 一致：`Model`（config）> `FallbackModel`（文件残留）> profile 默认。
- 提取方式：既有文件为 yaml frontmatter（`---` 块），解析出 `model` 键值；解析失败（无 frontmatter/无 model）视为无保留值，不报错。
- 优先级：`Model`（config）> `FallbackModel`（文件残留）> 无字段。config 已设时文件值被有意覆盖，无需提示。
- preserved 提示：采用 fallback 时向 stderr 输出一行（对齐 assets_deploy.go L292 风格）：`  info: preserved local model %q for %s (set agents.models in .flowforge/config.yaml to pin explicitly)`。
- codex 编译器无 model 字段，`FallbackModel` 被忽略（零行为差异）。

## Standards clauses

- must 无网络、无 LLM 调用，纯本地文件操作（源：`AGENTS.md` 核心设计原则 2，[Constraints]）。
- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- must 不改 CLI 接口签名与 Issue Schema 头规范（源：requirements 范围，[Constraints]）。
- must preserved 提示风格对齐 skills 部署既有消息（源：本设计 d-preserve-merge，[Conventions]）。

## 兼容与迁移

- 首次部署无既有文件 → 无 fallback → 行为不变。
- 既有部署已被抹掉的项目（如 tangram-v2 已升 v5.9）：一次性手工把 model 写回部署文件或 config `agents.models`，此后保留机制自动接管。
- config 设值项目：config 胜出语义不变（现行为即 config 产出）。

## 验证策略

- 单元：deploySubagents 三态表驱动（文件有+config 无→保留+提示；文件有+config 有→config；文件无→无）＋首次部署回归＋codex 不变回归。
- e2e：/tmp/opencode/upg-repro-loop.sh 场景 2（手编 model → init --force → 存活）作为回归回路转正为测试或手工演练步骤。
