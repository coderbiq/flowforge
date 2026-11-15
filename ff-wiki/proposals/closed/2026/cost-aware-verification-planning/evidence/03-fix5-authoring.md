---
flowforge:
  schema: 1
  role: evidence
---

# 03 Fix 5：作者预设与真实 red

2026-11-05；基线 HEAD `b3836b7e61db837da82d8b04f089a2e8281193e7`。消费03 Fix 5及Review主会话裁决、requirements rev2、test-authority design rev1；按本次明确派发在本文件记录新增作者事实，旧 `03-refinement.md` 与Changes 1–4历史不变。Node v22.16.0 / Pi 1.5.0。

## 预设与保留证明

仅在 `00-test-authority-probe.mjs` 原矩阵汇总前追加标记为 `Fix 5 additions begin/end` 的22行，使用既有 `observe`/`expect`、真实 `loadExtensions` / `ExtensionRunner`。五个目标各write/edit，共10观察：

| 目标 | 应有blocked | source / deployed实际blocked（write/edit相同） |
| --- | --- | --- |
| `..allowed_test.go`（精确授权） | false | true：2 red |
| `..fixtures/allowed_test.go`（精确授权） | false | true：2 red |
| `fixtures/allowed_test.go`（精确授权） | false | false：2 pass |
| `.flowforge/test-author-grants/..unused.json` | true | false：2 red |
| `../outside_test.go` | true | true：2 pass |

Python实际读取 `git show HEAD:ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs` 与工作树文件，仅移除上述新增块后断言bytes完全相等，exit 0。因此旧266项、断言、执行顺序及聚合失败判定逐字节保留。`node --check ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs` exit 0。

所有新目标在事件前后均断言不存在，两个普通父目录仍未创建，fixture grant内容未变。仅在临时目录准备fixture；未执行底层write/edit，finally清理临时目录并恢复cwd/env。未写本项目Go测试或授权记录，未修改产品/配置，未构建或部署。

## 真实命令与结果

```sh
node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent
node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent .pi/extensions/flowforge.ts
```

两次均exit 1，真实输出（各276项，旧266全通过）：

```text
GROUP dangling: {"total":32,"passed":32,"failed":0}
GROUP grant-zone: {"total":160,"passed":160,"failed":0}
GROUP UTC: {"total":36,"passed":36,"failed":0}
GROUP fix5-path-segments: {"total":10,"passed":4,"failed":6}
SUMMARY: original=38/38; repair=238; failures=6; passed=270; no built-in write/edit executed
AssertionError [ERR_ASSERTION]: repair boundary observations must all pass (aggregated failures above)
6 !== 0
```

原始完整日志保存在本机 `/tmp/flowforge-fix5-author-source.log` 与 `/tmp/flowforge-fix5-author-deployed.log`；临时日志不作为持久仓库产物，本文件保存准确结果与失败目标。

`shasum -a 256` 实测：

| 文件 | SHA256 |
| --- | --- |
| 基线probe | `970f7b68881b8a18a3be627192b6537cf8aba584df806082072cb79c2227a283` |
| 追加后probe | `8e28ffd552864a696e705860241cdaeb7aa60e8eba77eb1ad2cc9a258873c0f9` |
| `assets/pi/flowforge.ts` | `ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c` |
| `internal/command/assets/pi/flowforge.ts` | `ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c` |
| `.pi/extensions/flowforge.ts` | `ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c` |
| source日志 | `196f91d5a8623a01548830d808c6e8449d229dca560bb3e5998f1194b80b164d` |
| deployed日志 | `13d69b056c81b83919dbafeb9f15e91eb414743dc29814dac1b424b2fb91757a` |

`cmp assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts` 与 `cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts` 均exit 0。这是当前文件一致性，实际行为结果来自上述两次真实loader/runner，不作为新的native工具写入验收。

`GOPROXY=https://goproxy.cn,direct go test ./internal/...` 实跑exit 1；唯一失败为既有 `TestIssueTrackerDescribesRepairWorkflow`：`assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory`，仍归01 Change 8；其他包通过或无测试。日志 `/tmp/flowforge-fix5-author-go.log`。不宣称全仓绿。

`./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full` 与 `./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json` 均exit 0。03无execution-contract gap；01既有blocked-evidence warning、02既有execution-contract gap保留，frontier仅03，03→01→02阻塞链不变。`git diff --check` exit 0。

## 后继

Fix 5保持未完成；新鲜实现上下文仅按既定路径段判定Change执行，不修改probe。修后source/deployed各须276/276、0失败、exit 0；构建/部署、Go回归、native历史复用条件及完整范围独立审查仍按03合同。主会话保留收敛，01/02尚未交付。
