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
flowforge frontier [--dir <path>] [--json] [--quiet] [--strict] [--include-gaps]
```

先计算 DAG 的无阻塞 open ticket，再按内容诊断投影：

- clean ticket 默认输出；
- warning ticket 默认输出并保留诊断；
- gap ticket 默认排除，`--include-gaps` 可显式包含；
- blocker ticket 始终排除；
- `--strict` 只输出 clean ticket，并优先于 `--include-gaps`。

`--quiet` 只把可执行路径写到 stdout，诊断写到 stderr；`--json` 保留 clean、warning、gap、claimed、blocked 分组，适合 Agent 或自动化消费。

## 其他命令

- `flowforge status [--dir <path>]`：按 feature 汇总 ticket 生命周期和 DAG 状态。
- `flowforge config get|set|list`：读取或修改 `docs_dir`、version check 及兼容项目配置。
- `flowforge upgrade`：更新 CLI；已经是相同版本时仍同步当前项目的受管资产，降级返回独立错误。
- `flowforge version`：显示构建注入版本。

## 稳定边界

- CLI 不提供通过参数传递长正文的 create/update 接口。
- 图计算和诊断是确定性的；需求或设计充分性仍由对应 Skill 负责。
- warning、gap、blocker 不写回 readiness 状态。
- waiver 必须匹配一个诊断和精确目标并记录理由；`*` 式全局跳过无效。
