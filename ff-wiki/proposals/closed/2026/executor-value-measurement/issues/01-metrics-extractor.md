---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      executor-value-measurement-requirements: 1
    design:
      executor-value-measurement-design: 1
---

# 01: 指标提取器（extract → observations.md 幂等追加）

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

`scripts/executor_metrics.py extract`：只读打开 opencode DB，过滤目标项目的 `flowforge-implementer` 会话，按 d-extraction 六类指标逐会话输出一行，幂等追加到 `--out` 指向的 observations.md。

## Design context

三轮手工分析（2026-09-13）的 SQL/解析逻辑固化为一个零依赖脚本；幂等键 = session id，`--since` 支持增量。

See the design authority at [执行者价值度量方案](../design.md#executor-value-measurement-design)（d-extraction 节）. Requirement authority: [执行者价值度量需求](../requirements.md#executor-value-measurement-requirements)（目标 1/2，验收 1/2）.

## Touch points

- `scripts/executor_metrics.py` — 新建（python3 标准库；argparse 子命令 extract）
- `scripts/executor_metrics_test.py` — 新建（unittest：内存 fixture DB 回放已知形态）
- `docs/proposals/executor-value-measurement/observations.md` — 首次生成（含表头）

## Changes

- [x] 1. extract 子命令：会话过滤（`agent='flowforge-implementer'` AND `directory LIKE '%<project>'`）+ part 解析（steps = `step-start` 计数、tools = `type=tool` 计数、bash_n；命令规范化后计 rep_max；verdict 按 `BUILD SUCCESSFUL`/`BUILD FAILED`/`exit code N≠0` 判定，`state.metadata` 无 exitCode 不崩溃）+ token 列直读 + 末条 text part `STATUS:` 终态提取（缺失记 none）。
- [x] 2. 幂等追加：observations.md 已含 session id 的行跳过；文件不存在时写表头；`--since` 按 time_created 增量过滤。

## Execution detail

### Verified contracts

- opencode DB 实测 schema（2026-09-13，`sqlite_master` 取证）：`session(id, agent, model, directory, tokens_input, tokens_output, tokens_cache_read, tokens_cache_write, cost, time_created, time_updated, ...)`；`part(id, session_id, time_created, data)`；`part.data` 为 JSON，形态 `{"type":"text"|"step-start"|"tool","tool":<name>,"state":{"status","input","output","metadata","title"}}`。
- verdict 数据源事实：本机 opencode 4.x 的 bash part `state.metadata` 无 `exitCode` 键（实测 None）——verdict 必须由输出文本判定（`BUILD SUCCESSFUL` → ok、`BUILD FAILED` → fail、`exit code N` N≠0 → fail、其余 unknown）。
- 过滤实证：`agent='flowforge-implementer' AND directory LIKE '%tangram-v2%'` → 09-13 当日 18 会话；事故会话 2773 parts / 893 tools / 509 bash；STATUS 终态 BLOCKED/COMPLETED 均在末条 text part 出现。
- python3 标准库 `sqlite3`（URI 只读 `file:...?mode=ro`）、`json`、`argparse`、`re`、`datetime` 可用；`scripts/` 现有文件均为 shell（build.sh/install.sh/release.sh/install.ps1），新增 python 文件无命名冲突。

### Execution scenarios

- Success：对真实 DB 运行 extract，18 个 implementer 会话各产出 1 行；事故会话行 `steps=895`（step-start 计数）、`tools=893`、`rep_max≥184`、`status=BLOCKED`。
  - 更正（2026-09-13 实施取证）：事故会话（2773 parts/893 tools/509 bash/86.4min/rep 184，唯一命中 `ses_f66f341e5ffe5i4UyyAW1jrIoo`）末条 text part 为 `STATUS: COMPLETED` 终报；`BLOCKED` 仅出现于其**首条** text part（派发提示词引用的指令文本「歧义或规格冲突即 `STATUS: BLOCKED`」），把首条当末条是规划期手工 SQL 的 `LIKE '%"text"%'` + `ORDER BY time_created, id DESC` 优先级伪影。权威规则（末条 text part，d-extraction 与验收 2 一致钉死）不变，实测该行 `status=COMPLETED`；BLOCKED 终态在其余 4 个会话行出现。处置详见 Completion evidence 偏差条目。
- Success：重复运行同一命令，observations.md 行数不变（session id 幂等）；`--since` 传事故会话 time_created 后，仅输出其后的会话。
- Failure：DB 路径不存在或非 SQLite → 非零退出 + 明确错误信息，不创建 observations.md。
- Failure：part.data JSON 损坏（单行）→ 跳过该 part 并在 stderr 计数告警，会话行仍产出（不因单点脏数据丢整个会话）。

### Expected tests

- `scripts/executor_metrics_test.py`（unittest，内存 DB fixture 构造 session/part 行）：verdict 四分类、STATUS 提取（含缺失→none）、rep_max/fail_streak 计算、幂等（二次运行行数不变）、`--since` 边界（恰好等于 ts 含入）、损坏 JSON 跳过不崩溃。
- 验证命令：`python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` — 全部 PASS；`python3 scripts/executor_metrics.py extract --db <真实DB> --project tangram-v2 --out /tmp/opencode/obs-smoke.md` 冒烟非零退出码为 0 且行数 = 18。

### Generated artifacts

- `docs/proposals/executor-value-measurement/observations.md` — 首次生成：表头 + 09-13 全量 18 行基线（追加式，此后不手改）。

### Conventions

- python3 标准库零第三方依赖；DB 一律 URI `mode=ro` 只读打开。
- 命令正文仅保留规范化后前 120 字符（需求约束）。
- 不新增 flowforge CLI 子命令、不改 Go 代码（需求非目标）。
- 文件用 4 空格缩进、UTF-8；observations.md 行尾单一 `\n`（幂等 diff 友好）。

## Review rounds

### Round 1

- Fixed point: 工作树（新增 `scripts/executor_metrics.py`、`scripts/executor_metrics_test.py`；首次生成 `docs/proposals/executor-value-measurement/observations.md`；本票更新）。`flowforge check --dir docs/proposals/executor-value-measurement` 依赖图健康。
- Standards: 票面 Conventions 6 条逐条核验通过（python3 标准库零第三方依赖；DB 一律 `file:...?mode=ro`；命令正文规范化后 120 字符截断且不入输出行；不新增 flowforge CLI 子命令、不改 Go 代码；4 空格缩进 UTF-8；observations.md 行尾单一 `\n` 无 CR）。smell baseline 2 项 [Low] 均当场修正：测试 DDL 双副本（Duplicated Code → 提取共享 `SCHEMA` 常量）；`--out` 不可写路径裸 traceback（→ `OSError` 捕获为 `error: cannot write observations file` + exit 1，新增回归测试 `test_out_is_directory_errors_cleanly`）。其余 smell 无发现（`--epochs`/`--price-override` 为需求范围与设计运行手册明示接口，非投机泛化）。
- Spec: Changes 1/2、Expected tests 全部场景（verdict 四分类、STATUS 缺失→none、rep_max/fail_streak、幂等、`--since` 含入边界、损坏 JSON 跳过）、两个 Failure 场景（DB 缺失/非 SQLite/缺表 → exit 1 且不创建 observations.md；单行损坏 JSON → 跳过 + stderr 计数告警 + 会话行仍产出）、Generated artifacts（18 行基线）全部落地。1 项发现：[Medium→票面事实更正] Success 场景「事故会话行 status=BLOCKED」与权威规则及实测冲突（末条 text part 实为 COMPLETED 终报）——按「仓库事实更正过期位置/符号/命令细节」就地更正场景值并留痕，实现遵循 design d-extraction，非代码缺陷。
- Fix changes: none（2 项 Standards 发现当场修正并入交付；1 项 Spec 发现为票面事实更正）
- Design returns: none

## Completion evidence

- 交付行为：`scripts/executor_metrics.py extract`（python3 标准库零依赖）以 `file:...?mode=ro` 只读打开 opencode DB，按 `agent='flowforge-implementer'' AND directory LIKE '%<project>'` 过滤并按 time_created 升序逐会话输出 17 列 markdown 表行（d-extraction 15 字段 + `epoch`/`est_cost`，覆盖需求六类指标：身份/规模/循环信号/token/结局/成本），幂等键 session id；文件缺失或空时写表头；`--since`（epoch-ms 或 ISO8601，`>=` 含入）增量；`--epochs "name:<T,name:T1-T2,name:>T2"` 纪元标注（缺省 `-`）；`--price-override` 覆盖档位代理价（flash 0.30/2.50/0.075，旗舰 0.60/2.20/0.113）。verdict 四分类（`BUILD SUCCESSFUL`→ok、`BUILD FAILED`→fail、`exit code N≠0`→fail、其余 unknown 不计入也不重置 fail_streak）；bash 命令规范化（去首尾空白、压连续空白、截 120 字符）后计 rep_max；损坏 part.data JSON 跳过 + stderr 计数告警，会话行仍产出；DB 缺失/非 SQLite/缺 session|part 表 → exit 1 + 明确错误且不创建 observations.md；`--out` 不可写 → exit 1 干净报错。
- 验证命令与观测：`python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` → 26 tests OK。真实 DB 冒烟 `python3 scripts/executor_metrics.py extract --db /home/biqiang/.local/share/opencode/opencode.db --project tangram-v2 --out /tmp/opencode/obs-smoke.md` → exit 0、18 行；事故会话行 `dur_min=86.4 / steps=895 / tools=893 / bash_n=509 / rep_max=184 / fail_streak=20 / status=COMPLETED`。幂等：同命令重跑 md5 不变（appended 0, skipped 18）；`--since 1789274144000`（事故会话 time_created）→ 15 行且含该会话（含入边界实证）。`docs/proposals/executor-value-measurement/observations.md` 首次生成：27 行 = 9 行表头 + 09-13 全量 18 行基线（4 BLOCKED / 7 COMPLETED / 7 none）。`go build -trimpath ./cmd/flowforge` 通过（未改 Go 代码，健全性检查）。
- 双轴 review 与处置：Round 1 见上——Standards 2 项 [Low] 当场修正；Spec 1 项为票面事实更正；无 Critical/High，无未处置发现。
- 偏差：票面 Success 场景「事故会话行 status=BLOCKED」与权威规则及实测数据冲突，已就地更正并留痕（见 Execution scenarios 更正条目）：实现遵循 design d-extraction「末条 text part」规则，实测该行 `status=COMPLETED`，BLOCKED 终态在其余 4 会话行出现，与票面 Verified contracts「STATUS 终态 BLOCKED/COMPLETED 均在末条 text part 出现」一致。除此之外无权威偏差。
- 实现参考：本次提交（`scripts/executor_metrics.py`、`scripts/executor_metrics_test.py`、`docs/proposals/executor-value-measurement/observations.md`、本票），工作树 fixed point。
