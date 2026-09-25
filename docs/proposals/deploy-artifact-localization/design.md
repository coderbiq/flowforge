---
flowforge:
  schema: 1
  role: design
  id: deploy-artifact-localization-design
  revision: 3
  consumes:
    requirements:
      deploy-artifact-localization-requirements: 1
  areas:
    model-channels:
      revision: 1
      anchor: d-model-channels
    config-validation:
      revision: 1
      anchor: d-config-validation
    pi-model-injection:
      revision: 1
      anchor: d-pi-model-injection
    artifact-localization:
      revision: 1
      anchor: d-artifact-localization
---

<a id="deploy-artifact-localization-design"></a># 部署产物本地化与模型注入方案

依据：[部署产物本地化与模型注入需求](requirements.md#deploy-artifact-localization-requirements)。用户已裁决路线：「模型注入产物 + 产物本地化」（否决宿主用户级外置通道与项目级 `.pi/settings.json` 生成方案）。事实基础：[宿主模型外置能力调查](../../research/2026-09-20-host-model-portability.md）。2026-09-25：四项口味题（models_by_host 双键、层间优先级链、无 --untrack flag、格式级 fail-fast 校验）经用户逐条确认，设计定稿。

## <a id="d-model-channels"></a>d-model-channels：per-host 模型覆盖层与优先级链

**配置形状**：`AgentsConfig` 新增 `ModelHostOverrides map[string]map[string]string`（`yaml:"models_by_host,omitempty"`）。外层键 = 宿主 key（`opencode`/`claude`/`pi`；`codex` 是配置错误，见 d-config-validation）；内层键 = agent 名**或** profile 键（`tool-capable`/`tool-capable-read-only`），值 = 模型字符串：

```yaml
agents:
  models_by_host:
    pi:
      flowforge-investigator: cpa/deepseek-v4.1-flash   # name 键
      tool-capable-read-only: cpa/deepseek-v4.1-flash   # profile 键
    claude:
      flowforge-reviewer: opus
```

**优先级链**（需求"per-host 覆盖 > `agents.models`/`models_by_name` > 未配置"的机械展开；同层内 name 键 > profile 键，与既有全局层语义一致）：

1. `models_by_host.<host>[<agent-name>]`
2. `models_by_host.<host>[<profile-key>]`
3. `models_by_name.<agent-name>`（既有）
4. `models.<profile-key>`（既有）
5. preserve-merge 回填（既有部署文件 `model:`，仅 model 承载宿主）
6. 宿主默认（opencode/pi 省略字段；claude 输出 profile 别名档）

**替代方案（否决）**：
- *全局层 name 键跨层压过 per-host profile 键*（按"name 更具体"排序）——否决：需求明文把整个 per-host 层排在 `models`/`models_by_name` 之上；跨层混排使优先级不可局部推理，配置审查需全局计算。
- *扁平 host 前缀键*（`models: {"pi/tool-capable": ...}`）——否决：模型值本身含 `/`（`provider/model`），键再引入 `/` 语义产生解析歧义；嵌套 map 无歧义且 yaml 表达自然。
- *内层只收 name 键*——否决：全局层已有 profile 批量档位需求（验收 2 的"各自宿主的值"常按档位批量配），per-host 层镜像双键成本极低。

**解析缝**：`resolveCompileOptions(cfg, def)` 增加宿主参数（`resolveCompileOptions(cfg, def, hostKey)`），model 解析移入其中按上表取值；`deploySubagents`（agents.go 写文件循环）与 `computeSubagentStatus`（agents_status.go 期望编译循环）都在 `for _, h := range hosts` 层传 `h.key`，deploy/status 同一解析路径（既有约定 `TestStatusUsesSameCompileOptions` 的延伸）。缝在 command 包内部，CLI 接口签名不变。

## <a id="d-config-validation"></a>d-config-validation：模型值格式级校验（fail-fast）

新校验函数 `validateModelConfig(cfg, defs, enabledHosts)`，在 `deploySubagents` 与 `computeSubagentStatus` 中均于**任何目录创建/产物写入之前**调用（对齐 `validateModelOverrides` 先例），配置错误时两路径报同一错误。规则：

1. **键合法性**（全量校验，不受启用宿主影响——键错误是配置损坏信号）：`models` 键 ∈ `validModelProfileKeys`；`models_by_name` 内层键 ∈ 已发现 agent 名（profile 键属 `agents.models` 或 `agents.models_by_host.<host>`——rev 3 收敛裁决：惰性键是配置损坏信号）；`models_by_host.<host>` 内层键 ∈ 已发现 agent 名 ∪ `validModelProfileKeys`；`models_by_host` 外层键 ∈ 已知宿主集合。未知即错，错误消息沿用 `agents.models_by_name: unknown agent %q` 风格（`agents.models_by_host: unknown host %q` / `unknown key %q`）。
2. **codex 例外**：`models_by_host.codex` 存在即配置错误（codex 无 per-agent model 概念）；全局层值不校验 codex（codex 编译器天然丢弃 model）。
3. **值格式**（仅对**启用的 model 承载宿主**逐值校验；禁用宿主的值不阻塞 deploy）：
   - 通用：非空、去首尾空白后无空白/控制字符（防 YAML 注入与多值粘连）。
   - opencode / pi：必须形如 `provider/model`（含一个 `/`，两侧非空）——OpenCode 模型 ID 语法即 `provider/model-id`；pi 为 BYO-provider，同格式。
   - claude：单一 token 即合法（别名 `sonnet`/`opus`/`haiku`/`fable`、`inherit`、完整模型 ID 皆为合法 token）；不查存在性。
   - 全局层（`models`/`models_by_name`）的值会编译进**每个**启用 model 承载宿主 → 对每个启用宿主各跑一遍该宿主的格式规则；不满足即报错并提示改用 `models_by_host` 分宿主配置。
4. **preserve-merge 回填值不校验**：那是既有部署文件的用户手编值，不是本命令写出的配置（agent-model-preservation 显式通道语义不变）。

**替代方案（否决）**：全局层只按"最宽松宿主"校验或只警告——否决：验收 3 要求"不产出坏配置文件"，宽松校验会把 `model: sonnet` 写进 opencode 产物（opencode 视为非法模型 ID，运行期才爆）。格式错误早失败是 fail-fast 的最小实现；存在性校验被需求明文禁止（纯本地无网络）。

## <a id="d-pi-model-injection"></a>d-pi-model-injection：pi 宿主 model 注入与 preserve-merge 真实化

- `CompilePiWithOptions` frontmatter 增 `Model string \`yaml:"model,omitempty"\``（字段序固定插在 `description` 之后、`thinking` 之前），取值 `resolveModel(opts)`（Model > FallbackModel > 空）；空时省略键——未配置项目产物逐字节不变（回归零）。
- 该修订**覆盖** [PI 宿主集成方案 §一](../pi-host-integration/design.md#pi-host-integration-design) 的"`ModelProfile` → `thinking`：**不写 `model`**"一行（已按需求授权修订该行并将其 design revision 升至 2）：`thinking` 映射维持，`model` 从"永不写"改为"config/回填有值才写"。pi 语义：`model` 省略 = 继承父会话（原行为保留为默认态）。
- **preserve-merge 虚假保留随之消除**：deploy 循环的 preserve-merge（agents.go `opts.Model == ""` 分支）是宿主无关缝，pi 编译器此前忽略 `FallbackModel` 却打印 preserved 提示（真 bug）。`CompilePiWithOptions` 改为消费 `resolveModel(opts)` 后，pi 与 opencode/claude 同链（config > 回填 > 默认），提示变真话。**不做按宿主门控**——codex 由"frontmatter 无 `model` 键 → `deployedModel` 无值"天然短路（票 01 已验证零行为差异），四宿主收敛到同一语义，无需 host 能力表。
- preserved 提示文案同步微调：`(set agents.models in .flowforge/config.yaml to pin explicitly)` → `(set agents.models_by_name/models_by_host in .flowforge/config.yaml to pin explicitly)`（提示真实通道）。

**替代方案（否决）**：按宿主登记"model 承载能力"再门控 preserve-merge——否决：能力差异已由"编译器是否产出 `model` 键"自然编码，再建能力表是双重事实源；且 pi 修复后不存在"承诺保留却丢弃"的宿主。

## <a id="d-artifact-localization"></a>d-artifact-localization：产物本地化与存量迁移

**.gitignore 自动化**：新 helper（init 复用，幂等追加到项目根 `.gitignore`，逐条判存防重复；文件不存在则创建）。条目（需求钉定）：`.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/`、`.pi/extensions/`、`.agents/`、`.flowforge/config.yaml`。追加块带 marker 注释（`# flowforge: managed deploy artifacts (per-machine)`) 便于识别归属；用户已有同名条目视为已满足、不重复写。`upgrade` 的 `syncProjectAssets` 同样调用（升级路径保持守恒）。

**存量迁移（已跟踪受管产物）**：`init`/`upgrade` 在部署前用本地 `git ls-files --error-unmatch <受管路径>` 检测（仅本地 git 查询，无网络）；命中时向 stderr 输出逐条可复制的 `git rm --cached -r <path>` 指引 + 一句原因（"部署产物是每机器文件，不应随 git 旅行"），**绝不自动改动 git 索引**。不提供显式 flag（否决 `init --untrack`：一次性迁移不值得新旗标；自动改索引违反最小惊讶，指引已足够；未来有批量需求再加）。

**AGENTS 模板约定**：`assets/AGENTS.md` 增补一句：部署产物与 `.flowforge/config.yaml` 是每机器文件（init 自动 gitignore）；用户自写且想入库的 agent 用 `git add -f` 例外。属部署内容，符合 `assets/` 边界。

**替代方案（否决）**：`upgrade --untrack` 显式 flag 自动执行 `git rm --cached`——否决（见上）；`.gitignore` 只写目录不写 config 单文件——否决：config 携带本机模型 ID，入库即需求问题 2 的可移植性事故，必须显式忽略单文件（`.flowforge/` 其余内容如 `subagents/` 自定义源仍入库）。

## Standards clauses

- must 纯本地确定性文件操作，无网络、无 LLM 调用：模型校验仅格式级，git 检测仅本地 `git ls-files` 查询，绝不写 git 索引（源：[AGENTS.md](../../../AGENTS.md) 核心设计原则 2，[Constraints]）。
- must 变更后运行 `go test ./internal/...`（源：[AGENTS.md](../../../AGENTS.md) boundaries，[Conventions]）。
- must 不改 CLI 命令签名与 Issue Schema 头规范；`resolveCompileOptions` 加宿主参数属包内缝（源：[AGENTS.md](../../../AGENTS.md) Ask first 边界 + requirements 范围，[Constraints]）。
- must `assets/` 只放部署内容（AGENTS 模板 gitignore 约定属部署内容）（源：[AGENTS.md](../../../AGENTS.md) 🚫 Never，[Constraints]）。
- must config 显式配置总是压过 preserve-merge 回填值（源：[agent-model-preservation 需求](../agent-model-preservation/requirements.md#agent-model-preservation-requirements) 显式通道优先级，[Constraints]）。

## 兼容与迁移

- 未配置任何新键的项目：四宿主产物逐字节不变（pi 仍省略 `model`；claude 仍出 opus/sonnet；codex 不变）——回归零。
- 既有 `models`/`models_by_name` 配置语义不变（新层只在其上叠加）。
- 手编 `.pi/agents/<name>.md` 带 `model:` 的项目：行为从"虚假保留"变为"真实保留"（bug fix，验收 5）。
- [PI 宿主集成方案](../pi-host-integration/design.md) revision 升 2，§一 `model` 行改为条件注入语义。
- 本仓库存量一次性手工清理：对 `.claude/agents/`、`.codex/agents/`、`.pi/agents/`、`.pi/extensions/flowforge.ts` 执行 `git rm --cached`（按 d-artifact-localization 指引人工执行，不入 ticket 自动化）。

## 验证策略

- 单元：优先级链六级表驱动（含 per-host 双键、config 压回填、codex 短路）；校验矩阵（未知 host/未知 key/格式非法/禁用宿主不阻塞/全局值跨宿主不满足即错）；`CompilePiWithOptions` 三态（Model→写、仅 Fallback→写、无→省略）；pi preserve-merge 回归（手编 model → deploy 存活 + stderr 提示，且提示后值真实生效）；gitignore 幂等（重复 init 无重复行、已有条目不重写）；已跟踪检测（fixture git 仓库 → 输出 `git rm --cached` 指引、索引不变）。
- deploy/status 对称：同一 fixture 下 status 期望内容与 deploy 写出内容一致（扩展 `TestStatusUsesSameCompileOptions` 覆盖 per-host 层）。
- e2e 冒烟：临时项目双宿主（pi + claude）仅配 per-host 覆盖 → 两宿主产物各得其值（验收 2）；非法值 deploy/status 均 exit 1 且产物未写（验收 3）；全新 init 后 `.gitignore` 含全部条目且二次 init 无重复行（验收 4）。
