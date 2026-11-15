---
flowforge:
  schema: 1
  role: evidence
  id: cost-aware-verification-planning-03-native-authority
  revision: 1
---

# 03 修后 native 测试作者验收

**结论：NATIVE_VERIFIED**

本次是修后 `93f07a4dd60c85fb6721afa595699d0a17e4e4d5` 上重新执行的真实 native 工具链验收，不是旧 00 证据的复用。全程未改产品、Go 测试、模型或 guard 配置；未自行签发/撤销 grant，未用 bash 或其他工具绕过拒绝，未 push。主会话保留签发和撤销。

## 固定版本与身份

- HEAD：`93f07a4dd60c85fb6721afa595699d0a17e4e4d5`
- 实现提交标题：`fix(pi): close dangling-link, UTC expiry, and canonical grant-zone gaps`
- 本运行实际核对 source / embedded / deployed：

```text
ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c  assets/pi/flowforge.ts
ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c  internal/command/assets/pi/flowforge.ts
ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c  .pi/extensions/flowforge.ts
```

首次 `flowforge_test_authority` 原始结果：

```json
{
  "sessionId": "01a1ab74-027f-74bd-b8e4-953428b2aad3",
  "instanceId": "41bb5247-0360-4684-9877-39e9b72e1670",
  "projectRoot": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge",
  "grantPath": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.flowforge/test-author-grants/01a1ab74-027f-74bd-b8e4-953428b2aad3.json",
  "eligible": true,
  "grant": {
    "status": "missing",
    "role": "",
    "paths": [],
    "expiresAt": "",
    "reason": "no grant record for this session"
  }
}
```

## 1. 无 grant：正常 write/edit 均拒绝

验收前 `.flowforge/test-author-probe/allowed_test.go` 的实际内容为：

```go
package authorityprobe

// phase two
```

对该路径正常 `write`，目标内容 `package authorityprobe\n\n// phase one\n`，原始结果：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): no test-author grant for this session instance (no grant record for this session)
```

正常 `edit`（`// phase two` → `// phase one`）原始结果：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): no test-author grant for this session instance (no grant record for this session)
```

两次失败后读回仍为：

```go
package authorityprobe

// phase two
```

## 2. 主会话精确签发：正常 write/edit 成功

通过 `contact_supervisor(reason=need_decision)` 提交上述完整状态，请求仅对该实例和 `.flowforge/test-author-probe/allowed_test.go` 签发两小时 `plan` grant。主会话回复：

```text
已为该实际sessionId/instanceId签发2h、role=plan、唯一文件 .flowforge/test-author-probe/allowed_test.go 的授权记录，位置为你的grantPath。按既定native步骤继续正常write/edit、读回与越界/grant拒绝负例。完成后请求主会话撤销，不自行修改或删除grant，不写任何生产Go测试。
```

签发后 `flowforge_test_authority` 原始结果：

```json
{
  "sessionId": "01a1ab74-027f-74bd-b8e4-953428b2aad3",
  "instanceId": "41bb5247-0360-4684-9877-39e9b72e1670",
  "projectRoot": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge",
  "grantPath": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.flowforge/test-author-grants/01a1ab74-027f-74bd-b8e4-953428b2aad3.json",
  "eligible": true,
  "grant": {
    "status": "valid",
    "role": "plan",
    "paths": [
      ".flowforge/test-author-probe/allowed_test.go"
    ],
    "expiresAt": "2026-11-05T13:25:50.878Z",
    "reason": ""
  }
}
```

正常 `write` 的实际目标内容：

```go
package authorityprobe

// phase one
```

原始结果：

```text
Successfully wrote to .flowforge/test-author-probe/allowed_test.go
```

正常 `edit`（`// phase one` → `// phase two`）原始结果：

```text
Successfully replaced 1 block(s) in .flowforge/test-author-probe/allowed_test.go.
```

随后读回实际内容：

