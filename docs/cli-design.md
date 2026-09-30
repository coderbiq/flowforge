# FlowForge CLI 行为与策略

CLI 是 Markdown 工作流的确定性投影器，不是内容管理 API。所有需求、设计、ticket 和 evidence 都由文件工具直接编辑。

## 目录解析

`.flowforge/config.yaml` 的 `docs_dir` 指定文档根目录，默认为 `ff-wiki`，支持相对项目根目录和绝对路径。`check`、`frontier`、`status` 从当前目录向上定位项目配置；显式 `--dir` 则直接使用调用者给出的 proposals 路径。损坏的项目配置必须报错；没有 FlowForge 配置的普通目录回退到 `<startDir>/ff-wiki/proposals` 作为兼容默认路径。

`flowforge init [path]`（别名 `sync`，另接受 `-f`/`--force`）按顺序执行：

1. 创建 `.flowforge/config.yaml`，仅在文件缺失时写入 `version: 5.0.0`、`version_check: true`、`docs_dir: ff-wiki`；
2. 创建 `<docs_dir>/proposals/`、`<docs_dir>/adr/`，以及 `<docs_dir>/CONTEXT.md`（后两者仅在缺失时写入）；
3. 把受管 per-machine 部署产物（`.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/`、`.pi/extensions/`、`.agents/`、`.flowforge/config.yaml`）幂等写入项目根 `.gitignore`，不改 git index；已被 git 跟踪的受管路径只打印提示，不自动处理；
4. 部署受管资产（`.agents/skills/`、`<docs_dir>/agents/`、`AGENTS.md` 受管区块），随后立即复核，仍有缺失或漂移即报错；
5. 部署 subagent 到宿主目录。

已有 `.flowforge/config.yaml`、项目自定义文件与手工改写的模型配置始终保留；`upgrade` 内部以 `init <projectRoot> --force` 重新部署受管资产，`--force` 不改变上述保留语义（受管资产按幂等规则刷新）。

## Artifact Catalog

Catalog 扫描 proposals 下 Markdown，区分 requirement、design、spec、ticket、evidence、research 和 map。只有 `issues/*.md` 中的 ticket 能进入 DAG；旧 v5 ticket 没有 schema envelope 时仍兼容并产生 legacy warning。

确定性诊断覆盖：

- role 与物理位置冲突、无效或未来 schema；
- 重复 identity、缺失 authority、消费 revision 过旧或超前；
- 机器依赖缺少人类语义链接，或链接未记录 consumption；
- open item 缺少精确 scope/anchor，waiver 过宽、无效或过期；
- closed ticket 缺少非空 `Completion evidence`。

## `flowforge check`

```bash
flowforge check [--dir <path>] [--json] [--strict]
```

扫描 proposal 工件、构建 ticket 依赖 DAG，共检测 7 类问题：

1. 循环依赖（死锁）；
2. 悬空引用（blocked by 指向不存在的 ticket）；
3. 自依赖；
4. Catalog 诊断：工件元数据、authority、语义链接、waiver、scoped gap 与完成证据；
5. checked-change 证据四元组诊断（缺失、不完整、非零退出、artifact 缺失）；
6. 重复失败诊断（同一命令在单张票内报非零退出 3 次及以上）；
7. blocked evidence 诊断（open ticket 仍带 `## Blocked evidence` 小节，等待 refine-ticket 消费）。

默认 warning 与 gap 保持可见但不让 check 失败，blocker 始终失败；`--strict` 让未豁免 warning/gap 也失败。legacy proposal 可通过 `.flowforge/config.yaml` 的 `evidence.exempt_proposals` 退出四元组校验（`upgrade` 同步时会提示）。JSON 输出包含 issue/artifact 数量、DAG 结果和完整 diagnostics。

## `flowforge frontier`

```bash
flowforge frontier [--dir <path>] [--json] [--quiet] [--strict] [--include-gaps] [--pi-workflow]
```

先计算 DAG 的无阻塞 open ticket，再按内容诊断投影：

- clean ticket 默认输出；
- warning ticket 默认输出并保留诊断；
- gap ticket 默认排除，`--include-gaps` 可显式包含；
- blocker ticket 始终排除；
- `--strict` 只输出 clean ticket，并优先于 `--include-gaps`。

`--quiet` 只把可执行路径写到 stdout，诊断写到 stderr；`--json` 保留 clean、warning、gap、claimed、blocked 分组，适合 Agent 或自动化消费；`--pi-workflow` 把当前 ready 批次渲染为 pi-subagents workflowScript（顺序 fail-fast，每票一个 fresh 上下文的 `flowforge-implementer` 派发并携带 `flowforge check` gate），优先于 `--json`/`--quiet`，供已安装 pi-subagents 的 PI 会话消费。

## 其他命令

