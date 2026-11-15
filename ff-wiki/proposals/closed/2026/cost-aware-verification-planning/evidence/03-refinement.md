---
flowforge:
  schema: 1
  role: evidence
  id: cost-aware-verification-planning-03-refinement
  revision: 1
---

# 03 作者预设回归与执行合同核实

本次只完成03作者前置测试Change与五节Refine。消费 [requirements revision 2](../requirements.md#test-authority)、[test-authority design revision 1](../design/test-authority.md) 和 [03 repair](../issues/03-test-authority-boundary-repair.md)。主会话明确授权本运行作为测试作者；本mjs不匹配受保护glob，未把01被拒绝Go测试换格式，未应用01补丁、修改生产代码/Go测试、模型/guard配置、部署或push，也未向实际项目签发授权。

## 固定输入与覆盖

- 当前生产实现HEAD：`6c7c2667b59237d82319b3e8904e36053cd81607`。03票发布时的7cb27d8/6061b6c为历史版本；本次按实际HEAD执行，不冒称在旧checkout执行。完整最终审查范围仍为 `f1a51117b0baf5af6801238b2a12206ad6697b29...修后HEAD`，由主会话在实施后组织。
- 实跑环境：Node v22.16.0，安装Pi 1.5.0；包根 `/Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`。真实loader/runner只发tool_call事件，未调用write/edit底层实现；临时fixture含授权JSON及符号链接，finally清理并恢复cwd/env。
- source、embedded、deployed三个文件 `assets/pi/flowforge.ts`、`internal/command/assets/pi/flowforge.ts`、`.pi/extensions/flowforge.ts` sha256均为 `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153`，本次未改变。通过 `shasum -a 256` 实测；相同hash仅确定输入，不能代替native运行验收。
- 扩充probe sha256：`970f7b68881b8a18a3be627192b6537cf8aba584df806082072cb79c2227a283`；修前probe sha256：`8982aef59489971ce1d954f2b8a72910e3790e8a25948d744a5fd6d26d1559b7`。

| 矩阵 | source通过/总数 | deployed通过/总数 | 失败 |
| --- | --- | --- | --- |
| 原38观察 | 38/38 | 38/38 | 0 |
| dangling final/parent × 项目内/外 × 有/无grant × relative/absolute × write/edit | 16/32 | 16/32 | 各16 |
| 外部grant zone：四fixture × (18保护目标 × write/edit + 4配对观察) | 56/160 | 56/160 | 各104 |
| UTC：3个合法未来正例 + 15负例，均write/edit | 14/36 | 14/36 | 各22 |
| 总计 | 124/266 | 124/266 | 各142，exit 1 |

所有失败都是应当block=true而实际false；没有通过删除或翻转原断言获得red。原`expect`完全不变，新`observe`只捕获AssertionError以继续独立case；末尾`assert.equal(failures.length, 0)`抛出非零退出。loader/fixture错误仍立即失败。

Grant-zone四fixture覆盖 .flowforge 指向初始既有/尚未创建的项目外目录 × 有/无grant。无grant/missing整轮保持物理parent未创建；有grant/missing在start之后由fixture显式创建物理目录及记录，因此这一象限观察的是从dangling到存在的变化，不能称“缺失目录里存在有效记录”。18目标包括lexical自己的/其他missing记录、目录、父目录、项目根，以及physical自己的/其他missing记录、目录、父/祖父，zone/root/parent三层alias的记录、目录、父级。普通源码与精确未建测试正/负配对避免把一律拒绝当修复。

UTC用运行时未来年，包含完整秒Z、毫秒Z及合法闰日；拒绝year/date-only、两种缺Z、+08:00、+00:00、缺秒、Feb30、非闰Feb29、Apr31、month13、day0、hour24归一化、minute60与过期UTC。普通未建多层父目录在有效grant下仍应允许。

首次探索性执行在`zone-missing-grant`的fixture准备中遇到`mkdirSync(lexical/test-author-grants) ENOENT`，因为mkdir不能沿dangling链接创建目的地；那次未覆盖完矩阵，不计作完整red。只修正fixture顺序为先创建物理目的地再写grant，以下source/deployed均重新完整实跑至最终汇总和非零退出。

## 旧38保留证明

从上述固定HEAD `git show HEAD:<probe>` 读取原文件，比较从 `try {\n  const a = await start(join(scratch, 'a'));` 到旧最后的38项汇总console为止的整个块（含invalid数组及循环）在新文件逐字存在；另比较整个`expect`函数逐字存在。原38调用、断言及辅助状态检查未改写。运行时先执行原块并断言checks=38，才开始新矩阵。

实际Python字节比较及SHA256结果（exit0；同批 `node --check <probe>` exit0）：

```json
{
  "timestamp": "2026-11-05T11:05:38.734653+00:00",
  "original_block_bytes": 4257,
  "original_block_sha256": "846d97463f741f7a9e393e379a3c8e4406030f1278ccb7636a30abc3624a55d8",
  "expect_helper_bytes": 495,
  "expect_helper_sha256": "30cc47bdee207c5d3687eb1e2210da18e37dea4edb78a931116afe15c2a41948",
  "old_probe_sha256": "8982aef59489971ce1d954f2b8a72910e3790e8a25948d744a5fd6d26d1559b7",
  "new_probe_sha256": "970f7b68881b8a18a3be627192b6537cf8aba584df806082072cb79c2227a283",
  "original_block_and_helper_exact_match": true
}
```

可复核命令（固定提交，不依赖未来HEAD）：

```sh
python3 - <<'PY'
import subprocess, pathlib, hashlib
p = 'ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs'
old = subprocess.check_output(['git', 'show', '6c7c2667b59237d82319b3e8904e36053cd81607:' + p]).decode()
new = pathlib.Path(p).read_text()
start = "try {\n  const a = await start(join(scratch, 'a'));"
end = '  console.log(`PASS: ${checks} guard observations; no built-in write/edit executed`);'
block = old[old.index(start):old.index(end) + len(end)]
helper = old[old.index('async function expect('):old.index(start)]
assert block in new and helper in new
for name, data in [('original-block', block), ('expect-helper', helper), ('probe', new)]:
    print(name, len(data.encode()), hashlib.sha256(data.encode()).hexdigest())
PY
```

## 五节核实与派发边界

- 核实源码接缝：`resolveThroughLinks:196` realpath退回；`evaluateTestAuthorGrant:264`/`Date.parse:312` 宽格式；`testWriteDenial:333` 既有stat/resolved判定；`touchesGrantRecords:225` 与 factory `grantDir:443` 只保护lexical；`tool_call:508` 的异常拒绝出口。
- 生成/消费实链：`Makefile:15 dev` 复制assets、构建二进制；`internal/command/agents.go:358 deployPiExtension` 部署共享扩展；后续新背景运行消费。合同规定修完Changes1–3一次source验证、Change4部署后deployed及Go复验、保存生成差异和hash。未让实现者自行裁剪验证。
- 五节已填，模式仍lightweight、状态仍open、Changes未勾选；计划发布时“尚未Refine”的Planning evidence保留为历史快照，当前事实以五节和本文为准。未改写需求/设计/验收或Write set。
- 后续真实native作者握手/正常write/edit、主会话签发与撤销、完整授权实现Round0及Standards/Spec双轴仍待修后执行。本次只证明预设反例能暴露已定三项partial，不宣称修复或关闭03/00。

执行 `./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full` 后执行 `./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json`，exit0。03 gap已清除，01/02遗留项原样保留：

```text
Checked 4 issues in ff-wiki/proposals/cost-aware-verification-planning
warning: blocked-evidence-present ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md — Open ticket carries blocked evidence; run flowforge-refine-ticket to consume it before redispatch
gap: execution-contract-incomplete ff-wiki/proposals/cost-aware-verification-planning/issues/02-controlled-acceptance-and-deploy.md — Open ticket is missing a complete machine execution contract
✓ Dependency graph is healthy. No cycles or dangling references found.
{"blocked":[{"id":"01","feature":"cost-aware-verification-planning","title":"规则实现及必要保护性测试","waiting_on":["03"]},{"id":"02","feature":"cost-aware-verification-planning","title":"受控情景验收、源/嵌入/部署一致性及本项目部署","waiting_on":["01"]}],"claimed":null,"detail":"pass --full for complete bodies and diagnostics","diagnostics":{"counts":{"gap:execution-contract-incomplete":1,"warning:blocked-evidence-present":1},"blockers":null,"sample_warnings":[{"severity":"warning","code":"blocked-evidence-present","artifact":"ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md"},{"severity":"gap","code":"execution-contract-incomplete","artifact":"ff-wiki/proposals/cost-aware-verification-planning/issues/02-controlled-acceptance-and-deploy.md"}]},"gaps":null,"ready":[{"id":"03","feature":"cost-aware-verification-planning","file_path":"ff-wiki/proposals/cost-aware-verification-planning/issues/03-test-authority-boundary-repair.md"}],"ready_with_warnings":null}
```

初始 `check --strict` exit1，包含03/02两个execution-contract gap及01 blocked-evidence；Refine后03已清除，不能据此声称整个proposal strict通过。03→01→02关系未变，01原测试作者工作仍待后续授权。

## Go回归实际结果

命令：`GOPROXY=https://goproxy.cn,direct go test ./internal/...`，exit 1。唯一失败测试为既有 `TestIssueTrackerDescribesRepairWorkflow`，仍归01 Change8；不是全仓绿。其余package结果及错误摘录：

```text
--- FAIL: TestIssueTrackerDescribesRepairWorkflow (0.00s)
    assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory
FAIL	flowforge/internal/command	8.463s
ok  	flowforge/internal/config	(cached)
ok  	flowforge/internal/subagent	(cached)
ok  	flowforge/internal/tracker	0.594s
ok  	flowforge/internal/update	(cached)
?   	flowforge/internal/version	[no test files]
```

## Source 原始完整输出

命令：`node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`。exit 1；tool输出未截断。

<details>
<summary>266项观察与142项聚合失败</summary>

```text
PASS 1: default write denied
PASS 2: default edit denied
PASS 3: ordinary source preserved
PASS 4: child cannot create its own grant
PASS 5: child cannot overwrite grant directory
PASS 6: child cannot overwrite grant ancestor
PASS 7: authorized Plan write
PASS 8: authorized Plan edit
PASS 9: absolute exact target
PASS 10: unlisted protected target denied
PASS 11: outside project denied
PASS 12: authorized child cannot enlarge grant
PASS 13: authorized Refine write
PASS 14: implementer role denied
PASS 15: unknown field denied
PASS 16: empty paths denied
PASS 17: wildcard paths denied
PASS 18: absolute grant path denied
PASS 19: traversal grant denied
PASS 20: dot segment denied
PASS 21: wrong instance denied
PASS 22: wrong session denied
PASS 23: wrong project denied
PASS 24: expired denied
PASS 25: invalid expiry denied
PASS 26: version mismatch denied
PASS 27: malformed grant denied
PASS 28: revocation immediate
PASS 29: symlink parent denied
PASS 30: symlink target denied
PASS 31: non-test alias cannot reach protected target
PASS 32: grant directory alias denied
PASS 33: parallel environment changes do not alter A
PASS 34: parallel B cannot inherit broader rights
PASS 35: shutdown revokes instance
PASS 36: same session restart rejects old grant
PASS 37: controller can issue grant records
PASS 38: controller does not gain child-author exemption
PASS: 38 guard observations; no built-in write/edit executed
HOST: Node v22.16.0; Pi 1.5.0; source /Users/qiangbi/develop/projects/Syl/tangram/flowforge/assets/pi/flowforge.ts
PASS 39: 03/dangling/final/inside/no-grant/relative/write
PASS 40: 03/dangling/final/inside/no-grant/relative/edit
PASS 41: 03/dangling/final/inside/no-grant/absolute/write
PASS 42: 03/dangling/final/inside/no-grant/absolute/edit
FAIL: {"group":"dangling","label":"final/inside/grant/relative","toolName":"write","target":"dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/inside/grant/relative/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/inside/grant/relative","toolName":"edit","target":"dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/inside/grant/relative/edit\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/inside/grant/absolute","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/dangling-final-inside/dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/inside/grant/absolute/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/inside/grant/absolute","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/dangling-final-inside/dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/inside/grant/absolute/edit\n\nfalse !== true\n"}
PASS 43: 03/dangling/final/outside/no-grant/relative/write
PASS 44: 03/dangling/final/outside/no-grant/relative/edit
PASS 45: 03/dangling/final/outside/no-grant/absolute/write
PASS 46: 03/dangling/final/outside/no-grant/absolute/edit
FAIL: {"group":"dangling","label":"final/outside/grant/relative","toolName":"write","target":"dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/outside/grant/relative/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/outside/grant/relative","toolName":"edit","target":"dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/outside/grant/relative/edit\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/outside/grant/absolute","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/dangling-final-outside/dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/outside/grant/absolute/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/outside/grant/absolute","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/dangling-final-outside/dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/outside/grant/absolute/edit\n\nfalse !== true\n"}
PASS 47: 03/dangling/parent/inside/no-grant/relative/write
PASS 48: 03/dangling/parent/inside/no-grant/relative/edit
PASS 49: 03/dangling/parent/inside/no-grant/absolute/write
PASS 50: 03/dangling/parent/inside/no-grant/absolute/edit
FAIL: {"group":"dangling","label":"parent/inside/grant/relative","toolName":"write","target":"dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/inside/grant/relative/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/inside/grant/relative","toolName":"edit","target":"dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/inside/grant/relative/edit\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/inside/grant/absolute","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/dangling-parent-inside/dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/inside/grant/absolute/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/inside/grant/absolute","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/dangling-parent-inside/dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/inside/grant/absolute/edit\n\nfalse !== true\n"}
PASS 51: 03/dangling/parent/outside/no-grant/relative/write
PASS 52: 03/dangling/parent/outside/no-grant/relative/edit
PASS 53: 03/dangling/parent/outside/no-grant/absolute/write
PASS 54: 03/dangling/parent/outside/no-grant/absolute/edit
FAIL: {"group":"dangling","label":"parent/outside/grant/relative","toolName":"write","target":"dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/outside/grant/relative/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/outside/grant/relative","toolName":"edit","target":"dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/outside/grant/relative/edit\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/outside/grant/absolute","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/dangling-parent-outside/dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/outside/grant/absolute/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/outside/grant/absolute","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/dangling-parent-outside/dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/outside/grant/absolute/edit\n\nfalse !== true\n"}
PASS 55: 03/grant-zone/zone-existing-no-grant/lexical-record/write
PASS 56: 03/grant-zone/zone-existing-no-grant/lexical-record/edit
PASS 57: 03/grant-zone/zone-existing-no-grant/lexical-missing-record/write
PASS 58: 03/grant-zone/zone-existing-no-grant/lexical-missing-record/edit
PASS 59: 03/grant-zone/zone-existing-no-grant/lexical-zone/write
PASS 60: 03/grant-zone/zone-existing-no-grant/lexical-zone/edit
PASS 61: 03/grant-zone/zone-existing-no-grant/lexical-parent/write
PASS 62: 03/grant-zone/zone-existing-no-grant/lexical-parent/edit
PASS 63: 03/grant-zone/zone-existing-no-grant/lexical-project-ancestor/write
PASS 64: 03/grant-zone/zone-existing-no-grant/lexical-project-ancestor/edit
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent/config/test-author-grants/9ab1332b-dadd-4f6e-b5fd-3915b98da835.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent/config/test-author-grants/9ab1332b-dadd-4f6e-b5fd-3915b98da835.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-missing-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-missing-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-zone","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-zone","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-parent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-parent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-grandparent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-grandparent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-grandparent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-no-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-grandparent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-record","toolName":"write","target":"zone-alias/9ab1332b-dadd-4f6e-b5fd-3915b98da835.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-record","toolName":"edit","target":"zone-alias/9ab1332b-dadd-4f6e-b5fd-3915b98da835.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-missing-record","toolName":"write","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-missing-record","toolName":"edit","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-zone","toolName":"write","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-zone","toolName":"edit","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-record","toolName":"write","target":"root-alias/test-author-grants/9ab1332b-dadd-4f6e-b5fd-3915b98da835.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-record","toolName":"edit","target":"root-alias/test-author-grants/9ab1332b-dadd-4f6e-b5fd-3915b98da835.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-zone","toolName":"write","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-zone","toolName":"edit","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-parent","toolName":"write","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-parent","toolName":"edit","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/parent-alias-record","toolName":"write","target":"parent-alias/config/test-author-grants/9ab1332b-dadd-4f6e-b5fd-3915b98da835.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/parent-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/parent-alias-record","toolName":"edit","target":"parent-alias/config/test-author-grants/9ab1332b-dadd-4f6e-b5fd-3915b98da835.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/parent-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/parent-alias-ancestor","toolName":"write","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/parent-alias-ancestor/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/parent-alias-ancestor","toolName":"edit","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/parent-alias-ancestor/edit\n\nfalse !== true\n"}
PASS 65: 03/grant-zone/zone-existing-no-grant/ordinary-source/write
PASS 66: 03/grant-zone/zone-existing-no-grant/exact-missing-file/write
PASS 67: 03/grant-zone/zone-existing-no-grant/ordinary-source/edit
PASS 68: 03/grant-zone/zone-existing-no-grant/exact-missing-file/edit
PASS 69: 03/grant-zone/zone-existing-grant/lexical-record/write
PASS 70: 03/grant-zone/zone-existing-grant/lexical-record/edit
PASS 71: 03/grant-zone/zone-existing-grant/lexical-missing-record/write
PASS 72: 03/grant-zone/zone-existing-grant/lexical-missing-record/edit
PASS 73: 03/grant-zone/zone-existing-grant/lexical-zone/write
PASS 74: 03/grant-zone/zone-existing-grant/lexical-zone/edit
PASS 75: 03/grant-zone/zone-existing-grant/lexical-parent/write
PASS 76: 03/grant-zone/zone-existing-grant/lexical-parent/edit
PASS 77: 03/grant-zone/zone-existing-grant/lexical-project-ancestor/write
PASS 78: 03/grant-zone/zone-existing-grant/lexical-project-ancestor/edit
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent/config/test-author-grants/963b6f11-40b3-4b10-9ad2-95d80bc4e27c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent/config/test-author-grants/963b6f11-40b3-4b10-9ad2-95d80bc4e27c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-missing-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-missing-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-zone","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-zone","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-parent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-parent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-grandparent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-grandparent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-grandparent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-existing-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-grandparent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-record","toolName":"write","target":"zone-alias/963b6f11-40b3-4b10-9ad2-95d80bc4e27c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-record","toolName":"edit","target":"zone-alias/963b6f11-40b3-4b10-9ad2-95d80bc4e27c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-missing-record","toolName":"write","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-missing-record","toolName":"edit","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-zone","toolName":"write","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-zone","toolName":"edit","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-record","toolName":"write","target":"root-alias/test-author-grants/963b6f11-40b3-4b10-9ad2-95d80bc4e27c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-record","toolName":"edit","target":"root-alias/test-author-grants/963b6f11-40b3-4b10-9ad2-95d80bc4e27c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-zone","toolName":"write","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-zone","toolName":"edit","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-parent","toolName":"write","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-parent","toolName":"edit","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/parent-alias-record","toolName":"write","target":"parent-alias/config/test-author-grants/963b6f11-40b3-4b10-9ad2-95d80bc4e27c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/parent-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/parent-alias-record","toolName":"edit","target":"parent-alias/config/test-author-grants/963b6f11-40b3-4b10-9ad2-95d80bc4e27c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/parent-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/parent-alias-ancestor","toolName":"write","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/parent-alias-ancestor/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/parent-alias-ancestor","toolName":"edit","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/parent-alias-ancestor/edit\n\nfalse !== true\n"}
PASS 79: 03/grant-zone/zone-existing-grant/ordinary-source/write
PASS 80: 03/grant-zone/zone-existing-grant/exact-missing-file/write
PASS 81: 03/grant-zone/zone-existing-grant/ordinary-source/edit
PASS 82: 03/grant-zone/zone-existing-grant/exact-missing-file/edit
PASS 83: 03/grant-zone/zone-missing-no-grant/lexical-record/write
PASS 84: 03/grant-zone/zone-missing-no-grant/lexical-record/edit
PASS 85: 03/grant-zone/zone-missing-no-grant/lexical-missing-record/write
PASS 86: 03/grant-zone/zone-missing-no-grant/lexical-missing-record/edit
PASS 87: 03/grant-zone/zone-missing-no-grant/lexical-zone/write
PASS 88: 03/grant-zone/zone-missing-no-grant/lexical-zone/edit
PASS 89: 03/grant-zone/zone-missing-no-grant/lexical-parent/write
PASS 90: 03/grant-zone/zone-missing-no-grant/lexical-parent/edit
PASS 91: 03/grant-zone/zone-missing-no-grant/lexical-project-ancestor/write
PASS 92: 03/grant-zone/zone-missing-no-grant/lexical-project-ancestor/edit
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent/config/test-author-grants/4fb9e655-c674-4fb2-b15c-3891a25d385f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent/config/test-author-grants/4fb9e655-c674-4fb2-b15c-3891a25d385f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-missing-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-missing-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-zone","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-zone","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-parent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-parent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-grandparent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-grandparent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-grandparent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-no-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-grandparent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-record","toolName":"write","target":"zone-alias/4fb9e655-c674-4fb2-b15c-3891a25d385f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-record","toolName":"edit","target":"zone-alias/4fb9e655-c674-4fb2-b15c-3891a25d385f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-missing-record","toolName":"write","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-missing-record","toolName":"edit","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-zone","toolName":"write","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-zone","toolName":"edit","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-record","toolName":"write","target":"root-alias/test-author-grants/4fb9e655-c674-4fb2-b15c-3891a25d385f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-record","toolName":"edit","target":"root-alias/test-author-grants/4fb9e655-c674-4fb2-b15c-3891a25d385f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-zone","toolName":"write","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-zone","toolName":"edit","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-parent","toolName":"write","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-parent","toolName":"edit","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/parent-alias-record","toolName":"write","target":"parent-alias/config/test-author-grants/4fb9e655-c674-4fb2-b15c-3891a25d385f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/parent-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/parent-alias-record","toolName":"edit","target":"parent-alias/config/test-author-grants/4fb9e655-c674-4fb2-b15c-3891a25d385f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/parent-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/parent-alias-ancestor","toolName":"write","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/parent-alias-ancestor/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/parent-alias-ancestor","toolName":"edit","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/parent-alias-ancestor/edit\n\nfalse !== true\n"}
PASS 93: 03/grant-zone/zone-missing-no-grant/ordinary-source/write
PASS 94: 03/grant-zone/zone-missing-no-grant/exact-missing-file/write
PASS 95: 03/grant-zone/zone-missing-no-grant/ordinary-source/edit
PASS 96: 03/grant-zone/zone-missing-no-grant/exact-missing-file/edit
PASS 97: 03/grant-zone/zone-missing-grant/lexical-record/write
PASS 98: 03/grant-zone/zone-missing-grant/lexical-record/edit
PASS 99: 03/grant-zone/zone-missing-grant/lexical-missing-record/write
PASS 100: 03/grant-zone/zone-missing-grant/lexical-missing-record/edit
PASS 101: 03/grant-zone/zone-missing-grant/lexical-zone/write
PASS 102: 03/grant-zone/zone-missing-grant/lexical-zone/edit
PASS 103: 03/grant-zone/zone-missing-grant/lexical-parent/write
PASS 104: 03/grant-zone/zone-missing-grant/lexical-parent/edit
PASS 105: 03/grant-zone/zone-missing-grant/lexical-project-ancestor/write
PASS 106: 03/grant-zone/zone-missing-grant/lexical-project-ancestor/edit
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent/config/test-author-grants/c21a0c15-efab-4a17-9d10-7566c4b9d70f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent/config/test-author-grants/c21a0c15-efab-4a17-9d10-7566c4b9d70f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-missing-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-missing-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-zone","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-zone","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-parent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-parent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-grandparent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-grandparent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-grandparent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-4tgaA3/zone-missing-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-grandparent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-record","toolName":"write","target":"zone-alias/c21a0c15-efab-4a17-9d10-7566c4b9d70f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-record","toolName":"edit","target":"zone-alias/c21a0c15-efab-4a17-9d10-7566c4b9d70f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-missing-record","toolName":"write","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-missing-record","toolName":"edit","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-zone","toolName":"write","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-zone","toolName":"edit","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-record","toolName":"write","target":"root-alias/test-author-grants/c21a0c15-efab-4a17-9d10-7566c4b9d70f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-record","toolName":"edit","target":"root-alias/test-author-grants/c21a0c15-efab-4a17-9d10-7566c4b9d70f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-zone","toolName":"write","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-zone","toolName":"edit","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-parent","toolName":"write","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-parent","toolName":"edit","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/parent-alias-record","toolName":"write","target":"parent-alias/config/test-author-grants/c21a0c15-efab-4a17-9d10-7566c4b9d70f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/parent-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/parent-alias-record","toolName":"edit","target":"parent-alias/config/test-author-grants/c21a0c15-efab-4a17-9d10-7566c4b9d70f.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/parent-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/parent-alias-ancestor","toolName":"write","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/parent-alias-ancestor/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/parent-alias-ancestor","toolName":"edit","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/parent-alias-ancestor/edit\n\nfalse !== true\n"}
PASS 107: 03/grant-zone/zone-missing-grant/ordinary-source/write
PASS 108: 03/grant-zone/zone-missing-grant/exact-missing-file/write
PASS 109: 03/grant-zone/zone-missing-grant/ordinary-source/edit
PASS 110: 03/grant-zone/zone-missing-grant/exact-missing-file/edit
PASS 111: 03/UTC/utc-seconds/write
PASS 112: 03/UTC/utc-seconds/edit
PASS 113: 03/UTC/utc-milliseconds/write
PASS 114: 03/UTC/utc-milliseconds/edit
PASS 115: 03/UTC/utc-leap-date/write
PASS 116: 03/UTC/utc-leap-date/edit
PASS 117: 03/UTC/expired-utc/write
PASS 118: 03/UTC/expired-utc/edit
FAIL: {"group":"UTC","label":"year-only","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/year-only/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"year-only","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/year-only/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"date-only","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/date-only/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"date-only","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/date-only/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-Z","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-Z/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-Z","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-Z/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-Z-milliseconds","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-Z-milliseconds/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-Z-milliseconds","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-Z-milliseconds/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"non-UTC-offset","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/non-UTC-offset/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"non-UTC-offset","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/non-UTC-offset/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"zero-offset-without-Z","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/zero-offset-without-Z/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"zero-offset-without-Z","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/zero-offset-without-Z/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-seconds","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-seconds/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-seconds","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-seconds/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"impossible-Feb-30","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/impossible-Feb-30/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"impossible-Feb-30","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/impossible-Feb-30/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"non-leap-Feb-29","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/non-leap-Feb-29/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"non-leap-Feb-29","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/non-leap-Feb-29/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"impossible-Apr-31","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/impossible-Apr-31/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"impossible-Apr-31","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/impossible-Apr-31/edit\n\nfalse !== true\n"}
PASS 119: 03/UTC/impossible-month-13/write
PASS 120: 03/UTC/impossible-month-13/edit
PASS 121: 03/UTC/impossible-day-zero/write
PASS 122: 03/UTC/impossible-day-zero/edit
FAIL: {"group":"UTC","label":"normalized-hour-24","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/normalized-hour-24/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"normalized-hour-24","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/normalized-hour-24/edit\n\nfalse !== true\n"}
PASS 123: 03/UTC/impossible-minute-60/write
PASS 124: 03/UTC/impossible-minute-60/edit
GROUP dangling: {"total":32,"passed":16,"failed":16}
GROUP grant-zone: {"total":160,"passed":56,"failed":104}
GROUP UTC: {"total":36,"passed":14,"failed":22}
SUMMARY: original=38/38; repair=228; failures=142; passed=124; no built-in write/edit executed
node:internal/modules/run_main:123
    triggerUncaughtException(
    ^

AssertionError [ERR_ASSERTION]: repair boundary observations must all pass (aggregated failures above)

142 !== 0

    at runRepairMatrix (file:///Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs:213:10)
    at async file:///Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs:283:3 {
  generatedMessage: false,
  code: 'ERR_ASSERTION',
  actual: 142,
  expected: 0,
  operator: 'strictEqual'
}

Node.js v22.16.0
```

</details>

## Deployed 原始完整输出

命令：`node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent .pi/extensions/flowforge.ts`。exit 1；tool输出未截断。

<details>
<summary>266项观察与142项聚合失败</summary>

```text
PASS 1: default write denied
PASS 2: default edit denied
PASS 3: ordinary source preserved
PASS 4: child cannot create its own grant
PASS 5: child cannot overwrite grant directory
PASS 6: child cannot overwrite grant ancestor
PASS 7: authorized Plan write
PASS 8: authorized Plan edit
PASS 9: absolute exact target
PASS 10: unlisted protected target denied
PASS 11: outside project denied
PASS 12: authorized child cannot enlarge grant
PASS 13: authorized Refine write
PASS 14: implementer role denied
PASS 15: unknown field denied
PASS 16: empty paths denied
PASS 17: wildcard paths denied
PASS 18: absolute grant path denied
PASS 19: traversal grant denied
PASS 20: dot segment denied
PASS 21: wrong instance denied
PASS 22: wrong session denied
PASS 23: wrong project denied
PASS 24: expired denied
PASS 25: invalid expiry denied
PASS 26: version mismatch denied
PASS 27: malformed grant denied
PASS 28: revocation immediate
PASS 29: symlink parent denied
PASS 30: symlink target denied
PASS 31: non-test alias cannot reach protected target
PASS 32: grant directory alias denied
PASS 33: parallel environment changes do not alter A
PASS 34: parallel B cannot inherit broader rights
PASS 35: shutdown revokes instance
PASS 36: same session restart rejects old grant
PASS 37: controller can issue grant records
PASS 38: controller does not gain child-author exemption
PASS: 38 guard observations; no built-in write/edit executed
HOST: Node v22.16.0; Pi 1.5.0; source /Users/qiangbi/develop/projects/Syl/tangram/flowforge/.pi/extensions/flowforge.ts
PASS 39: 03/dangling/final/inside/no-grant/relative/write
PASS 40: 03/dangling/final/inside/no-grant/relative/edit
PASS 41: 03/dangling/final/inside/no-grant/absolute/write
PASS 42: 03/dangling/final/inside/no-grant/absolute/edit
FAIL: {"group":"dangling","label":"final/inside/grant/relative","toolName":"write","target":"dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/inside/grant/relative/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/inside/grant/relative","toolName":"edit","target":"dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/inside/grant/relative/edit\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/inside/grant/absolute","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/dangling-final-inside/dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/inside/grant/absolute/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/inside/grant/absolute","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/dangling-final-inside/dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/inside/grant/absolute/edit\n\nfalse !== true\n"}
PASS 43: 03/dangling/final/outside/no-grant/relative/write
PASS 44: 03/dangling/final/outside/no-grant/relative/edit
PASS 45: 03/dangling/final/outside/no-grant/absolute/write
PASS 46: 03/dangling/final/outside/no-grant/absolute/edit
FAIL: {"group":"dangling","label":"final/outside/grant/relative","toolName":"write","target":"dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/outside/grant/relative/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/outside/grant/relative","toolName":"edit","target":"dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/outside/grant/relative/edit\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/outside/grant/absolute","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/dangling-final-outside/dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/outside/grant/absolute/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"final/outside/grant/absolute","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/dangling-final-outside/dangling_test.go","expected":true,"actual":false,"message":"03/dangling/final/outside/grant/absolute/edit\n\nfalse !== true\n"}
PASS 47: 03/dangling/parent/inside/no-grant/relative/write
PASS 48: 03/dangling/parent/inside/no-grant/relative/edit
PASS 49: 03/dangling/parent/inside/no-grant/absolute/write
PASS 50: 03/dangling/parent/inside/no-grant/absolute/edit
FAIL: {"group":"dangling","label":"parent/inside/grant/relative","toolName":"write","target":"dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/inside/grant/relative/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/inside/grant/relative","toolName":"edit","target":"dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/inside/grant/relative/edit\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/inside/grant/absolute","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/dangling-parent-inside/dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/inside/grant/absolute/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/inside/grant/absolute","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/dangling-parent-inside/dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/inside/grant/absolute/edit\n\nfalse !== true\n"}
PASS 51: 03/dangling/parent/outside/no-grant/relative/write
PASS 52: 03/dangling/parent/outside/no-grant/relative/edit
PASS 53: 03/dangling/parent/outside/no-grant/absolute/write
PASS 54: 03/dangling/parent/outside/no-grant/absolute/edit
FAIL: {"group":"dangling","label":"parent/outside/grant/relative","toolName":"write","target":"dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/outside/grant/relative/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/outside/grant/relative","toolName":"edit","target":"dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/outside/grant/relative/edit\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/outside/grant/absolute","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/dangling-parent-outside/dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/outside/grant/absolute/write\n\nfalse !== true\n"}
FAIL: {"group":"dangling","label":"parent/outside/grant/absolute","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/dangling-parent-outside/dangling-parent/leaf_test.go","expected":true,"actual":false,"message":"03/dangling/parent/outside/grant/absolute/edit\n\nfalse !== true\n"}
PASS 55: 03/grant-zone/zone-existing-no-grant/lexical-record/write
PASS 56: 03/grant-zone/zone-existing-no-grant/lexical-record/edit
PASS 57: 03/grant-zone/zone-existing-no-grant/lexical-missing-record/write
PASS 58: 03/grant-zone/zone-existing-no-grant/lexical-missing-record/edit
PASS 59: 03/grant-zone/zone-existing-no-grant/lexical-zone/write
PASS 60: 03/grant-zone/zone-existing-no-grant/lexical-zone/edit
PASS 61: 03/grant-zone/zone-existing-no-grant/lexical-parent/write
PASS 62: 03/grant-zone/zone-existing-no-grant/lexical-parent/edit
PASS 63: 03/grant-zone/zone-existing-no-grant/lexical-project-ancestor/write
PASS 64: 03/grant-zone/zone-existing-no-grant/lexical-project-ancestor/edit
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent/config/test-author-grants/94e63fa4-8cf1-4b0e-a68f-ddef9c66b24c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent/config/test-author-grants/94e63fa4-8cf1-4b0e-a68f-ddef9c66b24c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-missing-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-missing-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-zone","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-zone","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-parent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-parent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-grandparent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-grandparent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/physical-grandparent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-no-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/physical-grandparent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-record","toolName":"write","target":"zone-alias/94e63fa4-8cf1-4b0e-a68f-ddef9c66b24c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-record","toolName":"edit","target":"zone-alias/94e63fa4-8cf1-4b0e-a68f-ddef9c66b24c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-missing-record","toolName":"write","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-missing-record","toolName":"edit","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-zone","toolName":"write","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/alias-zone","toolName":"edit","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-record","toolName":"write","target":"root-alias/test-author-grants/94e63fa4-8cf1-4b0e-a68f-ddef9c66b24c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-record","toolName":"edit","target":"root-alias/test-author-grants/94e63fa4-8cf1-4b0e-a68f-ddef9c66b24c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-zone","toolName":"write","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-zone","toolName":"edit","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-parent","toolName":"write","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/root-alias-parent","toolName":"edit","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/root-alias-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/parent-alias-record","toolName":"write","target":"parent-alias/config/test-author-grants/94e63fa4-8cf1-4b0e-a68f-ddef9c66b24c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/parent-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/parent-alias-record","toolName":"edit","target":"parent-alias/config/test-author-grants/94e63fa4-8cf1-4b0e-a68f-ddef9c66b24c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/parent-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/parent-alias-ancestor","toolName":"write","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/parent-alias-ancestor/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-no-grant/parent-alias-ancestor","toolName":"edit","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-no-grant/parent-alias-ancestor/edit\n\nfalse !== true\n"}
PASS 65: 03/grant-zone/zone-existing-no-grant/ordinary-source/write
PASS 66: 03/grant-zone/zone-existing-no-grant/exact-missing-file/write
PASS 67: 03/grant-zone/zone-existing-no-grant/ordinary-source/edit
PASS 68: 03/grant-zone/zone-existing-no-grant/exact-missing-file/edit
PASS 69: 03/grant-zone/zone-existing-grant/lexical-record/write
PASS 70: 03/grant-zone/zone-existing-grant/lexical-record/edit
PASS 71: 03/grant-zone/zone-existing-grant/lexical-missing-record/write
PASS 72: 03/grant-zone/zone-existing-grant/lexical-missing-record/edit
PASS 73: 03/grant-zone/zone-existing-grant/lexical-zone/write
PASS 74: 03/grant-zone/zone-existing-grant/lexical-zone/edit
PASS 75: 03/grant-zone/zone-existing-grant/lexical-parent/write
PASS 76: 03/grant-zone/zone-existing-grant/lexical-parent/edit
PASS 77: 03/grant-zone/zone-existing-grant/lexical-project-ancestor/write
PASS 78: 03/grant-zone/zone-existing-grant/lexical-project-ancestor/edit
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent/config/test-author-grants/a2e1246d-9af7-4dcf-8a47-0cd3fcf049a1.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent/config/test-author-grants/a2e1246d-9af7-4dcf-8a47-0cd3fcf049a1.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-missing-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-missing-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-zone","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-zone","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-parent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-parent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-grandparent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-grandparent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/physical-grandparent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-existing-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/physical-grandparent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-record","toolName":"write","target":"zone-alias/a2e1246d-9af7-4dcf-8a47-0cd3fcf049a1.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-record","toolName":"edit","target":"zone-alias/a2e1246d-9af7-4dcf-8a47-0cd3fcf049a1.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-missing-record","toolName":"write","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-missing-record","toolName":"edit","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-zone","toolName":"write","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/alias-zone","toolName":"edit","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-record","toolName":"write","target":"root-alias/test-author-grants/a2e1246d-9af7-4dcf-8a47-0cd3fcf049a1.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-record","toolName":"edit","target":"root-alias/test-author-grants/a2e1246d-9af7-4dcf-8a47-0cd3fcf049a1.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-zone","toolName":"write","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-zone","toolName":"edit","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-parent","toolName":"write","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/root-alias-parent","toolName":"edit","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/root-alias-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/parent-alias-record","toolName":"write","target":"parent-alias/config/test-author-grants/a2e1246d-9af7-4dcf-8a47-0cd3fcf049a1.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/parent-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/parent-alias-record","toolName":"edit","target":"parent-alias/config/test-author-grants/a2e1246d-9af7-4dcf-8a47-0cd3fcf049a1.json","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/parent-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/parent-alias-ancestor","toolName":"write","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/parent-alias-ancestor/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-existing-grant/parent-alias-ancestor","toolName":"edit","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-existing-grant/parent-alias-ancestor/edit\n\nfalse !== true\n"}
PASS 79: 03/grant-zone/zone-existing-grant/ordinary-source/write
PASS 80: 03/grant-zone/zone-existing-grant/exact-missing-file/write
PASS 81: 03/grant-zone/zone-existing-grant/ordinary-source/edit
PASS 82: 03/grant-zone/zone-existing-grant/exact-missing-file/edit
PASS 83: 03/grant-zone/zone-missing-no-grant/lexical-record/write
PASS 84: 03/grant-zone/zone-missing-no-grant/lexical-record/edit
PASS 85: 03/grant-zone/zone-missing-no-grant/lexical-missing-record/write
PASS 86: 03/grant-zone/zone-missing-no-grant/lexical-missing-record/edit
PASS 87: 03/grant-zone/zone-missing-no-grant/lexical-zone/write
PASS 88: 03/grant-zone/zone-missing-no-grant/lexical-zone/edit
PASS 89: 03/grant-zone/zone-missing-no-grant/lexical-parent/write
PASS 90: 03/grant-zone/zone-missing-no-grant/lexical-parent/edit
PASS 91: 03/grant-zone/zone-missing-no-grant/lexical-project-ancestor/write
PASS 92: 03/grant-zone/zone-missing-no-grant/lexical-project-ancestor/edit
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent/config/test-author-grants/2c0604f6-9636-44e0-898f-85a2c7b86306.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent/config/test-author-grants/2c0604f6-9636-44e0-898f-85a2c7b86306.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-missing-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-missing-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-zone","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-zone","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-parent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-parent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-grandparent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-grandparent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/physical-grandparent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-no-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/physical-grandparent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-record","toolName":"write","target":"zone-alias/2c0604f6-9636-44e0-898f-85a2c7b86306.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-record","toolName":"edit","target":"zone-alias/2c0604f6-9636-44e0-898f-85a2c7b86306.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-missing-record","toolName":"write","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-missing-record","toolName":"edit","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-zone","toolName":"write","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/alias-zone","toolName":"edit","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-record","toolName":"write","target":"root-alias/test-author-grants/2c0604f6-9636-44e0-898f-85a2c7b86306.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-record","toolName":"edit","target":"root-alias/test-author-grants/2c0604f6-9636-44e0-898f-85a2c7b86306.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-zone","toolName":"write","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-zone","toolName":"edit","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-parent","toolName":"write","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/root-alias-parent","toolName":"edit","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/root-alias-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/parent-alias-record","toolName":"write","target":"parent-alias/config/test-author-grants/2c0604f6-9636-44e0-898f-85a2c7b86306.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/parent-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/parent-alias-record","toolName":"edit","target":"parent-alias/config/test-author-grants/2c0604f6-9636-44e0-898f-85a2c7b86306.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/parent-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/parent-alias-ancestor","toolName":"write","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/parent-alias-ancestor/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-no-grant/parent-alias-ancestor","toolName":"edit","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-no-grant/parent-alias-ancestor/edit\n\nfalse !== true\n"}
PASS 93: 03/grant-zone/zone-missing-no-grant/ordinary-source/write
PASS 94: 03/grant-zone/zone-missing-no-grant/exact-missing-file/write
PASS 95: 03/grant-zone/zone-missing-no-grant/ordinary-source/edit
PASS 96: 03/grant-zone/zone-missing-no-grant/exact-missing-file/edit
PASS 97: 03/grant-zone/zone-missing-grant/lexical-record/write
PASS 98: 03/grant-zone/zone-missing-grant/lexical-record/edit
PASS 99: 03/grant-zone/zone-missing-grant/lexical-missing-record/write
PASS 100: 03/grant-zone/zone-missing-grant/lexical-missing-record/edit
PASS 101: 03/grant-zone/zone-missing-grant/lexical-zone/write
PASS 102: 03/grant-zone/zone-missing-grant/lexical-zone/edit
PASS 103: 03/grant-zone/zone-missing-grant/lexical-parent/write
PASS 104: 03/grant-zone/zone-missing-grant/lexical-parent/edit
PASS 105: 03/grant-zone/zone-missing-grant/lexical-project-ancestor/write
PASS 106: 03/grant-zone/zone-missing-grant/lexical-project-ancestor/edit
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent/config/test-author-grants/0aa95217-17c6-4376-96f4-d112876aea8c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent/config/test-author-grants/0aa95217-17c6-4376-96f4-d112876aea8c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-missing-record","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-missing-record","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent/config/test-author-grants/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-zone","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-zone","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent/config/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-parent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-parent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent/config","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-grandparent","toolName":"write","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-grandparent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/physical-grandparent","toolName":"edit","target":"/private/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-authority-probe-09GBef/zone-missing-grant-outside/parent","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/physical-grandparent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-record","toolName":"write","target":"zone-alias/0aa95217-17c6-4376-96f4-d112876aea8c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-record","toolName":"edit","target":"zone-alias/0aa95217-17c6-4376-96f4-d112876aea8c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-missing-record","toolName":"write","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-missing-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-missing-record","toolName":"edit","target":"zone-alias/new.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-missing-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-zone","toolName":"write","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/alias-zone","toolName":"edit","target":"zone-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-record","toolName":"write","target":"root-alias/test-author-grants/0aa95217-17c6-4376-96f4-d112876aea8c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-record","toolName":"edit","target":"root-alias/test-author-grants/0aa95217-17c6-4376-96f4-d112876aea8c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-zone","toolName":"write","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-zone/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-zone","toolName":"edit","target":"root-alias/test-author-grants","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-zone/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-parent","toolName":"write","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-parent/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/root-alias-parent","toolName":"edit","target":"root-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/root-alias-parent/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/parent-alias-record","toolName":"write","target":"parent-alias/config/test-author-grants/0aa95217-17c6-4376-96f4-d112876aea8c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/parent-alias-record/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/parent-alias-record","toolName":"edit","target":"parent-alias/config/test-author-grants/0aa95217-17c6-4376-96f4-d112876aea8c.json","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/parent-alias-record/edit\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/parent-alias-ancestor","toolName":"write","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/parent-alias-ancestor/write\n\nfalse !== true\n"}
FAIL: {"group":"grant-zone","label":"zone-missing-grant/parent-alias-ancestor","toolName":"edit","target":"parent-alias","expected":true,"actual":false,"message":"03/grant-zone/zone-missing-grant/parent-alias-ancestor/edit\n\nfalse !== true\n"}
PASS 107: 03/grant-zone/zone-missing-grant/ordinary-source/write
PASS 108: 03/grant-zone/zone-missing-grant/exact-missing-file/write
PASS 109: 03/grant-zone/zone-missing-grant/ordinary-source/edit
PASS 110: 03/grant-zone/zone-missing-grant/exact-missing-file/edit
PASS 111: 03/UTC/utc-seconds/write
PASS 112: 03/UTC/utc-seconds/edit
PASS 113: 03/UTC/utc-milliseconds/write
PASS 114: 03/UTC/utc-milliseconds/edit
PASS 115: 03/UTC/utc-leap-date/write
PASS 116: 03/UTC/utc-leap-date/edit
PASS 117: 03/UTC/expired-utc/write
PASS 118: 03/UTC/expired-utc/edit
FAIL: {"group":"UTC","label":"year-only","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/year-only/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"year-only","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/year-only/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"date-only","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/date-only/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"date-only","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/date-only/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-Z","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-Z/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-Z","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-Z/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-Z-milliseconds","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-Z-milliseconds/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-Z-milliseconds","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-Z-milliseconds/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"non-UTC-offset","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/non-UTC-offset/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"non-UTC-offset","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/non-UTC-offset/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"zero-offset-without-Z","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/zero-offset-without-Z/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"zero-offset-without-Z","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/zero-offset-without-Z/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-seconds","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-seconds/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"missing-seconds","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/missing-seconds/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"impossible-Feb-30","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/impossible-Feb-30/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"impossible-Feb-30","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/impossible-Feb-30/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"non-leap-Feb-29","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/non-leap-Feb-29/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"non-leap-Feb-29","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/non-leap-Feb-29/edit\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"impossible-Apr-31","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/impossible-Apr-31/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"impossible-Apr-31","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/impossible-Apr-31/edit\n\nfalse !== true\n"}
PASS 119: 03/UTC/impossible-month-13/write
PASS 120: 03/UTC/impossible-month-13/edit
PASS 121: 03/UTC/impossible-day-zero/write
PASS 122: 03/UTC/impossible-day-zero/edit
FAIL: {"group":"UTC","label":"normalized-hour-24","toolName":"write","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/normalized-hour-24/write\n\nfalse !== true\n"}
FAIL: {"group":"UTC","label":"normalized-hour-24","toolName":"edit","target":"ordinary/new/allowed_test.go","expected":true,"actual":false,"message":"03/UTC/normalized-hour-24/edit\n\nfalse !== true\n"}
PASS 123: 03/UTC/impossible-minute-60/write
PASS 124: 03/UTC/impossible-minute-60/edit
GROUP dangling: {"total":32,"passed":16,"failed":16}
GROUP grant-zone: {"total":160,"passed":56,"failed":104}
GROUP UTC: {"total":36,"passed":14,"failed":22}
SUMMARY: original=38/38; repair=228; failures=142; passed=124; no built-in write/edit executed
node:internal/modules/run_main:123
    triggerUncaughtException(
    ^

AssertionError [ERR_ASSERTION]: repair boundary observations must all pass (aggregated failures above)

142 !== 0

    at runRepairMatrix (file:///Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs:213:10)
    at async file:///Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs:283:3 {
  generatedMessage: false,
  code: 'ERR_ASSERTION',
  actual: 142,
  expected: 0,
  operator: 'strictEqual'
}

Node.js v22.16.0
```

</details>
