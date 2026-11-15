# 01 实施前 lint 基线与工具

记录者：主会话。时点：00/03完成、01尚未落地Go测试；Git基线 `16cce8e`（与 `2a924c1` 的 cmd/internal 无差分）。本记录不是 lint 通过证明。

## 实际结果

- `golangci-lint` 初始不在 PATH。主会话从官方 GitHub release 安装临时验证工具 `/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint`；`version` 输出 `2.15.0 built with go1.27.0 from 25c5369a on 2026-10-16T11:03:09Z`。下载归档 `golangci-lint-2.15.0-darwin-arm64.tar.gz` 的 SHA256 `d1e409481a8ca846c4bc97e781eb5cb1bd5c23353a0b3977a1cbda1b1b6447d5` 与同一官方release checksums一致。工具不进入仓库。
- `/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 45s ./...`：exit 1。`ff-wiki/proposals/documentation-contract-refinement/parser-prototype/` 中两个 throwaway Go 程序的 `envelope`、`main` 重复声明，且其中一个 `md.Areas` 不存在；typecheck阻断。不是01引入。
- 同工具 `run --timeout 45s ./cmd/... ./internal/...`：exit 1，默认输出限制下显示40项既有报告（errcheck26、ineffassign2、staticcheck7、unused5），不是完整无上限数量或通过证明。涉及command/config/update中的既有未检查清理/输出错误、无效赋值、静态简化和未用函数。没有01新测试文件，不能把这些输出归为新测试引入。
- `run --help` 确认 `--new-from-rev`、`--whole-files`、`--max-issues-per-linter`、`--max-same-issues` 可用；测试默认纳入分析。不使用 `--disable`、`--fix` 或改变 issues exit code。

## 主会话设计裁决

01仍须完整内部Go回归全部通过，旧docs路径失败由合法作者执行既有Change8修正，保留repair断言。lint采用覆盖所有新增/修改Go文件**整文件**的门禁，避免把本需求扩展为全仓历史清理：

```
/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...
```

该命令必须exit 0。作者测试必须已提交，执行者/审查者须检查Git差分列出的全部新增/修改Go文件都在cmd/internal范围内；不允许遗漏未跟踪新文件、关闭linter或忽略新失败。若实际写集新增其他Go包，返回权威扩大范围。基线通过与否、全仓失败和本次修改文件的结果分开记录，不能宣布全仓lint绿。此裁决在01 G1开始前由设计/规划所有者发布，不交执行者临场选测。

两个拟新增Go测试保护现有资产交付接缝和Catalog兼容性。00/03的make dev已消除旧资产漂移，它们落地后可能直接green；应记录实际基线，不故意破坏实现制造red，也不把旧docs失败叫作新增测试的red。新规则语义由01完整双轴与02真实受控情景证明。
