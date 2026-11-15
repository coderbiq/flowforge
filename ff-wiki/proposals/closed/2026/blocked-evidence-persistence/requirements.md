---
flowforge:
  schema: 1
  role: requirement
  id: blocked-evidence-persistence-requirements
  revision: 1
---

<a id="blocked-evidence-persistence-requirements"></a>
# BLOCKED 证据工件化需求

## 问题

executor-loop-hardening 交付后，弱执行者的失败被及时截断（steps 预算 + fail-fast 契约 + 事后诊断），但失败知识的流向存在结构缺口：

1. `STATUS: BLOCKED` 是**会话返回约定**：逐字错误、已尝试路径、失败假设只存在于派发方对话里，不落任何工件。
2. 执行单元政策要求"跨上下文状态经工件传递"（ticket 文件 + STATUS 契约 + frontier），但重新派发的全新执行上下文对上次失败一无所知——除非派发方人工转录，这与政策的自洽性矛盾。
3. 失败知识没有结构化归宿：refine-ticket 的 Verified contracts 只承载成功契约事实，"什么命令以何种方式失败过"无处沉淀，重派发的新上下文可能原样重蹈覆辙（tangram-v2 票 06 的 290 次重试若有失败记录在前，本可在第 2 次尝试后换路）。
4. 派发方无法从 frontier/check 机械看出"这张票是被阻断过的、带失败上下文"——阻断票与全新就绪票在图上不可区分。

## 目标

1. 执行者阻断即落工件：返回 `STATUS: BLOCKED` 的同一动作必须先在票内追加 `## Blocked evidence` 小节（逐字错误、已尝试命令与退出码、下一步假设），失败知识随票持久化。
2. 阻断状态机器可见：open 票携带 `## Blocked evidence` → `flowforge check` 报 `blocked-evidence-present` warning；frontier 将其归入 ready_with_warnings（可派发但带上下文标记），派发方无需读对话即可识别阻断史。
3. 失败知识结构化再利用：flowforge-refine-ticket 消费 `## Blocked evidence`，把失败事实转写为 Verified contracts（何命令、为何失败、正确调用方式），随后移除该小节（瞬态状态消费即清除），诊断随之消失——重派发上下文从工件继承教训。

## 范围与约束

- 不新增 CLI 子命令与接口签名；诊断接入现有 check/frontier 管线与 warning/strict 语义。
- 不修改 Issue Schema 头规范；`## Blocked evidence` 是正文小节约定（与 evidence 四元组同层的 Markdown 标记）。
- 固化 prompt 的 Non-negotiables 增补第四句（block-and-record），既有三锚点不动。
- 消费循环（refine 转写+移除）是方法论层流程，不做 CLI 强制（诊断的存在已提供机器可见性；移除由 refine 流程完成）。
- BLOCKED 的会话返回通道保留（双通道：会话返回 + 工件持久化），工件是权威。

## 可观察验收

1. 构造 open 票含 `## Blocked evidence` 小节 → `flowforge check` 报 `blocked-evidence-present` warning；`--strict` 退出非零。
2. 同票 `**Status:** closed` 后不再报（阻断史在关票后不构成诊断）。
3. 小节被移除（消费）后诊断消失；refine 后重跑 check 恢复 clean。
4. `flowforge frontier` 将携带该小节的 open 就绪票归入 ready_with_warnings 展示，不进 clean ready。
5. `assets/skills/flowforge-implement/SKILL.md` 的 BLOCKED 路径含"先写 `## Blocked evidence` 再返回"的指令；结构断言锁定。
6. `assets/subagents/flowforge-implementer.md` 的 Non-negotiables 含第四句 block-and-record 锚点；部署产物同步；既有三锚点测试不受影响。
7. `assets/skills/flowforge-refine-ticket/SKILL.md` 含消费步骤（转写 Verified contracts + 移除小节 + 复跑 check）；结构断言锁定。
8. 存量票（无该小节）零影响；本仓自身 proposal 检查不因此产生新 warning。

## 术语

- **Blocked evidence**：执行者阻断时写入票内的结构化失败记录（逐字错误、已尝试命令与退出码、下一步假设），瞬态工件——被 refine 消费后移除。
- **消费（consumption）**：refine-ticket 把 Blocked evidence 的事实转写进 Verified contracts 并移除原小节的流程；失败知识从瞬态记录转为持久契约。
- **block-and-record**：固化 prompt 的第四句契约：返回 STATUS: BLOCKED 前必须先落 `## Blocked evidence`。