- `flowforge status [--dir <path>]`：按 feature 汇总 ticket 生命周期和 DAG 状态。
- `flowforge init [path]`（别名 `sync`）：初始化或同步本地 tracker、wiki 结构与受管资产（逐条副作用见「目录解析」）。
- `flowforge agents deploy [name]`：将权威 subagent 定义编译并部署到宿主原生目录（`.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/`）。
- `flowforge agents remove <name>`：移除 subagent 部署文件；内置角色写入 `.flowforge/config.yaml` 的 `agents.disabled` 持久化停用，自定义角色删除源文件。
- `flowforge agents status [--json]`：报告各 subagent 在各宿主目录中的状态（`current`、`missing`、`drifted`、`project-owned`）。
- `flowforge assets verify [project] [--json]`：比较运行中二进制内嵌的 Skills 与 agent 规则和所选项目，不修改文件（细节见下文「Managed asset verification」）。
- `flowforge config get <key>` / `flowforge config set <key> <value> [--dry-run]` / `flowforge config list`：读取、修改或列出配置；键名与取值范围见下文「配置键」，`--dry-run` 只预览不落盘。
- `flowforge model-set list|show [name]|use <name>`：管理命名 agent 模型方案并切换。`list` 列出已声明方案并标出激活项；`show` 打印某方案的逐宿主/角色的生效 model、effort 表（缺省为激活方案，`default` 显示基础层；两个独立 source 列标注 base/set、层级和完整配置路径，不读取本地产物充当配置来源）；`use <name>` 切换激活方案（也可用 `default` 回到基础层）并重新部署 subagent，重新部署失败则激活指针回滚到原值并报错；配置校验、读取和编译失败发生在产物写入前，中途写入 I/O 失败仍可能留下部分更新的产物。方案在 `.flowforge/config.yaml` 的 `agents.model_sets` 下声明，激活指针存于 `.flowforge/model-set.active`，缺失即基础层。
- `flowforge completion bash|fish|powershell|zsh`：为指定 shell 生成自动补全脚本。
- `flowforge upgrade [--dry-run] [--version <v>]`：更新 CLI；`--dry-run` 只显示可用更新而不安装，`--version <v>` 升级到指定版本；已经是相同版本时仍同步当前项目的受管资产与 subagent，降级返回独立错误。
- `flowforge version`：显示构建注入版本。

### 配置键

`flowforge config` 可读写的键：

- `docs_dir`（别名 `docsDir`；空值被拒绝）；
- `version_check`（`true`/`false`，也接受 `1`/`0`/`yes`/`no`；键名是下划线形式，不是 `version check`）；
- `standards.guide`（空值被拒绝）；
- `project.<id>.srcDirs`。

仅由配置文件消费、`flowforge config` 不暴露的字段：`version`、`projects[]`（`id`、`srcDirs`）、`knowledge_sources`、`evidence.exempt_proposals`，以及 `agents.disabled`、`agents.hosts`、`agents.max_steps`、`agents.models`、`agents.models_by_name`、`agents.models_by_host`、`agents.model_sets`、`agents.test_file_globs`、`agents.disable_test_guard`。

