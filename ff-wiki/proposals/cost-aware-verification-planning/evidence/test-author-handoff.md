# 合法测试作者交接

[可审阅测试补丁](test-authoring.patch) 是唯一源码提案，未应用。新增 Go 内容已通过 gofmt 的 stdin 解析/格式化，未写入任何测试路径；`git apply --check` 只检查补丁适用性，不算测试编译或运行。

宿主拒绝正常写入 `internal/command/cost_aware_assets_test.go`。主会话核实 `.pi/extensions/flowforge.ts:219-235` 的保护无角色判定，要求不关闭保护、不通过 bash/重命名/其他工具绕过。测试未落地，新测试 red/green 均未发生，不能输出 PLAN_READY。

| 拟测试/修正 | 断言与权威 | 现有接缝 |
| --- | --- | --- |
| `TestCostAwareRulePackageSourceEmbeddedDeployed` | design#surfaces / #verification：七个规则面源→嵌入和嵌入→新隔离部署分别逐字节相同；只证明交付一致性，不证明模型行为 | `embeddedAssets.ReadFile`、`deployManagedAssets`、真实文件读取 |
| `TestCostAwareSharedRunPreservesEvidenceGate` | design#schedule / #evidence、R8：同一实跑四元组可被两个 Change 引用；新 checked Fix 缺证据仍报 evidence-missing；schedule 正文不豁免原证据门禁 | `writeEvidenceTicket`、`DiscoverArtifacts`、`evidenceCodes`，有/无 schedule × 有/无缺证据 Fix 四组合 |
| `TestIssueTrackerDescribesRepairWorkflow` 路径 | 仅将 `filepath.Join("..", "..", "docs", "agents", "issue-tracker.md")` 改为 `filepath.Join("..", "..", "ff-wiki", "agents", "issue-tracker.md")`；保留原 `strings.Contains(string(data), "repair")` 及 Fatal 消息 | 原 `internal/command/assets_deploy_test.go` 用例；01 Change 8 明确授权 |

合法作者后续运行入口：

- `go test ./internal/command -run '^TestCostAwareRulePackageSourceEmbeddedDeployed$' -count=1 -v`：源未同步时预期共享契约、implement、review 三子项 red；这只是基线推导，须用实际输出确认。
- `go test ./internal/tracker -run '^TestCostAwareSharedRunPreservesEvidenceGate$' -count=1 -v`：预期四组合通过，当前尚未运行。
- `go test ./internal/command -run '^TestIssueTrackerDescribesRepairWorkflow$' -count=1 -v`：本轮已有旧路径失败实证，修正后的成功尚未发生。

恢复条件是主会话提供合法作者执行入口并实际落地/记录测试，再由 Refine 消费 01 的 Blocked evidence、验证命令/断言，按当前文件事实重算 frontier。结构校验通过不能替代这个条件。
