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

# 02: 纪元 × 难度档分组与 report 聚合

**Blocked by:** 01
**Status:** open
**Mode:** full

## Delivery

`scripts/executor_metrics.py report`：读 observations.md，按 `--epochs` 时间戳映射纪元、按票文件客观代理判 S/M/L 档，输出纪元 × 档位聚合表（n、dur 中位数、单票 est_cost 中位数、cacheR 占比、失控事件数）到 report.md（当日替换制）。

## Design context

纪元三段（pre-hardening / provider-switch / hardened）显式传参不自动猜；难度档读票文件 Changes 行数与 evidence cmd 类型；est_cost 用 d-cost 档位代理价。

See the design authority at [执行者价值度量方案](../design.md#executor-value-measurement-design)（d-epoch / d-strata / d-cost 节）. Requirement authority: [执行者价值度量需求](../requirements.md#executor-value-measurement-requirements)（目标 3，验收 3/4/5）.

## Touch points

- `scripts/executor_metrics.py` — report 子命令 + epoch/strata/cost 模块
- `scripts/executor_metrics_test.py` — 分组与聚合用例
- `docs/proposals/executor-value-measurement/report.md` — 首次生成

## Changes

- [ ] 1. 纪元映射：`--epochs "name:<ts,name:ts-ts,name:>ts"` 三形态解析；会话按 time_created 落位；未覆盖区间记 `unclassified`。
- [ ] 2. 难度档：读 `--project-root` 下 `**/issues/*.md` 票文件，统计 `- [ ]`/`- [x]` Changes 行数 + evidence cmd 类型 → S（≤2 且单测级）/ L（≥4 或 integrationTest/E2E/--rerun）/ M（其余）；票文件缺失记 `unknown` 档并在报告计数。
- [ ] 3. report 输出：纪元 × 档位聚合 + est_cost（P_in/P_out/P_cacheR 常量，`--price-override` 覆盖）+ cacheR 占比 + 失控事件数（steps>150 / rep_max>10 / dur>30min 任一）。

## Execution detail

### Verified contracts

- 纪元时间戳来源（本 proposal design 运行手册钉定，2026-09-13 部署 mtime 实测）：pre-hardening `< 1789274100000`（17:25 前，含两次事故）；provider-switch `1789274100000-1789283100000`（17:25–19:25，无 steps/Non-negotiables 的重部署）；hardened `> 1789283100000`（19:25 后，steps: 200 + Non-negotiables 生效）。
- 票文件路径实证：implementer 会话首条 text part 含 `proposals/<name>/issues/<id>.md` 绝对路径（01 票提取的 ticket 字段）；tangram-v2 票根目录 `ff-wiki-v5/proposals/*/issues/`。
- 中位数定义：偶数 n 取两中间值均值（statistics.median 标准行为）。
- est_cost 公式与常量钉定于 design d-cost：`in/1e6*0.30 + out/1e6*2.50 + cacheR/1e6*0.075`（flash）与 `0.60/2.20/0.113`（旗舰），`--price-override json` 整表替换。

### Execution scenarios

- Success：09-13 基线回放——事故会话（12:35）落 pre-hardening 且失控计数 ≥1；17:38–18:53 的 9 会话落 provider-switch 且失控 0；est_cost 单票中位数量级与手工核算一致（pre-hardening ≈$3.1/票，provider-switch ≈$0.2/票）。
- Success：flash/旗舰混合 observations 下，同档分组各自聚合，档内对比行输出两者中位数与比值。
- Failure：`--epochs` 语法非法（缺冒号/时间戳非数字）→ 非零退出 + 指出非法片段，不产出 report。
- Failure：`--project-root` 无匹配票文件 → 全部记 unknown 档，report 顶部 warning 行提示，不中断。

### Expected tests

- unittest 用例：epoch 三形态边界（ts 恰好等于边界归右闭区间）、unclassified 落位、strata S/L 边界（2 Changes 单测级 → S；4 Changes → L；integrationTest → L 无论行数）、票文件缺失 → unknown、est_cost 计算与 price-override、失控判定三条件各自触发、report 表含全部纪元行。
- 验证命令：`python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` 全 PASS；对 09-13 真实 observations.md 运行 report 冒烟，产出三纪元分组。

### Generated artifacts

- `docs/proposals/executor-value-measurement/report.md` — 首次生成：纪元 × 档位聚合 + 09-13 基线数据行（当日替换制模板）。

### Conventions

- 沿用 01 票 Conventions（零依赖、只读、UTF-8、单 `\n`）。
- report 为覆盖写（替换制），observations 为追加制——两者不得混淆。
- 判定阈值（失控 150/10/30、S/M/L 边界）作为模块常量集中定义，不在多处散落字面量。

## Completion evidence

(pending)
