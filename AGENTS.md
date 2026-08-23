# FlowForge Agent 配置

> 本文档约束 Agent 如何在 FlowForge 本地优先 Issue Tracker & DAG 引擎体系下协作。

## Commands

- Build (当前平台): `make dev` 或 `go build -trimpath -o bin/flowforge ./cmd/flowforge`
- Build (所有平台): `make build` 或 `./scripts/build.sh <version> all`
- Test: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`
- Lint: `golangci-lint run ./...`

## 核心设计原则

1. **文件负责内容（零摩擦）**：
   - 所有 Spec、Ticket、Map 均直接通过 Markdown 文件读写（`.scratch/`）。
   - 严禁设计通过 CLI 传大文本长文本的 API。
2. **CLI 负责图计算（高确定性）**：
   - 使用 `flowforge frontier` 获取无阻塞就绪任务队列。
   - 使用 `flowforge check` 进行 DAG 依赖与循环死锁检查。
3. **方法论原装集成**：
   - 100% 采用 `mattpocock/skills`，不自造 prompt 敏捷规则。

## boundaries

- ✅ **Always**: 直接使用文件工具操作 `.scratch/` 下的 Markdown；使用 `flowforge frontier` 校验执行顺序；变更后运行 `go test ./internal/...`
- ⚠️ **Ask first**: 修改 Issue Schema 头规范、变更 CLI 接口签名
- 🚫 **Never**: 引入通过 CLI 传长文本的接口；在 `assets/` 中放不部署的内容

<!-- FLOWFORGE:START -->
## Agent skills

When asked to work on a feature, bug, refactor, or complex task, invoke the appropriate skill:

| Phase | Skill | Role & Responsibility |
|:---|:---|:---|
| **Triage** | `/triage` | Categorize, prioritize and prepare tasks into crisp actionable briefs |
| **Align** | `/grill-with-docs` | Relentless frontier grilling; inline sync with `CONTEXT.md` & `docs/adr/` |
| **Spec** | `/to-spec` | Synthesize current consensus into unambiguous specification (`spec.md`) |
| **Plan** | `/to-tickets` | Polymorphic vertical slicing with explicit DAG blocking edges (`issues/`) |
| **Implement** | `/implement` | TDD delivery on pre-agreed seams; close out with dual-axis code review |
| **Wayfinding** | `/wayfinder` | Fog-of-war decision mapping (`map.md`) for high-uncertainty efforts |
| **Review** | `/code-review` | Dual-axis (Standards vs Spec) parallel sub-agent code inspection |
| **Handoff** | `/handoff` | Compact session memory into cross-agent handoff artifact |
<!-- FLOWFORGE:END -->
