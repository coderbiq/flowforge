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
**Status:** closed
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

- [x] 1. 纪元映射：`--epochs "name:<ts,name:ts-ts,name:>ts"` 三形态解析；会话按 time_created 落位；未覆盖区间记 `unclassified`。
- [x] 2. 难度档：读 `--project-root` 下 `**/issues/*.md` 票文件，统计 `- [ ]`/`- [x]` Changes 行数 + evidence cmd 类型 → S（≤2 且单测级）/ L（≥4 或 integrationTest/E2E/--rerun）/ M（其余）；票文件缺失记 `unknown` 档并在报告计数。
- [x] 3. report 输出：纪元 × 档位聚合 + est_cost（P_in/P_out/P_cacheR 常量，`--price-override` 覆盖）+ cacheR 占比 + 失控事件数（steps>150 / rep_max>10 / dur>30min 任一）。

## Execution detail

### Verified contracts

- 纪元时间戳来源（本 proposal design 运行手册钉定，2026-09-13 部署 mtime 实测）：pre-hardening `< 1789291500000`（17:25 前，含两次事故）；provider-switch `1789291500000-1789298700000`（17:25–19:25，无 steps/Non-negotiables 的重部署）；hardened `> 1789298700000`（19:25 后，steps: 200 + Non-negotiables 生效）。
  - 更正（2026-09-13 实施取证）：原钉定 `1789274100000/1789283100000` 换算实为 12:35/15:05 +08:00（12:35 恰为事故会话始发分钟——规划期把事故时间误作边界编码），与 d-epoch 规范表（17:25/19:25）、本票回放场景（事故会话须落 pre-hardening，而 12:35:44 > 12:35:00 会落 provider-switch）、手工成本核算（pre ≈$3.1/票 = flash 5 会话 est_cost 均值 3.107）及部署产物 mtime 四路证据矛盾。实证：`tangram-v2/.opencode/agent/*.md` mtime = 2026-09-13 19:25:19 +0800 且内容含 `steps: 200` 与 `## Non-negotiables`（hardened 判据）；17:25+08:00 = 1789291500000 ms，19:25+08:00 = 1789298700000 ms。design.md 运行手册示例命令中同两个字面量需设计 owner 同步更正（本票写集外，留待处置）。
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

## Review rounds

### Round 1

- Fixed point: 工作树（`scripts/executor_metrics.py` 新增 report 子命令 + epoch/strata/cost 模块；`scripts/executor_metrics_test.py` 扩展 38 用例；`docs/proposals/executor-value-measurement/report.md` 首次生成；本票纪元字面量更正留痕）。`flowforge check --dir docs/proposals/executor-value-measurement` 依赖图健康。
- Standards: 票面 Conventions 7 条逐条核验通过（python3 标准库零第三方依赖；observations 对 report 只读、report 覆盖写与追加制不混淆且有测试锚定；不新增 flowforge CLI 子命令、不改 Go 代码；4 空格缩进 UTF-8 单 `\n`；阈值 150/10/30 与 S/M/L 边界为单处模块常量）。smell baseline 2 项 [Low] Duplicated Code 当场修正：est_cost 公式双写（extract 内联 vs report est_cost_for → 提取 `est_cost_tokens` 单点，extract 改调）；价格表合并双写（collect_metrics vs cmd_report → 提取 `effective_prices`）。其余 smell 无发现（`--out` 默认值与 `price_note`/`generated` 渲染参数均被生产路径消费，非投机泛化）。
- Spec: Changes 1/2/3、4 个 Execution scenarios、Expected tests 全部枚举项（三形态边界归右闭、unclassified、S/L 边界与 integrationTest 无论行数、票缺失→unknown、est_cost 重算与 override、失控三条件各自触发、report 含全部纪元行）、Generated artifacts 全部落地（真实 09-13 冒烟 + 独立核算逐值吻合，见 Completion evidence）。零 finding。偏差 3 项均已处置（纪元字面量更正 / design.md 示例待 owner 同步 / price-override 三元组覆盖语义钉定，见 Completion evidence 偏差条目）。
- Fix changes: none（2 项 Standards [Low] 当场修正并入交付）
- Design returns: none（design.md 运行手册示例字面量为同义事实修正——规范表 17:25/19:25 含义不变——不涉责任/接口/缝变更，已记录待设计 owner 执行，不阻塞本票）

## Completion evidence

