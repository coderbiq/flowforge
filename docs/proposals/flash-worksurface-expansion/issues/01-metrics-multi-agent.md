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

# 01: 度量多角色提取（extract --agents + report 按 agent 分组）

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

`scripts/executor_metrics.py` 的 extract 支持 `--agents`（逗号分隔，默认 `flowforge-implementer`），observations 表新增 `agent` 列；report 按 agent 分组聚合并逐角色出 G1 失控门。

## Design context

度量先于迁移（需求约束）：P1 迁移 investigator/explore 前必须能观测这两个角色。本 proposal 使用独立观测文件 `docs/proposals/flash-worksurface-expansion/observations.md`，不迁移 executor-value-measurement 的既有文件（其表无 agent 列且观察窗运行中，09-20 收口）。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-phase-1 度量先行节）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 1/4，验收 1）.

## Touch points

- `scripts/executor_metrics.py` — extract 子命令（SQL 过滤、表头构造、幂等追加）与 report 子命令（聚合、门判定）
- `scripts/executor_metrics_test.py` — unittest fixture

## Changes

- [x] 1. extract 增加 `--agents` 参数（默认 `flowforge-implementer`）：SQL 过滤 `agent='flowforge-implementer'` 改为 `agent IN (?,...)`；表头在 `session` 列后插入 `agent` 列，数据行同步。
  - cmd: python3 scripts/executor_metrics.py extract --db /home/biqiang/.local/share/opencode/opencode.db --project tangram-v2 --agents flowforge-investigator,explore --out /tmp/opencode/obs-multi.md
  - exit: 0
  - output: "extract: 11 session(s) matched, appended 11, skipped 0 existing"；表头 "| session | agent | ts | ..."；行含 3 × "| flowforge-investigator |"、8 × "| explore |"
  - artifact: scripts/executor_metrics.py
- [x] 2. report 解析 observations：按 `agent` 列分组聚合（epoch × agent × stratum）；agent 列缺失的旧行按 `flowforge-implementer` 兼容；G1/G4 门逐 agent 输出（阈值沿用）。
  - cmd: python3 scripts/executor_metrics.py report --obs /tmp/opencode/obs-multi.md --project-root /vol3/1000/develop/tangram-v2 --out /tmp/opencode/report-multi.md
  - exit: 0
  - output: "report: 11 session(s), 2 agent(s), 0 epoch(s), strata unknown=11"；4 个 section 各含 "### explore" / "### flowforge-investigator" 子节
  - artifact: scripts/executor_metrics.py
- [x] 3. 测试：多 agent 提取（两 agent 各出各行）、幂等不变、旧表无 agent 列的 report 兼容、按 agent 分组的门判定各 1 例。
  - cmd: python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v
  - exit: 0
  - output: "Ran 108 tests in 0.258s" / "OK"（现有 97 例零回归 + 新增 11 例）
  - artifact: scripts/executor_metrics_test.py

## Constraints

- 不修改既有 `docs/proposals/executor-value-measurement/observations.md` 与 report.md（运行中的观察窗产物）。
- DB 一律 URI `mode=ro` 只读打开；python3 标准库零第三方依赖。
- Write set: `scripts/executor_metrics.py`, `scripts/executor_metrics_test.py`, `docs/proposals/flash-worksurface-expansion/`

## Done and verify

- 多角色提取冒烟: `python3 scripts/executor_metrics.py extract --db /home/biqiang/.local/share/opencode/opencode.db --project tangram-v2 --agents flowforge-investigator,explore --out /tmp/opencode/obs-multi.md` — 退出码 0，行含 investigator 与 explore 会话，表头含 agent 列。
- 全部测试通过: `python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` — 全部 PASS（含新增 3 类用例）。
- 既有观测不被破坏: `python3 scripts/executor_metrics.py extract --db /home/biqiang/.local/share/opencode/opencode.db --project tangram-v2 --out /tmp/opencode/obs-impl.md` 后 `/tmp/opencode/obs-impl.md` 行数与 docs/proposals/executor-value-measurement/observations.md 一致（42→72 会话幂等语义不变，仅新文件含 agent 列）。

---

## Execution detail

### Verified contracts