```go
package authorityprobe

// phase two
```

## 3. 未列出文件：正常 write/edit 均拒绝，目标不存在

对未授权的 `.flowforge/test-author-probe/denied_test.go` 正常 `write`，目标内容 `package authorityprobe\n\n// must not be created\n`，原始结果：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): the test-author grant does not list this file
```

正常 `edit`（`// must not be created` → `// must remain absent`）原始结果：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): the test-author grant does not list this file
```

两次失败后读回的原始结果：

```text
ENOENT: no such file or directory, access '/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.flowforge/test-author-probe/denied_test.go'
```

## 4. 授权记录：原字节 write 与同值字段 edit 均拒绝

尝试前正常 `read` 得到 grant 实际内容：

```json
{
  "version": 1,
  "sessionId": "01a1ab74-027f-74bd-b8e4-953428b2aad3",
  "instanceId": "41bb5247-0360-4684-9877-39e9b72e1670",
  "projectRoot": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge",
  "role": "plan",
  "paths": [
    ".flowforge/test-author-probe/allowed_test.go"
  ],
  "expiresAt": "2026-11-05T13:25:50.878Z"
}
```

将上述原字节（末尾换行相同）正常 `write` 回 grantPath，原始结果：

```text
protected test-authority record (.flowforge/test-author-grants): grant records are issued by the main session and a background session cannot write them, their directory, or its parents
```

正常 `edit`，同值字段替换 `"role": "plan"` → `"role": "plan"`，原始结果：

```text
protected test-authority record (.flowforge/test-author-grants): grant records are issued by the main session and a background session cannot write them, their directory, or its parents
```

失败后再次读回，逐字仍为上面的 JSON，字段和末尾换行未变。

## 5. 主会话撤销：原允许文件 write/edit 再拒绝

通过 `contact_supervisor(reason=need_decision)` 请求主会话撤销同一精确 grant。主会话回复：

```text
已独立核对allowed为phase two、denied不存在，并删除此精确grant文件；未修改其他配置。继续正常write/edit验证撤销后的拒绝，保留原始结果、读回和修后版本，提交03-native-authority.md后返回。不要删除临时测试fixture或更改其他文件，后续独立审查由既定workflow继续。
```

撤销后 `flowforge_test_authority` 原始结果：

```json
{
  "sessionId": "01a1ab74-027f-74bd-b8e4-953428b2aad3",
  "instanceId": "41bb5247-0360-4684-9877-39e9b72e1670",
  "projectRoot": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge",
  "grantPath": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.flowforge/test-author-grants/01a1ab74-027f-74bd-b8e4-953428b2aad3.json",
  "eligible": true,
  "grant": {
    "status": "missing",
    "role": "",
    "paths": [],
    "expiresAt": "",
    "reason": "no grant record for this session"
  }
}
```

再次正常 `write` 原允许文件，目标内容 `package authorityprobe\n\n// phase one\n`，原始结果：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): no test-author grant for this session instance (no grant record for this session)
```

再次正常 `edit`（`// phase two` → `// phase one`）原始结果：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): no test-author grant for this session instance (no grant record for this session)
```

撤销后两次失败均没有改变文件；最终读回：

```go
package authorityprobe

// phase two
```

临时 fixture 按主会话要求保留；真实 grant 已由主会话删除。

## 证据范围

本报告只证明上述真实 native `read`/`write`/`edit`/`flowforge_test_authority` 调用及主会话签发/撤销链。未授权、精确授权、未列出文件、grant 原字节及同值字段、撤销后拒绝均得到预期结果，未发现异常。

03 修后 source/deployed 内存边界矩阵各 266/266 的记录见 [03-implementation.md](03-implementation.md)（dangling 32、grant-zone 160、UTC 36，含原 38）。本 native 验收未逐项重跑该内存矩阵，不将其冒充 native 已执行项目，也不以旧 00 native 记录替代本次运行。