- 交付行为：`scripts/executor_metrics.py report`（python3 标准库零依赖）：只读解析 observations.md（表头驱动、17 列强校验），按 `--epochs` 三形态以 ts（会话 time_created）落纪元、未覆盖记 `unclassified`；按 `--project-root` 遍历 `**/issues/*.md`（剪枝 node_modules/.git 等）后缀匹配票文件并施客观档位代理——L = Changes ≥ 4 或命令含 integrationTest/E2E/--rerun，S = ≤ 2 且全部命令 span 单测级，M 其余，票缺失或无 `-` 记 `unknown` 并在报告顶部 warning 计数；est_cost 自 token 列按档位代理价重算（`--price-override` 整三元组逐档覆盖，未提及档位保留常量）；输出纪元 × 档位聚合（n、dur 中位数、单票 est_cost 中位数、cacheR 成本占比、失控数=steps>150∨rep_max>10∨dur>30min 单处常量）+ 同档 flash vs 旗舰中位数与比值表 + 分纪元 status 计数（含 `none` 无终态计数）；report 为覆盖写（默认写到 obs 同目录 report.md），observations 不被触碰；`--epochs` 语法非法 → exit 1 指名非法片段且不产出/截断 report。
- 验证命令与观测：`python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` → 64 tests OK（#01 的 26 个保留 + 新增 38 个，零删改）。真实冒烟 `report --obs .../observations.md --epochs "pre-hardening:<1789291500000,provider-switch:1789291500000-1789298700000,hardened:>1789298700000" --project-root <tangram-v2>` → exit 0，18 会话全落 L 档（5 张 timeline 票按客观代理解析成功）：pre-hardening L n=8 / dur_med 16.6 / est_cost_med $0.7374 / runaway 2（事故会话 12:35 steps=895 与 14:13 steps=527）；provider-switch L n=10 / dur_med 1.45 / est_cost_med $0.0975 / runaway 0（17:38–18:53 九会话全部落位且失控 0）；hardened 零行在场；同档对比 pre flash $0.8635(n=5) vs 旗舰 $0.6113(n=3) 比值 1.41、provider flash $0.1051(n=9) vs $0.0239(n=1) 比值 4.40；est_cost 量级与手工核算一致（手工 ≈$3.1/$0.2 系 flash 均值 3.107/0.210，报告 flash 中位 0.8635/0.1051 同为 $10^0/$10^-1 量级）；status 计 pre 3C/4B/1none、provider 4C/0B/6none（与 #01 的 18 行 4 BLOCKED/7 COMPLETED/7 none 总量交叉吻合）。独立核算脚本（不经被测代码）逐值一致。失败路径：`--epochs` 缺冒号/时间戳非数字 → exit 1 且错误指名片段、不产出 report（既有 report 不被截断）；空 `--project-root` → 全 unknown + 顶部 warning 行、exit 0 不中断。幂等安全：report 运行前后 observations.md 字节不变（git diff 空）。Go 侧健全性：`go vet ./cmd/flowforge ./internal/...` 通过、`go test ./internal/...` 全 ok（本票零 Go 代码改动）；`go build -o bin/flowforge` 在会话末因宿主环境对临时可执行文件创建的限制（`/tmp/go-build*/exe/a.out: permission denied`，会话初 `go run` 同源构建曾成功）未能完成链接，与本交付无关。
- 双轴 review 与处置：Round 1 见上——Standards 2 项 [Low] 当场修正；Spec 零 finding、偏差 3 项均处置；无 Critical/High，无未处置发现。
- 偏差：(1) 票面钉定纪元字面量 `1789274100000/1789283100000` 实为 12:35/15:05 +08:00（12:35 恰为事故会话始发分钟），与 d-epoch 规范表（17:25/19:25）、本票回放场景、手工成本核算、部署产物 mtime（`.opencode/agent/*.md` = 2026-09-13 19:25:19 +0800，含 `steps: 200` 与 `## Non-negotiables`）四路证据矛盾——按仓库事实就地更正为 `<1789291500000` / `1789291500000-1789298700000` / `>1789298700000` 并在 Verified contracts 留痕，authority 含义（墙钟边界）不变。(2) design.md 运行手册示例命令携带同一错误字面量，在本票写集外，待设计 owner 同步更正（d-epoch 规范表为准，不阻塞）。(3) `--price-override`「整表替换」钉定为完整 [in,out,cacheR] 三元组逐档覆盖、未提及档位保留代理常量——与 #01 已交付已 review 语义及本票 Change 3「覆盖」措辞一致。除此之外无权威偏差。
- 实现参考：本次提交（`scripts/executor_metrics.py`、`scripts/executor_metrics_test.py`、`docs/proposals/executor-value-measurement/report.md`、本票），工作树 fixed point。