- `scripts/executor_metrics.py:98` — `COLUMNS` 元组（17 列，无 agent 列）；`:102` `HEADER_LINES` 由 COLUMNS 生成表头（含写死的一句 "Append-only per-session metrics for `flowforge-implementer` sessions"）。
- `scripts/executor_metrics.py:965` — `collect_metrics(conn, project, since_ms=None, prices=None, epochs=None)`；`:968` SQL 硬编码 `WHERE agent='flowforge-implementer' AND directory LIKE ?`。
- `scripts/executor_metrics.py:287` — `parse_obs_rows(text)` 按表头解析并做严格列校验（`missing = [c for c in COLUMNS if c not in header]` 非空即抛 `MetricsError`）：agent 列必须作为**可选列**处理（缺失时不进 missing 校验，行级回填 `flowforge-implementer`），否则旧格式 observations 文件全部炸。
- `scripts/executor_metrics.py:695` — `render_report(rows, epochs, strata, prices, *, obs_name, epochs_spec, warnings, price_note)`；门判定（G1/G4）在此函数内。`:373` `resolve_strata(rows, project_root)` 按 ticket 路径解析 stratum——investigator/explore 会话首条 text part 可能无 `proposals/<p>/issues/` 路径，落入 stratum `unknown` 是预期行为，不修。
- `scripts/executor_metrics.py:1014/:1031` — argparse 子命令；extract 现有参数 `--db --project --out --since --epochs --price-override`。
- `scripts/executor_metrics_test.py` — `SessionMetricsTests.test_filter_agent_directory_and_order` 直接调 `em.collect_metrics(conn, "tangram-v2")`：`collect_metrics` 新增 `agents` 参数必须带默认值 `("flowforge-implementer",)` 保持该用例语义；CLI 级用例参考 `CliExtractTests.test_writes_header_and_rows` / `test_idempotent_rerun_and_incremental_append` 的 fixture 写法（`make_conn` / `add_session` / `add_part` / `dispatch_and_report_parts` helper 已存在，`add_session` 已有 `agent=` 可选参数）。

### Execution scenarios

- Success：`extract --agents flowforge-investigator,explore` → 仅输出这两个 agent 的会话行，每行 agent 列有值；表头含 agent 列（session 之后）。
- Success：对含 agent 列的新格式 observations 运行 report → 聚合表与 G1/G4 门按 agent 分组输出（epoch × agent × stratum）。
- Failure（兼容路径）：对无 agent 列的旧格式文件（executor-value-measurement/observations.md 副本到 /tmp）运行 report → 不报错，所有行按 `flowforge-implementer` 处理，产出与旧版一致的 implementer 聚合。
- Failure：`--agents` 含无会话的 agent 名 → 0 行追加、退出码 0（空结果不是错误）。

### Expected tests

- `python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` — 全部 PASS（现有 97 例零回归 + 新增：多 agent SQL 过滤各出各行、agent 列写入、旧格式 report 兼容回填、按 agent 分组的门判定）。

### Generated artifacts

- `docs/proposals/flash-worksurface-expansion/observations.md` — producer：首次 `extract --agents ... --out` 运行；consumer：`report --obs`。同步断言：文件含 agent 列表头 + 首批历史会话行（investigator 3 个 glm 会话入库为 baseline 纪元）。

### Conventions

- DB 一律 URI `mode=ro` 只读打开；python3 标准库零第三方依赖；行尾单一 `\n`。
- 不修改 `docs/proposals/executor-value-measurement/` 下任何文件（运行中的观察窗产物，Constraints 转录）。
- 门阈值沿用现有 G1 形状（steps>150 / rep_max>10 / dur>30min），首轮观察后调整需修订 design.md。

## Implementation note

**做了什么**（full mode，TDD：先新增 11 个失败用例再实现）：

1. `scripts/executor_metrics.py`：
   - `DEFAULT_AGENT`/`OPTIONAL_COLUMNS`/`REQUIRED_COLUMNS` 常量；`COLUMNS` 在 `session` 后插入 `agent`（17→18 列）；`HEADER_LINES` 同步且描述句不再写死 implementer（新增"旧行按 implementer 读"说明）。
   - `parse_agents()`（逗号分隔、去空格去重、空 spec=空集）；`_agent_filter()` 生成 `agent IN (?,?)`（空集退化为恒假子句 `0`，空结果不是错误）；`collect_metrics(..., agents=("flowforge-implementer",))` 保持既有默认调用语义。
   - `_session_metrics` SELECT 追加 `agent` 列并写入每行 metrics。
   - `parse_obs_rows`：agent 为可选列（缺失不进严格列校验），行级缺失/空值回填 `flowforge-implementer`。
   - `render_report`：先按 agent 切分 rows，逐 agent 做 `aggregate_cells` + `evaluate_gates`（G1/G4 阈值常量未动）；多 agent 时各 section 出 `### <agent>` 子节，ALARM 行加 `[agent]` 前缀；单 agent 保持与旧版逐字节一致的表体。`cmd_report` 摘要行加 agent 计数。
   - **旧观察窗保护（Constraints 推导）**：`_file_columns()` + `format_row(metrics, columns)`——向已有表头的文件追加行时按该文件自己的列格式输出，默认 agent 的每日 extract 不会向运行中的 17 列旧表追加 18 列行（否则 09-20 收口 report 会因列数不匹配炸掉）。测试 `test_extract_appends_in_legacy_files_own_columns` 锁定该行为。
2. `scripts/executor_metrics_test.py`：`rrow` 加可选 `agent=`（缺省不写入键，保持旧 fixture 形状）；`test_writes_header_and_rows` 表头断言随 Change 1 更新为 18 列（Change 显式改表头，预设断言随之修订）；新增 4 个测试类 11 例（多 agent SQL 过滤/CLI 提取/幂等/空 agent、旧格式 parse 回填 + report 兼容 + 追加列格式保护、按 agent 的 G1 隔离判定与告警前缀）。
3. `docs/proposals/flash-worksurface-expansion/observations.md`：首次 extract 生成，11 行（investigator 3 个 glm-5.3 会话 = baseline 纪元 + explore 8 个），重跑幂等（appended 0, skipped 11）。