模型与推理强度配置依据[兼容配置设计 revision 3](proposals/subagent-model-reasoning/design.md#配置表达)。值兼容原字符串及 `{model: ..., reasoning_effort: ...}` 对象，两个字段均可独立声明；对象至少含一项，仅支持这两键，字段须为非空字符串。统一层 `models[profile]` / `models_by_name[name]` 可供多个宿主使用，`models_by_host` 是稀疏覆盖；只声明 Codex 即可，PI/OpenCode 无需重复配置：

```yaml
agents:
  hosts: [codex, pi, opencode]
  models_by_name:
    flowforge-investigator: {model: cpa/deepseek-v4.1-flash, reasoning_effort: high}
  models_by_host:
    codex:
      flowforge-investigator: gpt-6-luna
  model_sets:
    quick:
      models_by_host:
        codex:
          flowforge-investigator: {reasoning_effort: low}
```

模型与 effort 分别按 `by_host[name] > by_host[profile] > by_name[name] > models[profile] > 本地保留值 > 宿主默认` 解析。方案只覆盖已声明字段：换 model 保留原 effort，effort-only 方案不改 model。`reasoning_effort: inherit` 终止查找，屏蔽低层、本地值和 profile 默认，省略原生字段并交宿主继承。

原生字段为 Codex TOML `model` / `model_reasoning_effort`、Claude YAML `model` / `effort`、OpenCode YAML `model` / `reasoningEffort`、PI YAML `model` / `thinking`。Claude effort 接受 low/medium/high/xhigh/max；PI 接受 off/minimal/low/medium/high/xhigh/max；Codex/OpenCode 做 token 检查，provider/client/model 能力由宿主运行时校验。OpenCode/PI 的最终 model 仍须 `provider/model`，Claude/Codex 接受单 token。所有声明检查结构、键和 token，宿主规则只检查启用宿主、启用角色的最终使用字段，覆盖后的统一值不会被该宿主误拒绝。

`agents deploy` 与 `agents status` 应用同一 active-set、按字段保留 YAML frontmatter/TOML 顶层的本地 model/effort；显式配置优先。无显式 effort 且无本地字段时，Codex/PI 保留 profile 的 high/medium 默认，Claude/OpenCode 省略字段。无配置无本地值的默认产物保持原字节。损坏或不可读的本地文件报路径且阻止写入。解析、原生字段、保留与切换边界见[独立解析与清除](proposals/subagent-model-reasoning/design.md#独立解析与清除)、[原生编译与校验](proposals/subagent-model-reasoning/design.md#原生编译与校验)及[部署与观察](proposals/subagent-model-reasoning/design.md#部署保留与可观察性)。

已弃用且被忽略（仅告警）：`wiki.root` 与 `projects[].wikiRoot`；wiki 根只由 `docs_dir` 决定。

PI 宿主说明：

- 前提：PI 核心不含子代理，需先安装 pi-subagents 扩展（`pi install npm:pi-subagents`），PI 才会发现 `.pi/agents/` 下的原生 agent 文件。
- 部署产物：`agents deploy` 在 `pi` 宿主启用（`agents.hosts` 含 `pi`，默认启用）时写入两处——`.pi/agents/<name>.md`（子代理定义，含 `thinking` 档位、只读角色工具白名单、`skills` 绑定）与 `.pi/extensions/flowforge.ts`（项目级扩展，提供三项能力：① 拦截对 `agents.test_file_globs` 匹配文件的 write/edit；② 注册 `flowforge_frontier`/`flowforge_check` 原生工具（PATH 优先，再回退 `bin/flowforge`）；③ raw-script 委派模型保护——拦截经 `workflowScript`/`workflowScriptPath` 发起、且未显式给 `model` 的 `subagent` 派发，这类派发会静默继承编排会话的模型并绕过逐角色模型钉定）。扩展是宿主级受管资源，随 `pi` 宿主选中与否收敛（移出 `agents.hosts` 后再次 deploy 即删除），不随单个 subagent 的 remove 变化。
- 作用范围：前台（`async: false`）子代理不加载 ambient extensions（pi-subagents "Ambient extensions depend on where the child runs"），写拦截与 `flowforge_frontier`/`flowforge_check` 在主会话与后台/workflow 子代理（默认路径）生效。
- 逃生阀：两个键都只被扩展读取，且不在 `flowforge config` 的键集内，需直接编辑 `.flowforge/config.yaml`。`agents.disable_test_guard: true` 停用写拦截（与 opencode 宿主同语义）；`agents.disable_raw_script_model_guard: true` 停用 raw-script 委派模型保护。该保护也可逐次规避：改用具名 `agent:` 派发（带模型钉定），或在调用与脚本内每次 `runs.run()` 显式传 `model`。
- 手工冒烟：`pi -e ./assets/pi/flowforge.ts` 会话中尝试 write 任意 `*_test.go` 应被 block 并显示原因；修改扩展文件后用 `/reload` 重新加载。

## 稳定边界

- CLI 不提供通过参数传递长正文的 create/update 接口。
- 图计算和诊断是确定性的；需求或设计充分性仍由对应 Skill 负责。
- warning、gap、blocker 不写回 readiness 状态。
- waiver 必须匹配一个诊断和精确目标并记录理由；`*` 式全局跳过无效。
## Managed asset verification

`flowforge assets verify [project]` 比较运行中二进制内嵌的 Skills 与 agent rules 和所选项目，不修改任何文件。它把每个文件报告为 `current`、`missing`、`drifted` 或 `project-owned`；`--json` 为工具消费提供同样的事实。缺失或漂移的受管文件返回非零；`project-owned` 文件只作信息提示，验证永不覆盖它们。

`flowforge init` 与 `flowforge upgrade` 显式同步受管资产，然后用同一比较判定同步是否成功。若仍有受管文件不一致，它们打印其路径并提示用户运行 `flowforge assets verify`。


<a id="agent-model-state-restoration"></a>
### 方案切换与部署字段恢复（2026-09-30）

每机器内部快照 `.flowforge/agent-model-state.json` 按产物路径、model/effort 字段分别保存最后生成值和真实本地配置备份。实际字段等于最后生成值时继续使用备份；不同则捕获手改，包括字段删除。显式配置仍优先，临时方案退出或配置字段移除后恢复本地备份或目标 profile/default；`inherit` 暂时省略 effort，保留备份供退出后恢复。删除整文件清掉两字段备份。同值手改无法观察，继续视为生成值。

旧项目缺少快照时按已有保留规则迁移，不猜测旧字段来源；首次成功部署建立快照。快照精确列入 managed gitignore，`.flowforge/subagents/` 仍可提交。坏快照、未知版本或本地字段解析失败带路径报错，deploy/status 共用准备判断，status 完全只读且不会创建快照。配置、读取、编译及快照序列化完成后才写产物。

agent、PI extension 和宿主清理全部成功后，以同目录临时文件加 rename 原子提交快照。单角色部署保持其他记录；实际清理的宿主产物移除记录。最终提交失败保留旧快照、清理临时文件并报错，但产物可能已经更新；沿用非事务部署范围。`model-set use` 失败恢复激活指针，不承诺产物整体回滚。
