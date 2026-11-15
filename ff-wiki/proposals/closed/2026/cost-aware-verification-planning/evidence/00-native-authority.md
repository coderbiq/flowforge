# 00 native 测试作者验收

结论：**NATIVE_VERIFIED（本文执行的 native 握手与精确路径矩阵）**。**00 整体仍 BLOCKED**：主会话另行报告三项内存反例，见末节；本文不关闭票、不替代独立审计。

## 权威与运行上下文

- 消费 requirements.md revision 2 / R10、design/test-authority.md revision 1、issues/00-test-authority.md。
- 独立新 native flowforge-planner；本轮未写生产代码、未应用 01 测试补丁、未改配置、未关闭 guard、未委派、未 push。
- cwd / canonical projectRoot：`/Users/qiangbi/develop/projects/Syl/tangram/flowforge`。
- 验收源码 HEAD：`7cb27d84b60ed3462f7b3582feefe45e8cdf107d`。
- 最后文件核对时刻：`2026-11-05T10:35:57.544427+00:00`。
- source `assets/pi/flowforge.ts`、embedded `internal/command/assets/pi/flowforge.ts`、deployed `.pi/extensions/flowforge.ts` 实际 SHA-256 均为 `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153`。哈希用于定位版本；native 允许/拒绝以以下实际工具结果证明。
- 正常工具发现：`searchTools('flowforge_test_authority', {limit:5})` 返回只读状态工具。随后直接调用，未从提示词推定资格。
- 开始及撤销后 `git status --short` 均为空；临时 `.flowforge/` 材料被忽略。

## 1. 未授权负例与真实身份

`functions.flowforge_test_authority({})` 原始结果：

```json
{
  "sessionId": "01a1ab42-6dd8-73b3-8558-6554da00c9ef",
  "instanceId": "87f4cabd-7d97-4ffa-985d-7f9034e31ea7",
  "projectRoot": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge",
  "grantPath": "/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.flowforge/test-author-grants/01a1ab42-6dd8-73b3-8558-6554da00c9ef.json",
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

正常 `functions.write`：path `.flowforge/test-author-probe/allowed_test.go`，content 精确为 `"package authorityprobe\n\n// phase one\n"`。原始结果：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): no test-author grant for this session instance (no grant record for this session)
```

通过 `contact_supervisor(reason="need_decision")` 发送上述完整身份、状态与原始拒绝，请主会话仅签发该实例、该文件的 2h plan grant，并等待答复。主会话回复：

> 已按用户授权为该真实session/instance签发2h、role=plan、唯一paths=[.flowforge/test-author-probe/allowed_test.go]记录，位于你报告的grantPath。现在用正常工具验证允许write/edit及读回；配对未列出文件和授权记录write/edit拒绝，保留原始结果。不改任何生产测试或保护配置。完成正负例后再contact_supervisor请求我撤销；撤销后确认允许文件也重新拒绝。

授权后、首次允许写之前，只读 Python `os.path.lexists` 核对结果：

```json
{
  ".flowforge/test-author-probe/allowed_test.go": false,
  ".flowforge/test-author-probe/denied_test.go": false
}
```

## 2. 授权 write/edit 正例

再次调用状态工具，sessionId、instanceId、projectRoot、grantPath、eligible 与第 1 节逐项相同；原始 grant 状态片段：

```json
{
  "status": "valid",
  "role": "plan",
  "paths": [".flowforge/test-author-probe/allowed_test.go"],
  "expiresAt": "2026-11-05T12:31:44.779Z",
  "reason": ""
}
```

只读 `read` 取得主会话签发记录供原字节回写负例使用；未自行签发或更改任何字段。本文不保存可复用的完整授权记录。

正常 `functions.write` 参数与第 1 节完全相同，原始结果：

```text
Successfully wrote to .flowforge/test-author-probe/allowed_test.go
```

正常 `functions.edit`：path 相同，`oldText="// phase one"`、`newText="// phase two"`。原始结果：

```text
Successfully replaced 1 block(s) in .flowforge/test-author-probe/allowed_test.go.
```

正常 `functions.read` 读回：

```go
package authorityprobe

// phase two
```

精确字符串为 `"package authorityprobe\n\n// phase two\n"`，包括末尾换行。

## 3. 授权范围外与自改授权负例

1. 正常 `functions.write`：path `.flowforge/test-author-probe/denied_test.go`，content `"package authorityprobe\n\n// denied probe\n"`。
2. 正常 `functions.edit`：相同 denied path，`oldText="// denied probe"`、`newText="// denied probe"`。文件不存在；guard 在底层文件操作之前拒绝，本步骤不声称底层 edit 能成功。

