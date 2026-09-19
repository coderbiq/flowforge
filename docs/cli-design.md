# FlowForge CLI 行为与策略

CLI 是 Markdown 工作流的确定性投影器，不是内容管理 API。所有需求、设计、ticket 和 evidence 都由文件工具直接编辑。

## 目录解析

`.flowforge/config.yaml` 的 `docs_dir` 指定文档根目录，默认为 `docs`，支持相对项目根目录和绝对路径。`check`、`frontier`、`status` 从当前目录向上定位项目配置；显式 `--dir` 则直接使用调用者给出的 proposals 路径。损坏的项目配置必须报错。

`flowforge init [path]` 创建配置、`<docs_dir>/CONTEXT.md`、`adr/` 和 `proposals/`，部署 `<docs_dir>/agents/`、`.agents/skills/` 并维护 `AGENTS.md` 中的受管区块。已有配置始终保留；`--force` 只强制刷新受管资产。

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

检查循环依赖、悬空依赖、自依赖和 Catalog 诊断。默认 warning 与 gap 保持可见但不让 check 失败，blocker 始终失败；`--strict` 让未豁免 warning/gap 也失败。JSON 输出包含 issue/artifact 数量、DAG 结果和完整 diagnostics。

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
- `flowforge agents deploy [name]`：将权威 subagent 定义编译并部署到宿主原生目录（`.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/`）。
- `flowforge agents remove <name>`：移除 subagent 部署文件；内置角色写入 `.flowforge/config.yaml` 的 `agents.disabled` 持久化停用，自定义角色删除源文件。
- `flowforge agents status [--json]`：报告各 subagent 在各宿主目录中的状态（`current`、`missing`、`drifted`、`project-owned`）。
- `flowforge config get|set|list`：读取或修改 `docs_dir`、`standards.guide`、version check 及兼容项目配置。
- `flowforge upgrade`：更新 CLI；已经是相同版本时仍同步当前项目的受管资产与 subagent，降级返回独立错误。
- `flowforge version`：显示构建注入版本。

PI 宿主说明：

- 前提：PI 核心不含子代理，需先安装 pi-subagents 扩展（`pi install npm:pi-subagents`），PI 才会发现 `.pi/agents/` 下的原生 agent 文件。
- 部署产物：`agents deploy` 在 `pi` 宿主启用（`agents.hosts` 含 `pi`，默认启用）时写入两处——`.pi/agents/<name>.md`（子代理定义，含 `thinking` 档位、只读角色工具白名单、`skills` 绑定）与 `.pi/extensions/flowforge.ts`（项目级扩展：拦截对 `agents.test_file_globs` 匹配文件的 write/edit，注册 `flowforge_frontier`/`flowforge_check` 原生工具，PATH 优先回退 `bin/flowforge`）。扩展是宿主级受管资源，随 `pi` 宿主选中与否收敛（移出 `agents.hosts` 后再次 deploy 即删除），不随单个 subagent 的 remove 变化。
- 作用范围：前台（`async: false`）子代理不加载 ambient extensions（pi-subagents "Ambient extensions depend on where the child runs"），写拦截与 `flowforge_frontier`/`flowforge_check` 在主会话与后台/workflow 子代理（默认路径）生效。
- 逃生阀：`.flowforge/config.yaml` 设 `agents.disable_test_guard: true` 时扩展不注册写拦截（与 opencode 宿主同语义）。
- 手工冒烟：`pi -e ./assets/pi/flowforge.ts` 会话中尝试 write 任意 `*_test.go` 应被 block 并显示原因；修改扩展文件后用 `/reload` 重新加载。

## 稳定边界

- CLI 不提供通过参数传递长正文的 create/update 接口。
- 图计算和诊断是确定性的；需求或设计充分性仍由对应 Skill 负责。
- warning、gap、blocker 不写回 readiness 状态。
- waiver 必须匹配一个诊断和精确目标并记录理由；`*` 式全局跳过无效。
# Managed asset verification

`flowforge assets verify [project]` compares the running binary's embedded Skills and agent rules with the selected project without changing files. It reports each file as `current`, `missing`, `drifted`, or `project-owned`; `--json` provides the same facts for tooling. Missing or drifted managed files return non-zero. Project-owned files are informational and are never overwritten by verification.

`flowforge init` and `flowforge upgrade` explicitly synchronize managed assets, then use this same comparison before reporting synchronization success. If a managed file remains divergent, they show its path and direct the user to `flowforge assets verify`.
