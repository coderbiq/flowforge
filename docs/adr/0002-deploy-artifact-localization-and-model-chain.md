# ADR 0002: 部署产物本地化与六级模型优先级链

日期：2026-09-25 · 状态：已接受 · 来源提案：[deploy-artifact-localization](../proposals/deploy-artifact-localization/design.md)

## 背景

两个并存问题：其一，subagent/AGENTS 部署产物被项目 git 跟踪，携带本机模型 ID 随仓旅行——换机 clone 即配置事故（[宿主模型外置能力调查](../research/2026-09-20-host-model-portability.md)）；其二，模型钉扎只有全局层（`agents.models`/`models_by_name`），pi 宿主产物无 `model` 字段（继承会话模型），per-host 差异配置无通道。

## 决策

1. **产物本地化**：部署产物与 `.flowforge/config.yaml` 是**每机器文件**——init/upgrade 自动 gitignore（幂等追加；config 必须单文件条目而非 `.flowforge/` 目录，自定义源仍入库）；存量已跟踪产物经 `git ls-files` 检测后输出可复制 `git rm --cached -r` 指引，**绝不自动改 git 索引**。
2. **`models_by_host`** 嵌套双键 map（外层宿主键 → 内层 agent 名 | profile 档位键）；**per-host 层整体压过全局层**，同层内 name 键 > profile 键——六级链：`by_host[name] > by_host[profile] > by_name[name] > models[profile] > preserve-merge 回填 > 宿主默认`。
3. **pi 注入**：pi frontmatter 增 `model` 字段（Description/Thinking 之间），取链解析值，空则省略键——未配置项目产物逐字节不变；preserve-merge 对 pi 真实化，提示文案指向 `models_by_name/models_by_host` 真实通道。
4. **格式级 fail-fast 校验**（`validateModelConfig`，deploy/status 两路径任何写入之前、同一报错）：键全量校验；`models_by_name` 只收 agent 名（惰性键=配置损坏信号）；`models_by_host.codex` 即错（编译器丢弃 model）；值按启用宿主规则（opencode/pi 必须 `provider/model`，claude 单 token）。

## 取舍与被拒替代

- **拒绝外置通道**（宿主用户级配置 / 生成 `.pi/settings.json`）：调查证实宿主外置能力不一且引入第二真相源；产物内注入单源可验。
- **拒绝 `--untrack` flag**：一次性迁移不值新旗标；自动改索引违反最小惊讶，指引已足够。
- **拒绝扁平 host 前缀键**（`models: {"pi/tool-capable": …}`）：模型值本身含 `/`，键再引入 `/` 产生解析歧义。
- **拒绝跨层按具体度混排**：层间整体排序可局部推理，混排使配置审查需全局计算。
- **拒绝宽松校验**：会把 `model: sonnet` 写进 opencode 产物运行期才爆；存在性校验被纯本地红线禁止。

## 后果

- config 是每机器文件：入库即需求问题 2 的可移植性事故，gitignore 是契约不是便利。
- `by_name` 写 profile 键直接报错（收敛裁决，rev 3）；codex 无 per-host 键。
- 六级链被表驱动测试钉死；未配置项目经双二进制 sha256 清单保证逐字节回归零。
- GIIS 形态实证：既有全局 profile 钉扎**零 config 变更**自动作用于 pi 宿主。

## 验证

票 01-03 证据四元组（六级链/校验矩阵/幂等 gitignore/索引不变/逐字节回归）+ GIIS 实弹（v5.11.0 deploy：investigator / reviewer-lite / implementer 的 pi 产物均带 `model: cpa/deepseek-v4.1-flash`）。