**验证结果**（全部通过）：

- `python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` → Ran 108 tests, OK（97 现有零回归 + 11 新增）。
- 多 agent extract 冒烟 → exit 0，11 行（3 investigator + 8 explore），表头含 agent 列。
- 多 agent report 冒烟（--project-root tangram-v2）→ exit 0，四个 section 均按 `### explore`/`### flowforge-investigator` 分组。
- 旧格式兼容：旧表副本 report → exit 0，72 行全部按 implementer 处理；用运行窗 report 自身参数（5 epochs + project-root）重跑，聚合/对比/状态/门四张表与 `executor-value-measurement/report.md` 逐行一致（仅节尾空行差异）。
- 既有观测不被破坏：默认 extract 到新文件 → 72 行与旧文件会话数一致，仅新文件含 agent 列。
- `flowforge check --dir docs/proposals/flash-worksurface-expansion` → exit 0，依赖图健康（02/03/05/06 的 execution-contract gap 为既有诊断，非本票范围）。

**环境注记**：`executor-value-measurement/*.md` 目录内文件为 mode 000（观察窗锁定惯例）；旧格式兼容冒烟的副本需 `chmod 644` 于 /tmp 副本（源文件未动，diff 实验中 report.md 的 mode 已恢复 000）。

**Write set 合规**：全部修改限于 `scripts/executor_metrics.py`、`scripts/executor_metrics_test.py`、`docs/proposals/flash-worksurface-expansion/`（ticket + observations.md）。工作树中其他改动（Go 文件、04 票、executor-value-measurement 运行窗产物）为本任务之前已存在，未触碰。

## Review rounds

### Round 0

- Gaps: none（逐 Change 对照 diff：1/2/3 均 delivered，无 missing/partial/contradicts/unrequested-without-disposition）
- Disposition: none
- Escalated to dual axes: yes（单会话环境无 subagent 派发工具，按协议内联分轴执行并如实记录）

### Round 1

- Fixed point: 447f033（工作树 diff：scripts/executor_metrics.py、scripts/executor_metrics_test.py、docs/proposals/flash-worksurface-expansion/observations.md + 本票）
- Standards: 1 × [Low] Duplicated Code（render_report 四处重复的 `### <agent>` 子节发射形状）→ 当轮修复为 `agent_heading()` 闭包 helper，108 测试复跑全绿；其余 smell 候选（Data Clumps=既有形状、Speculative Generality=format_row 列参数有两个真实消费方、Primitive Obsession=脚本级设计取舍）均已证伪并记录。文档化条款（mode=ro、stdlib-only、单 `\n`、不改 executor-value-measurement、不动 Go/*_test.go）全部符合。
- Spec: none blocking。5 项超字面变更均记为 disposition：① 旧文件按自身列格式追加（Constraints"不改运行窗产物"的必要推导，有测试锁定）；② report 头部 agents bullet 与 stdout agent 计数；③ ALARM `[agent]` 前缀（d-decision-gates 按角色回滚的作用域要求）；④ G2/G3 亦逐 agent 判定（"按 agent 分组聚合"语义内，investigator 的 baseline/flash 对比正是后续用途）；⑤ 单 agent 不出 `###` 子节（兼容场景"与旧版一致的聚合"已逐行验证）。
- Fix changes: none
- Design returns: none
- Repair: none

## Completion evidence

- 交付行为：`extract --agents a,b` 多角色提取（默认 implementer，幂等键仍为 session id）；observations 表 `session` 后含 `agent` 列；report 按 epoch × agent × stratum 聚合，G1/G4（及 G2/G3）逐 agent 判定、阈值沿用；无 agent 列的旧表 report 兼容回填 implementer，向旧表的追加保持其 17 列格式，运行中的 executor-value-measurement 观察窗不受影响。
- 命令与观测：见 Changes 证据四元组与 Implementation note 验证结果（108 测试 OK；三组冒烟 exit 0；旧格式表体与运行窗 report 逐行一致；72 会话计数一致；`flowforge check` exit 0）。
- 双轴与处置：Round 0 零 gap；Round 1 Standards 1 Low（已修复）+ Spec 0 blocking（5 项 disposition 记录于 Review rounds）。
- 偏差与处置：超字面变更 5 项全部在 spec 语义或 Constraints 推导范围内，无一需要 authority 裁定；表头断言的预设测试更新源于 Change 1 显式改表头。
- 实现参考：工作树 diff vs 447f033（scripts/executor_metrics.py +297−102 中主体、scripts/executor_metrics_test.py 新增 4 测试类 11 例、docs/proposals/flash-worksurface-expansion/observations.md 新建 11 行）；closeout 提交随本票一并落库。