两次原始结果均为：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): the test-author grant does not list this file
```

3. 正常 `functions.write`：第 1 节绝对 grantPath，content 为刚读出的完整原字节（含原缩进和末尾换行），没有扩充字段、path 或有效期。
4. 正常 `functions.edit`：相同 grant 的项目相对路径，`oldText="\"version\": 1"`、`newText="\"version\": 1"`，同值尝试。

两次原始结果均为：

```text
protected test-authority record (.flowforge/test-author-grants): grant records are issued by the main session and a background session cannot write them, their directory, or its parents
```

正常 read 再次读取 grant，与最初 read 返回的字符串比较，工具结果：

```json
{"grantBytesUnchanged":true}
```

只读 Python 核对实际文件与 grant 哈希，原始 stdout：

```text
{
  ".flowforge/test-author-probe/allowed_test.go": {
    "exists": true,
    "content": "package authorityprobe\n\n// phase two\n"
  },
  ".flowforge/test-author-probe/denied_test.go": {
    "exists": false,
    "content": null
  }
}
grant sha256: 7be17a1da2ad1167098d58618a5e3ce2d2165a0049b642b91a164dd9a84b1bfd
```

## 4. 主会话撤销与撤销后负例

通过 `contact_supervisor(reason="need_decision")` 发送正负例原始结果、内容及未写目标核对，请主会话撤销并等待答复。主会话回复已独立读到 phase two、确认 denied 不存在，且已删除该精确 grant。并要求保留 fixture，由主会话处理。

随后正常状态工具返回值与第 1 节原始 JSON 完全相同：身份未变，`eligible=true`，`grant.status="missing"`，reason 为 `no grant record for this session`。

- 正常 `functions.write` 允许文件，content `"package authorityprobe\n\n// revoked write attempt\n"`。
- 正常 `functions.edit` 允许文件，`oldText="// phase two"`、`newText="// revoked edit attempt"`。

两次原始结果均为：

```text
protected test file (matched **/*_test.go, source: FlowForge defaults): no test-author grant for this session instance (no grant record for this session)
```

正常 read 仍读到 `"package authorityprobe\n\n// phase two\n"`。最后只读 `os.path.lexists` 结果：

```json
{
  ".flowforge/test-author-probe/allowed_test.go": true,
  ".flowforge/test-author-probe/denied_test.go": false,
  ".flowforge/test-author-grants/01a1ab42-6dd8-73b3-8558-6554da00c9ef.json": false
}
```

**授权已撤销**。未保留有效 grant 文件供复用；临时允许 fixture 由主会话后续清理。本验收所有受保护写尝试均为正常 write/edit；bash/Python 仅做只读核对，未 shell 写测试。

## 5. 主会话补充的独立内存反例：00 整体阻塞

以下是主会话在撤销回复中提供的实际执行报告，**不是本文作者执行的 native 步骤，本文未自行复跑**。记录来源为 `contact_supervisor` 回复；主会话负责将其完整原始执行记录交独立审计。

主会话报告源码版本 `7cb27d8`，使用真实 `loadExtensions` + `ExtensionRunner`，有效 plan grant `paths=['allowed_test.go']`：

- 创建该路径的 symlink，指向临时根的 `outside/missing_test.go`（目标不存在）；`emitToolCall(write)` 返回：

```text
{case:'dangling symlink to nonexistent outside target',expectedBlocked:true,actualBlocked:false}
```

- 移除 link 后，把 grant `expiresAt` 设置为 `'2099'`；返回：

```text
{case:'expiry is year-only, not UTC timestamp',expectedBlocked:true,actualBlocked:false}
```

主会话说明：前两例均未执行底层 write，临时 fixture 已删除；具体触发为 `resolveThroughLinks` 失败回退 ancestor 与 `Date.parse` 宽松校验。

主会话随后通过 mid-run steering 补充同一实际源码的第三项内存反例：临时 `project/.flowforge` 是指向临时 `outside` 的 symlink；背景 eligible 真实 loader/runner，未签发任何 grant，`emitToolCall(write, path=outside/test-author-grants/<当前sessionId>.json)` 返回：

```text
{case:'child writes grant via physical path when .flowforge is symlink',expectedBlocked:true,actualBlocked:false}
```

主会话说明第三例同样未执行底层 write，fixture 已删除；待独立审计核实的原因是 `grantDir` 为 lexical project 路径，resolved 候选没有与 canonical grantDir 比较。本文作者未复跑此例，不将原因假说视为已经完成独立审计。

三项报告涉及设计中的 symlink 拒绝、UTC 时间戳及授权目录别名保护约束，应一并交 workflow 审计/修复；本作者不修改生产实现或降低要求。

本 native 验收只证明缺 grant 拒绝、精确授权后 write/edit 成功、非列出路径拒绝、原字节自改授权拒绝、撤销后 write/edit 拒绝及实际目标状态。未在 native 会话验证过期、重启、并发环境或 symlink 组合；不得把本矩阵通过扩大为所有 R10 条款通过。00 不可依据本证据单独关闭。

## 6. 仓库规定的补充检查

- `git diff --check`：exit 0；该次运行时证据仍是 untracked，另在暂存后检查新增证据。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...`：本轮实际运行 exit 1。唯一失败测试 `TestIssueTrackerDescribesRepairWorkflow`：`assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory`；其他 internal 包通过或无测试。与票 00 已记录的 01 待修 docs 路径失败一致，未应用补丁，不声称全仓通过。
- `flowforge_frontier({feature:"cost-aware-verification-planning"})`：00 ready、01 waiting_on 00、02 waiting_on 01；02 带 execution-contract-incomplete，01 带 blocked-evidence-present，另有全局旧文档警告。这里只记录当前图，不改票状态，也不把图上的 ready 当作语义验收通过。
- 暂存检查发现 EOF 多余空行，已移除；重新运行 `git diff --cached --check` 验证格式。
